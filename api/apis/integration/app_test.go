package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"code.cloudfoundry.org/cf-k8s-controllers/api/repositories/fake"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	hnsv1alpha2 "sigs.k8s.io/hierarchical-namespaces/api/v1alpha2"

	workloads "code.cloudfoundry.org/cf-k8s-controllers/controllers/apis/workloads/v1alpha1"

	"code.cloudfoundry.org/cf-k8s-controllers/api/actions"
	. "code.cloudfoundry.org/cf-k8s-controllers/api/apis"
	"code.cloudfoundry.org/cf-k8s-controllers/api/payloads"
	"code.cloudfoundry.org/cf-k8s-controllers/api/repositories"
)

var _ = FDescribe("App Handler", func() {
	var (
		apiHandler *AppHandler
		org, space *hnsv1alpha2.SubnamespaceAnchor
		spaceGUID  string
	)

	BeforeEach(func() {
		appRepo := repositories.NewAppRepo(k8sClient, clientFactory, nsPermissions)
		dropletRepo := repositories.NewDropletRepo(k8sClient, clientFactory)
		processRepo := repositories.NewProcessRepo(k8sClient, clientFactory)
		routeRepo := repositories.NewRouteRepo(k8sClient, clientFactory)
		domainRepo := repositories.NewDomainRepo(k8sClient, clientFactory)
		podRepo := repositories.NewPodRepo(clientFactory, new(fake.MetricsFetcherFn).Spy)
		orgRepo := repositories.NewOrgRepo(rootNamespace, k8sClient, clientFactory, nsPermissions, time.Minute, true)
		scaleProcess := actions.NewScaleProcess(processRepo).Invoke
		scaleAppProcess := actions.NewScaleAppProcess(appRepo, processRepo, scaleProcess).Invoke
		decoderValidator, err := NewDefaultDecoderValidator()
		Expect(err).NotTo(HaveOccurred())

		apiHandler = NewAppHandler(
			logf.Log.WithName("integration tests"),
			*serverURL,
			appRepo,
			dropletRepo,
			processRepo,
			routeRepo,
			domainRepo,
			podRepo,
			orgRepo,
			scaleAppProcess,
			decoderValidator,
		)
		apiHandler.RegisterRoutes(router)

		org = createOrgAnchorAndNamespace(ctx, rootNamespace, generateGUID())
		space = createSpaceAnchorAndNamespace(ctx, org.Name, "spacename-"+generateGUID())
		spaceGUID = space.Name
	})

	Describe("POST /v3/apps endpoint", func() {
		When("on the happy path", func() {
			const (
				appName = "my-test-app"
			)
			var testEnvironmentVariables map[string]string

			BeforeEach(func() {
				createRoleBinding(ctx, userName, spaceDeveloperRole.Name, spaceGUID)

				testEnvironmentVariables = map[string]string{"foo": "foo", "bar": "bar"}
				envJSON, _ := json.Marshal(&testEnvironmentVariables)
				requestBody := fmt.Sprintf(`{
                    "name": %q,
                    "relationships": {
                        "space": {
                            "data": {
                                "guid": %q
                            }
                        }
                    },
                    "environment_variables": %s
                }`, appName, spaceGUID, envJSON)

				var err error
				req, err = http.NewRequestWithContext(ctx, "POST", serverURI("/v3/apps"), strings.NewReader(requestBody))
				Expect(err).NotTo(HaveOccurred())

				req.Header.Add("Content-type", "application/json")
			})

			JustBeforeEach(func() {
				router.ServeHTTP(rr, req)
			})

			It("creates a CFApp and Secret, returns 201 and an App object as JSON", func() {
				Expect(rr.Code).To(Equal(201))

				var parsedBody map[string]interface{}
				body, err := ioutil.ReadAll(rr.Body)
				Expect(err).NotTo(HaveOccurred())

				Expect(
					json.Unmarshal(body, &parsedBody),
				).To(Succeed())

				Expect(parsedBody).To(MatchKeys(IgnoreExtras, Keys{
					"guid":       Not(BeEmpty()),
					"name":       Equal("my-test-app"),
					"state":      Equal("STOPPED"),
					"created_at": Not(BeEmpty()),
					"relationships": Equal(map[string]interface{}{
						"space": map[string]interface{}{
							"data": map[string]interface{}{
								"guid": spaceGUID,
							},
						},
					}),
				}))

				appGUID := parsedBody["guid"].(string)
				appNSName := types.NamespacedName{
					Name:      appGUID,
					Namespace: spaceGUID,
				}
				var appRecord workloads.CFApp
				Eventually(func() error {
					return k8sClient.Get(ctx, appNSName, &appRecord)
				}).Should(Succeed())

				Expect(appRecord.Spec.Name).To(Equal("my-test-app"))
				Expect(appRecord.Spec.DesiredState).To(BeEquivalentTo("STOPPED"))
				Expect(appRecord.Spec.EnvSecretName).NotTo(BeEmpty())

				secretNSName := types.NamespacedName{
					Name:      appRecord.Spec.EnvSecretName,
					Namespace: spaceGUID,
				}
				var secretCR corev1.Secret
				Eventually(func() error {
					return k8sClient.Get(ctx, secretNSName, &secretCR)
				}).Should(Succeed())

				Expect(secretCR.Data).To(MatchAllKeys(Keys{
					"foo": BeEquivalentTo(testEnvironmentVariables["foo"]),
					"bar": BeEquivalentTo(testEnvironmentVariables["bar"]),
				}))
			})
		})
	})

	Describe("app sub-resources", func() {
		var (
			app         *workloads.CFApp
			dropletGUID string
		)

		BeforeEach(func() {
			appGUID := generateGUID()
			dropletGUID = generateGUID()

			app = &workloads.CFApp{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appGUID,
					Namespace: spaceGUID,
				},
				Spec: workloads.CFAppSpec{
					Name:         generateGUID(),
					DesiredState: "STOPPED",
					Lifecycle: workloads.Lifecycle{
						Type: "buildpack",
					},
					CurrentDropletRef: corev1.LocalObjectReference{
						Name: dropletGUID,
					},
				},
			}

			Expect(k8sClient.Create(ctx, app)).To(Succeed())
		})

		Describe("get processes", func() {
			JustBeforeEach(func() {
				var err error
				req, err = http.NewRequestWithContext(ctx, http.MethodGet, serverURI("/v3/apps/"+app.Name+"/processes"), nil)
				Expect(err).NotTo(HaveOccurred())

				router.ServeHTTP(rr, req)
			})

			When("the user is not authorized in the space", func() {
				It("returns a not found status", func() {
					Expect(rr).To(HaveHTTPStatus(http.StatusNotFound))
				})
			})

			When("the user is a space developer", func() {
				BeforeEach(func() {
					createRoleBinding(ctx, userName, spaceDeveloperRole.Name, spaceGUID)
				})

				It("returns the (empty) process list", func() {
					Expect(rr).To(HaveHTTPStatus(http.StatusOK))
					Expect(rr).To(HaveHTTPBody(ContainSubstring(`"resources":[]`)), rr.Body.String())
				})
			})
		})

		Describe("get routes", func() {
			JustBeforeEach(func() {
				var err error
				req, err = http.NewRequestWithContext(ctx, http.MethodGet, serverURI("/v3/apps/"+app.Name+"/routes"), nil)
				Expect(err).NotTo(HaveOccurred())

				router.ServeHTTP(rr, req)
			})

			When("the user is not authorized in the space", func() {
				It("returns a not found status", func() {
					Expect(rr).To(HaveHTTPStatus(http.StatusNotFound))
					Expect(rr).To(HaveHTTPBody(ContainSubstring("App not found")), rr.Body.String())
				})
			})

			When("the user is a space developer", func() {
				BeforeEach(func() {
					createRoleBinding(ctx, userName, spaceDeveloperRole.Name, spaceGUID)
				})

				It("returns the (empty) route list", func() {
					Expect(rr).To(HaveHTTPStatus(http.StatusOK))
					Expect(rr).To(HaveHTTPBody(ContainSubstring(`"resources":[]`)), rr.Body.String())
				})
			})
		})

		Describe("restart app", func() {
			JustBeforeEach(func() {
				var err error
				req, err = http.NewRequestWithContext(ctx, http.MethodPost, serverURI("/v3/apps/"+app.Name+"/actions/restart"), nil)
				Expect(err).NotTo(HaveOccurred())

				router.ServeHTTP(rr, req)
			})

			When("the user is not authorized in the space", func() {
				It("returns a not found status", func() {
					Expect(rr).To(HaveHTTPStatus(http.StatusNotFound))
					Expect(rr).To(HaveHTTPBody(ContainSubstring("App not found")), rr.Body.String())
				})
			})

			When("the user has readonly access to the app", func() {
				BeforeEach(func() {
					createRoleBinding(ctx, userName, spaceManagerRole.Name, spaceGUID)
				})

				It("returns a forbidden error", func() {
					Expect(rr).To(HaveHTTPStatus(http.StatusForbidden))
				})
			})

			When("the user is a space developer", func() {
				BeforeEach(func() {
					createRoleBinding(ctx, userName, spaceDeveloperRole.Name, spaceGUID)
				})

				It("restarts the app", func() {
					Expect(rr).To(HaveHTTPStatus(http.StatusOK))
					Expect(rr).To(HaveHTTPBody(ContainSubstring(`"state":"STARTED"`)), rr.Body.String())
				})
			})
		})

		Describe("stop app", func() {
			BeforeEach(func() {
				stoppedApp := app.DeepCopy()
				stoppedApp.Spec.DesiredState = "STARTED"
				Expect(k8sClient.Patch(ctx, stoppedApp, client.MergeFrom(app))).To(Succeed())
			})

			JustBeforeEach(func() {
				var err error
				req, err = http.NewRequestWithContext(ctx, http.MethodPost, serverURI("/v3/apps/"+app.Name+"/actions/stop"), nil)
				Expect(err).NotTo(HaveOccurred())

				router.ServeHTTP(rr, req)
			})

			When("the user is not authorized in the space", func() {
				It("returns a not found status", func() {
					Expect(rr).To(HaveHTTPStatus(http.StatusNotFound))
					Expect(rr).To(HaveHTTPBody(ContainSubstring("App not found")), rr.Body.String())
				})
			})

			When("the user has readonly access to the app", func() {
				BeforeEach(func() {
					createRoleBinding(ctx, userName, spaceManagerRole.Name, spaceGUID)
				})

				It("returns a forbidden error", func() {
					Expect(rr).To(HaveHTTPStatus(http.StatusForbidden))
				})
			})

			When("the user is a space developer", func() {
				BeforeEach(func() {
					createRoleBinding(ctx, userName, spaceDeveloperRole.Name, spaceGUID)
				})

				It("stops the app", func() {
					Expect(rr).To(HaveHTTPStatus(http.StatusOK))
					Expect(rr).To(HaveHTTPBody(ContainSubstring(`"state":"STOPPED"`)), rr.Body.String())
				})
			})
		})

		Describe("droplets", func() {
			var droplet *workloads.CFBuild

			BeforeEach(func() {
				droplet = &workloads.CFBuild{
					ObjectMeta: metav1.ObjectMeta{
						Name:      dropletGUID,
						Namespace: spaceGUID,
					},
					Spec: workloads.CFBuildSpec{
						AppRef: corev1.LocalObjectReference{
							Name: app.Name,
						},
						Lifecycle: workloads.Lifecycle{
							Type: "buildpack",
						},
					},
				}
				Expect(k8sClient.Create(ctx, droplet)).To(Succeed())
				droplet.Status = workloads.CFBuildStatus{
					Conditions: []metav1.Condition{
						{
							Type:               "Staging",
							Status:             metav1.ConditionFalse,
							Reason:             "foo",
							LastTransitionTime: metav1.NewTime(time.Now()),
						},
						{
							Type:               "Succeeded",
							Status:             metav1.ConditionTrue,
							Reason:             "foo",
							LastTransitionTime: metav1.NewTime(time.Now()),
						},
					},
					BuildDropletStatus: &workloads.BuildDropletStatus{
						ProcessTypes: []workloads.ProcessType{},
						Ports:        []int32{},
					},
				}
				Expect(k8sClient.Status().Update(ctx, droplet)).To(Succeed())
			})

			Describe("get current droplet", func() {
				JustBeforeEach(func() {
					var err error
					req, err = http.NewRequestWithContext(ctx, http.MethodGet, serverURI("/v3/apps/"+app.Name+"/droplets/current"), nil)
					Expect(err).NotTo(HaveOccurred())

					router.ServeHTTP(rr, req)
				})

				When("having the space developer role", func() {
					BeforeEach(func() {
						createRoleBinding(ctx, userName, spaceDeveloperRole.Name, spaceGUID)
					})

					It("gets the droplet", func() {
						Expect(rr).To(HaveHTTPStatus(http.StatusOK))
					})
				})

				When("not authorized to get app", func() {
					It("returns a 404", func() {
						Expect(rr).To(HaveHTTPStatus(http.StatusNotFound))
						Expect(rr).To(HaveHTTPBody(ContainSubstring("App not found")))
					})
				})
			})

			Describe("set current droplet", func() {
				var payload payloads.AppSetCurrentDroplet

				BeforeEach(func() {
					payload = payloads.AppSetCurrentDroplet{
						Relationship: payloads.Relationship{
							Data: &payloads.RelationshipData{
								GUID: droplet.Name,
							},
						},
					}
				})

				JustBeforeEach(func() {
					payloadJSON, err := json.Marshal(payload)
					Expect(err).NotTo(HaveOccurred())

					req, err = http.NewRequestWithContext(ctx, http.MethodPatch, serverURI("/v3/apps/"+app.Name+"/relationships/current_droplet"), bytes.NewReader(payloadJSON))
					Expect(err).NotTo(HaveOccurred())

					router.ServeHTTP(rr, req)
				})

				When("having the space developer role", func() {
					BeforeEach(func() {
						createRoleBinding(ctx, userName, spaceDeveloperRole.Name, spaceGUID)
					})

					It("sets the droplet", func() {
						Expect(rr).To(HaveHTTPStatus(http.StatusOK))
					})
				})

				When("no access to app", func() {
					It("returns a 404", func() {
						Expect(rr).To(HaveHTTPStatus(http.StatusNotFound))
						Expect(rr).To(HaveHTTPBody(ContainSubstring("App not found")), rr.Body.String())
					})
				})

				When("access to app but no write permissions", func() {
					BeforeEach(func() {
						createRoleBinding(ctx, userName, spaceManagerRole.Name, spaceGUID)
					})

					It("returns a 403", func() {
						Expect(rr).To(HaveHTTPStatus(http.StatusForbidden))
						Expect(rr).To(HaveHTTPBody(ContainSubstring("CF-NotAuthorized")), rr.Body.String())
					})
				})

				When("the droplet does not exist", func() {
					BeforeEach(func() {
						createRoleBinding(ctx, userName, spaceDeveloperRole.Name, spaceGUID)
						payload.Data.GUID = "not-a-real-guid"
					})

					It("returns unprocessable entity", func() {
						Expect(rr).To(HaveHTTPStatus(http.StatusUnprocessableEntity))
					})
				})
			})
		})
	})
})

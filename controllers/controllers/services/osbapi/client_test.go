package osbapi_test

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"code.cloudfoundry.org/korifi/controllers/controllers/services/osbapi"
	"code.cloudfoundry.org/korifi/model/services"
	"code.cloudfoundry.org/korifi/tests/helpers/broker"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("OSBAPI Client", func() {
	var (
		brokerClient *osbapi.Client
		brokerServer *broker.BrokerServer
	)

	BeforeEach(func() {
		brokerServer = broker.NewServer()
	})

	JustBeforeEach(func() {
		brokerServer.Start()
		DeferCleanup(func() {
			brokerServer.Stop()
		})

		brokerClient = osbapi.NewClient(osbapi.Broker{
			URL:      brokerServer.URL(),
			Username: "broker-user",
			Password: "broker-password",
		}, &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //#nosec G402
		}})
	})

	Describe("GetCatalog", func() {
		var (
			catalog       osbapi.Catalog
			getCatalogErr error
		)

		BeforeEach(func() {
			brokerServer.WithResponse(
				"/v2/catalog",
				map[string]any{
					"services": []map[string]any{
						{
							"id":          "123456",
							"name":        "test-service",
							"description": "test service description",
							"bindable":    true,
						},
					},
				},
				http.StatusOK,
			)
		})

		JustBeforeEach(func() {
			catalog, getCatalogErr = brokerClient.GetCatalog(ctx)
		})

		It("gets the catalog", func() {
			Expect(getCatalogErr).NotTo(HaveOccurred())
			Expect(catalog).To(Equal(osbapi.Catalog{
				Services: []osbapi.Service{{
					ID:          "123456",
					Name:        "test-service",
					Description: "test service description",
					BrokerCatalogFeatures: services.BrokerCatalogFeatures{
						Bindable: true,
					},
				}},
			}))
		})

		It("sends a sync request", func() {
			servedRequests := brokerServer.ServedRequests()
			Expect(servedRequests).To(HaveLen(1))
			Expect(servedRequests[0].Method).To(Equal(http.MethodGet))
			Expect(servedRequests[0].URL.Query().Get("accepts_incomplete")).To(BeEmpty())
		})

		It("sends broker credentials in the Authorization request header", func() {
			Expect(brokerServer.ServedRequests()).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
				"Header": HaveKeyWithValue(
					"Authorization", ConsistOf("Basic "+base64.StdEncoding.EncodeToString([]byte("broker-user:broker-password"))),
				),
			}))))
		})

		It("sends OSBAPI version request header", func() {
			Expect(brokerServer.ServedRequests()).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
				"Header": HaveKeyWithValue(
					"X-Broker-Api-Version", ConsistOf("2.17"),
				),
			}))))
		})

		When("getting the catalog fails", func() {
			BeforeEach(func() {
				brokerServer = brokerServer.WithResponse("/v2/catalog", nil, http.StatusTeapot)
			})

			It("returns an error", func() {
				Expect(getCatalogErr).To(MatchError(ContainSubstring(strconv.Itoa(http.StatusTeapot))))
			})
		})
	})

	Describe("Instances", func() {
		Describe("Provision", func() {
			var (
				provisionResp osbapi.ServiceInstanceOperationResponse
				provisionErr  error
			)

			BeforeEach(func() {
				brokerServer = brokerServer.WithResponse(
					"/v2/service_instances/{id}",
					nil,
					http.StatusCreated,
				)
			})

			JustBeforeEach(func() {
				provisionResp, provisionErr = brokerClient.Provision(ctx, osbapi.InstanceProvisionPayload{
					InstanceID: "my-service-instance",
					InstanceProvisionRequest: osbapi.InstanceProvisionRequest{
						ServiceId: "service-guid",
						PlanID:    "plan-guid",
						SpaceGUID: "space-guid",
						OrgGUID:   "org-guid",
						Parameters: map[string]any{
							"foo": "bar",
						},
					},
				})
			})

			It("sends async provision request to broker", func() {
				Expect(provisionErr).NotTo(HaveOccurred())
				requests := brokerServer.ServedRequests()

				Expect(requests).To(HaveLen(1))

				Expect(requests[0].Method).To(Equal(http.MethodPut))
				Expect(requests[0].URL.Path).To(Equal("/v2/service_instances/my-service-instance"))

				Expect(requests[0].URL.Query().Get("accepts_incomplete")).To(Equal("true"))
			})

			It("sends correct request body", func() {
				Expect(provisionErr).NotTo(HaveOccurred())
				requests := brokerServer.ServedRequests()

				Expect(requests).To(HaveLen(1))

				requestBytes, err := io.ReadAll(requests[0].Body)
				Expect(err).NotTo(HaveOccurred())
				requestBody := map[string]any{}
				Expect(json.Unmarshal(requestBytes, &requestBody)).To(Succeed())

				Expect(requestBody).To(MatchAllKeys(Keys{
					"service_id":        Equal("service-guid"),
					"plan_id":           Equal("plan-guid"),
					"space_guid":        Equal("space-guid"),
					"organization_guid": Equal("org-guid"),
					"parameters": MatchAllKeys(Keys{
						"foo": Equal("bar"),
					}),
				}))
			})

			It("provisions the service synchronously", func() {
				Expect(provisionErr).NotTo(HaveOccurred())
				Expect(provisionResp).To(Equal(osbapi.ServiceInstanceOperationResponse{}))
			})

			When("the broker accepts the provision request", func() {
				BeforeEach(func() {
					brokerServer = brokerServer.WithResponse(
						"/v2/service_instances/{id}",
						map[string]any{
							"operation": "provision_op1",
						},
						http.StatusAccepted,
					)
				})

				It("provisions the service asynchronously", func() {
					Expect(provisionErr).NotTo(HaveOccurred())
					Expect(provisionResp).To(Equal(osbapi.ServiceInstanceOperationResponse{
						IsAsync:   true,
						Operation: "provision_op1",
					}))
				})
			})

			When("the provision request fails", func() {
				BeforeEach(func() {
					brokerServer = brokerServer.WithResponse("/v2/service_instances/{id}", nil, http.StatusTeapot)
				})

				It("returns an error", func() {
					Expect(provisionErr).To(MatchError(ContainSubstring("provision request failed")))
				})
			})
		})

		Describe("Deprovision", func() {
			var (
				deprovisionResp osbapi.ServiceInstanceOperationResponse
				deprovisionErr  error
			)

			BeforeEach(func() {
				brokerServer.WithResponse(
					"/v2/service_instances/{id}",
					map[string]any{
						"operation": "provision_op1",
					},
					http.StatusOK,
				)
			})

			JustBeforeEach(func() {
				deprovisionResp, deprovisionErr = brokerClient.Deprovision(ctx, osbapi.InstanceDeprovisionPayload{
					ID: "my-service-instance",
					InstanceDeprovisionRequest: osbapi.InstanceDeprovisionRequest{
						ServiceId: "service-guid",
						PlanID:    "plan-guid",
					},
				})
			})

			It("deprovisions the service", func() {
				Expect(deprovisionErr).NotTo(HaveOccurred())
				Expect(deprovisionResp).To(Equal(osbapi.ServiceInstanceOperationResponse{
					Operation: "provision_op1",
				}))
			})

			It("sends async deprovision request to broker", func() {
				Expect(deprovisionErr).NotTo(HaveOccurred())
				requests := brokerServer.ServedRequests()

				Expect(requests).To(HaveLen(1))

				Expect(requests[0].Method).To(Equal(http.MethodDelete))
				Expect(requests[0].URL.Path).To(Equal("/v2/service_instances/my-service-instance"))

				Expect(requests[0].URL.Query().Get("accepts_incomplete")).To(Equal("true"))
			})

			It("sends correct request body", func() {
				Expect(deprovisionErr).NotTo(HaveOccurred())
				requests := brokerServer.ServedRequests()

				Expect(requests).To(HaveLen(1))

				requestBytes, err := io.ReadAll(requests[0].Body)
				Expect(err).NotTo(HaveOccurred())
				requestBody := map[string]any{}
				Expect(json.Unmarshal(requestBytes, &requestBody)).To(Succeed())

				Expect(requestBody).To(MatchAllKeys(Keys{
					"service_id": Equal("service-guid"),
					"plan_id":    Equal("plan-guid"),
				}))
			})

			When("the deprovision request fails", func() {
				BeforeEach(func() {
					brokerServer = brokerServer.WithResponse(
						"/v2/service_instances/{id}",
						nil,
						http.StatusTeapot,
					)
				})

				It("returns an error", func() {
					Expect(deprovisionErr).To(MatchError(ContainSubstring("deprovision request failed")))
				})
			})
		})

		Describe("GetServiceInstanceLastOperation", func() {
			var (
				lastOpResp           osbapi.LastOperationResponse
				lastOpErr            error
				lastOperationRequest osbapi.GetServiceInstanceLastOperationRequest
			)

			BeforeEach(func() {
				brokerServer.WithResponse(
					"/v2/service_instances/{id}/last_operation",
					map[string]any{
						"state":       "in-progress",
						"description": "provisioning",
					},
					http.StatusOK,
				)

				lastOperationRequest = osbapi.GetServiceInstanceLastOperationRequest{
					InstanceID: "my-service-instance",
					GetLastOperationRequestParameters: osbapi.GetLastOperationRequestParameters{
						ServiceId: "service-guid",
						PlanID:    "plan-guid",
						Operation: "op-guid",
					},
				}
			})

			JustBeforeEach(func() {
				lastOpResp, lastOpErr = brokerClient.GetServiceInstanceLastOperation(ctx, lastOperationRequest)
			})

			It("gets the last operation", func() {
				Expect(lastOpErr).NotTo(HaveOccurred())
				Expect(lastOpResp).To(Equal(osbapi.LastOperationResponse{
					State:       "in-progress",
					Description: "provisioning",
				}))
			})

			It("sends correct request to broker", func() {
				Expect(lastOpErr).NotTo(HaveOccurred())
				requests := brokerServer.ServedRequests()

				Expect(requests).To(HaveLen(1))

				Expect(requests[0].Method).To(Equal(http.MethodGet))
				Expect(requests[0].URL.Path).To(Equal("/v2/service_instances/my-service-instance/last_operation"))
				Expect(requests[0].URL.Query()).To(BeEquivalentTo(map[string][]string{
					"service_id": {"service-guid"},
					"plan_id":    {"plan-guid"},
					"operation":  {"op-guid"},
				}))

				requestBytes, err := io.ReadAll(requests[0].Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(requestBytes).To(BeEmpty())
			})

			When("request parameters are not specified", func() {
				BeforeEach(func() {
					lastOperationRequest = osbapi.GetServiceInstanceLastOperationRequest{
						InstanceID: "my-service-instance",
					}
				})

				It("does not specify http request query parameters", func() {
					Expect(lastOpErr).NotTo(HaveOccurred())
					requests := brokerServer.ServedRequests()

					Expect(requests).To(HaveLen(1))
					Expect(requests[0].URL.Query()).To(BeEmpty())
				})
			})

			When("getting the last operation request fails", func() {
				BeforeEach(func() {
					brokerServer = brokerServer.WithResponse(
						"/v2/service_instances/{id}/last_operation",
						nil,
						http.StatusTeapot,
					)
				})

				It("returns an error", func() {
					Expect(lastOpErr).To(MatchError(ContainSubstring("last operation request failed")))
				})
			})

			When("getting the last operation request fails with 410 Gone", func() {
				BeforeEach(func() {
					brokerServer = brokerServer.WithResponse(
						"/v2/service_instances/{id}/last_operation",
						nil,
						http.StatusGone,
					)
				})

				It("returns a gone error", func() {
					Expect(lastOpErr).To(BeAssignableToTypeOf(osbapi.GoneError{}))
				})
			})
		})
	})

	Describe("Bindings", func() {
		Describe("Bind", func() {
			var (
				bindResp osbapi.BindResponse
				bindErr  error
			)

			BeforeEach(func() {
				brokerServer.WithResponse(
					"/v2/service_instances/{instance_id}/service_bindings/{binding_id}",
					map[string]any{
						"credentials": map[string]string{
							"foo": "bar",
						},
					},
					http.StatusCreated,
				)
			})

			JustBeforeEach(func() {
				bindResp, bindErr = brokerClient.Bind(ctx, osbapi.BindPayload{
					InstanceID: "instance-id",
					BindingID:  "binding-id",
					BindRequest: osbapi.BindRequest{
						ServiceId: "service-guid",
						PlanID:    "plan-guid",
						AppGUID:   "app-guid",
						BindResource: osbapi.BindResource{
							AppGUID: "app-guid",
						},
						Parameters: map[string]any{
							"foo": "bar",
						},
					},
				})
			})

			It("sends async bind request to broker", func() {
				Expect(bindErr).NotTo(HaveOccurred())
				requests := brokerServer.ServedRequests()

				Expect(requests).To(HaveLen(1))

				Expect(requests[0].Method).To(Equal(http.MethodPut))
				Expect(requests[0].URL.Path).To(Equal("/v2/service_instances/instance-id/service_bindings/binding-id"))

				Expect(requests[0].URL.Query().Get("accepts_incomplete")).To(Equal("true"))
			})

			It("sends correct request to broker", func() {
				Expect(bindErr).NotTo(HaveOccurred())
				requests := brokerServer.ServedRequests()

				Expect(requests).To(HaveLen(1))

				Expect(requests[0].Method).To(Equal(http.MethodPut))
				Expect(requests[0].URL.Path).To(Equal("/v2/service_instances/instance-id/service_bindings/binding-id"))

				requestBytes, err := io.ReadAll(requests[0].Body)
				Expect(err).NotTo(HaveOccurred())
				requestBody := map[string]any{}
				Expect(json.Unmarshal(requestBytes, &requestBody)).To(Succeed())

				Expect(requestBody).To(MatchAllKeys(Keys{
					"service_id": Equal("service-guid"),
					"plan_id":    Equal("plan-guid"),
					"app_guid":   Equal("app-guid"),
					"bind_resource": MatchAllKeys(Keys{
						"app_guid": Equal("app-guid"),
					}),
					"parameters": MatchAllKeys(Keys{
						"foo": Equal("bar"),
					}),
				}))
			})

			It("binds the service", func() {
				Expect(bindErr).NotTo(HaveOccurred())
				Expect(bindResp).To(Equal(osbapi.BindResponse{
					Credentials: map[string]any{
						"foo": "bar",
					},
					Complete: true,
				}))
			})

			When("bind is asynchronous", func() {
				BeforeEach(func() {
					brokerServer.WithResponse(
						"/v2/service_instances/{instance_id}/service_bindings/{binding_id}",
						map[string]any{
							"operation": "bind_op1",
						},
						http.StatusAccepted,
					)
				})

				It("binds the service asynchronously", func() {
					Expect(bindErr).NotTo(HaveOccurred())
					Expect(bindResp).To(Equal(osbapi.BindResponse{
						Operation: "bind_op1",
						Complete:  false,
					}))
				})
			})

			When("binding request fails", func() {
				BeforeEach(func() {
					brokerServer = brokerServer.WithResponse(
						"/v2/service_instances/{instance_id}/service_bindings/{binding_id}",
						nil,
						http.StatusTeapot,
					)
				})

				It("returns an error", func() {
					Expect(bindErr).To(MatchError(ContainSubstring("binding request failed")))
				})
			})

			When("binding request fails with 409 Confilct", func() {
				BeforeEach(func() {
					brokerServer = brokerServer.WithResponse(
						"/v2/service_instances/{instance_id}/service_bindings/{binding_id}",
						nil,
						http.StatusConflict,
					)
				})

				It("returns a confilct error", func() {
					Expect(bindErr).To(BeAssignableToTypeOf(osbapi.ConflictError{}))
				})
			})
		})

		Describe("GetServiceBindingLastOperation", func() {
			var (
				lastOpResp osbapi.LastOperationResponse
				lastOpErr  error
			)

			BeforeEach(func() {
				brokerServer.WithResponse(
					"/v2/service_instances/{instance_id}/service_bindings/{binding_id}/last_operation",
					map[string]any{
						"state":       "in-progress",
						"description": "provisioning",
					},
					http.StatusOK,
				)
			})

			JustBeforeEach(func() {
				lastOpResp, lastOpErr = brokerClient.GetServiceBindingLastOperation(ctx, osbapi.GetServiceBindingLastOperationRequest{
					InstanceID: "my-service-instance",
					BindingID:  "my-binding-id",
					GetLastOperationRequestParameters: osbapi.GetLastOperationRequestParameters{
						ServiceId: "service-guid",
						PlanID:    "plan-guid",
						Operation: "op-guid",
					},
				})
			})

			It("gets the last operation", func() {
				Expect(lastOpErr).NotTo(HaveOccurred())
				Expect(lastOpResp).To(Equal(osbapi.LastOperationResponse{
					State:       "in-progress",
					Description: "provisioning",
				}))
			})

			It("sends correct request to broker", func() {
				Expect(lastOpErr).NotTo(HaveOccurred())
				requests := brokerServer.ServedRequests()

				Expect(requests).To(HaveLen(1))

				Expect(requests[0].Method).To(Equal(http.MethodGet))
				Expect(requests[0].URL.Path).To(Equal("/v2/service_instances/my-service-instance/service_bindings/my-binding-id/last_operation"))

				requestBytes, err := io.ReadAll(requests[0].Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(requestBytes).To(BeEmpty())
			})

			When("getting the last operation request fails", func() {
				BeforeEach(func() {
					brokerServer = brokerServer.WithResponse(
						"/v2/service_instances/{instance_id}/service_bindings/{binding_id}/last_operation",
						nil,
						http.StatusTeapot,
					)
				})

				It("returns an error", func() {
					Expect(lastOpErr).To(MatchError(ContainSubstring("last operation request failed")))
				})
			})

			When("getting the last operation request fails with 410 Gone", func() {
				BeforeEach(func() {
					brokerServer = brokerServer.WithResponse(
						"/v2/service_instances/{instance_id}/service_bindings/{binding_id}/last_operation",
						nil,
						http.StatusGone,
					)
				})

				It("returns a gone error", func() {
					Expect(lastOpErr).To(BeAssignableToTypeOf(osbapi.GoneError{}))
				})
			})
		})

		Describe("GetServiceBinding", func() {
			var (
				getBindingResponse osbapi.GetBindingResponse
				getBindingErr      error
			)

			BeforeEach(func() {
				brokerServer.WithResponse(
					"/v2/service_instances/{instance_id}/service_bindings/{binding_id}",
					map[string]any{
						"credentials": map[string]any{
							"credentialKey": "credentialValue",
						},
					},
					http.StatusOK,
				)
			})

			JustBeforeEach(func() {
				getBindingResponse, getBindingErr = brokerClient.GetServiceBinding(ctx, osbapi.GetServiceBindingRequest{
					InstanceID: "my-service-instance",
					BindingID:  "my-binding-id",
					ServiceId:  "service-guid",
					PlanID:     "plan-guid",
				})
			})

			It("gets the binding", func() {
				Expect(getBindingErr).NotTo(HaveOccurred())
				Expect(getBindingResponse).To(Equal(osbapi.GetBindingResponse{
					Credentials: map[string]any{
						"credentialKey": "credentialValue",
					},
				}))
			})

			It("sends correct request to broker", func() {
				Expect(getBindingErr).NotTo(HaveOccurred())
				requests := brokerServer.ServedRequests()

				Expect(requests).To(HaveLen(1))

				Expect(requests[0].Method).To(Equal(http.MethodGet))
				Expect(requests[0].URL.Path).To(Equal("/v2/service_instances/my-service-instance/service_bindings/my-binding-id"))

				Expect(requests[0].URL.Query()).To(BeEquivalentTo(map[string][]string{
					"service_id": {"service-guid"},
					"plan_id":    {"plan-guid"},
				}))

				requestBytes, err := io.ReadAll(requests[0].Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(requestBytes).To(BeEmpty())
			})

			When("getting the binding fails", func() {
				BeforeEach(func() {
					brokerServer = brokerServer.WithResponse(
						"/v2/service_instances/{instance_id}/service_bindings/{binding_id}",
						nil,
						http.StatusTeapot,
					)
				})

				It("returns an error", func() {
					Expect(getBindingErr).To(MatchError(ContainSubstring("get binding request failed")))
				})
			})
		})
	})

	Describe("Unbind", func() {
		var (
			unbindResp osbapi.UnbindResponse
			unbindErr  error
		)

		BeforeEach(func() {
			brokerServer.WithResponse(
				"/v2/service_instances/{instance_id}/service_bindings/{binding_id}",
				nil,
				http.StatusOK,
			)
		})

		JustBeforeEach(func() {
			unbindResp, unbindErr = brokerClient.Unbind(ctx, osbapi.UnbindPayload{
				InstanceID: "instance-id",
				BindingID:  "binding-id",
				UnbindRequestParameters: osbapi.UnbindRequestParameters{
					ServiceId: "service-guid",
					PlanID:    "plan-guid",
				},
			})
		})

		It("sends an unbind request to broker", func() {
			Expect(unbindErr).NotTo(HaveOccurred())
			requests := brokerServer.ServedRequests()

			Expect(requests).To(HaveLen(1))

			Expect(requests[0].Method).To(Equal(http.MethodDelete))
			Expect(requests[0].URL.Path).To(Equal("/v2/service_instances/instance-id/service_bindings/binding-id"))

			Expect(requests[0].URL.Query()).To(BeEquivalentTo(map[string][]string{
				"service_id":         {"service-guid"},
				"plan_id":            {"plan-guid"},
				"accepts_incomplete": {"true"},
			}))
		})

		It("responds synchronously", func() {
			Expect(unbindErr).NotTo(HaveOccurred())
			Expect(unbindResp.IsComplete()).To(BeTrue())
		})

		When("broker return 202 Accepted", func() {
			BeforeEach(func() {
				brokerServer = brokerServer.WithResponse(
					"/v2/service_instances/{instance_id}/service_bindings/{binding_id}",
					map[string]any{
						"operation": "operation-id",
					},
					http.StatusAccepted,
				)
			})

			It("responds asynchronously", func() {
				Expect(unbindErr).NotTo(HaveOccurred())
				Expect(unbindResp.IsComplete()).To(BeFalse())
				Expect(unbindResp.Operation).To(Equal("operation-id"))
			})
		})

		When("the binding does not exist", func() {
			BeforeEach(func() {
				brokerServer = brokerServer.WithResponse(
					"/v2/service_instances/{instance_id}/service_bindings/{binding_id}",
					nil,
					http.StatusGone,
				)
			})

			It("returns a gone error", func() {
				Expect(unbindErr).To(BeAssignableToTypeOf(osbapi.GoneError{}))
			})
		})

		When("the unbind request fails", func() {
			BeforeEach(func() {
				brokerServer = brokerServer.WithResponse(
					"/v2/service_instances/{instance_id}/service_bindings/{binding_id}",
					nil,
					http.StatusTeapot,
				)
			})

			It("returns an error", func() {
				Expect(unbindErr).To(MatchError(ContainSubstring("unbind request failed")))
			})
		})
	})
})

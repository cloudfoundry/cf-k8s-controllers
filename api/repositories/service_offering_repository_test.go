package repositories_test

import (
	"context"
	"errors"
	"fmt"

	apierrors "code.cloudfoundry.org/korifi/api/errors"
	"code.cloudfoundry.org/korifi/api/repositories"
	korifiv1alpha1 "code.cloudfoundry.org/korifi/controllers/api/v1alpha1"
	"code.cloudfoundry.org/korifi/model/services"
	"code.cloudfoundry.org/korifi/tools"
	. "github.com/onsi/gomega/gstruct"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ServiceOfferingRepo", func() {
	var (
		repo  *repositories.ServiceOfferingRepo
		org   *korifiv1alpha1.CFOrg
		space *korifiv1alpha1.CFSpace
	)

	BeforeEach(func() {
		repo = repositories.NewServiceOfferingRepo(
			userClientFactory,
			rootNamespace,
			repositories.NewServiceBrokerRepo(
				userClientFactory,
				rootNamespace,
			),
			nsPerms,
		)

		org = createOrgWithCleanup(ctx, uuid.NewString())
		space = createSpaceWithCleanup(ctx, org.Name, uuid.NewString())
	})

	Describe("Get", func() {
		var (
			offeringGUID    string
			broker          *korifiv1alpha1.CFServiceBroker
			desiredOffering repositories.ServiceOfferingRecord
			getErr          error
		)

		BeforeEach(func() {
			offeringGUID = uuid.NewString()

			broker = &korifiv1alpha1.CFServiceBroker{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: rootNamespace,
					Name:      uuid.NewString(),
				},
				Spec: korifiv1alpha1.CFServiceBrokerSpec{
					ServiceBroker: services.ServiceBroker{
						Name: uuid.NewString(),
					},
				},
			}
			Expect(k8sClient.Create(ctx, broker)).To(Succeed())

			Expect(k8sClient.Create(ctx, &korifiv1alpha1.CFServiceOffering{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: rootNamespace,
					Name:      offeringGUID,
					Labels: map[string]string{
						korifiv1alpha1.RelServiceBrokerGUIDLabel: broker.Name,
						korifiv1alpha1.RelServiceBrokerNameLabel: broker.Spec.Name,
					},
					Annotations: map[string]string{
						"annotation": "annotation-value",
					},
				},
				Spec: korifiv1alpha1.CFServiceOfferingSpec{
					ServiceOffering: services.ServiceOffering{
						Name:             "my-offering",
						Description:      "my offering description",
						Tags:             []string{"t1"},
						Requires:         []string{"r1"},
						DocumentationURL: tools.PtrTo("https://my.offering.com"),
						BrokerCatalog: services.ServiceBrokerCatalog{
							ID: "offering-catalog-guid",
							Metadata: &runtime.RawExtension{
								Raw: []byte(`{"offering-md": "offering-md-value"}`),
							},
							Features: services.BrokerCatalogFeatures{
								PlanUpdateable:       true,
								Bindable:             true,
								InstancesRetrievable: true,
								BindingsRetrievable:  true,
								AllowContextUpdates:  true,
							},
						},
					},
				},
			})).To(Succeed())
		})

		JustBeforeEach(func() {
			desiredOffering, getErr = repo.GetServiceOffering(ctx, authInfo, offeringGUID)
		})

		It("gets the service offering", func() {
			Expect(getErr).NotTo(HaveOccurred())
			Expect(desiredOffering).To(
				MatchFields(IgnoreExtras, Fields{
					"ServiceOffering": MatchFields(IgnoreExtras, Fields{
						"Name":             Equal("my-offering"),
						"Description":      Equal("my offering description"),
						"Tags":             ConsistOf("t1"),
						"Requires":         ConsistOf("r1"),
						"DocumentationURL": PointTo(Equal("https://my.offering.com")),
						"BrokerCatalog": MatchFields(IgnoreExtras, Fields{
							"ID": Equal("offering-catalog-guid"),
							"Metadata": PointTo(MatchFields(IgnoreExtras, Fields{
								"Raw": MatchJSON(`{"offering-md": "offering-md-value"}`),
							})),
							"Features": MatchFields(IgnoreExtras, Fields{
								"PlanUpdateable":       BeTrue(),
								"Bindable":             BeTrue(),
								"InstancesRetrievable": BeTrue(),
								"BindingsRetrievable":  BeTrue(),
								"AllowContextUpdates":  BeTrue(),
							}),
						}),
					}),
					"CFResource": MatchFields(IgnoreExtras, Fields{
						"GUID":      Equal(offeringGUID),
						"CreatedAt": Not(BeZero()),
						"UpdatedAt": BeNil(),
						"Metadata": MatchAllFields(Fields{
							"Labels":      HaveKeyWithValue(korifiv1alpha1.RelServiceBrokerGUIDLabel, broker.Name),
							"Annotations": HaveKeyWithValue("annotation", "annotation-value"),
						}),
					}),
					"ServiceBrokerGUID": Equal(broker.Name),
				}),
			)
		})

		When("the service offering does not exist", func() {
			BeforeEach(func() {
				offeringGUID = "does-not-exist"
			})
			It("returns a not found error", func() {
				notFoundError := apierrors.NotFoundError{}
				Expect(errors.As(getErr, &notFoundError)).To(BeTrue())
			})
		})
	})

	Describe("List", func() {
		var (
			offeringGUID        string
			anotherOfferingGUID string
			broker              *korifiv1alpha1.CFServiceBroker
			listedOfferings     []repositories.ServiceOfferingRecord
			message             repositories.ListServiceOfferingMessage
			listErr             error
		)

		BeforeEach(func() {
			offeringGUID = uuid.NewString()
			anotherOfferingGUID = uuid.NewString()

			broker = &korifiv1alpha1.CFServiceBroker{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: rootNamespace,
					Name:      uuid.NewString(),
				},
				Spec: korifiv1alpha1.CFServiceBrokerSpec{
					ServiceBroker: services.ServiceBroker{
						Name: uuid.NewString(),
					},
				},
			}
			Expect(k8sClient.Create(ctx, broker)).To(Succeed())

			Expect(k8sClient.Create(ctx, &korifiv1alpha1.CFServiceOffering{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: rootNamespace,
					Name:      offeringGUID,
					Labels: map[string]string{
						korifiv1alpha1.RelServiceBrokerGUIDLabel: broker.Name,
						korifiv1alpha1.RelServiceBrokerNameLabel: broker.Spec.Name,
					},
					Annotations: map[string]string{
						"annotation": "annotation-value",
					},
				},
				Spec: korifiv1alpha1.CFServiceOfferingSpec{
					ServiceOffering: services.ServiceOffering{
						Name:             "my-offering",
						Description:      "my offering description",
						Tags:             []string{"t1"},
						Requires:         []string{"r1"},
						DocumentationURL: tools.PtrTo("https://my.offering.com"),
						BrokerCatalog: services.ServiceBrokerCatalog{
							ID: "offering-catalog-guid",
							Metadata: &runtime.RawExtension{
								Raw: []byte(`{"offering-md": "offering-md-value"}`),
							},
							Features: services.BrokerCatalogFeatures{
								PlanUpdateable:       true,
								Bindable:             true,
								InstancesRetrievable: true,
								BindingsRetrievable:  true,
								AllowContextUpdates:  true,
							},
						},
					},
				},
			})).To(Succeed())

			Expect(k8sClient.Create(ctx, &korifiv1alpha1.CFServiceOffering{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: rootNamespace,
					Name:      anotherOfferingGUID,
					Labels: map[string]string{
						korifiv1alpha1.RelServiceBrokerGUIDLabel: "another-broker",
						korifiv1alpha1.RelServiceBrokerNameLabel: "another-broker-name",
					},
				},
				Spec: korifiv1alpha1.CFServiceOfferingSpec{
					ServiceOffering: services.ServiceOffering{
						Name: "another-offering",
					},
				},
			})).To(Succeed())

			message = repositories.ListServiceOfferingMessage{}
		})

		JustBeforeEach(func() {
			listedOfferings, listErr = repo.ListOfferings(ctx, authInfo, message)
		})

		It("lists service offerings", func() {
			Expect(listErr).NotTo(HaveOccurred())
			Expect(listedOfferings).To(ConsistOf(
				MatchFields(IgnoreExtras, Fields{
					"ServiceOffering": MatchFields(IgnoreExtras, Fields{
						"Name":             Equal("my-offering"),
						"Description":      Equal("my offering description"),
						"Tags":             ConsistOf("t1"),
						"Requires":         ConsistOf("r1"),
						"DocumentationURL": PointTo(Equal("https://my.offering.com")),
						"BrokerCatalog": MatchFields(IgnoreExtras, Fields{
							"ID": Equal("offering-catalog-guid"),
							"Metadata": PointTo(MatchFields(IgnoreExtras, Fields{
								"Raw": MatchJSON(`{"offering-md": "offering-md-value"}`),
							})),
							"Features": MatchFields(IgnoreExtras, Fields{
								"PlanUpdateable":       BeTrue(),
								"Bindable":             BeTrue(),
								"InstancesRetrievable": BeTrue(),
								"BindingsRetrievable":  BeTrue(),
								"AllowContextUpdates":  BeTrue(),
							}),
						}),
					}),
					"CFResource": MatchFields(IgnoreExtras, Fields{
						"GUID":      Equal(offeringGUID),
						"CreatedAt": Not(BeZero()),
						"UpdatedAt": BeNil(),
						"Metadata": MatchAllFields(Fields{
							"Labels":      HaveKeyWithValue(korifiv1alpha1.RelServiceBrokerGUIDLabel, broker.Name),
							"Annotations": HaveKeyWithValue("annotation", "annotation-value"),
						}),
					}),
					"ServiceBrokerGUID": Equal(broker.Name),
				}),
				MatchFields(IgnoreExtras, Fields{
					"CFResource": MatchFields(IgnoreExtras, Fields{
						"GUID": Equal(anotherOfferingGUID),
					}),
				}),
			))
		})

		When("filtering by name", func() {
			BeforeEach(func() {
				message.Names = []string{"my-offering"}
			})

			It("returns the matching offerings", func() {
				Expect(listErr).NotTo(HaveOccurred())
				Expect(listedOfferings).To(ConsistOf(MatchFields(IgnoreExtras, Fields{
					"ServiceOffering": MatchFields(IgnoreExtras, Fields{
						"Name": Equal("my-offering"),
					}),
				})))
			})
		})

		When("filtering by broker name", func() {
			BeforeEach(func() {
				message.BrokerNames = []string{broker.Spec.Name}
			})

			It("returns the matching offerings", func() {
				Expect(listErr).NotTo(HaveOccurred())
				Expect(listedOfferings).To(ConsistOf(MatchFields(IgnoreExtras, Fields{
					"ServiceBrokerGUID": Equal(broker.Name),
				})))
			})
		})

		When("filtering by guid", func() {
			BeforeEach(func() {
				message.GUIDs = []string{offeringGUID}
			})

			It("returns the matching offerings", func() {
				Expect(listErr).NotTo(HaveOccurred())
				Expect(listedOfferings).To(ConsistOf(MatchFields(IgnoreExtras, Fields{
					"CFResource": MatchFields(IgnoreExtras, Fields{
						"GUID": Equal(offeringGUID),
					}),
				})))
			})
		})
	})

	Describe("Delete", func() {
		var (
			plan      *korifiv1alpha1.CFServicePlan
			offering  *korifiv1alpha1.CFServiceOffering
			instance  *korifiv1alpha1.CFServiceInstance
			binding   *korifiv1alpha1.CFServiceBinding
			message   repositories.DeleteServiceOfferingMessage
			deleteErr error
		)

		BeforeEach(func() {
			offering = &korifiv1alpha1.CFServiceOffering{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: rootNamespace,
					Name:      uuid.NewString(),
				},
				Spec: korifiv1alpha1.CFServiceOfferingSpec{
					ServiceOffering: services.ServiceOffering{
						Name:        "my-offering",
						Description: "my offering description",
						Tags:        []string{"t1"},
						Requires:    []string{"r1"},
					},
				},
			}
			Expect(k8sClient.Create(ctx, offering)).To(Succeed())

			plan = &korifiv1alpha1.CFServicePlan{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: rootNamespace,
					Name:      uuid.NewString(),
					Labels: map[string]string{
						korifiv1alpha1.RelServiceOfferingGUIDLabel: offering.Name,
					},
				},
				Spec: korifiv1alpha1.CFServicePlanSpec{
					ServicePlan: services.ServicePlan{
						Name:        "my-service-plan",
						Free:        true,
						Description: "service plan description",
					},
					Visibility: korifiv1alpha1.ServicePlanVisibility{
						Type: korifiv1alpha1.PublicServicePlanVisibilityType,
					},
				},
			}
			Expect(k8sClient.Create(ctx, plan)).To(Succeed())

			instance = createServiceInstanceCR(ctx, k8sClient, uuid.NewString(), space.Name, "my-service-instance", "secret-name")
			instance.Spec.PlanGUID = plan.Name
			instance.Finalizers = append(instance.Finalizers, korifiv1alpha1.CFManagedServiceInstanceFinalizerName)

			Expect(k8sClient.Update(ctx, instance)).To(Succeed())

			binding = &korifiv1alpha1.CFServiceBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name:      uuid.NewString(),
					Namespace: space.Name,
					Labels: map[string]string{
						korifiv1alpha1.PlanGUIDLabelKey: plan.Name,
					},
				},
				Spec: korifiv1alpha1.CFServiceBindingSpec{
					Service: corev1.ObjectReference{
						Kind:       "CFServiceInstance",
						APIVersion: korifiv1alpha1.SchemeGroupVersion.Identifier(),
						Name:       instance.Name,
					},
					AppRef: corev1.LocalObjectReference{
						Name: "some-app-guid",
					},
				},
			}

			binding.Finalizers = append(binding.Finalizers, korifiv1alpha1.CFServiceBindingFinalizerName)
			Expect(k8sClient.Create(ctx, binding)).To(Succeed())

			message = repositories.DeleteServiceOfferingMessage{GUID: offering.Name}
		})

		JustBeforeEach(func() {
			createRoleBinding(ctx, userName, spaceDeveloperRole.Name, space.Name)
			deleteErr = repo.DeleteOffering(ctx, authInfo, message)
		})

		It("successfully deletes the offering", func() {
			Expect(deleteErr).ToNot(HaveOccurred())

			namespacedName := types.NamespacedName{
				Name:      offering.Name,
				Namespace: rootNamespace,
			}

			err := k8sClient.Get(context.Background(), namespacedName, &korifiv1alpha1.CFServiceOffering{})
			Expect(k8serrors.IsNotFound(err)).To(BeTrue(), fmt.Sprintf("error: %+v", err))
		})

		When("the service offering does not exist", func() {
			BeforeEach(func() {
				message.GUID = "does-not-exist"
			})

			It("returns a error", func() {
				Expect(errors.As(deleteErr, &apierrors.NotFoundError{})).To(BeTrue())
			})
		})

		When("Purge is set to true", func() {
			BeforeEach(func() {
				message.Purge = true
			})
			It("successfully deletes the offering and all related resources", func() {
				Expect(deleteErr).ToNot(HaveOccurred())

				err := k8sClient.Get(context.Background(), types.NamespacedName{Name: offering.Name, Namespace: rootNamespace}, &korifiv1alpha1.CFServiceOffering{})
				Expect(k8serrors.IsNotFound(err)).To(BeTrue(), fmt.Sprintf("error: %+v", err))

				err = k8sClient.Get(context.Background(), types.NamespacedName{Name: plan.Name, Namespace: rootNamespace}, &korifiv1alpha1.CFServicePlan{})
				Expect(k8serrors.IsNotFound(err)).To(BeTrue(), fmt.Sprintf("error: %+v", err))

				err = k8sClient.Get(context.Background(), types.NamespacedName{Name: instance.Name, Namespace: space.Name}, &korifiv1alpha1.CFServiceInstance{})
				Expect(k8serrors.IsNotFound(err)).To(BeTrue(), fmt.Sprintf("error: %+v", err))

				serviceBinding := new(korifiv1alpha1.CFServiceBinding)
				err = k8sClient.Get(context.Background(), types.NamespacedName{Name: binding.Name, Namespace: space.Name}, serviceBinding)

				Expect(err).ToNot(HaveOccurred())
				Expect(serviceBinding.Finalizers).To(BeEmpty())
			})
		})
	})
})

package workloads_test

import (
	"context"
	"errors"

	workloadsv1alpha1 "code.cloudfoundry.org/cf-k8s-controllers/apis/workloads/v1alpha1"
	cfconfig "code.cloudfoundry.org/cf-k8s-controllers/config/cf"
	. "code.cloudfoundry.org/cf-k8s-controllers/controllers/workloads"
	"code.cloudfoundry.org/cf-k8s-controllers/controllers/workloads/fake"
	. "code.cloudfoundry.org/cf-k8s-controllers/controllers/workloads/testutils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	buildv1alpha1 "github.com/pivotal/kpack/pkg/apis/build/v1alpha1"
	"github.com/pivotal/kpack/pkg/apis/core/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var _ = Describe("CFBuildReconciler", func() {
	const (
		defaultNamespace        = "default"
		stagingConditionType    = "Staging"
		readyConditionType      = "Ready"
		succeededConditionType  = "Succeeded"
		kpackReadyConditionType = "Ready"
	)

	var (
		fakeClient       *fake.CFClient
		fakeStatusWriter *fake.StatusWriter

		cfAppGUID      string
		cfPackageGUID  string
		cfBuildGUID    string
		kpackImageGUID string

		cfBuild           *workloadsv1alpha1.CFBuild
		cfBuildError      error
		cfApp             *workloadsv1alpha1.CFApp
		cfAppError        error
		cfPackage         *workloadsv1alpha1.CFPackage
		cfPackageError    error
		kpackImage        *buildv1alpha1.Image
		kpackImageError   error
		cfBuildReconciler *CFBuildReconciler
		req               ctrl.Request
		ctx               context.Context

		reconcileResult ctrl.Result
		reconcileErr    error
	)

	BeforeEach(func() {
		fakeClient = new(fake.CFClient)

		cfAppGUID = "cf-app-guid"
		cfPackageGUID = "cf-package-guid"
		cfBuildGUID = "cf-build-guid"
		kpackImageGUID = cfBuildGUID

		cfBuild = MockCFBuildObject(cfBuildGUID, defaultNamespace, cfPackageGUID, cfAppGUID)
		cfBuildError = nil
		cfApp = MockAppCRObject(cfAppGUID, defaultNamespace)
		cfAppError = nil
		cfPackage = MockPackageCRObject(cfPackageGUID, defaultNamespace, cfAppGUID)
		cfPackageError = nil
		kpackImage = mockKpackImageObject(kpackImageGUID, defaultNamespace)
		kpackImageError = apierrors.NewNotFound(schema.GroupResource{}, cfBuild.Name)

		fakeClient.GetStub = func(_ context.Context, _ types.NamespacedName, obj client.Object) error {
			// cast obj to find its kind
			switch obj.(type) {
			case *workloadsv1alpha1.CFBuild:
				cfBuild.DeepCopyInto(obj.(*workloadsv1alpha1.CFBuild))
				return cfBuildError
			case *workloadsv1alpha1.CFApp:
				cfApp.DeepCopyInto(obj.(*workloadsv1alpha1.CFApp))
				return cfAppError
			case *workloadsv1alpha1.CFPackage:
				cfPackage.DeepCopyInto(obj.(*workloadsv1alpha1.CFPackage))
				return cfPackageError
			case *buildv1alpha1.Image:
				kpackImage.DeepCopyInto(obj.(*buildv1alpha1.Image))
				return kpackImageError
			default:
				panic("test Client Get provided a weird obj")
			}
		}

		// Configure mock status update to succeed
		fakeStatusWriter = &fake.StatusWriter{}
		fakeClient.StatusReturns(fakeStatusWriter)

		// configure a CFBuildReconciler with the client
		Expect(workloadsv1alpha1.AddToScheme(scheme.Scheme)).To(Succeed())
		Expect(buildv1alpha1.AddToScheme(scheme.Scheme)).To(Succeed())
		cfBuildReconciler = &CFBuildReconciler{
			Client:           fakeClient,
			Scheme:           scheme.Scheme,
			Log:              zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)),
			ControllerConfig: &cfconfig.ControllerConfig{KpackImageTag: "image/registry/tag"},
		}
		req = ctrl.Request{
			NamespacedName: types.NamespacedName{
				Namespace: defaultNamespace,
				Name:      cfBuildGUID,
			},
		}
		ctx = context.Background()
	})

	When("CFBuild status conditions are unknown and", func() {
		When("on the happy path", func() {

			BeforeEach(func() {
				reconcileResult, reconcileErr = cfBuildReconciler.Reconcile(ctx, req)
			})

			It("does not return an error", func() {
				Expect(reconcileErr).NotTo(HaveOccurred())
			})

			It("returns an empty result", func() {
				Expect(reconcileResult).To(Equal(ctrl.Result{}))
			})

			It("should create kpack image with the same GUID as the CF Build", func() {
				Expect(fakeClient.CreateCallCount()).To(Equal(1), "fakeClient Create was not called 1 time")
				_, kpackImage, _ := fakeClient.CreateArgsForCall(0)
				Expect(kpackImage.GetName()).To(Equal(cfBuildGUID))
			})
			It("should update the status conditions on CFBuild", func() {
				Expect(fakeClient.StatusCallCount()).To(Equal(1))
			})

		})
		When("on unhappy path", func() {
			When("fetch CFBuild returns an error", func() {
				BeforeEach(func() {
					cfBuildError = errors.New("failing on purpose")
					reconcileResult, reconcileErr = cfBuildReconciler.Reconcile(ctx, req)
				})
				It("should return an error", func() {
					Expect(reconcileErr).To(HaveOccurred())
				})
			})
			When("fetch CFBuild returns a NotFoundError", func() {
				BeforeEach(func() {
					cfBuildError = apierrors.NewNotFound(schema.GroupResource{}, cfBuild.Name)
					reconcileResult, reconcileErr = cfBuildReconciler.Reconcile(ctx, req)
				})
				It("should NOT return any error", func() {
					Expect(reconcileErr).NotTo(HaveOccurred())
				})
			})
			When("fetch CFApp returns an error", func() {
				BeforeEach(func() {
					cfAppError = errors.New("failing on purpose")
					reconcileResult, reconcileErr = cfBuildReconciler.Reconcile(ctx, req)
				})
				It("should return an error", func() {
					Expect(reconcileErr).To(HaveOccurred())
				})
			})
			When("fetch CFPackage returns an error", func() {
				BeforeEach(func() {
					cfPackageError = errors.New("failing on purpose")
					reconcileResult, reconcileErr = cfBuildReconciler.Reconcile(ctx, req)
				})
				It("should return an error", func() {
					Expect(reconcileErr).To(HaveOccurred())
				})
			})
			When("create Kpack Image returns an error", func() {
				BeforeEach(func() {
					fakeClient.CreateReturns(errors.New("failing on purpose"))
					reconcileResult, reconcileErr = cfBuildReconciler.Reconcile(ctx, req)
				})
				It("should return an error", func() {
					Expect(reconcileErr).To(HaveOccurred())
				})
			})
			When("update status conditions returns an error", func() {
				BeforeEach(func() {
					fakeStatusWriter.UpdateReturns(errors.New("failing on purpose"))
					reconcileResult, reconcileErr = cfBuildReconciler.Reconcile(ctx, req)
				})
				It("should return an error", func() {
					Expect(reconcileErr).To(HaveOccurred())
				})
			})
		})

	})
	When("CFBuild status conditions for Staging is True and others are unknown", func() {
		BeforeEach(func() {
			SetStatusCondition(&cfBuild.Status.Conditions, stagingConditionType, metav1.ConditionTrue)
			SetStatusCondition(&cfBuild.Status.Conditions, succeededConditionType, metav1.ConditionUnknown)
		})
		When("on the happy path", func() {
			BeforeEach(func() {
				kpackImageError = nil
				setKpackImageStatus(kpackImage, kpackReadyConditionType, "True")
				reconcileResult, reconcileErr = cfBuildReconciler.Reconcile(ctx, req)
			})
			It("does not return an error", func() {
				Expect(reconcileErr).NotTo(HaveOccurred())
			})

			It("returns an empty result", func() {
				Expect(reconcileResult).To(Equal(ctrl.Result{}))
			})

			It("does not create a new kpack image", func() {
				Expect(fakeClient.CreateCallCount()).To(Equal(0), "fakeClient Create was called When It should not have been")
			})

			It("should update the status conditions on CFBuild", func() {
				Expect(fakeClient.StatusCallCount()).To(Equal(1))
			})
		})
		When("on the unhappy path", func() {
			When("fetch KpackImage returns an error", func() {
				BeforeEach(func() {
					kpackImageError = errors.New("failing on purpose")
					reconcileResult, reconcileErr = cfBuildReconciler.Reconcile(ctx, req)
				})
				It("should return an error", func() {
					Expect(reconcileErr).To(HaveOccurred())
				})
			})
			When("fetch KpackImage returns a NotFoundError", func() {
				BeforeEach(func() {
					kpackImageError = apierrors.NewNotFound(schema.GroupResource{}, cfBuild.Name)
					reconcileResult, reconcileErr = cfBuildReconciler.Reconcile(ctx, req)
				})
				It("should NOT return any error", func() {
					Expect(reconcileErr).NotTo(HaveOccurred())
				})
			})
			When("kpack image status conditions for Type Succeeded is nil", func() {
				BeforeEach(func() {
					reconcileResult, reconcileErr = cfBuildReconciler.Reconcile(ctx, req)
				})
				It("does not return an error", func() {
					Expect(reconcileErr).NotTo(HaveOccurred())
				})

				It("returns an empty result", func() {
					Expect(reconcileResult).To(Equal(ctrl.Result{}))
				})
			})
			When("kpack image status conditions for Type Succeeded is UNKNOWN", func() {
				BeforeEach(func() {
					kpackImageError = nil
					setKpackImageStatus(kpackImage, kpackReadyConditionType, "Unknown")
					reconcileResult, reconcileErr = cfBuildReconciler.Reconcile(ctx, req)
				})
				It("does not return an error", func() {
					Expect(reconcileErr).NotTo(HaveOccurred())
				})

				It("returns an empty result", func() {
					Expect(reconcileResult).To(Equal(ctrl.Result{}))
				})
			})
			When("kpack image status conditions for Type Succeeded is FALSE", func() {
				BeforeEach(func() {
					kpackImageError = nil
					setKpackImageStatus(kpackImage, kpackReadyConditionType, "False")
					reconcileResult, reconcileErr = cfBuildReconciler.Reconcile(ctx, req)
				})
				It("does not return an error", func() {
					Expect(reconcileErr).NotTo(HaveOccurred())
				})
				It("updates status conditions on CFBuild", func() {
					Expect(fakeClient.StatusCallCount()).To(Equal(1))
				})
				When("update status conditions returns an error", func() {
					BeforeEach(func() {
						fakeStatusWriter.UpdateReturns(errors.New("failing on purpose"))
						reconcileResult, reconcileErr = cfBuildReconciler.Reconcile(ctx, req)
					})
					It("should return an error", func() {
						Expect(reconcileErr).To(HaveOccurred())
					})
				})
			})
			When("kpack image status conditions for Type Succeeded is True", func() {
				When("update status conditions returns an error", func() {
					BeforeEach(func() {
						kpackImageError = nil
						setKpackImageStatus(kpackImage, kpackReadyConditionType, "True")
						fakeStatusWriter.UpdateReturns(errors.New("failing on purpose"))
						reconcileResult, reconcileErr = cfBuildReconciler.Reconcile(ctx, req)
					})
					It("should return an error", func() {
						Expect(reconcileErr).To(HaveOccurred())
					})
				})
			})
		})
	})
})

func mockKpackImageObject(guid string, namespace string) *buildv1alpha1.Image {
	return &buildv1alpha1.Image{
		ObjectMeta: metav1.ObjectMeta{
			Name:      guid,
			Namespace: namespace,
		},
		Spec: buildv1alpha1.ImageSpec{
			Tag:            "test-tag",
			Builder:        corev1.ObjectReference{},
			ServiceAccount: "test-service-account",
			Source: buildv1alpha1.SourceConfig{
				Registry: &buildv1alpha1.Registry{
					Image:            "image-path",
					ImagePullSecrets: nil,
				},
			},
		},
	}
}

func setKpackImageStatus(kpackImage *buildv1alpha1.Image, conditionType string, conditionStatus string) {
	kpackImage.Status.Conditions = append(kpackImage.Status.Conditions, v1alpha1.Condition{
		Type:   v1alpha1.ConditionType(conditionType),
		Status: corev1.ConditionStatus(conditionStatus),
	})
}

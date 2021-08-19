package controllers_test

import (
	"code.cloudfoundry.org/cf-k8s-controllers/api/v1alpha1"
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"time"
)

var _ = Describe("CFAppReconciler", func() {
	When("a new record is created", func() {
		const (
			cfAppGUID = "test-app-guid"
			namespace = "default"
		)
		It("sets its status.conditions", func() {
			ctx := context.Background()
			cfApp := &v1alpha1.CFApp{
				TypeMeta: metav1.TypeMeta{
					Kind:       "CFApp",
					APIVersion: v1alpha1.GroupVersion.Identifier(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      cfAppGUID,
					Namespace: namespace,
				},
				Spec: v1alpha1.CFAppSpec{
					Name:         "test-app",
					DesiredState: "STOPPED",
					Lifecycle: v1alpha1.Lifecycle{
						Type: "buildpack",
						Data: v1alpha1.LifecycleData{
							Buildpacks: []string{}, // TODO: field can't be null - may want to change this
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, cfApp)).To(Succeed())

			cfAppLookupKey := types.NamespacedName{Name: cfAppGUID, Namespace: namespace}
			createdCFApp := new(v1alpha1.CFApp)

			Eventually(func() []metav1.Condition {
				err := k8sClient.Get(ctx, cfAppLookupKey, createdCFApp)
				if err != nil {
					return nil
				}
				return createdCFApp.Status.Conditions
			}, 10*time.Second, 250*time.Millisecond).ShouldNot(BeEmpty())

			runningConditionFalse := meta.IsStatusConditionFalse(createdCFApp.Status.Conditions, "Running")
			Expect(runningConditionFalse).To(BeTrue())

			restartingConditionFalse := meta.IsStatusConditionFalse(createdCFApp.Status.Conditions, "Restarting")
			Expect(restartingConditionFalse).To(BeTrue())
		})
	})
})

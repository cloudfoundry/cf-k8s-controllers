package conditions_test

import (
	"sync"
	"time"

	"code.cloudfoundry.org/korifi/api/repositories/conditions"
	korifiv1alpha1 "code.cloudfoundry.org/korifi/controllers/api/v1alpha1"
	"code.cloudfoundry.org/korifi/tools/k8s"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("ConditionAwaiter", func() {
	var (
		awaiter     *conditions.Awaiter[*korifiv1alpha1.CFTask, korifiv1alpha1.CFTaskList, *korifiv1alpha1.CFTaskList]
		task        *korifiv1alpha1.CFTask
		awaitedTask *korifiv1alpha1.CFTask
		awaitErr    error
	)

	asyncPatchTask := func(patchTask func(*korifiv1alpha1.CFTask)) {
		wg := &sync.WaitGroup{}
		wg.Add(1)

		go func() {
			defer GinkgoRecover()
			defer wg.Done()

			patchedTask := task.DeepCopy()
			Expect(k8s.Patch(ctx, k8sClient, patchedTask, func() {
				patchTask(patchedTask)
			})).To(Succeed())
		}()

		DeferCleanup(func() {
			wg.Wait()
		})
	}

	BeforeEach(func() {
		awaiter = conditions.NewConditionAwaiter[*korifiv1alpha1.CFTask, korifiv1alpha1.CFTaskList](time.Second)
		awaitedTask = nil
		awaitErr = nil

		task = &korifiv1alpha1.CFTask{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespace,
				Name:      "my-task",
			},
		}

		Expect(k8sClient.Create(ctx, task)).To(Succeed())
	})

	JustBeforeEach(func() {
		awaitedTask, awaitErr = awaiter.AwaitCondition(ctx, k8sClient, task, korifiv1alpha1.TaskInitializedConditionType)
	})

	It("returns an error", func() {
		Expect(awaitErr).To(MatchError(ContainSubstring("condition Initialized not set yet")))
	})

	When("the condition becomes false", func() {
		BeforeEach(func() {
			asyncPatchTask(func(cfTask *korifiv1alpha1.CFTask) {
				meta.SetStatusCondition(&cfTask.Status.Conditions, metav1.Condition{
					Type:               korifiv1alpha1.TaskInitializedConditionType,
					Status:             metav1.ConditionFalse,
					Reason:             "initialized",
					ObservedGeneration: task.Generation,
				})
			})
		})

		It("returns an error", func() {
			Expect(awaitErr).To(MatchError(ContainSubstring("expected the Initialized condition to be true")))
		})
	})

	When("the condition becomes true", func() {
		BeforeEach(func() {
			asyncPatchTask(func(cfTask *korifiv1alpha1.CFTask) {
				meta.SetStatusCondition(&cfTask.Status.Conditions, metav1.Condition{
					Type:               korifiv1alpha1.TaskInitializedConditionType,
					Status:             metav1.ConditionTrue,
					Reason:             "initialized",
					ObservedGeneration: task.Generation,
				})
			})
		})

		It("succeeds and returns the updated object", func() {
			Expect(awaitErr).NotTo(HaveOccurred())
			Expect(awaitedTask).NotTo(BeNil())

			Expect(awaitedTask.Name).To(Equal(task.Name))
			Expect(meta.IsStatusConditionTrue(awaitedTask.Status.Conditions, korifiv1alpha1.TaskInitializedConditionType)).To(BeTrue())
		})
	})

	When("the condition becomes true but is outdated", func() {
		BeforeEach(func() {
			Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(task), task)).To(Succeed())
			asyncPatchTask(func(cfTask *korifiv1alpha1.CFTask) {
				meta.SetStatusCondition(&cfTask.Status.Conditions, metav1.Condition{
					Type:               korifiv1alpha1.TaskInitializedConditionType,
					Status:             metav1.ConditionTrue,
					Reason:             "initialized",
					ObservedGeneration: task.Generation - 1,
				})
			})
		})

		It("returns an error", func() {
			Expect(awaitErr).To(MatchError(ContainSubstring("condition Initialized is outdated")))
		})
	})
})

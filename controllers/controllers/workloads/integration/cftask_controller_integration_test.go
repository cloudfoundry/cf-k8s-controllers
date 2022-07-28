package integration_test

import (
	"context"

	eiriniv1 "code.cloudfoundry.org/eirini-controller/pkg/apis/eirini/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	korifiv1alpha1 "code.cloudfoundry.org/korifi/controllers/api/v1alpha1"
	"code.cloudfoundry.org/korifi/controllers/controllers/workloads/testutils"
)

var _ = Describe("CFTaskReconciler Integration Tests", func() {
	var (
		ctx context.Context
		ns  string

		cfTask    *korifiv1alpha1.CFTask
		cfApp     *korifiv1alpha1.CFApp
		cfDroplet *korifiv1alpha1.CFBuild
		envSecret *corev1.Secret
	)

	BeforeEach(func() {
		ctx = context.Background()
		ns = testutils.PrefixedGUID("namespace")
		Expect(k8sClient.Create(ctx, &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: ns,
			},
		})).To(Succeed())

		cfDroplet = &korifiv1alpha1.CFBuild{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: ns,
				Name:      testutils.PrefixedGUID("droplet"),
			},
			Spec: korifiv1alpha1.CFBuildSpec{
				Lifecycle: korifiv1alpha1.Lifecycle{Type: "buildpack"},
			},
		}
		Expect(k8sClient.Create(ctx, cfDroplet)).To(Succeed())

		cfDropletCopy := cfDroplet.DeepCopy()
		cfDropletCopy.Status.Droplet = &korifiv1alpha1.BuildDropletStatus{
			Registry: korifiv1alpha1.Registry{Image: "registry.io/my/image"},
			ProcessTypes: []korifiv1alpha1.ProcessType{{
				Type:    "web",
				Command: "cmd",
			}},
			Ports: []int32{8080},
		}
		meta.SetStatusCondition(&cfDropletCopy.Status.Conditions, metav1.Condition{
			Type:   "type",
			Status: "Unknown",
			Reason: "reason",
		})
		Expect(k8sClient.Status().Patch(ctx, cfDropletCopy, client.MergeFrom(cfDroplet))).To(Succeed())

		envSecret = &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testutils.PrefixedGUID("env-secret"),
				Namespace: ns,
			},
			StringData: map[string]string{
				"BOB":  "flemming",
				"FAST": "show",
			},
			Type: corev1.SecretTypeOpaque,
		}
		Expect(k8sClient.Create(ctx, envSecret)).To(Succeed())

		cfApp = &korifiv1alpha1.CFApp{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: ns,
				Name:      testutils.PrefixedGUID("app"),
			},
			Spec: korifiv1alpha1.CFAppSpec{
				Lifecycle: korifiv1alpha1.Lifecycle{Type: "buildpack"},
				CurrentDropletRef: corev1.LocalObjectReference{
					Name: cfDroplet.Name,
				},
				DesiredState:  "STOPPED",
				DisplayName:   "app",
				EnvSecretName: envSecret.Name,
			},
		}
		Expect(k8sClient.Create(ctx, cfApp)).To(Succeed())

		cfProcess := &korifiv1alpha1.CFProcess{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: ns,
				Name:      testutils.PrefixedGUID("web-process"),
				Labels: map[string]string{
					korifiv1alpha1.CFProcessTypeLabelKey: "web",
					korifiv1alpha1.CFAppGUIDLabelKey:     cfApp.Name,
				},
			},
			Spec: korifiv1alpha1.CFProcessSpec{
				AppRef:      corev1.LocalObjectReference{Name: cfApp.Name},
				ProcessType: "web",
				Command:     "echo hello",
				MemoryMB:    768,
				HealthCheck: korifiv1alpha1.HealthCheck{
					Type: "process",
				},
				Ports: []int32{8080},
			},
		}
		Expect(k8sClient.Create(ctx, cfProcess)).To(Succeed())

		cfTask = &korifiv1alpha1.CFTask{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: ns,
				Name:      testutils.PrefixedGUID("cftask"),
			},
			Spec: korifiv1alpha1.CFTaskSpec{
				Command: "echo hello",
				AppRef: corev1.LocalObjectReference{
					Name: cfApp.Name,
				},
			},
		}
	})

	Describe("CFTask creation", func() {
		JustBeforeEach(func() {
			Expect(k8sClient.Create(ctx, cfTask)).To(Succeed())
		})

		When("the task gets initialized", func() {
			var task *korifiv1alpha1.CFTask

			BeforeEach(func() {
				task = &korifiv1alpha1.CFTask{}
			})

			JustBeforeEach(func() {
				Eventually(func(g Gomega) {
					g.Expect(k8sClient.Get(ctx, types.NamespacedName{Namespace: ns, Name: cfTask.Name}, task)).To(Succeed())
					g.Expect(meta.IsStatusConditionTrue(task.Status.Conditions, korifiv1alpha1.TaskInitializedConditionType)).To(BeTrue(), "task did not become initialized")
				}).Should(Succeed())
			})

			It("populates the Status of the CFTask", func() {
				Expect(task.Status.SequenceID).NotTo(BeZero())
				Expect(task.Status.MemoryMB).To(Equal(cfProcessDefaults.MemoryMB))
				Expect(task.Status.DiskQuotaMB).To(Equal(cfProcessDefaults.DiskQuotaMB))
				Expect(task.Status.DropletRef.Name).To(Equal(cfDroplet.Name))
			})
		})

		It("SequenceID does not change on task update", func() {
			var task korifiv1alpha1.CFTask

			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(ctx, types.NamespacedName{Namespace: ns, Name: cfTask.Name}, &task)).To(Succeed())
				g.Expect(task.Status.SequenceID).NotTo(BeZero())
			}).Should(Succeed())

			seqId := task.Status.SequenceID

			updatedTask := task.DeepCopy()
			updatedTask.Spec.Command = "foo bar"
			Expect(k8sClient.Patch(ctx, updatedTask, client.MergeFrom(&task))).To(Succeed())

			Consistently(func(g Gomega) {
				g.Expect(k8sClient.Get(ctx, types.NamespacedName{Namespace: ns, Name: cfTask.Name}, &task)).To(Succeed())
				g.Expect(task.Status.SequenceID).To(Equal(seqId))
			}).Should(Succeed())
		})

		It("creates an eirini.Task", func() {
			var tasks eiriniv1.TaskList
			Eventually(func() ([]eiriniv1.Task, error) {
				err := k8sClient.List(
					ctx,
					&tasks,
					client.InNamespace(ns),
					client.MatchingLabels{korifiv1alpha1.CFTaskGUIDLabelKey: cfTask.Name},
				)
				return tasks.Items, err
			}).Should(HaveLen(1))

			Expect(tasks.Items[0].Name).To(Equal(cfTask.Name))
			Expect(tasks.Items[0].Spec.GUID).To(Equal(cfTask.Name))
			Expect(tasks.Items[0].Spec.Command).To(Equal([]string{"/cnb/lifecycle/launcher", "echo hello"}))
			Expect(tasks.Items[0].Spec.Image).To(Equal("registry.io/my/image"))
			Expect(tasks.Items[0].Spec.MemoryMB).To(Equal(cfProcessDefaults.MemoryMB))
			Expect(tasks.Items[0].Spec.DiskMB).To(Equal(cfProcessDefaults.DiskQuotaMB))
			Expect(tasks.Items[0].Spec.CPUMillis).To(BeEquivalentTo(75))

			// refresh the VCAPServicesSecretName
			Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(cfApp), cfApp)).To(Succeed())

			Expect(tasks.Items[0].Spec.Environment).To(ConsistOf(
				corev1.EnvVar{
					Name: "BOB",
					ValueFrom: &corev1.EnvVarSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: envSecret.Name,
							},
							Key: "BOB",
						},
					},
				},
				corev1.EnvVar{
					Name: "FAST",
					ValueFrom: &corev1.EnvVarSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: envSecret.Name,
							},
							Key: "FAST",
						},
					},
				},
				corev1.EnvVar{
					Name: "VCAP_SERVICES",
					ValueFrom: &corev1.EnvVarSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: cfApp.Status.VCAPServicesSecretName,
							},
							Key: "VCAP_SERVICES",
						},
					},
				},
			))
		})

		When("the eirini task status condition changes", func() {
			JustBeforeEach(func() {
				var tasks eiriniv1.TaskList
				Eventually(func() ([]eiriniv1.Task, error) {
					err := k8sClient.List(
						ctx,
						&tasks,
						client.InNamespace(ns),
						client.MatchingLabels{korifiv1alpha1.CFTaskGUIDLabelKey: cfTask.Name},
					)
					return tasks.Items, err
				}).Should(HaveLen(1))

				modifiedEiriniTask := tasks.Items[0].DeepCopy()
				meta.SetStatusCondition(&modifiedEiriniTask.Status.Conditions, metav1.Condition{
					Type:    eiriniv1.TaskStartedConditionType,
					Status:  metav1.ConditionTrue,
					Reason:  "task_started",
					Message: "task started",
				})

				Expect(k8sClient.Status().Patch(ctx, modifiedEiriniTask, client.MergeFrom(&tasks.Items[0]))).To(Succeed())
			})

			It("reflects the status in the korifi task", func() {
				Eventually(func(g Gomega) {
					var task korifiv1alpha1.CFTask

					g.Expect(k8sClient.Get(ctx, types.NamespacedName{Namespace: ns, Name: cfTask.Name}, &task)).To(Succeed())
					g.Expect(meta.IsStatusConditionTrue(task.Status.Conditions, korifiv1alpha1.TaskStartedConditionType)).To(BeTrue())
				}).Should(Succeed())
			})
		})
	})

	Describe("CFTask Cancellation", func() {
		BeforeEach(func() {
			Expect(k8sClient.Create(ctx, cfTask)).To(Succeed())
		})

		When("spec.canceled is set to true", func() {
			BeforeEach(func() {
				cfTask.Spec.Canceled = true
				Expect(k8sClient.Update(ctx, cfTask)).To(Succeed())
			})

			It("sets the canceled status condition", func() {
				Eventually(func(g Gomega) {
					g.Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(cfTask), cfTask)).To(Succeed())
					g.Expect(meta.IsStatusConditionTrue(cfTask.Status.Conditions, korifiv1alpha1.TaskCanceledConditionType)).To(BeTrue())
				})
			})
		})
	})

	Describe("CFTask TTL", func() {
		BeforeEach(func() {
			Expect(k8sClient.Create(ctx, cfTask)).To(Succeed())

			updatedTask := cfTask.DeepCopy()
			meta.SetStatusCondition(&updatedTask.Status.Conditions, metav1.Condition{
				Type:               korifiv1alpha1.TaskSucceededConditionType,
				Status:             metav1.ConditionTrue,
				LastTransitionTime: metav1.Now(),
				Reason:             "succeeded",
				Message:            "succeeded",
			})
			Expect(k8sClient.Status().Patch(ctx, updatedTask, client.MergeFrom(cfTask))).To(Succeed())
		})

		It("it can get the task shortly after completion", func() {
			Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(cfTask), cfTask)).To(Succeed())
		})

		It("deletes the task after it expires", func() {
			Eventually(func(g Gomega) {
				err := k8sClient.Get(ctx, client.ObjectKeyFromObject(cfTask), cfTask)
				g.Expect(k8serrors.IsNotFound(err)).To(BeTrue())
			}).Should(Succeed())
		})
	})
})

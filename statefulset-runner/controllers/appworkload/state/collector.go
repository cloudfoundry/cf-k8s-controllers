package state

import (
	"context"
	"fmt"
	"time"

	korifiv1alpha1 "code.cloudfoundry.org/korifi/controllers/api/v1alpha1"
	"code.cloudfoundry.org/korifi/statefulset-runner/controllers"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type AppWorkloadStateCollector struct {
	client client.Client
}

func NewAppWorkloadStateCollector(client client.Client) *AppWorkloadStateCollector {
	return &AppWorkloadStateCollector{
		client: client,
	}
}

func (c *AppWorkloadStateCollector) CollectState(ctx context.Context, appWorkloadGUID string) (map[string]korifiv1alpha1.InstanceStatus, error) {
	workloadPods := &corev1.PodList{}
	err := c.client.List(ctx, workloadPods,
		client.MatchingLabels{
			controllers.LabelGUID: appWorkloadGUID,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list pods for workload %q: %w", appWorkloadGUID, err)
	}

	result := map[string]korifiv1alpha1.InstanceStatus{}

	for _, pod := range workloadPods.Items {
		result[pod.Labels["apps.kubernetes.io/pod-index"]] = getPodState(pod)
	}

	return result, nil
}

// Logic from Kubernetes in Action 2nd Edition - Ch 6.
// DOWN => !pod || !pod.conditions.PodScheduled
// CRASHED => any(pod.ContainerStatuses.State isA Terminated)
// RUNNING => pod.conditions.Ready
// STARTING => default
func getPodState(pod corev1.Pod) korifiv1alpha1.InstanceStatus {
	// return running when all containers are ready
	if podConditionStatus(pod, corev1.PodReady) {
		return korifiv1alpha1.InstanceStatus{
			State:  korifiv1alpha1.InstanceStateRunning,
			Uptime: getPodUptime(pod),
		}
	}

	if !podConditionStatus(pod, corev1.PodScheduled) {
		return korifiv1alpha1.InstanceStatus{
			State: korifiv1alpha1.InstanceStateDown,
		}
	}

	if podHasCrashedContainer(pod) {
		return korifiv1alpha1.InstanceStatus{
			State: korifiv1alpha1.InstanceStateCrashed,
		}
	}

	return korifiv1alpha1.InstanceStatus{
		State: korifiv1alpha1.InstanceStateStarting,
	}
}

func podHasCrashedContainer(pod corev1.Pod) bool {
	for _, cond := range pod.Status.ContainerStatuses {
		if cond.State.Waiting != nil && cond.State.Waiting.Reason == "CrashLoopBackOff" {
			return true
		}
	}

	return false
}

func podConditionStatus(pod corev1.Pod, conditionType corev1.PodConditionType) bool {
	for _, cond := range pod.Status.Conditions {
		if cond.Type == conditionType {
			return cond.Status == corev1.ConditionTrue
		}
	}

	return false
}

func getPodUptime(pod corev1.Pod) int {
	for _, cond := range pod.Status.Conditions {
		if cond.Type == corev1.PodReady {
			return int(time.Since(cond.LastTransitionTime.Time).Seconds())
		}
	}

	return 0
}

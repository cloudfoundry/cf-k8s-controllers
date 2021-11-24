package repositories_test

import (
	. "code.cloudfoundry.org/cf-k8s-controllers/api/repositories"
	workloadsv1alpha1 "code.cloudfoundry.org/cf-k8s-controllers/controllers/apis/workloads/v1alpha1"

	"github.com/google/uuid"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	cfAppGUIDLabelKey = "workloads.cloudfoundry.org/app-guid"
)

func generateGUID() string {
	return uuid.NewString()
}

func initializeAppCR(appName string, appGUID string, spaceGUID string) *workloadsv1alpha1.CFApp {
	return &workloadsv1alpha1.CFApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      appGUID,
			Namespace: spaceGUID,
		},
		Spec: workloadsv1alpha1.CFAppSpec{
			Name:         appName,
			DesiredState: "STOPPED",
			Lifecycle: workloadsv1alpha1.Lifecycle{
				Type: "buildpack",
				Data: workloadsv1alpha1.LifecycleData{
					Buildpacks: []string{},
					Stack:      "",
				},
			},
		},
	}
}

func initializeProcessCR(processGUID, spaceGUID, appGUID string) *workloadsv1alpha1.CFProcess {
	return &workloadsv1alpha1.CFProcess{
		ObjectMeta: metav1.ObjectMeta{
			Name:      processGUID,
			Namespace: spaceGUID,
			Labels: map[string]string{
				cfAppGUIDLabelKey: appGUID,
			},
		},
		Spec: workloadsv1alpha1.CFProcessSpec{
			AppRef: corev1.LocalObjectReference{
				Name: appGUID,
			},
			ProcessType: "web",
			Command:     "",
			HealthCheck: workloadsv1alpha1.HealthCheck{
				Type: "process",
				Data: workloadsv1alpha1.HealthCheckData{
					InvocationTimeoutSeconds: 0,
					TimeoutSeconds:           0,
				},
			},
			DesiredInstances: 1,
			MemoryMB:         500,
			DiskQuotaMB:      512,
			Ports:            []int32{8080},
		},
	}
}

func initializeDropletCR(dropletGUID, appGUID, spaceGUID string) workloadsv1alpha1.CFBuild {
	return workloadsv1alpha1.CFBuild{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dropletGUID,
			Namespace: spaceGUID,
		},
		Spec: workloadsv1alpha1.CFBuildSpec{
			AppRef: corev1.LocalObjectReference{Name: appGUID},
			Lifecycle: workloadsv1alpha1.Lifecycle{
				Type: "buildpack",
			},
		},
	}
}

func initializeAppCreateMessage(appName string, spaceGUID string) AppCreateMessage {
	return AppCreateMessage{
		Name:      appName,
		SpaceGUID: spaceGUID,
		State:     "STOPPED",
		Lifecycle: Lifecycle{
			Type: "buildpack",
			Data: LifecycleData{
				Buildpacks: []string{},
				Stack:      "cflinuxfs3",
			},
		},
	}
}

func generateAppEnvSecretName(appGUID string) string {
	return appGUID + "-env"
}

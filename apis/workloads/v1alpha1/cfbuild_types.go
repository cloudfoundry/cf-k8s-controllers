/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CFBuildSpec defines the desired state of CFBuild
type CFBuildSpec struct {
	// Specifies the CFPackage associated with this build
	PackageRef v1.LocalObjectReference `json:"packageRef"`
	// Specifies the CFApp associated with this build
	AppRef v1.LocalObjectReference `json:"appRef"`

	// Specifies the memory request for the staging image
	StagingMemoryMB int `json:"stagingMemoryMB"`
	// Specifies the disk request for the staging image - Do we need this?
	StagingDiskMB int `json:"stagingDiskMB"`

	// Specifies the buildpacks and stack for the build
	Lifecycle Lifecycle `json:"lifecycle"`
}

// CFBuildStatus defines the observed state of CFBuild
type CFBuildStatus struct {
	BuildDropletStatus *BuildDropletStatus `json:"droplet,omitempty"`
	// Conditions capture the current status of the Build
	Conditions []metav1.Condition `json:"conditions"`
}

// BuildDropletStatus defines the observed state of the CFBuild's Droplet or runnable image
type BuildDropletStatus struct {
	// Specifies the Container registry image, and secrets to access
	Registry Registry `json:"registry"`

	// Specifies the process types and associated start commands for the Droplet
	ProcessTypes []ProcessType `json:"processTypes"`

	// Specifies the exposed ports for the application
	Ports []int32 `json:"ports"`
}

// ProcessType is a map of process names and associated start commands for the Droplet
type ProcessType struct {
	Type    string `json:"type"`
	Command string `json:"command"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// CFBuild is the Schema for the cfbuilds API
type CFBuild struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CFBuildSpec   `json:"spec,omitempty"`
	Status CFBuildStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CFBuildList contains a list of CFBuild
type CFBuildList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CFBuild `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CFBuild{}, &CFBuildList{})
}

/*
Copyright 2017 The Kubernetes Authors.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Log is a specification for a Log resource
type Log struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              LogSpec   `json:"spec"`
	Status            LogStatus `json:"status"`
}

// LogSpec is the spec for a Log resource
type LogSpec struct {
	Prometheus Prometheus `json:"prometheus"`
	Warning    Warning    `json: warning`
}

type Prometheus struct {
	Name     string `json:"name"`
	Protocol string `json:"protocol"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	// unit second
	Period int64 `json:"period"`
}

type Warning struct {
	Sustained Sustained `json: "sustained"`
}

type Sustained struct {
	Step         int64 `json: "step"`
	Range        int64 `json: "range"`
	WarningValue int64 `json: warningValue`
}

// LogStatus is the status for a Log resource
type LogStatus struct {
	AvailableReplicas int32 `json:"availableReplicas"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LogList is a list of Log resources
type LogList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Log `json:"items"`
}

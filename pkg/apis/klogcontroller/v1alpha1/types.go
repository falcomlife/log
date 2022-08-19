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

// KLog is a specification for a KLog resource
type Klog struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              KlogSpec   `json:"spec"`
	Status            KlogStatus `json:"status"`
}

// KLogSpec is the spec for a KLog resource
type KlogSpec struct {
	Prometheus Prometheus `json:"prometheus"`
	Template   Template   `json:"template"`
	Warning    Warning    `json:"warning"`
	Context    Context    `json:"context"`
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
	Sustained          Sustained          `json:"sustained"`
	ExtremePointMedian ExtremePointMedian `json:"extremePointMedian"`
}
type Context struct {
	EnvoyFilters []string `json:"envoyFilters"`
}
type Sustained struct {
	Cpu    Cpu    `json:"cpu"`
	Memory Memory `json:"memory"`
	Disk   Disk   `json:"disk"`
}

type ExtremePointMedian struct {
	Cpu    Ecpu    `json:"cpu"`
	Memory Ememory `json:"memory"`
}

type Cpu struct {
	Step         int64 `json:"step"`
	Range        int64 `json:"range"`
	WarningValue int64 `json:"warningValue"`
	LeftTime     int   `json:"leftTime"`
}

type Memory struct {
	Step         int64 `json:"step"`
	Range        int64 `json:"range"`
	WarningValue int64 `json:"warningValue"`
	LeftTime     int   `json:"leftTime"`
}

type Disk struct {
	Range    int64 `json:"range"`
	LeftTime int   `json:"leftTime"`
}

type Ecpu struct {
	WarningValue int64 `json:"warningValue"`
}

type Ememory struct {
	WarningValue int64 `json:"warningValue"`
}

type Template struct {
	Registry         Registry `json:"registry"`
	Env              string   `json:"env"`
	WebRealmName     string   `json:"webRealmName"`
	GatewayRealmName string   `json:"gatewayRealmName"`
	Address          string   `json:"address"`
	Username         string   `json:"username"`
	Password         string   `json:"password"`
}

type Registry struct {
	Java Java `json:"java"`
	Npm  Npm  `json:"npm"`
}

type Java struct {
	PackageImage string `json:"packageImage"`
	ReleaseImage string `json:"releaseImage"`
	DeployImage  string `json:"deployImage"`
}
type Npm struct {
	PackageImage string `json:"packageImage"`
	ReleaseImage string `json:"releaseImage"`
	DeployImage  string `json:"deployImage"`
}

// KLogStatus is the status for a KLog resource
type KlogStatus struct {
	AvailableReplicas int32 `json:"availableReplicas"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KLogList is a list of KLog resources
type KlogList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Klog `json:"items"`
}

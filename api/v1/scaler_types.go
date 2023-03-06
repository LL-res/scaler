/*
Copyright 2023.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
type Metric struct {
	Name      string `json:"name"`
	Threshold string `json:"threshold"`
}
type PodMetric struct {
	Metric []Metric `json:"metric"`
}
type NodeMetric struct {
	Metric []Metric `json:"metric"`
}
type Predictor struct {
	Name      string `json:"name"`
	Algorithm string `json:"algorithm"`
	Step      string `json:"step"`
	Image     string `json:"image"`
	Replica   int    `json:"replica"`
	Port      int32  `json:"port"`
}
type Update struct {
	On bool `json:"on"`
}
type Application struct {
	Name    string `json:"name"`
	Image   string `json:"image"`
	Replica int    `json:"replica"`
	Ports   []Port `json:"ports"`
}
type Port struct {
	Name string `json:"name,omitempty"`
	Port int32  `json:"port"`
}

// ScalerSpec defines the desired state of Scaler
type ScalerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	PodMetric  PodMetric  `json:"pod_metric,omitempty"`
	NodeMetric NodeMetric `json:"node_metric,omitempty"`
	Predictor  Predictor  `json:"algorithm,omitempty"`
	Update     Update     `json:"update,omitempty"`
	// Foo is an example field of Scaler. Edit scaler_types.go to remove/update

}

// ScalerStatus defines the observed state of Scaler
type ScalerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	MonitorHealth     bool `json:"monitor_health"`
	ApplicationHealth bool `json:"application_health"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Scaler is the Schema for the scalers API
type Scaler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ScalerSpec   `json:"spec,omitempty"`
	Status ScalerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ScalerList contains a list of Scaler
type ScalerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Scaler `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Scaler{}, &ScalerList{})
}

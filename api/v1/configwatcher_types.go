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

// ConfigWatcherSpec defines the desired state of ConfigWatcher
type ConfigWatcherSpec struct {

	//Name of the ConfigMap to monitor for changes
	ConfigMap string `json:"configMap,omitempty"`
	// PodSelector defines the label selector for the pods to delete, if the given configmap changes
	PodSelector map[string]string `json:"podSelector,omitempty"`
}

// ConfigWatcherStatus defines the observed state of ConfigWatcher
type ConfigWatcherStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ConfigWatcher is the Schema for the configwatchers API
type ConfigWatcher struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ConfigWatcherSpec   `json:"spec,omitempty"`
	Status ConfigWatcherStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ConfigWatcherList contains a list of ConfigWatcher
type ConfigWatcherList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ConfigWatcher `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ConfigWatcher{}, &ConfigWatcherList{})
}

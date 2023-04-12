// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	healthcheckconfigv1alpha1 "github.com/gardener/gardener/extensions/pkg/apis/config/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Configuration contains information about the network traffic gauger configuration.
type Configuration struct {
	metav1.TypeMeta `json:",inline"`

	// NetworkTrafficGauger contains the configuration for the network traffic gauger
	// +optional
	NetworkTrafficGauger *NetworkTrafficGauger `json:"networkTrafficGauger,omitempty"`

	// HealthCheckConfig is the config for the health check controller.
	// +optional
	HealthCheckConfig *healthcheckconfigv1alpha1.HealthCheckConfig `json:"healthCheckConfig,omitempty"`
}

// NetworkTrafficGauger contains the configuration for the network traffic gauger.
type NetworkTrafficGauger struct {
}

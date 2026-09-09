// SPDX-FileCopyrightText: Contributors to the Gardener project
//
// SPDX-License-Identifier: Apache-2.0

package constants

const (
	// ExtensionType is the name of the extension type.
	ExtensionType = "shoot-networking-traffic-gauger"
	// ServiceName is the name of the service.
	ServiceName = ExtensionType

	// ApplicationName is the name for resource describing the components deployed by the extension controller.
	ApplicationName = "network-traffic-gauger"

	// AgentImageName is the image name for network traffic gauger agent
	AgentImageName = "network-traffic-gauger"

	extensionServiceName = "extension-" + ServiceName
	// NamespaceKubeSystem kube-system namespace
	NamespaceKubeSystem = "kube-system"
	// ManagedResourceNamesSeed is the name used to describe the managed seed resources.
	ManagedResourceNamesSeed = extensionServiceName + "-seed"
	// ManagedResourceNamesAgentShoot is the name used to describe the managed shoot resources for the agents.
	ManagedResourceNamesAgentShoot = extensionServiceName + "-agent-shoot"

	// NetworkTrafficGaugerChartNameSeed is the chart name for network traffic gauger resources in the seed.
	NetworkTrafficGaugerChartNameSeed = "shoot-network-traffic-gauger-seed"
)

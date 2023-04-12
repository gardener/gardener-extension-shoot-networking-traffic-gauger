// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"github.com/gardener/gardener-extension-shoot-networking-traffic-gauger/pkg/apis/config"
)

// Config contains configuration for the network traffic gauger.
type Config struct {
	config.Configuration
}

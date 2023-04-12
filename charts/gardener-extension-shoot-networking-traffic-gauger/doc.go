// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

//go:generate sh -c "../../vendor/github.com/gardener/gardener/hack/generate-controller-registration.sh shoot-networking-traffic-gauger . $(cat ../../VERSION) ../../example/controller-registration.yaml Extension:shoot-networking-traffic-gauger"

// Package chart enables go:generate support for generating the correct controller registration.
package chart

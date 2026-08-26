// SPDX-FileCopyrightText: Copyright Contributors to the Gardener project
//
// SPDX-License-Identifier: Apache-2.0

// +k8s:deepcopy-gen=package
// +k8s:conversion-gen=github.com/gardener/gardener-extension-shoot-networking-traffic-gauger/pkg/apis/config
// +k8s:defaulter-gen=TypeMeta
// +k8s:openapi-gen=true

//go:generate gen-crd-api-reference-docs -api-dir . -config ../../../../hack/api-reference/api.json -template-dir $GARDENER_HACK_DIR/api-reference/template -out-file ../../../../hack/api-reference/api.md

// Package v1alpha1 contains the shoot networking filter extension configuration.
// +groupName=shoot-networking-traffic-gauger.extensions.config.gardener.cloud
package v1alpha1

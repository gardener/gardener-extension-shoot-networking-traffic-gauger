// SPDX-FileCopyrightText: Contributors to the Gardener project
//
// SPDX-License-Identifier: Apache-2.0

package lifecycle

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	networkingv1 "k8s.io/api/networking/v1"
)

var _ = Describe("activator methods", func() {
	It("#buildDaemonSet", func() {
		obj := buildDaemonSet("image:tag", "serviceAccount")
		Expect(obj).NotTo(BeNil())
		var ds *appsv1.DaemonSet
		switch v := obj.(type) {
		case *appsv1.DaemonSet:
			ds = v
		}
		Expect(ds).NotTo(BeNil())
		Expect(len(ds.Spec.Template.Spec.Containers)).To(Equal(1))
		Expect(ds.Spec.Template.Spec.Containers[0].Image).To(Equal("image:tag"))
	})

	It("#buildNetworkPolicy", func() {
		obj := buildNetworkPolicy()
		Expect(obj).NotTo(BeNil())
		var np *networkingv1.NetworkPolicy
		switch v := obj.(type) {
		case *networkingv1.NetworkPolicy:
			np = v
		}
		Expect(np).NotTo(BeNil())
		Expect(np.Name).To(Equal("gardener.cloud--allow-to-from-network-traffic-gauger"))
	})
})

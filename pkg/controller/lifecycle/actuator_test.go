// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package lifecycle

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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

	It("#buildPodSecurityPolicy", func() {
		objs := buildPodSecurityPolicy("serviceAccount")
		Expect(len(objs)).NotTo(BeZero())
		var cr *rbacv1.ClusterRole
		var crb *rbacv1.ClusterRoleBinding
		var sa *corev1.ServiceAccount
		var psp *policyv1beta1.PodSecurityPolicy
		for _, obj := range objs {
			switch v := obj.(type) {
			case *rbacv1.ClusterRole:
				cr = v
			case *rbacv1.ClusterRoleBinding:
				crb = v
			case *corev1.ServiceAccount:
				sa = v
			case *policyv1beta1.PodSecurityPolicy:
				psp = v
			}
		}
		Expect(cr).NotTo(BeNil())
		Expect(cr.Name).To(Equal("gardener.cloud:psp:kube-system:network-traffic-gauger"))
		Expect(crb).NotTo(BeNil())
		Expect(crb.Name).To(Equal("gardener.cloud:psp:kube-system:network-traffic-gauger"))
		Expect(sa).NotTo(BeNil())
		Expect(sa.Name).To(Equal("serviceAccount"))
		Expect(psp).NotTo(BeNil())
		Expect(cr.Name).To(Equal("gardener.kube-system.network-traffic-gauger"))
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

// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package lifecycle

import (
	"context"
	_ "embed"
	"fmt"
	"path/filepath"
	"time"

	"github.com/gardener/gardener/extensions/pkg/controller"
	"github.com/gardener/gardener/extensions/pkg/controller/extension"
	corev1betaconstants "github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	"github.com/gardener/gardener/pkg/chartrenderer"
	"github.com/gardener/gardener/pkg/client/kubernetes"
	managedresources "github.com/gardener/gardener/pkg/utils/managedresources"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/rest"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/gardener/gardener-extension-shoot-networking-traffic-gauger/charts"
	"github.com/gardener/gardener-extension-shoot-networking-traffic-gauger/imagevector"
	"github.com/gardener/gardener-extension-shoot-networking-traffic-gauger/pkg/apis/config"
	"github.com/gardener/gardener-extension-shoot-networking-traffic-gauger/pkg/constants"
)

const (
	labelKeyK8sApp = "k8s-app"
	metricsPort    = 16160
)

// NewActuator returns an actuator responsible for Extension resources.
func NewActuator(config config.Configuration) extension.Actuator {
	return &actuator{
		serviceConfig: config,
	}
}

type actuator struct {
	client        client.Client
	config        *rest.Config
	decoder       runtime.Decoder
	serviceConfig config.Configuration
}

// Reconcile the Extension resource.
func (a *actuator) Reconcile(ctx context.Context, log logr.Logger, ex *extensionsv1alpha1.Extension) error {
	namespace := ex.GetNamespace()

	cluster, err := controller.GetCluster(ctx, a.client, namespace)
	if err != nil {
		return err
	}

	if !controller.IsHibernated(cluster) {
		if err := a.createShootResources(ctx, cluster, namespace); err != nil {
			return err
		}
	}

	if err := a.createSeedResources(ctx, log, cluster, namespace); err != nil {
		return err
	}

	return nil
}

func (a *actuator) createSeedResources(ctx context.Context, log logr.Logger, cluster *controller.Cluster, namespace string) error {
	renderer, err := chartrenderer.NewForConfig(a.config)
	if err != nil {
		return fmt.Errorf("could not create chart renderer: %w", err)
	}

	log.Info("Component is being applied", "component", "network-traffic-gauger-seed", "namespace", namespace)

	return a.createManagedResource(ctx, namespace, constants.ManagedResourceNamesSeed, "seed", renderer, constants.NetworkTrafficGaugerChartNameSeed, namespace, nil, nil)
}

func (a *actuator) createShootResources(ctx context.Context, cluster *controller.Cluster, namespace string) error {
	shootResources, err := a.getShootAgentResources()
	if err != nil {
		return err
	}
	return managedresources.CreateForShoot(ctx, a.client, namespace, constants.ManagedResourceNamesAgentShoot, "gardener-extension-shoot-networking-traffic-gauger", false, shootResources)
}

func (a *actuator) createManagedResource(ctx context.Context, namespace, name, class string, renderer chartrenderer.Interface, chartName, chartNamespace string, chartValues map[string]interface{}, injectedLabels map[string]string) error {
	chartPath := filepath.Join(charts.ChartsPath, chartName)
	chart, err := renderer.RenderEmbeddedFS(charts.Internal, chartPath, chartName, chartNamespace, chartValues)
	if err != nil {
		return err
	}

	data := map[string][]byte{chartName: chart.Manifest()}
	keepObjects := false
	forceOverwriteAnnotations := false
	return managedresources.Create(ctx, a.client, namespace, name, nil, false, class, data, &keepObjects, injectedLabels, &forceOverwriteAnnotations)
}

// Delete the Extension resource.
func (a *actuator) Delete(ctx context.Context, log logr.Logger, ex *extensionsv1alpha1.Extension) error {
	namespace := ex.GetNamespace()

	err := a.deleteShootResources(ctx, log, namespace, false)
	if err != nil {
		return err
	}

	return a.deleteSeedResources(ctx, log, namespace)
}

// ForceDelete the Extension resource.
func (a *actuator) ForceDelete(ctx context.Context, log logr.Logger, ex *extensionsv1alpha1.Extension) error {
	namespace := ex.GetNamespace()

	err := a.deleteShootResources(ctx, log, namespace, true)
	if err != nil {
		return err
	}

	return a.deleteSeedResources(ctx, log, namespace)
}

func (a *actuator) deleteShootResources(ctx context.Context, log logr.Logger, namespace string, forceDelete bool) error {
	log.Info("Deleting managed resource for shoot", "namespace", namespace)
	if err := managedresources.DeleteForShoot(ctx, a.client, namespace, constants.ManagedResourceNamesAgentShoot); err != nil {
		return err
	}

	// We don't need to wait for the shoot managed resource deletion because managed resources are finalized by gardenlet
	// in later step in the Shoot force deletion flow.
	if forceDelete {
		return nil
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()
	if err := managedresources.WaitUntilDeleted(timeoutCtx, a.client, namespace, constants.ManagedResourceNamesAgentShoot); err != nil {
		return err
	}
	return nil
}

func (a *actuator) deleteSeedResources(ctx context.Context, log logr.Logger, namespace string) error {
	log.Info("Deleting managed resource for seed", "namespace", namespace)

	if err := managedresources.Delete(ctx, a.client, namespace, constants.ManagedResourceNamesSeed, false); err != nil {
		return err
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()
	return managedresources.WaitUntilDeleted(timeoutCtx, a.client, namespace, constants.ManagedResourceNamesSeed)
}

// Restore the Extension resource.
func (a *actuator) Restore(ctx context.Context, log logr.Logger, ex *extensionsv1alpha1.Extension) error {
	return a.Reconcile(ctx, log, ex)
}

// Migrate the Extension resource.
func (a *actuator) Migrate(ctx context.Context, log logr.Logger, ex *extensionsv1alpha1.Extension) error {
	// Keep objects for shoot managed resources so that they are not deleted from the shoot during the migration
	if err := managedresources.SetKeepObjects(ctx, a.client, ex.GetNamespace(), constants.ManagedResourceNamesAgentShoot, true); err != nil {
		return err
	}

	return a.Delete(ctx, log, ex)
}

// InjectConfig injects the rest config to this actuator.
func (a *actuator) InjectConfig(config *rest.Config) error {
	a.config = config
	return nil
}

// InjectClient injects the controller runtime client into the reconciler.
func (a *actuator) InjectClient(client client.Client) error {
	a.client = client
	return nil
}

// InjectScheme injects the given scheme into the reconciler.
func (a *actuator) InjectScheme(scheme *runtime.Scheme) error {
	a.decoder = serializer.NewCodecFactory(scheme, serializer.EnableStrict).UniversalDecoder()
	return nil
}

func (a *actuator) getShootAgentResources() (map[string][]byte, error) {
	shootRegistry := managedresources.NewRegistry(kubernetes.ShootScheme, kubernetes.ShootCodec, kubernetes.ShootSerializer)

	image, err := imagevector.ImageVector().FindImage(constants.AgentImageName)
	if err != nil {
		return nil, err
	}

	objects := []client.Object{}
	serviceAccountName := ""

	daemonset := buildDaemonSet(image.String(), serviceAccountName)
	networkPolicy := buildNetworkPolicy()
	objects = append(objects, daemonset, networkPolicy)

	shootResources, err := shootRegistry.AddAllAndSerialize(objects...)
	if err != nil {
		return nil, err
	}
	return shootResources, nil
}

func buildDaemonSet(image string, serviceAccountName string) client.Object {
	var (
		requestCPU                   = resource.MustParse("10m")
		requestMemory                = resource.MustParse("32Mi")
		limitMemory                  = resource.MustParse("64Mi")
		automountServiceAccountToken = false
		fileStoreDirectory           = "/var/log/net-gauger"
		hostPathType                 = corev1.HostPathDirectoryOrCreate
		labels                       = map[string]string{
			labelKeyK8sApp:        constants.ApplicationName,
			"gardener.cloud/role": constants.ApplicationName,
		}
	)

	return &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      constants.ApplicationName,
			Namespace: metav1.NamespaceSystem,
		},
		Spec: appsv1.DaemonSetSpec{
			RevisionHistoryLimit: ptr.To[int32](5),
			Selector:             &metav1.LabelSelector{MatchLabels: labels},
			UpdateStrategy: appsv1.DaemonSetUpdateStrategy{
				Type: appsv1.RollingUpdateDaemonSetStrategyType,
				RollingUpdate: &appsv1.RollingUpdateDaemonSet{
					MaxUnavailable: &intstr.IntOrString{Type: intstr.String, StrVal: "100%"},
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					HostNetwork:                   true,
					PriorityClassName:             corev1betaconstants.PriorityClassNameShootSystem900,
					TerminationGracePeriodSeconds: ptr.To[int64](0),
					Tolerations: []corev1.Toleration{
						{
							Effect:   corev1.TaintEffectNoSchedule,
							Operator: corev1.TolerationOpExists,
						},
						{
							Effect:   corev1.TaintEffectNoExecute,
							Operator: corev1.TolerationOpExists,
						},
					},
					AutomountServiceAccountToken: &automountServiceAccountToken,
					ServiceAccountName:           serviceAccountName,
					InitContainers: []corev1.Container{{
						Name:            constants.ApplicationName + "-init",
						Image:           image,
						ImagePullPolicy: corev1.PullIfNotPresent,
						Command: []string{
							"/net-gauger",
							"setup",
						},
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    requestCPU,
								corev1.ResourceMemory: requestMemory,
							},
							Limits: corev1.ResourceList{
								corev1.ResourceMemory: limitMemory,
							},
						},
						SecurityContext: &corev1.SecurityContext{
							Privileged: ptr.To(true),
							Capabilities: &corev1.Capabilities{
								Add: []corev1.Capability{"NET_ADMIN"},
							},
						},
					}},
					Containers: []corev1.Container{{
						Name:            constants.ApplicationName,
						Image:           image,
						ImagePullPolicy: corev1.PullIfNotPresent,
						Command: []string{
							"/net-gauger",
							"run-agent",
							"--dump-period=1m",
							"--event-channel-buffer-size=1024",
							"--event-receive-buffer-size=1048576",
							fmt.Sprintf("--file-store-directory=%s", fileStoreDirectory),
							"--file-store-channel-buffer-size=1024",
							fmt.Sprintf("--metrics-port=%d", metricsPort),
							"--enable-service-metrics=true",
							"--enable-byte-metrics=true",
							"--enable-packet-metrics=true",
							"--enable-flow-count-metrics=true",
						},
						LivenessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								TCPSocket: &corev1.TCPSocketAction{
									Port: intstr.FromInt(metricsPort),
								},
							},
						},
						ReadinessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								TCPSocket: &corev1.TCPSocketAction{
									Port: intstr.FromInt(metricsPort),
								},
							},
						},
						Ports: []corev1.ContainerPort{
							{
								Name:          "metrics",
								ContainerPort: int32(metricsPort),
								Protocol:      "TCP",
							},
						},
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    requestCPU,
								corev1.ResourceMemory: requestMemory,
							},
							Limits: corev1.ResourceList{
								corev1.ResourceMemory: limitMemory,
							},
						},
						SecurityContext: &corev1.SecurityContext{
							Capabilities: &corev1.Capabilities{
								Add: []corev1.Capability{"NET_ADMIN"},
							},
						},
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      "storage",
								ReadOnly:  false,
								MountPath: fileStoreDirectory,
							},
						},
					}},
					Volumes: []corev1.Volume{
						{
							Name: "storage",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: fileStoreDirectory,
									Type: &hostPathType,
								},
							},
						},
					},
					SecurityContext: &corev1.PodSecurityContext{
						SeccompProfile: &corev1.SeccompProfile{
							Type: corev1.SeccompProfileTypeRuntimeDefault,
						},
					},
				},
			},
		},
	}
}

func buildNetworkPolicy() client.Object {
	tcp := corev1.ProtocolTCP
	port := intstr.FromInt(metricsPort)
	return &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "gardener.cloud--allow-to-from-network-traffic-gauger",
			Namespace: constants.NamespaceKubeSystem,
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					labelKeyK8sApp: constants.ApplicationName,
				},
			},
			Egress: []networkingv1.NetworkPolicyEgressRule{},
			Ingress: []networkingv1.NetworkPolicyIngressRule{
				{
					Ports: []networkingv1.NetworkPolicyPort{
						{
							Protocol: &tcp,
							Port:     &port,
						},
					},
				},
			},
			PolicyTypes: []networkingv1.PolicyType{networkingv1.PolicyTypeEgress, networkingv1.PolicyTypeIngress},
		},
	}
}

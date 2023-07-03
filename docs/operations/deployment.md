# Gardener Networking Traffic Gauger for Shoots

## Introduction
Gardener allows shoot clusters to add network traffice observability using the network traffic gauger.
To support this the Gardener must be installed with the `shoot-networking-traffic-gauger` extension.

## Configuration

To generally enable the networking traffic gauger for shoot objects the `shoot-networking-traffic-gauger` extension must be registered by providing an appropriate [extension registration](https://github.com/gardener/gardener-extension-shoot-networking-traffic-gauger/blob/master/example/controller-registration.yaml) in the garden cluster.

Here it is possible to decide whether the extension should be always available for all shoots or whether the extension must be separately enabled per shoot.

If the extension should be used for all shoots the `globallyEnabled` flag should be set to `true`.

```yaml
apiVersion: core.gardener.cloud/v1beta1
kind: ControllerRegistration
...
spec:
  resources:
    - kind: Extension
      type: shoot-networking-traffic-gauger
      globallyEnabled: true
```

### ControllerRegistration
An example of a `ControllerRegistration` for the `shoot-networking-traffic-gauger` can be found here: https://github.com/gardener/gardener-extension-shoot-networking-traffic-gauger/blob/master/example/controller-registration.yaml

The `ControllerRegistration` contains a Helm chart which eventually deploys the `shoot-networking-traffic-gauger` to seed clusters.

### Enablement for a Shoot

If the shoot network traffic gauger is not globally enabled by default (depends on the extension registration on the garden cluster), it can be enabled per shoot. To enable the service for a shoot, the shoot manifest must explicitly add the `shoot-networking-traffic-gauger` extension.

```yaml
apiVersion: core.gardener.cloud/v1beta1
kind: Shoot
...
spec:
  extensions:
    - type: shoot-networking-traffic-gauger
...
```

If the shoot network traffic gauger is globally enabled by default, it can be disabled per shoot. To disable the service for a shoot, the shoot manifest must explicitly state it.

```yaml
apiVersion: core.gardener.cloud/v1beta1
kind: Shoot
...
spec:
  extensions:
    - type: shoot-networking-traffic-gauger
      disabled: true
...
```

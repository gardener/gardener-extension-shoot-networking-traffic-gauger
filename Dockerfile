# SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

############# builder
FROM golang:1.22.1 AS builder

WORKDIR /go/src/github.com/gardener/gardener-extension-shoot-networking-traffic-gauger
COPY . .
RUN make install

############# gardener-extension-shoot-networking-traffic-gauger
FROM gcr.io/distroless/static-debian11:nonroot AS gardener-extension-shoot-networking-traffic-gauger

WORKDIR /
COPY charts /charts
COPY --from=builder /go/bin/gardener-extension-shoot-networking-traffic-gauger /gardener-extension-shoot-networking-traffic-gauger
ENTRYPOINT ["/gardener-extension-shoot-networking-traffic-gauger"]

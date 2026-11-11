FROM golang:1.27.4@sha256:08bcb99a38f85a265f9a93d1d3323e2440eec333380f3ea26ee0dec9264bb41c AS builder

WORKDIR /go/src/github.com/gardener/gardener-extension-shoot-networking-traffic-gauger

# Copy go mod and sum files
COPY go.mod go.sum ./
# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

COPY . .

ARG EFFECTIVE_VERSION
RUN make install EFFECTIVE_VERSION=$EFFECTIVE_VERSION

############# base
FROM gcr.io/distroless/static-debian11:nonroot@sha256:63ebe035fbdd056ed682e6a87b286d07d3f05f12cb46f26b2b44fc10fc4a59ed AS base

############# gardener-extension-shoot-networking-traffic-gauger
FROM base AS gardener-extension-shoot-networking-traffic-gauger
WORKDIR /

COPY charts /charts
COPY --from=builder /go/bin/gardener-extension-shoot-networking-traffic-gauger /gardener-extension-shoot-networking-traffic-gauger
ENTRYPOINT ["/gardener-extension-shoot-networking-traffic-gauger"]

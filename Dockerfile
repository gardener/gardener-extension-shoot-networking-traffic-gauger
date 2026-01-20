FROM golang:1.25.6 AS builder

WORKDIR /go/src/github.com/gardener/gardener-extension-shoot-networking-traffic-gauger

# Copy go mod and sum files
COPY go.mod go.sum ./
# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

COPY . .

ARG EFFECTIVE_VERSION
RUN make install EFFECTIVE_VERSION=$EFFECTIVE_VERSION

############# base
FROM gcr.io/distroless/static-debian11:nonroot AS base

############# gardener-extension-shoot-networking-traffic-gauger
FROM base AS gardener-extension-shoot-networking-traffic-gauger
WORKDIR /

COPY charts /charts
COPY --from=builder /go/bin/gardener-extension-shoot-networking-traffic-gauger /gardener-extension-shoot-networking-traffic-gauger
ENTRYPOINT ["/gardener-extension-shoot-networking-traffic-gauger"]

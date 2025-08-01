FROM golang:1.23-alpine AS build-env

# Customize to your build env

# TARGETPLATFORM should be one of linux/amd64 or linux/arm64.
ARG TARGETPLATFORM="linux/amd64"
# Version to build. Default is empty
ARG VERSION="v3.0.0-snapshot.250729"
ARG BABYLON_BUILD_OPTIONS=""
ARG COSMOS_BUILD_OPTIONS=""

# Use muslc for static libs
ARG BUILD_TAGS="muslc"
ARG LEDGER_ENABLED="false"


# Install cli tools for building and final image
RUN apk add --update --no-cache make git bash gcc linux-headers eudev-dev ncurses-dev openssh curl jq
RUN apk add --no-cache musl-dev

# Build
WORKDIR /go/src/github.com/babylonlabs-io/babylon
# Clone repo
RUN git clone https://github.com/babylonlabs-io/babylon.git /go/src/github.com/babylonlabs-io/babylon
# Checkout version (can be tag or commit hash)
RUN if [ -n "${VERSION}" ]; then \
    if git rev-parse --verify "refs/tags/${VERSION}" >/dev/null 2>&1; then \
    # If it's a tag
    git fetch origin tag ${VERSION} --no-tags; \
    else \
    # If it's a commit
    git fetch origin ${VERSION}; \
    fi; \
    git checkout -f ${VERSION}; \
    fi

# Cosmwasm - Download correct libwasmvm version
RUN WASMVM_VERSION=$(go list -m github.com/CosmWasm/wasmvm/v2 | cut -d ' ' -f 2) && \
    wget https://github.com/CosmWasm/wasmvm/releases/download/$WASMVM_VERSION/libwasmvm_muslc.$(uname -m).a \
    -O /lib/libwasmvm_muslc.$(uname -m).a && \
    # verify checksum
    wget https://github.com/CosmWasm/wasmvm/releases/download/$WASMVM_VERSION/checksums.txt -O /tmp/checksums.txt && \
    sha256sum /lib/libwasmvm_muslc.$(uname -m).a | grep $(cat /tmp/checksums.txt | grep libwasmvm_muslc.$(uname -m) | cut -d ' ' -f 1)

RUN LEDGER_ENABLED=$LEDGER_ENABLED \
    BABYLON_BUILD_OPTIONS=$BABYLON_BUILD_OPTIONS \
    COSMOS_BUILD_OPTIONS=$COSMOS_BUILD_OPTIONS \
    BUILD_TAGS=$BUILD_TAGS \
    LINK_STATICALLY=true \
    make build

FROM alpine:3.14 AS run
# Create a user
RUN addgroup --gid 1137 -S babylon && adduser --uid 1137 -S babylon -G babylon
RUN apk add bash curl jq

# Label should match your github repo
ARG VERSION
LABEL org.opencontainers.image.source="https://github.com/babylonlabs-io/babylond:${VERSION}"

# Install Libraries
# COPY --from=build-env /usr/lib/libgcc_s.so.1 /lib/
# COPY --from=build-env /lib/ld-musl*.so.1* /lib

COPY --from=build-env /go/src/github.com/babylonlabs-io/babylon/build/babylond /bin/babylond

# Set home directory and user
WORKDIR /home/babylon
RUN chown -R babylon /home/babylon
USER babylon

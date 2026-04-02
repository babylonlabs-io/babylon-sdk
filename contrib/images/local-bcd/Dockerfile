FROM debian:bookworm-slim AS build-env

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates=20230311+deb12u1 \
    gcc=4:12.2.0-3 \
    git=1:2.39.5-0+deb12u3 \
    libc6-dev=2.36-9+deb12u13 \
    make=4.3-4.1 \
    wget=1.21.3-1+deb12u1 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /work

ARG TARGETARCH

# Download and install Go
ENV GOLANG_VERSION=1.21.4
RUN wget -q https://golang.org/dl/go${GOLANG_VERSION}.linux-${TARGETARCH}.tar.gz && \
    tar -C /usr/local -xzf go${GOLANG_VERSION}.linux-${TARGETARCH}.tar.gz && \
    rm go${GOLANG_VERSION}.linux-${TARGETARCH}.tar.gz
# Set Go environment variables
ENV PATH=/usr/local/go/bin:$PATH
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH

ENV GO111MODULE=on

# Install bcd
COPY . /work
RUN make install

FROM debian:bookworm-slim AS run

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates=20230311+deb12u1 \
    curl=7.88.1-10+deb12u14 \
    jq=1.6-2.1+deb12u1 \
    wget=1.21.3-1+deb12u1 \
    && rm -rf /var/lib/apt/lists/*

# Install libraries
# Cosmwasm - Download correct libwasmvm version
COPY --from=build-env /work/demo/go.mod /tmp
RUN WASMVM_VERSION=$(grep github.com/CosmWasm/wasmvm /tmp/go.mod | cut -d' ' -f2) && \
    wget -q "https://github.com/CosmWasm/wasmvm/releases/download/${WASMVM_VERSION}/libwasmvm.$(uname -m).so" \
    -O "/lib/libwasmvm.$(uname -m).so" && \
    # verify checksum
    wget -q "https://github.com/CosmWasm/wasmvm/releases/download/${WASMVM_VERSION}/checksums.txt" -O /tmp/checksums.txt && \
    sha256sum "/lib/libwasmvm.$(uname -m).so" | grep "$(cat /tmp/checksums.txt | grep "libwasmvm.$(uname -m)" | cut -d ' ' -f 1)"
RUN rm -f /tmp/go.mod

# Install binaries
COPY --from=build-env /go/bin/bcd /usr/bin/bcd

WORKDIR /bcd
COPY --from=build-env /work/contrib/images/local-bcd/wrapper.sh /bcd
COPY --from=build-env /work/contrib/images/local-bcd/setup-bcd.sh /bcd
COPY --from=build-env /work/tests/testdata/babylon_contract.wasm /bcd
COPY --from=build-env /work/tests/testdata/btc_staking.wasm /bcd

ENV CONSUMER_CONF=/data/bcd

ENTRYPOINT ["/bcd/wrapper.sh"]
CMD []
STOPSIGNAL SIGTERM

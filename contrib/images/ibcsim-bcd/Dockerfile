FROM debian:bullseye-slim AS build-env

RUN apt-get update && apt-get install -y git make gcc wget

WORKDIR /work

ARG TARGETARCH

# Download and install Go
ENV GOLANG_VERSION=1.23.4
RUN wget -q https://golang.org/dl/go${GOLANG_VERSION}.linux-${TARGETARCH}.tar.gz && \
    tar -C /usr/local -xzf go${GOLANG_VERSION}.linux-${TARGETARCH}.tar.gz && \
    rm go${GOLANG_VERSION}.linux-${TARGETARCH}.tar.gz
# Set Go environment variables
ENV PATH=/usr/local/go/bin:$PATH
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH

ENV GO111MODULE=on
ENV RELAYER_TAG=v2.5.3

# Install the relayer
RUN git clone https://github.com/cosmos/relayer.git
RUN cd relayer && git fetch origin && git checkout ${RELAYER_TAG} && make install && cd -

# build bcd
COPY . /work/babylon-sdk
RUN cd babylon-sdk && \
    make build && \
    cd -

FROM debian:bullseye-slim AS run

RUN apt-get update && apt-get install -y bash curl jq wget

# Install libraries
# Cosmwasm - Download correct libwasmvm version
COPY --from=build-env /work/babylon-sdk/demo/go.mod /tmp
RUN WASMVM_VERSION=$(grep github.com/CosmWasm/wasmvm /tmp/go.mod | cut -d' ' -f2) && \
    wget https://github.com/CosmWasm/wasmvm/releases/download/$WASMVM_VERSION/libwasmvm.$(uname -m).so \
    -O /lib/libwasmvm.$(uname -m).so && \
    # verify checksum
    wget https://github.com/CosmWasm/wasmvm/releases/download/$WASMVM_VERSION/checksums.txt -O /tmp/checksums.txt && \
    sha256sum /lib/libwasmvm.$(uname -m).so | grep $(cat /tmp/checksums.txt | grep libwasmvm.$(uname -m) | cut -d ' ' -f 1)
RUN rm -f /tmp/go.mod

# Copy binaries
COPY --from=build-env /go/bin/rly /usr/bin/rly
COPY --from=build-env /work/babylon-sdk/demo/build/bcd /usr/bin/bcd

WORKDIR /ibcsim-bcd
COPY contrib/images/ibcsim-bcd/wrapper.sh /ibcsim-bcd/wrapper.sh
COPY contrib/images/ibcsim-bcd/setup-bcd.sh /ibcsim-bcd/setup-bcd.sh

# Make all scripts executable
RUN chmod +x /ibcsim-bcd/*.sh
COPY --from=build-env /work/babylon-sdk/tests/testdata/babylon_contract.wasm /ibcsim-bcd/babylon_contract.wasm
COPY --from=build-env /work/babylon-sdk/tests/testdata/btc_light_client.wasm /ibcsim-bcd/btc_light_client.wasm
COPY --from=build-env /work/babylon-sdk/tests/testdata/btc_staking.wasm /ibcsim-bcd/btc_staking.wasm
COPY --from=build-env /work/babylon-sdk/tests/testdata/btc_finality.wasm /ibcsim-bcd/btc_finality.wasm

ENV BABYLON_HOME=/data/node1/babylond
ENV BABYLON_NODE_RPC="http://babylondnode1:26657"
ENV RELAYER_CONF_DIR=/data/relayer
ENV CONSUMER_CONF=/data/bcd
ENV UPDATE_CLIENTS_INTERVAL=20s

ENTRYPOINT ["/ibcsim-bcd/wrapper.sh"]
CMD []
STOPSIGNAL SIGTERM

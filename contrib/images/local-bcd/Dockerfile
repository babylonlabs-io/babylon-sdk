FROM debian:bullseye-slim AS build-env

RUN apt-get update && apt-get install -y git make gcc wget

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

FROM debian:bullseye-slim AS run

RUN apt-get update && apt-get install -y bash curl jq wget

# Install libraries
# Cosmwasm - Download correct libwasmvm version
COPY --from=build-env /work/demo/go.mod /tmp
RUN WASMVM_VERSION=$(grep github.com/CosmWasm/wasmvm /tmp/go.mod | cut -d' ' -f2) && \
    wget https://github.com/CosmWasm/wasmvm/releases/download/$WASMVM_VERSION/libwasmvm.$(uname -m).so \
    -O /lib/libwasmvm.$(uname -m).so && \
    # verify checksum
    wget https://github.com/CosmWasm/wasmvm/releases/download/$WASMVM_VERSION/checksums.txt -O /tmp/checksums.txt && \
    sha256sum /lib/libwasmvm.$(uname -m).so | grep $(cat /tmp/checksums.txt | grep libwasmvm.$(uname -m) | cut -d ' ' -f 1)
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

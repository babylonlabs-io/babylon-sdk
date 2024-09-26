FROM golang:1.22.7-alpine AS go-builder

RUN apk add --no-cache ca-certificates build-base git

WORKDIR /code

# Copy over code
COPY . /code

#ADD https://github.com/CosmWasm/wasmvm/releases/download/v$wasmvm/libwasmvm_muslc.$arch.a /lib/libwasmvm_muslc.$arch.a
## Download
RUN WASMVM_VERSION=$(grep github.com/CosmWasm/wasmvm demo/go.mod | cut -d' ' -f2) && \
  wget https://github.com/CosmWasm/wasmvm/releases/download/$WASMVM_VERSION/libwasmvm_muslc.$(uname -m).a \
  -O /lib/libwasmvm_muslc.$(uname -m).a && \
  # verify checksum
  wget https://github.com/CosmWasm/wasmvm/releases/download/$WASMVM_VERSION/checksums.txt -O /tmp/checksums.txt && \
  sha256sum /lib/libwasmvm_muslc.$(uname -m).a | grep $(cat /tmp/checksums.txt | grep libwasmvm_muslc.$(uname -m) | cut -d ' ' -f 1)

# force it to use static lib (from above) not standard libgo_cosmwasm.so file
# then log output of file /code/bin/bcd
# then ensure static linking
RUN LEDGER_ENABLED=false BUILD_TAGS=muslc LINK_STATICALLY=true make build \
  && file /code/demo/build/bcd \
  && echo "Ensuring binary is statically linked ..." \
  && (file /code/demo/build/bcd | grep "statically linked")

# --------------------------------------------------------
FROM alpine:3.17

COPY --from=go-builder /code/demo/build/bcd /usr/bin/bcd

# Install dependencies used for Starship
RUN apk add --no-cache curl make bash jq sed

WORKDIR /opt

# rest server, tendermint p2p, tendermint rpc
EXPOSE 1317 26656 26657

CMD ["/usr/bin/bcd", "version"]

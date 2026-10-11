FROM golang:1.26-bookworm@sha256:a688600ca24f8a4d3ca77f95b0dd40704a9fc787c826660eb7ba0b641b8b175d AS build

ENV GOTOOLCHAIN=local \
    CGO_ENABLED=0

COPY bin/install-toolchain /tmp/install-toolchain
RUN sh /tmp/install-toolchain /usr/local/bin

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY contracts/foundry.toml contracts/soldeer.lock contracts/remappings.txt contracts/
RUN cd contracts && forge soldeer install

COPY . .

RUN sh bin/generate-bindings \
    && sh bin/check-boundaries \
    && go build -trimpath -ldflags="-s -w" -o /out/notary-verify ./cmd/notary-verify

FROM debian:bookworm-slim@sha256:3783cc01769c7b2b1b83a5c5ad96c815348e28ed7da68e2e3687004faa906251 AS final

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/* \
    && useradd --system --uid 10001 --no-create-home notary

COPY --from=build /out/notary-verify /usr/local/bin/notary-verify

USER notary

ENTRYPOINT ["notary-verify"]

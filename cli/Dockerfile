FROM --platform=$BUILDPLATFORM golang:1.27-bookworm@sha256:d9b6a1beb9474adc2a1ce842573639b44b2a0dd49b6669f51d3821baf9ec37c9 AS build

ENV GOTOOLCHAIN=local \
    CGO_ENABLED=0

RUN apt-get update \
    && apt-get install -y --no-install-recommends unzip \
    && rm -rf /var/lib/apt/lists/*

COPY bin/install-toolchain /tmp/install-toolchain
RUN sh /tmp/install-toolchain /usr/local/bin foundry abigen protobuf sqlc \
    && go clean -cache -modcache

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY contracts/foundry.toml contracts/soldeer.lock contracts/remappings.txt contracts/
RUN cd contracts && forge soldeer install

COPY . .

RUN sh bin/generate

ARG TARGETOS
ARG TARGETARCH
RUN --mount=type=cache,target=/root/.cache/go-build \
    GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -ldflags="-s -w" -o /out/notary-verify ./cmd/notary-verify \
    && GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -ldflags="-s -w" -o /out/notary-sponsor ./cmd/notary-sponsor

FROM gcr.io/distroless/static-debian12:nonroot@sha256:afa5c872c891853ca7fcf1f12c3edb23f7eeef36189728842dd51042ff57f7ab

LABEL org.opencontainers.image.source="https://github.com/WildanFrananda/kinetix-notary-service" \
      org.opencontainers.image.description="notary-verify and notary-sponsor: check a Kinetix proof of liabilities against the chain, and send a holder's court calls through the paymaster"
COPY --from=build /out/notary-verify /out/notary-sponsor /usr/local/bin/

CMD ["notary-verify"]

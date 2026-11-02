FROM golang:1.27-bookworm@sha256:d9b6a1beb9474adc2a1ce842573639b44b2a0dd49b6669f51d3821baf9ec37c9 AS build

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

ARG VERSION=unknown
RUN --mount=type=cache,target=/root/.cache/go-build \
    sh bin/generate \
    && sh bin/check-boundaries \
    && go build -trimpath -ldflags="-s -w -X main.version=${VERSION}" -o /out/notaryd ./cmd/notaryd \
    && go build -trimpath -ldflags="-s -w" -o /out/notary-migrate ./cmd/notary-migrate \
    && go build -trimpath -ldflags="-s -w" -o /out/notary-verify ./cmd/notary-verify \
    && go build -trimpath -ldflags="-s -w -X main.version=${VERSION}" -o /out/notary-watch ./cmd/notary-watch \
    && go build -trimpath -ldflags="-s -w" -o /out/notary-sponsor ./cmd/notary-sponsor \
    && go build -trimpath -ldflags="-s -w" -o /out/notary-prove ./cmd/notary-prove

FROM scratch AS srs

ADD --checksum=sha256:63cdfe7a45db5ac0bf760fb0245abf5dde16fa9b745636d80449e8c50f2f784d \
    https://github.com/WildanFrananda/kinetix-notary-service/releases/download/srs-ignition-d8/ignition-d8.srs \
    /ignition-d8.srs

FROM debian:bookworm-slim@sha256:99c0b25c46d0214d5a2091092e7908907dcaecd6a962cde0ac700acee222b6f3 AS final

RUN apt-get update \
    && apt-get upgrade -y --no-install-recommends \
    && apt-get install -y --no-install-recommends ca-certificates curl \
    && rm -rf /var/lib/apt/lists/* \
    && groupadd --system --gid 10001 notary \
    && useradd --system --uid 10001 --gid 10001 --no-create-home notary \
    && install -d -m 0755 /usr/local/share/kinetix-notary

COPY --from=build /out/notaryd /out/notary-migrate /out/notary-verify /out/notary-watch /out/notary-sponsor \
    /out/notary-prove /usr/local/bin/

COPY --from=srs --chmod=0444 /ignition-d8.srs /usr/local/share/kinetix-notary/ignition-d8.srs

USER notary

EXPOSE 8010

HEALTHCHECK --interval=15s --timeout=5s --start-period=60s --retries=3 \
    CMD curl -fsS "http://127.0.0.1:${KINETIX_NOTARY_HTTP_PORT}/health/ready" > /dev/null || exit 1

CMD ["notaryd"]

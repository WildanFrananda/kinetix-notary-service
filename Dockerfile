FROM golang:1.26-bookworm@sha256:d9c68c2c51161e12fd77e4c6320687c9cd86e1af1e3ad6e6cd63ff970641453c AS build

ENV GOTOOLCHAIN=local \
    CGO_ENABLED=0

COPY bin/install-toolchain /tmp/install-toolchain
RUN sh /tmp/install-toolchain /usr/local/bin foundry abigen

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY contracts/foundry.toml contracts/soldeer.lock contracts/remappings.txt contracts/
RUN cd contracts && forge soldeer install

COPY . .

RUN sh bin/generate-bindings \
    && sh bin/check-boundaries \
    && go build -trimpath -ldflags="-s -w" -o /out/notary-verify ./cmd/notary-verify

FROM debian:bookworm-slim@sha256:7c7b2c966bc9ee8cedfeef67e0e279108992c77681fa595db4a9d65c06ccc587 AS final

RUN apt-get update \
    && apt-get upgrade -y --no-install-recommends \
    && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/* \
    && useradd --system --uid 10001 --no-create-home notary

COPY --from=build /out/notary-verify /usr/local/bin/notary-verify

USER notary

ENTRYPOINT ["notary-verify"]

FROM golang:1.27-bookworm@sha256:d9b6a1beb9474adc2a1ce842573639b44b2a0dd49b6669f51d3821baf9ec37c9

ENV GOTOOLCHAIN=local \
    CGO_ENABLED=0

RUN apt-get update \
    && apt-get install -y --no-install-recommends python3-venv unzip xz-utils \
    && rm -rf /var/lib/apt/lists/*

COPY bin/install-toolchain /opt/toolchain/bin/install-toolchain
COPY toolchain/slither-requirements.txt /opt/toolchain/toolchain/slither-requirements.txt
RUN sh /opt/toolchain/bin/install-toolchain /usr/local/bin foundry abigen protobuf sqlc lint analysis bun \
    && go clean -cache -modcache

WORKDIR /workspace

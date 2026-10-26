FROM golang:1.27-bookworm@sha256:972b4110ba0f47539a975c31e52c64765e1ff19f2d74469ae4a595c6676e80a5

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

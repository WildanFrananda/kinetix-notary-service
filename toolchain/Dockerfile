FROM golang:1.27-bookworm@sha256:dc0963f32053e4bc451d24d87e3892d42d01021b61f093e98eb95c0352bb84b0

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

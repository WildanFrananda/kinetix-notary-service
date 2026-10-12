FROM golang:1.26-bookworm@sha256:d9c68c2c51161e12fd77e4c6320687c9cd86e1af1e3ad6e6cd63ff970641453c

ENV GOTOOLCHAIN=local \
    CGO_ENABLED=0

RUN apt-get update \
    && apt-get install -y --no-install-recommends python3-venv unzip xz-utils \
    && rm -rf /var/lib/apt/lists/*

COPY bin/install-toolchain /opt/toolchain/bin/install-toolchain
COPY toolchain/slither-requirements.txt /opt/toolchain/toolchain/slither-requirements.txt
RUN sh /opt/toolchain/bin/install-toolchain /usr/local/bin foundry abigen protobuf sqlc lint analysis

WORKDIR /workspace

-- +goose Up
CREATE TABLE evidence_cursors (
    stream              SMALLINT    PRIMARY KEY CHECK (stream BETWEEN 1 AND 3),
    through_transaction BIGINT      NOT NULL CHECK (through_transaction >= 0),
    through_sequence    BIGINT      NOT NULL CHECK (through_sequence >= 0),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO evidence_cursors (stream, through_transaction, through_sequence) VALUES (1, 0, 0), (2, 0, 0), (3, 0, 0);

CREATE TABLE evidence_batches (
    stream              SMALLINT    NOT NULL CHECK (stream BETWEEN 1 AND 3),
    number              BIGINT      NOT NULL CHECK (number >= 0),
    state               TEXT        NOT NULL CHECK (state IN ('built', 'submitted', 'mined', 'finalized', 'abandoned')),
    chain_id            BIGINT      NOT NULL CHECK (chain_id > 0),
    registry_address    BYTEA       NOT NULL CHECK (length(registry_address) = 20),
    scheme              INTEGER     NOT NULL CHECK (scheme BETWEEN 1 AND 65535),
    previous_digest     BYTEA       NOT NULL CHECK (length(previous_digest) = 32),
    root                BYTEA       NOT NULL CHECK (length(root) = 32),
    digest              BYTEA       NOT NULL UNIQUE CHECK (length(digest) = 32),
    leaf_count          BIGINT      NOT NULL CHECK (leaf_count > 0),
    through_transaction BIGINT      NOT NULL CHECK (through_transaction >= 0),
    through_sequence    BIGINT      NOT NULL CHECK (through_sequence >= 0),
    closed_at           BIGINT      NOT NULL CHECK (closed_at > 0),
    abandoned_reason    TEXT,
    built_at            TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (stream, number),
    CONSTRAINT evidence_batches_abandoned_has_reason CHECK ((state = 'abandoned') = (abandoned_reason IS NOT NULL))
);

CREATE UNIQUE INDEX evidence_batches_one_in_flight ON evidence_batches ((true)) WHERE state IN ('built', 'submitted', 'mined');

CREATE TABLE evidence_leaves (
    stream         SMALLINT    NOT NULL CHECK (stream BETWEEN 1 AND 3),
    transaction_id BIGINT      NOT NULL CHECK (transaction_id >= 0),
    sequence       BIGINT      NOT NULL CHECK (sequence >= 0),
    salt           BYTEA       NOT NULL CHECK (length(salt) = 32),
    record_digest  BYTEA       NOT NULL CHECK (length(record_digest) = 32),
    leaf           BYTEA       NOT NULL CHECK (length(leaf) = 32),
    batch_number   BIGINT,
    leaf_index     BIGINT      CHECK (leaf_index >= 0),
    collected_at   TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (stream, transaction_id, sequence),
    FOREIGN KEY (stream, batch_number) REFERENCES evidence_batches (stream, number),
    UNIQUE (stream, batch_number, leaf_index),
    CONSTRAINT evidence_leaves_batched_with_index CHECK ((batch_number IS NULL) = (leaf_index IS NULL))
);

CREATE INDEX evidence_leaves_unbatched ON evidence_leaves (stream, transaction_id, sequence) WHERE batch_number IS NULL;

CREATE TABLE lifecycle_records (
    stream                SMALLINT NOT NULL DEFAULT 1 CHECK (stream = 1),
    transaction_id        BIGINT   NOT NULL,
    sequence              BIGINT   NOT NULL,
    order_number          TEXT     NOT NULL CHECK (octet_length(order_number) BETWEEN 1 AND 65536),
    kind                  BIGINT   NOT NULL CHECK (kind BETWEEN 1 AND 4294967295),
    from_status           BIGINT   NOT NULL CHECK (from_status BETWEEN 0 AND 4294967295),
    to_status             BIGINT   NOT NULL CHECK (to_status BETWEEN 0 AND 4294967295),
    return_number         TEXT     NOT NULL CHECK (octet_length(return_number) <= 65536),
    return_status         BIGINT   NOT NULL CHECK (return_status BETWEEN 0 AND 4294967295),
    actor_kind            BIGINT   NOT NULL CHECK (actor_kind BETWEEN 1 AND 4294967295),
    actor_principal_id    TEXT     NOT NULL CHECK (octet_length(actor_principal_id) <= 65536),
    buyer_principal_id    TEXT     NOT NULL CHECK (octet_length(buyer_principal_id) BETWEEN 1 AND 65536),
    merchant_principal_id TEXT     NOT NULL CHECK (octet_length(merchant_principal_id) BETWEEN 1 AND 65536),
    occurred_seconds      BIGINT   NOT NULL CHECK (occurred_seconds > 0),
    occurred_nanos        INTEGER  NOT NULL CHECK (occurred_nanos BETWEEN 0 AND 999999999),

    PRIMARY KEY (stream, transaction_id, sequence),
    FOREIGN KEY (stream, transaction_id, sequence) REFERENCES evidence_leaves (stream, transaction_id, sequence)
);

CREATE INDEX lifecycle_records_by_order ON lifecycle_records (order_number, transaction_id, sequence);

CREATE TABLE delivery_records (
    stream              SMALLINT NOT NULL DEFAULT 2 CHECK (stream = 2),
    transaction_id      BIGINT   NOT NULL,
    sequence            BIGINT   NOT NULL,
    order_number        TEXT     NOT NULL CHECK (octet_length(order_number) <= 65536),
    driver_principal_id TEXT     NOT NULL CHECK (octet_length(driver_principal_id) BETWEEN 1 AND 65536),
    photo_sha256        BYTEA    NOT NULL CHECK (length(photo_sha256) = 32),
    signature_sha256    BYTEA    CHECK (length(signature_sha256) = 32),
    delivered_seconds   BIGINT   NOT NULL CHECK (delivered_seconds > 0),
    delivered_nanos     INTEGER  NOT NULL CHECK (delivered_nanos BETWEEN 0 AND 999999999),

    PRIMARY KEY (stream, transaction_id, sequence),
    FOREIGN KEY (stream, transaction_id, sequence) REFERENCES evidence_leaves (stream, transaction_id, sequence)
);

CREATE INDEX delivery_records_by_order ON delivery_records (order_number, transaction_id, sequence) WHERE order_number <> '';

CREATE TABLE review_records (
    stream                SMALLINT NOT NULL DEFAULT 3 CHECK (stream = 3),
    transaction_id        BIGINT   NOT NULL,
    sequence              BIGINT   NOT NULL,
    review_id             TEXT     NOT NULL CHECK (octet_length(review_id) BETWEEN 1 AND 65536),
    kind                  BIGINT   NOT NULL CHECK (kind BETWEEN 1 AND 4294967295),
    order_id              TEXT     NOT NULL CHECK (octet_length(order_id) BETWEEN 1 AND 65536),
    subject_id            TEXT     NOT NULL CHECK (octet_length(subject_id) BETWEEN 1 AND 65536),
    author_principal_id   TEXT     NOT NULL CHECK (octet_length(author_principal_id) BETWEEN 1 AND 65536),
    merchant_principal_id TEXT     NOT NULL CHECK (octet_length(merchant_principal_id) <= 65536),
    rating                BIGINT   NOT NULL CHECK (rating BETWEEN 1 AND 4294967295),
    comment_sha256        BYTEA    CHECK (length(comment_sha256) = 32),
    created_seconds       BIGINT   NOT NULL CHECK (created_seconds > 0),
    created_nanos         INTEGER  NOT NULL CHECK (created_nanos BETWEEN 0 AND 999999999),

    PRIMARY KEY (stream, transaction_id, sequence),
    FOREIGN KEY (stream, transaction_id, sequence) REFERENCES evidence_leaves (stream, transaction_id, sequence)
);

CREATE INDEX review_records_by_review ON review_records (review_id, transaction_id, sequence);

CREATE TABLE evidence_attempts (
    tx_hash                  BYTEA       PRIMARY KEY CHECK (length(tx_hash) = 32),
    stream                   SMALLINT    NOT NULL,
    batch_number             BIGINT      NOT NULL,
    nonce                    BIGINT      NOT NULL CHECK (nonce >= 0),
    raw_tx                   BYTEA       NOT NULL CHECK (length(raw_tx) > 0),
    gas_limit                BIGINT      NOT NULL CHECK (gas_limit > 0),
    max_fee_per_gas          BIGINT      NOT NULL CHECK (max_fee_per_gas > 0),
    max_priority_fee_per_gas BIGINT      NOT NULL CHECK (max_priority_fee_per_gas >= 0),
    outcome                  TEXT        NOT NULL CHECK (outcome IN ('pending', 'included', 'superseded', 'reverted')),
    sent_at                  TIMESTAMPTZ NOT NULL DEFAULT now(),
    block_number             BIGINT      CHECK (block_number >= 0),
    gas_used                 BIGINT      CHECK (gas_used >= 0),
    effective_gas_price      BIGINT      CHECK (effective_gas_price >= 0),
    l1_fee_wei               BIGINT      CHECK (l1_fee_wei >= 0),

    FOREIGN KEY (stream, batch_number) REFERENCES evidence_batches (stream, number),
    CONSTRAINT evidence_attempts_settled_has_receipt CHECK (
        (outcome IN ('included', 'reverted')) = (block_number IS NOT NULL AND gas_used IS NOT NULL
            AND effective_gas_price IS NOT NULL AND l1_fee_wei IS NOT NULL)
    )
);

CREATE INDEX evidence_attempts_by_batch ON evidence_attempts (stream, batch_number, nonce, sent_at);

-- +goose Down
DROP TABLE evidence_attempts;
DROP TABLE review_records;
DROP TABLE delivery_records;
DROP TABLE lifecycle_records;
DROP TABLE evidence_leaves;
DROP TABLE evidence_batches;
DROP TABLE evidence_cursors;

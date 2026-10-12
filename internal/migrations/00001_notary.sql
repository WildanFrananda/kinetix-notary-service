-- +goose Up
CREATE TABLE writer_lease (
    id         SMALLINT    PRIMARY KEY CHECK (id = 1),
    holder     TEXT        NOT NULL,
    token      BIGINT      NOT NULL CHECK (token >= 0),
    expires_at TIMESTAMPTZ NOT NULL
);

INSERT INTO writer_lease (id, holder, token, expires_at) VALUES (1, '', 0, '1970-01-01T00:00:00Z');

CREATE TABLE epochs (
    number           BIGINT      PRIMARY KEY CHECK (number >= 0),
    state            TEXT        NOT NULL CHECK (state IN ('built', 'submitted', 'mined', 'finalized', 'abandoned')),
    chain_id         BIGINT      NOT NULL CHECK (chain_id > 0),
    ledger_address   BYTEA       NOT NULL CHECK (length(ledger_address) = 20),
    scheme           INTEGER     NOT NULL CHECK (scheme BETWEEN 1 AND 65535),
    previous_digest  BYTEA       NOT NULL CHECK (length(previous_digest) = 32),
    root             BYTEA       NOT NULL CHECK (length(root) = 32),
    digest           BYTEA       NOT NULL UNIQUE CHECK (length(digest) = 32),
    total            BIGINT      NOT NULL CHECK (total >= 0),
    leaf_count       BIGINT      NOT NULL CHECK (leaf_count >= 0),
    as_of            BIGINT      NOT NULL CHECK (as_of > 0),
    currency         CHAR(3)     NOT NULL,
    abandoned_reason TEXT,
    built_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT epochs_abandoned_has_reason CHECK ((state = 'abandoned') = (abandoned_reason IS NOT NULL))
);

CREATE UNIQUE INDEX epochs_one_in_flight ON epochs ((true)) WHERE state IN ('built', 'submitted', 'mined');

CREATE TABLE accounts (
    epoch_number BIGINT   NOT NULL REFERENCES epochs (number),
    position     BIGINT   NOT NULL CHECK (position >= 0),
    holder       TEXT     NOT NULL CHECK (octet_length(holder) BETWEEN 1 AND 255),
    kind         SMALLINT NOT NULL CHECK (kind BETWEEN 1 AND 7),
    amount       BIGINT   NOT NULL CHECK (amount >= 0),
    salt         BYTEA    NOT NULL CHECK (length(salt) = 32),
    commitment   BYTEA    NOT NULL CHECK (length(commitment) = 32),

    PRIMARY KEY (epoch_number, position),
    UNIQUE (epoch_number, holder, kind),
    UNIQUE (epoch_number, commitment)
);

CREATE INDEX accounts_by_holder ON accounts (holder, epoch_number);

CREATE TABLE attempts (
    tx_hash                  BYTEA       PRIMARY KEY CHECK (length(tx_hash) = 32),
    epoch_number             BIGINT      NOT NULL REFERENCES epochs (number),
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

    CONSTRAINT attempts_settled_has_receipt CHECK (
        (outcome IN ('included', 'reverted')) = (block_number IS NOT NULL AND gas_used IS NOT NULL
            AND effective_gas_price IS NOT NULL AND l1_fee_wei IS NOT NULL)
    )
);

CREATE INDEX attempts_by_epoch ON attempts (epoch_number, nonce, sent_at);

-- +goose Down
DROP TABLE attempts;
DROP TABLE accounts;
DROP TABLE epochs;
DROP TABLE writer_lease;

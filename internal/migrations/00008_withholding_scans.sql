-- +goose Up
CREATE TABLE withholding_scans (
    chain_id   BIGINT      NOT NULL CHECK (chain_id > 0),
    court      BYTEA       NOT NULL CHECK (length(court) = 20),
    next_block BIGINT      NOT NULL CHECK (next_block >= 0),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (chain_id, court)
);

CREATE TABLE withholding_requests (
    chain_id     BIGINT      NOT NULL,
    court        BYTEA       NOT NULL,
    registration BYTEA       NOT NULL CHECK (length(registration) = 32),
    epoch_number BIGINT      NOT NULL CHECK (epoch_number >= 0),
    account      BYTEA       NOT NULL CHECK (length(account) = 32),
    holder       BYTEA       NOT NULL CHECK (length(holder) = 20),
    due_at       BIGINT      NOT NULL CHECK (due_at > 0),
    seen_at      TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (chain_id, court, registration, epoch_number),
    FOREIGN KEY (chain_id, court) REFERENCES withholding_scans (chain_id, court)
);

-- +goose Down
DROP TABLE withholding_requests;
DROP TABLE withholding_scans;

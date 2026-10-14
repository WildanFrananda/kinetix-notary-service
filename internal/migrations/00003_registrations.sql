-- +goose Up
CREATE TABLE registrations (
    digest         BYTEA       PRIMARY KEY CHECK (length(digest) = 32),
    holder_address BYTEA       NOT NULL CHECK (length(holder_address) = 20),
    account        BYTEA       NOT NULL CHECK (length(account) = 32),
    from_epoch     BIGINT      NOT NULL REFERENCES epochs (number),
    principal      TEXT        NOT NULL CHECK (octet_length(principal) BETWEEN 1 AND 255),
    kind           SMALLINT    NOT NULL CHECK (kind BETWEEN 1 AND 7),
    signature      BYTEA       NOT NULL CHECK (length(signature) = 65),
    issued_at      TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT registrations_name_a_kept_account
        FOREIGN KEY (principal, kind, account) REFERENCES account_secrets (holder, kind, account)
);

CREATE INDEX registrations_by_account ON registrations (account, from_epoch);

-- +goose Down
DROP TABLE registrations;

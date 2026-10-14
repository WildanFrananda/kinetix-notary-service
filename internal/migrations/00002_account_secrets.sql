-- +goose Up
CREATE TABLE account_secrets (
    holder     TEXT        NOT NULL CHECK (octet_length(holder) BETWEEN 1 AND 255),
    kind       SMALLINT    NOT NULL CHECK (kind BETWEEN 1 AND 7),
    secret     BYTEA       NOT NULL CHECK (length(secret) = 32),
    account    BYTEA       NOT NULL UNIQUE CHECK (length(account) = 32),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (holder, kind),
    UNIQUE (holder, kind, account)
);

ALTER TABLE accounts ADD COLUMN account BYTEA CHECK (length(account) = 32);

ALTER TABLE accounts ADD CONSTRAINT accounts_account_has_secret
    FOREIGN KEY (holder, kind, account) REFERENCES account_secrets (holder, kind, account);

CREATE INDEX accounts_by_account ON accounts (account, epoch_number) WHERE account IS NOT NULL;

-- +goose Down
DROP INDEX accounts_by_account;
ALTER TABLE accounts DROP CONSTRAINT accounts_account_has_secret;
ALTER TABLE accounts DROP COLUMN account;
DROP TABLE account_secrets;

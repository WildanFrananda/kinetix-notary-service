-- +goose Up
CREATE TABLE statement_progress (
    id         SMALLINT    PRIMARY KEY CHECK (id = 1),
    last_epoch BIGINT      NOT NULL REFERENCES epochs (number),
    last_day   DATE        NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE statement_deliveries (
    holder          TEXT        NOT NULL CHECK (octet_length(holder) BETWEEN 1 AND 255),
    epoch_number    BIGINT      NOT NULL REFERENCES epochs (number),
    day             DATE        NOT NULL,
    fingerprint     BYTEA       NOT NULL CHECK (length(fingerprint) = 32),
    state           TEXT        NOT NULL DEFAULT 'pending'
        CHECK (state IN ('pending', 'sent', 'refused', 'superseded')),
    attempts        INTEGER     NOT NULL DEFAULT 0 CHECK (attempts >= 0),
    next_attempt_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    notification_id TEXT        CHECK (octet_length(notification_id) BETWEEN 1 AND 64),
    refused_reason  TEXT,
    planned_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (holder, epoch_number),
    UNIQUE (holder, day),
    CONSTRAINT statement_deliveries_sent_names_its_notification CHECK ((state = 'sent') = (notification_id IS NOT NULL)),
    CONSTRAINT statement_deliveries_refused_has_reason CHECK ((state = 'refused') = (refused_reason IS NOT NULL))
);

CREATE INDEX statement_deliveries_due ON statement_deliveries (next_attempt_at) WHERE state = 'pending';

CREATE INDEX statement_deliveries_latest ON statement_deliveries (holder, epoch_number DESC) WHERE state <> 'superseded';

-- +goose Down
DROP TABLE statement_deliveries;
DROP TABLE statement_progress;

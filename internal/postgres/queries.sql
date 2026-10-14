-- name: AcquireLease :one
UPDATE writer_lease
SET holder = sqlc.arg(instance)::text,
    token = CASE WHEN holder = sqlc.arg(instance)::text AND expires_at > now() THEN token ELSE token + 1 END,
    expires_at = now() + sqlc.arg(ttl_millis)::bigint * interval '1 millisecond'
WHERE id = 1 AND (holder = sqlc.arg(instance)::text OR expires_at <= now())
RETURNING holder, token, expires_at;

-- name: FenceLease :one
SELECT token
FROM writer_lease
WHERE id = 1 AND holder = sqlc.arg(instance)::text AND token = sqlc.arg(token)::bigint AND expires_at > now()
FOR SHARE;

-- name: InFlightEpoch :one
SELECT * FROM epochs WHERE state IN ('built', 'submitted', 'mined');

-- name: LastFinalizedEpoch :one
SELECT * FROM epochs WHERE state = 'finalized' ORDER BY number DESC LIMIT 1;

-- name: InsertEpoch :exec
INSERT INTO epochs (
    number,
    state,
    chain_id,
    ledger_address,
    scheme,
    previous_digest,
    root,
    digest,
    total,
    leaf_count,
    as_of,
    currency,
    snapshot_id
) VALUES (
    $1, 'built', $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
);

-- name: InsertAccounts :copyfrom
INSERT INTO accounts (epoch_number, position, holder, kind, amount, salt, commitment, account)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: AccountSecretsFor :many
SELECT holder, kind, secret
FROM account_secrets
WHERE (holder, kind) IN (
    SELECT unnest(sqlc.arg(holders)::text[]), unnest(sqlc.arg(kinds)::smallint[])
);

-- name: InsertAccountSecrets :exec
INSERT INTO account_secrets (holder, kind, secret, account)
SELECT
    unnest(sqlc.arg(holders)::text[]),
    unnest(sqlc.arg(kinds)::smallint[]),
    unnest(sqlc.arg(secrets)::bytea[]),
    unnest(sqlc.arg(accounts)::bytea[])
ON CONFLICT (holder, kind) DO NOTHING;

-- name: TransitionEpoch :execrows
UPDATE epochs
SET state = sqlc.arg(to_state)::text, abandoned_reason = sqlc.narg(reason), updated_at = now()
WHERE number = sqlc.arg(number) AND state = sqlc.arg(from_state)::text;

-- name: InsertAttempt :exec
INSERT INTO attempts (
    tx_hash, epoch_number, nonce, raw_tx, gas_limit, max_fee_per_gas, max_priority_fee_per_gas, outcome, sent_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, 'pending', $8
);

-- name: EpochAttempts :many
SELECT * FROM attempts WHERE epoch_number = $1 ORDER BY sent_at, tx_hash;

-- name: SettleAttempt :execrows
UPDATE attempts
SET outcome = sqlc.arg(outcome)::text,
    block_number = sqlc.narg(block_number),
    gas_used = sqlc.narg(gas_used),
    effective_gas_price = sqlc.narg(effective_gas_price),
    l1_fee_wei = sqlc.narg(l1_fee_wei)
WHERE tx_hash = sqlc.arg(tx_hash);

-- name: AbandonedCount :one
SELECT count(*) FROM epochs WHERE state = 'abandoned';

-- name: RecentIncludedReceipts :many
SELECT a.tx_hash, a.block_number, a.gas_used, a.effective_gas_price, a.l1_fee_wei
FROM attempts a
JOIN epochs e ON e.number = a.epoch_number
WHERE e.state = 'finalized' AND a.outcome = 'included'
ORDER BY e.number DESC
LIMIT sqlc.arg(row_limit);

-- name: EpochByNumber :one
SELECT * FROM epochs WHERE number = $1;

-- name: EpochAccounts :many
SELECT a.position, a.holder, a.kind, a.amount, a.salt, a.commitment, a.account, s.secret
FROM accounts a
LEFT JOIN account_secrets s ON s.account = a.account
WHERE a.epoch_number = $1
ORDER BY a.position;

-- name: InsertRegistrations :exec
INSERT INTO registrations (digest, holder_address, account, from_epoch, principal, kind, signature)
SELECT
    unnest(sqlc.arg(digests)::bytea[]),
    unnest(sqlc.arg(holder_addresses)::bytea[]),
    unnest(sqlc.arg(accounts)::bytea[]),
    unnest(sqlc.arg(from_epochs)::bigint[]),
    unnest(sqlc.arg(principals)::text[]),
    unnest(sqlc.arg(kinds)::smallint[]),
    unnest(sqlc.arg(signatures)::bytea[])
ON CONFLICT (digest) DO NOTHING;

-- name: RegistrationByDigest :one
SELECT digest, holder_address, account, from_epoch, kind, signature
FROM registrations
WHERE digest = $1;

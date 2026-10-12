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
    number, state, chain_id, ledger_address, scheme, previous_digest, root, digest, total, leaf_count, as_of, currency
) VALUES (
    $1, 'built', $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
);

-- name: InsertAccounts :copyfrom
INSERT INTO accounts (epoch_number, position, holder, kind, amount, salt, commitment)
VALUES ($1, $2, $3, $4, $5, $6, $7);

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

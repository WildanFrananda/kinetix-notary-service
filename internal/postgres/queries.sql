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

-- name: EvidenceCursor :one
SELECT through_transaction, through_sequence FROM evidence_cursors WHERE stream = $1;

-- name: AdvanceEvidenceCursor :execrows
UPDATE evidence_cursors
SET through_transaction = sqlc.arg(to_transaction)::bigint,
    through_sequence = sqlc.arg(to_sequence)::bigint,
    updated_at = now()
WHERE stream = sqlc.arg(stream)::smallint
  AND through_transaction = sqlc.arg(from_transaction)::bigint
  AND through_sequence = sqlc.arg(from_sequence)::bigint;

-- name: InsertEvidenceLeaves :copyfrom
INSERT INTO evidence_leaves (stream, transaction_id, sequence, salt, record_digest, leaf)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: InsertLifecycleRecords :copyfrom
INSERT INTO lifecycle_records (
    transaction_id,
    sequence,
    order_number,
    kind,
    from_status,
    to_status,
    return_number,
    return_status,
    actor_kind,
    actor_principal_id,
    buyer_principal_id,
    merchant_principal_id,
    occurred_seconds,
    occurred_nanos
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
);

-- name: InsertDeliveryRecords :copyfrom
INSERT INTO delivery_records (
    transaction_id,
    sequence,
    order_number,
    driver_principal_id,
    photo_sha256,
    signature_sha256,
    delivered_seconds,
    delivered_nanos
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
);

-- name: InsertReviewRecords :copyfrom
INSERT INTO review_records (
    transaction_id,
    sequence,
    review_id,
    kind,
    order_id,
    subject_id,
    author_principal_id,
    merchant_principal_id,
    rating,
    comment_sha256,
    created_seconds,
    created_nanos
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
);

-- name: UnbatchedLifecycleRecords :many
SELECT l.transaction_id, l.sequence, l.salt, l.record_digest, l.leaf, l.batch_number, l.leaf_index,
    r.order_number, r.kind, r.from_status, r.to_status, r.return_number, r.return_status, r.actor_kind,
    r.actor_principal_id, r.buyer_principal_id, r.merchant_principal_id, r.occurred_seconds, r.occurred_nanos
FROM evidence_leaves l
JOIN lifecycle_records r ON r.stream = l.stream AND r.transaction_id = l.transaction_id AND r.sequence = l.sequence
WHERE l.stream = 1 AND l.batch_number IS NULL
ORDER BY l.transaction_id, l.sequence
LIMIT sqlc.arg(row_limit);

-- name: LifecycleRecordsOfOrder :many
SELECT l.transaction_id, l.sequence, l.salt, l.record_digest, l.leaf, l.batch_number, l.leaf_index,
    r.order_number, r.kind, r.from_status, r.to_status, r.return_number, r.return_status, r.actor_kind,
    r.actor_principal_id, r.buyer_principal_id, r.merchant_principal_id, r.occurred_seconds, r.occurred_nanos
FROM evidence_leaves l
JOIN lifecycle_records r ON r.stream = l.stream AND r.transaction_id = l.transaction_id AND r.sequence = l.sequence
WHERE l.stream = 1 AND r.order_number = sqlc.arg(order_number)::text
ORDER BY l.transaction_id, l.sequence;

-- name: UnbatchedDeliveryRecords :many
SELECT l.transaction_id, l.sequence, l.salt, l.record_digest, l.leaf, l.batch_number, l.leaf_index,
    r.order_number, r.driver_principal_id, r.photo_sha256, r.signature_sha256, r.delivered_seconds, r.delivered_nanos
FROM evidence_leaves l
JOIN delivery_records r ON r.stream = l.stream AND r.transaction_id = l.transaction_id AND r.sequence = l.sequence
WHERE l.stream = 2 AND l.batch_number IS NULL
ORDER BY l.transaction_id, l.sequence
LIMIT sqlc.arg(row_limit);

-- name: DeliveryRecordsOfOrder :many
SELECT l.transaction_id, l.sequence, l.salt, l.record_digest, l.leaf, l.batch_number, l.leaf_index,
    r.order_number, r.driver_principal_id, r.photo_sha256, r.signature_sha256, r.delivered_seconds, r.delivered_nanos
FROM evidence_leaves l
JOIN delivery_records r ON r.stream = l.stream AND r.transaction_id = l.transaction_id AND r.sequence = l.sequence
WHERE l.stream = 2 AND r.order_number <> '' AND r.order_number = sqlc.arg(order_number)::text
ORDER BY l.transaction_id, l.sequence;

-- name: UnbatchedReviewRecords :many
SELECT l.transaction_id, l.sequence, l.salt, l.record_digest, l.leaf, l.batch_number, l.leaf_index,
    r.review_id, r.kind, r.order_id, r.subject_id, r.author_principal_id, r.merchant_principal_id, r.rating,
    r.comment_sha256, r.created_seconds, r.created_nanos
FROM evidence_leaves l
JOIN review_records r ON r.stream = l.stream AND r.transaction_id = l.transaction_id AND r.sequence = l.sequence
WHERE l.stream = 3 AND l.batch_number IS NULL
ORDER BY l.transaction_id, l.sequence
LIMIT sqlc.arg(row_limit);

-- name: ReviewRecordsOfReview :many
SELECT l.transaction_id, l.sequence, l.salt, l.record_digest, l.leaf, l.batch_number, l.leaf_index,
    r.review_id, r.kind, r.order_id, r.subject_id, r.author_principal_id, r.merchant_principal_id, r.rating,
    r.comment_sha256, r.created_seconds, r.created_nanos
FROM evidence_leaves l
JOIN review_records r ON r.stream = l.stream AND r.transaction_id = l.transaction_id AND r.sequence = l.sequence
WHERE l.stream = 3 AND r.review_id = sqlc.arg(review_id)::text
ORDER BY l.transaction_id, l.sequence;

-- name: UnbatchedEvidenceCount :one
SELECT count(*) FROM evidence_leaves WHERE stream = $1 AND batch_number IS NULL;

-- name: InsertEvidenceBatch :exec
INSERT INTO evidence_batches (
    stream,
    number,
    state,
    chain_id,
    registry_address,
    scheme,
    previous_digest,
    root,
    digest,
    leaf_count,
    through_transaction,
    through_sequence,
    closed_at
) VALUES (
    $1, $2, 'built', $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
);

-- name: AssignEvidenceLeaves :execrows
UPDATE evidence_leaves l
SET batch_number = sqlc.arg(batch_number)::bigint, leaf_index = a.leaf_index
FROM (
    SELECT
        unnest(sqlc.arg(transaction_ids)::bigint[]) AS transaction_id,
        unnest(sqlc.arg(sequences)::bigint[]) AS sequence,
        unnest(sqlc.arg(leaf_indexes)::bigint[]) AS leaf_index,
        unnest(sqlc.arg(leaves)::bytea[]) AS leaf
) a
WHERE l.stream = sqlc.arg(stream)::smallint
  AND l.transaction_id = a.transaction_id
  AND l.sequence = a.sequence
  AND l.leaf = a.leaf
  AND l.batch_number IS NULL;

-- name: InFlightEvidenceBatch :one
SELECT * FROM evidence_batches WHERE state IN ('built', 'submitted', 'mined');

-- name: LastEvidenceBatch :one
SELECT * FROM evidence_batches WHERE stream = $1 ORDER BY number DESC LIMIT 1;

-- name: LastFinalizedEvidenceBatch :one
SELECT * FROM evidence_batches WHERE stream = $1 AND state = 'finalized' ORDER BY number DESC LIMIT 1;

-- name: EvidenceBatchByNumber :one
SELECT * FROM evidence_batches WHERE stream = $1 AND number = $2;

-- name: EvidenceBatchLeaves :many
SELECT leaf FROM evidence_leaves WHERE stream = $1 AND batch_number = $2 ORDER BY leaf_index;

-- name: TransitionEvidenceBatch :execrows
UPDATE evidence_batches
SET state = sqlc.arg(to_state)::text, abandoned_reason = sqlc.narg(reason), updated_at = now()
WHERE stream = sqlc.arg(stream)::smallint AND number = sqlc.arg(number) AND state = sqlc.arg(from_state)::text;

-- name: InsertEvidenceAttempt :exec
INSERT INTO evidence_attempts (
    tx_hash,
    stream,
    batch_number,
    nonce,
    raw_tx,
    gas_limit,
    max_fee_per_gas,
    max_priority_fee_per_gas,
    outcome,
    sent_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, 'pending', $9
);

-- name: EvidenceBatchAttempts :many
SELECT * FROM evidence_attempts WHERE stream = $1 AND batch_number = $2 ORDER BY sent_at, tx_hash;

-- name: SettleEvidenceAttempt :execrows
UPDATE evidence_attempts
SET outcome = sqlc.arg(outcome)::text,
    block_number = sqlc.narg(block_number),
    gas_used = sqlc.narg(gas_used),
    effective_gas_price = sqlc.narg(effective_gas_price),
    l1_fee_wei = sqlc.narg(l1_fee_wei)
WHERE tx_hash = sqlc.arg(tx_hash);

-- name: AbandonedEvidenceBatches :one
SELECT count(*) FROM evidence_batches WHERE state = 'abandoned';

-- name: RecentEvidenceReceipts :many
SELECT a.tx_hash, a.block_number, a.gas_used, a.effective_gas_price, a.l1_fee_wei
FROM evidence_attempts a
JOIN evidence_batches b ON b.stream = a.stream AND b.number = a.batch_number
WHERE b.state = 'finalized' AND a.outcome = 'included'
ORDER BY b.updated_at DESC
LIMIT sqlc.arg(row_limit);

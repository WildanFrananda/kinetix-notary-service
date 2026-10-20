-- +goose Up
ALTER TABLE epochs
    ADD COLUMN tree_root BYTEA CHECK (length(tree_root) = 32),
    ADD COLUMN total_commitment BYTEA CHECK (length(total_commitment) = 32),
    ADD COLUMN blinding BYTEA CHECK (length(blinding) = 31),
    ADD COLUMN depth SMALLINT CHECK (depth BETWEEN 0 AND 16),
    ADD COLUMN proof BYTEA CHECK (length(proof) BETWEEN 1 AND 65536),
    ADD CONSTRAINT epochs_scheme_two_is_whole CHECK (
        (scheme = 2) = (tree_root IS NOT NULL AND total_commitment IS NOT NULL AND blinding IS NOT NULL
            AND depth IS NOT NULL AND proof IS NOT NULL)
    ),
    ADD CONSTRAINT epochs_scheme_two_hides_its_total CHECK (scheme <> 2 OR total = 0);

ALTER TABLE attempts ADD COLUMN kind TEXT NOT NULL DEFAULT 'anchor' CHECK (kind IN ('anchor', 'proof'));

-- +goose Down
ALTER TABLE attempts DROP COLUMN kind;

ALTER TABLE epochs
    DROP CONSTRAINT epochs_scheme_two_hides_its_total,
    DROP CONSTRAINT epochs_scheme_two_is_whole,
    DROP COLUMN proof,
    DROP COLUMN depth,
    DROP COLUMN blinding,
    DROP COLUMN total_commitment,
    DROP COLUMN tree_root;

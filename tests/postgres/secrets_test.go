package postgres_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

func keysOf(accounts []notary.Account) []notary.AccountKey {
	keys := make([]notary.AccountKey, len(accounts))
	for i, account := range accounts {
		keys[i] = notary.AccountKey{Holder: account.Liability.Holder, Kind: account.Liability.Kind}
	}

	return keys
}

func abandon(t *testing.T, d *db, lease notary.Lease, epoch uint64) {
	t.Helper()

	require.NoError(
		t,
		d.store.Transition(context.Background(), lease, epoch, notary.StateBuilt, notary.StateAbandoned, "test"),
	)
}

func TestAnAccountKeepsOneSecretAcrossEpochs(t *testing.T) {
	d := open(t)
	ctx := context.Background()
	lease := d.lease(t, "writer")
	first := plan(t, 0, notary.Digest{}, 6)
	_, err := d.store.SaveBuilt(ctx, lease, ledgerID, first)
	require.NoError(t, err)
	abandon(t, d, lease, 0)

	known, err := d.store.AccountSecrets(ctx, keysOf(first.Accounts))
	require.NoError(t, err)
	require.Len(t, known, 6)

	snapshot := notarytest.Snapshot(t, 1760000001, notarytest.Liabilities(t, 6))
	second, err := notary.PlanEpoch(
		keccak.NewScheme(),
		&notarytest.Salts{Seed: 1},
		&notarytest.Secrets{Seed: 99},
		known,
		1,
		notary.Digest{},
		snapshot,
	)
	require.NoError(t, err)
	_, err = d.store.SaveBuilt(ctx, lease, ledgerID, second)
	require.NoError(t, err)

	before, err := d.store.Accounts(ctx, 0)
	require.NoError(t, err)
	after, err := d.store.Accounts(ctx, 1)
	require.NoError(t, err)

	byKey := map[notary.AccountKey]notary.Account{}
	for _, account := range before {
		byKey[notary.AccountKey{Holder: account.Liability.Holder, Kind: account.Liability.Kind}] = account
	}

	for _, account := range after {
		earlier := byKey[notary.AccountKey{Holder: account.Liability.Holder, Kind: account.Liability.Kind}]
		require.False(t, account.Secret.IsZero())
		require.Equal(t, earlier.Secret, account.Secret)
		require.Equal(t, earlier.ID, account.ID)
		require.NotEqual(t, earlier.Commitment, account.Commitment)
	}

	require.Equal(t, int64(6), d.count(t, "SELECT count(*) FROM account_secrets"))
}

func TestAnotherSecretForAKeptAccountIsRefusedAndNothingIsWritten(t *testing.T) {
	d := open(t)
	ctx := context.Background()
	lease := d.lease(t, "writer")
	_, err := d.store.SaveBuilt(ctx, lease, ledgerID, plan(t, 0, notary.Digest{}, 3))
	require.NoError(t, err)
	abandon(t, d, lease, 0)

	snapshot := notarytest.Snapshot(t, 1760000001, notarytest.Liabilities(t, 3))
	other, err := notary.PlanEpoch(
		keccak.NewScheme(),
		&notarytest.Salts{Seed: 1},
		&notarytest.Secrets{Seed: 99},
		nil,
		1,
		notary.Digest{},
		snapshot,
	)
	require.NoError(t, err)

	_, err = d.store.SaveBuilt(ctx, lease, ledgerID, other)

	require.Equal(t, notary.KindDisagreement, notary.KindOf(err))
	require.Equal(t, int64(1), d.count(t, "SELECT count(*) FROM epochs"))
	require.Equal(t, int64(0), d.count(t, "SELECT count(*) FROM accounts WHERE epoch_number = 1"))
}

func TestAnEpochBuiltBeforeSecretsReadsBackWithoutOne(t *testing.T) {
	d := open(t)
	ctx := context.Background()
	lease := d.lease(t, "writer")
	scheme := keccak.NewScheme()
	p := plan(t, 0, notary.Digest{}, 4)

	entries := make([]notary.Entry, len(p.Accounts))
	for i := range p.Accounts {
		account := &p.Accounts[i]
		account.Secret = notary.AccountSecret{}
		account.ID = notary.Digest{}
		account.Commitment = scheme.Commitment(account.Liability.Holder, account.Liability.Kind, account.Salt)
		entries[i] = notary.Entry{Commitment: account.Commitment, Amount: account.Liability.Amount}
	}

	tree, err := notary.BuildTree(scheme, 0, entries)
	require.NoError(t, err)
	p.Header, err = tree.Header(notary.Digest{}, p.Header.AsOf)
	require.NoError(t, err)

	p.Digest = scheme.EpochDigest(p.Header)

	_, err = d.store.SaveBuilt(ctx, lease, ledgerID, p)
	require.NoError(t, err)

	accounts, err := d.store.Accounts(ctx, 0)
	require.NoError(t, err)
	require.Len(t, accounts, 4)

	for _, account := range accounts {
		require.True(t, account.Secret.IsZero())
		require.True(t, account.ID == notary.Digest{})
	}

	require.Equal(t, int64(0), d.count(t, "SELECT count(*) FROM account_secrets"))
	require.Equal(t, int64(4), d.count(t, "SELECT count(*) FROM accounts WHERE account IS NULL"))
}

func TestTheDatabaseRefusesAnAccountRowWithoutItsSecret(t *testing.T) {
	d := open(t)
	ctx := context.Background()
	lease := d.lease(t, "writer")
	_, err := d.store.SaveBuilt(ctx, lease, ledgerID, plan(t, 0, notary.Digest{}, 1))
	require.NoError(t, err)

	_, err = d.pool.Exec(
		ctx,
		"INSERT INTO accounts (epoch_number, position, holder, kind, amount, salt, commitment, account) "+
			"VALUES (0, 7, 'stray', 1, 0, $1, $2, $3)",
		make([]byte, 32),
		make([]byte, 32),
		make([]byte, 32),
	)

	require.ErrorContains(t, err, "accounts_account_has_secret")
}

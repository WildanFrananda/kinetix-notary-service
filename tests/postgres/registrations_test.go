package postgres_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

func TestARegistrationIsKeptOnceAndReadBackByItsDigest(t *testing.T) {
	d := open(t)
	ctx := context.Background()
	lease := d.lease(t, "writer")
	p := plan(t, 0, notary.Digest{}, 3)
	_, err := d.store.SaveBuilt(ctx, lease, ledgerID, p)
	require.NoError(t, err)

	account := p.Accounts[1]
	issued := notary.IssuedRegistration{
		Registration: notary.Registration{Holder: notary.Address{0xa1}, Account: account.ID, FromEpoch: 0},
		Kind:         account.Liability.Kind,
		Digest:       notary.Digest{0xdd},
		Signature:    notary.Signature{0x5e, 0x01},
	}

	require.NoError(t, d.store.SaveRegistrations(ctx, account.Liability.Holder, []notary.IssuedRegistration{issued}))
	require.NoError(t, d.store.SaveRegistrations(ctx, account.Liability.Holder, []notary.IssuedRegistration{issued}))
	require.Equal(t, int64(1), d.count(t, "SELECT count(*) FROM registrations"))

	read, found, err := d.store.RegistrationByDigest(ctx, issued.Digest)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, issued, read)

	_, found, err = d.store.RegistrationByDigest(ctx, notary.Digest{0xee})
	require.NoError(t, err)
	require.False(t, found)
}

func TestARegistrationMustNameAnAccountThePrincipalHolds(t *testing.T) {
	d := open(t)
	ctx := context.Background()
	lease := d.lease(t, "writer")
	p := plan(t, 0, notary.Digest{}, 3)
	_, err := d.store.SaveBuilt(ctx, lease, ledgerID, p)
	require.NoError(t, err)

	account := p.Accounts[0]
	issued := notary.IssuedRegistration{
		Registration: notary.Registration{Holder: notary.Address{0xa1}, Account: account.ID, FromEpoch: 0},
		Kind:         account.Liability.Kind,
		Digest:       notary.Digest{0xdd},
		Signature:    notary.Signature{0x5e},
	}

	err = d.store.SaveRegistrations(ctx, notarytest.Holder(t, "someone-else"), []notary.IssuedRegistration{issued})

	require.Equal(t, notary.KindStoreUnavailable, notary.KindOf(err))
	require.Equal(t, int64(0), d.count(t, "SELECT count(*) FROM registrations"))
}

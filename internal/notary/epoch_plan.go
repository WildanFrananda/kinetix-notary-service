package notary

import (
	"bytes"
	"sort"
)

type EpochPlan struct {
	Header   EpochHeader
	Digest   Digest
	Currency string
	Accounts []Account
}

func PlanEpoch(
	scheme Scheme,
	salts SaltSource,
	secrets SecretSource,
	known map[AccountKey]AccountSecret,
	number uint64,
	previous Digest,
	snapshot Snapshot,
) (EpochPlan, error) {
	const op = "notary.PlanEpoch"

	liabilities := snapshot.Liabilities()

	accounts := make([]Account, len(liabilities))
	for i, liability := range liabilities {
		salt, err := salts.Salt()
		if err != nil {
			return EpochPlan{}, Errf(KindUnknown, op, err, "drawing a salt")
		}

		if salt.IsZero() {
			return EpochPlan{}, Errf(KindInvalidArgument, op, nil, "the salt source returned an all-zero salt")
		}

		secret, err := secretFor(secrets, known, AccountKey{Holder: liability.Holder, Kind: liability.Kind})
		if err != nil {
			return EpochPlan{}, err
		}

		id := scheme.AccountID(liability.Holder, liability.Kind, secret)
		accounts[i] = Account{
			Liability:  liability,
			Salt:       salt,
			Secret:     secret,
			ID:         id,
			Commitment: scheme.AccountCommitment(id, salt),
		}
	}

	sort.Slice(accounts, func(a, b int) bool {
		return bytes.Compare(accounts[a].Commitment[:], accounts[b].Commitment[:]) < 0
	})

	entries := make([]Entry, len(accounts))
	for i := range accounts {
		if i > 0 && accounts[i].Commitment == accounts[i-1].Commitment {
			return EpochPlan{}, Errf(
				KindInvalidArgument,
				op,
				nil,
				"two accounts of epoch %d share a commitment",
				number,
			)
		}

		accounts[i].Position = uint64(i)
		entries[i] = Entry{Commitment: accounts[i].Commitment, Amount: accounts[i].Liability.Amount}
	}

	tree, err := BuildTree(scheme, number, entries)
	if err != nil {
		return EpochPlan{}, err
	}

	header, err := tree.Header(previous, snapshot.AsOf())
	if err != nil {
		return EpochPlan{}, err
	}

	return EpochPlan{
		Header:   header,
		Digest:   scheme.EpochDigest(header),
		Currency: snapshot.Currency(),
		Accounts: accounts,
	}, nil
}

func secretFor(secrets SecretSource, known map[AccountKey]AccountSecret, key AccountKey) (AccountSecret, error) {
	const op = "notary.PlanEpoch"

	if secret, ok := known[key]; ok {
		if secret.IsZero() {
			return AccountSecret{}, Errf(
				KindInvalidArgument,
				op,
				nil,
				"%s of %s has an all-zero secret",
				key.Kind,
				key.Holder,
			)
		}

		return secret, nil
	}

	secret, err := secrets.AccountSecret()
	if err != nil {
		return AccountSecret{}, Errf(KindUnknown, op, err, "drawing an account secret")
	}

	if secret.IsZero() {
		return AccountSecret{}, Errf(KindInvalidArgument, op, nil, "the secret source returned an all-zero secret")
	}

	return secret, nil
}

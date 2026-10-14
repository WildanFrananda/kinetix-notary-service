package notary

import (
	"bytes"
	"sort"
)

type EpochPlan struct {
	Header     EpochHeader
	Digest     Digest
	Currency   string
	SnapshotID SnapshotID
	Accounts   []Account
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

		accounts[i] = OpenAccount(scheme, liability, salt, secret)
	}

	plan, err := AssembleEpoch(scheme, number, previous, snapshot.AsOf(), snapshot.Currency(), accounts)
	if err != nil {
		return EpochPlan{}, err
	}

	plan.SnapshotID = snapshot.Record().ID

	return plan, nil
}

func OpenAccount(scheme Scheme, liability Liability, salt Salt, secret AccountSecret) Account {
	id := scheme.AccountID(liability.Holder, liability.Kind, secret)

	return Account{
		Liability:  liability,
		Salt:       salt,
		Secret:     secret,
		ID:         id,
		Commitment: scheme.AccountCommitment(id, salt),
	}
}

func AssembleEpoch(
	scheme Scheme,
	number uint64,
	previous Digest,
	asOf UnixTime,
	currency string,
	accounts []Account,
) (EpochPlan, error) {
	const op = "notary.AssembleEpoch"

	ordered := make([]Account, len(accounts))
	copy(ordered, accounts)
	sort.Slice(ordered, func(a, b int) bool {
		return bytes.Compare(ordered[a].Commitment[:], ordered[b].Commitment[:]) < 0
	})

	entries := make([]Entry, len(ordered))
	for i := range ordered {
		if i > 0 && ordered[i].Commitment == ordered[i-1].Commitment {
			return EpochPlan{}, Errf(
				KindInvalidArgument,
				op,
				nil,
				"two accounts of epoch %d share a commitment",
				number,
			)
		}

		ordered[i].Position = uint64(i)
		entries[i] = Entry{Commitment: ordered[i].Commitment, Amount: ordered[i].Liability.Amount}
	}

	tree, err := BuildTree(scheme, number, entries)
	if err != nil {
		return EpochPlan{}, err
	}

	header, err := tree.Header(previous, asOf)
	if err != nil {
		return EpochPlan{}, err
	}

	return EpochPlan{
		Header:   header,
		Digest:   scheme.EpochDigest(header),
		Currency: currency,
		Accounts: ordered,
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

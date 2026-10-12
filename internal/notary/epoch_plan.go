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

func PlanEpoch(scheme Scheme, salts SaltSource, number uint64, previous Digest, snapshot Snapshot) (EpochPlan, error) {
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

		accounts[i] = Account{
			Liability:  liability,
			Salt:       salt,
			Commitment: scheme.Commitment(liability.Holder, liability.Kind, salt),
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

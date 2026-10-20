package notary

func PlanConfidentialEpoch(
	accounts Scheme,
	scheme ConfidentialScheme,
	depth int,
	salts SaltSource,
	secrets SecretSource,
	blindings BlindingSource,
	known map[AccountKey]AccountSecret,
	number uint64,
	previous Digest,
	snapshot Snapshot,
) (EpochPlan, error) {
	const op = "notary.PlanConfidentialEpoch"

	opened, err := openAccounts(accounts, salts, secrets, known, snapshot.Liabilities())
	if err != nil {
		return EpochPlan{}, err
	}

	blinding, err := blindings.Blinding()
	if err != nil {
		return EpochPlan{}, Errf(KindUnknown, op, err, "drawing a blinding")
	}

	plan, err := AssembleConfidentialEpoch(
		scheme,
		depth,
		blinding,
		number,
		previous,
		snapshot.AsOf(),
		snapshot.Currency(),
		opened,
	)
	if err != nil {
		return EpochPlan{}, err
	}

	plan.SnapshotID = snapshot.Record().ID

	return plan, nil
}

func AssembleConfidentialEpoch(
	scheme ConfidentialScheme,
	depth int,
	blinding Blinding,
	number uint64,
	previous Digest,
	asOf UnixTime,
	currency string,
	accounts []Account,
) (EpochPlan, error) {
	ordered, entries, err := orderAccounts(number, accounts)
	if err != nil {
		return EpochPlan{}, err
	}

	tree, err := BuildConfidentialTree(scheme, number, depth, entries)
	if err != nil {
		return EpochPlan{}, err
	}

	header, err := tree.Header(scheme, previous, asOf, blinding)
	if err != nil {
		return EpochPlan{}, err
	}

	return EpochPlan{
		Header:   header.Header,
		Digest:   scheme.EpochDigest(header.Header),
		Currency: currency,
		Accounts: ordered,
		Confidential: &ConfidentialPlan{
			Scheme:          scheme.Version(),
			TreeRoot:        header.TreeRoot,
			TotalCommitment: header.TotalCommitment,
			Blinding:        blinding,
			Depth:           depth,
		},
	}, nil
}

package notary

func distinctParties(candidates ...string) []string {
	parties := make([]string, 0, len(candidates))

	seen := make(map[string]struct{}, len(candidates))
	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}

		if _, ok := seen[candidate]; ok {
			continue
		}

		seen[candidate] = struct{}{}

		parties = append(parties, candidate)
	}

	return parties
}

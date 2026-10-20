package plonk

import (
	"math/big"

	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/permutation/poseidon2"
	"github.com/consensys/gnark/std/rangecheck"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type LiabilityCircuit struct {
	Epoch           frontend.Variable `gnark:",public"`
	TreeRoot        frontend.Variable `gnark:",public"`
	TotalCommitment frontend.Variable `gnark:",public"`
	LeafCount       frontend.Variable `gnark:",public"`
	High            []frontend.Variable
	Low             []frontend.Variable
	Amounts         []frontend.Variable
	Blinding        frontend.Variable
}

func NewLiabilityCircuit(depth int) *LiabilityCircuit {
	capacity := 1 << depth

	return &LiabilityCircuit{
		High:    make([]frontend.Variable, capacity),
		Low:     make([]frontend.Variable, capacity),
		Amounts: make([]frontend.Variable, capacity),
	}
}

func (c *LiabilityCircuit) Define(api frontend.API) error {
	permutation, err := poseidon2.NewPoseidon2FromParameters(api, 2, 6, 50)
	if err != nil {
		return notary.Errf(
			notary.KindInvalidArgument,
			"plonk.LiabilityCircuit.Define",
			err,
			"the Poseidon2 permutation",
		)
	}

	checker := rangecheck.New(api)
	capacity := len(c.Amounts)

	checker.Check(c.Epoch, 64)
	api.AssertIsLessOrEqual(c.LeafCount, capacity)

	empty := permutation.Compress(0, tagged(emptyTag))
	head := api.Add(tagged(leafTag), api.Mul(c.Epoch, shifted(64)))

	level := make([]frontend.Variable, capacity)
	total := frontend.Variable(0)
	inside := frontend.Variable(1)

	for i := range capacity {
		inside = api.Mul(inside, api.Sub(1, api.IsZero(api.Sub(c.LeafCount, i))))

		checker.Check(c.Amounts[i], 64)
		checker.Check(c.High[i], 128)
		checker.Check(c.Low[i], 128)
		api.AssertIsEqual(api.Mul(api.Sub(1, inside), c.Amounts[i]), 0)

		state := permutation.Compress(0, api.Add(head, c.Amounts[i]))
		state = permutation.Compress(state, c.High[i])
		state = permutation.Compress(state, c.Low[i])

		level[i] = api.Select(inside, state, empty)
		total = api.Add(total, c.Amounts[i])
	}

	checker.Check(total, 64)

	for len(level) > 1 {
		above := make([]frontend.Variable, len(level)/2)
		for i := range above {
			above[i] = permutation.Compress(level[2*i], level[2*i+1])
		}

		level = above
	}

	api.AssertIsEqual(level[0], c.TreeRoot)

	commitment := permutation.Compress(permutation.Compress(0, api.Add(tagged(totalTag), total)), c.Blinding)
	api.AssertIsEqual(commitment, c.TotalCommitment)

	return nil
}

const (
	leafTag  = 1
	emptyTag = 2
	totalTag = 5
)

func tagged(tag int64) *big.Int {
	return new(big.Int).Lsh(big.NewInt(tag), 128)
}

func shifted(bits uint) *big.Int {
	return new(big.Int).Lsh(big.NewInt(1), bits)
}

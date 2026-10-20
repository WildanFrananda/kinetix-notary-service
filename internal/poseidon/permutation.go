package poseidon

import (
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/poseidon2"
)

type permutation struct {
	fullRounds    int
	partialRounds int
	roundKeys     [][]fr.Element
}

func defaultPermutation() permutation {
	parameters := poseidon2.GetDefaultParameters()

	return permutation{
		fullRounds:    parameters.NbFullRounds,
		partialRounds: parameters.NbPartialRounds,
		roundKeys:     parameters.RoundKeys,
	}
}

func (p permutation) compress(left, right fr.Element) fr.Element {
	state := [2]fr.Element{left, right}
	p.permute(&state)

	var out fr.Element
	out.Add(&state[1], &right)

	return out
}

func (p permutation) permute(state *[2]fr.Element) {
	external(state)

	half := p.fullRounds / 2
	for round := range half {
		p.addRoundKeys(round, state)
		sbox(&state[0])
		sbox(&state[1])
		external(state)
	}

	for round := half; round < half+p.partialRounds; round++ {
		p.addRoundKeys(round, state)
		sbox(&state[0])
		internal(state)
	}

	for round := half + p.partialRounds; round < p.fullRounds+p.partialRounds; round++ {
		p.addRoundKeys(round, state)
		sbox(&state[0])
		sbox(&state[1])
		external(state)
	}
}

func (p permutation) addRoundKeys(round int, state *[2]fr.Element) {
	for i := range p.roundKeys[round] {
		state[i].Add(&state[i], &p.roundKeys[round][i])
	}
}

func sbox(x *fr.Element) {
	var squared, fourth fr.Element
	squared.Square(x)
	fourth.Square(&squared)
	x.Mul(x, &fourth)
}

func external(state *[2]fr.Element) {
	var sum fr.Element
	sum.Add(&state[0], &state[1])
	state[0].Add(&state[0], &sum)
	state[1].Add(&state[1], &sum)
}

func internal(state *[2]fr.Element) {
	var sum fr.Element
	sum.Add(&state[0], &state[1])
	state[0].Add(&state[0], &sum)
	state[1].Double(&state[1])
	state[1].Add(&state[1], &sum)
}

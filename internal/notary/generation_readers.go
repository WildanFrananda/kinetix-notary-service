package notary

type GenerationReaders interface {
	Ledger(generation Generation) Ledger
	Court(generation Generation) Court
}

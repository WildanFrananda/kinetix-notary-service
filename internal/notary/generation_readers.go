package notary

type GenerationReaders interface {
	Ledger(generation Generation) Ledger
	Confidential(generation Generation) ConfidentialLedger
	Court(generation Generation) Court
}

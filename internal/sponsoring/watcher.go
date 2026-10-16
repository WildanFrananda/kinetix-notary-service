package sponsoring

import (
	"context"
	"log/slog"
	"math"
	"time"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Watcher struct {
	paymaster notary.Paymaster
	court     notary.Court
	expected  notary.Generation
	observer  Observer
	log       *slog.Logger
	interval  time.Duration
}

func NewWatcher(
	paymaster notary.Paymaster,
	court notary.Court,
	expected notary.Generation,
	observer Observer,
	log *slog.Logger,
	interval time.Duration,
) *Watcher {
	return &Watcher{
		paymaster: paymaster,
		court:     court,
		expected:  expected,
		observer:  observer,
		log:       log,
		interval:  interval,
	}
}

func (w *Watcher) Check(ctx context.Context) error {
	const op = "sponsoring.Watcher.Check"

	identity, err := w.paymaster.Identity(ctx)
	if err != nil {
		return err
	}

	if identity.ChainID != w.expected.ChainID ||
		identity.Court != w.expected.Court ||
		identity.Withholding != w.expected.Withholding {
		return notary.Errf(
			notary.KindNotSponsored,
			op,
			nil,
			"the paymaster %s sponsors the courts %s and %s on chain %d; this notary's are %s and %s on chain %d",
			identity.Address.Hex(),
			identity.Court.Hex(),
			identity.Withholding.Hex(),
			identity.ChainID,
			w.expected.Court.Hex(),
			w.expected.Withholding.Hex(),
			w.expected.ChainID,
		)
	}

	signer, err := w.court.SignerOf(ctx, math.MaxUint64)
	if err != nil {
		return err
	}

	if identity.Signer != signer {
		return notary.Errf(
			notary.KindNotSponsored,
			op,
			nil,
			"the paymaster %s accepts statements signed by %s, the court's signer is now %s: deploy a paymaster for it",
			identity.Address.Hex(),
			identity.Signer.Hex(),
			signer.Hex(),
		)
	}

	return nil
}

func (w *Watcher) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		w.Look(ctx)

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (w *Watcher) Look(ctx context.Context) {
	funds, err := w.paymaster.Funds(ctx)
	if err != nil {
		w.log.WarnContext(
			ctx,
			"the paymaster's deposit and stake could not be read",
			"kind",
			notary.KindOf(err).String(),
			"err",
			err.Error(),
		)
	} else {
		w.observer.PaymasterDeposit(funds.Deposit)
		w.observer.PaymasterStake(funds.LockedStake())
	}

	err = w.Check(ctx)
	w.observer.PaymasterCurrent(err == nil)

	if err != nil {
		w.log.ErrorContext(
			ctx,
			"holders' calls are not being sponsored",
			"kind",
			notary.KindOf(err).String(),
			"err",
			err.Error(),
		)
	}
}

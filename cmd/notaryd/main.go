package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/dig"

	"github.com/WildanFrananda/kinetix-notary-service/internal/anchoring"
	"github.com/WildanFrananda/kinetix-notary-service/internal/answering"
	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/epochwatch"
	"github.com/WildanFrananda/kinetix-notary-service/internal/evidencing"
	"github.com/WildanFrananda/kinetix-notary-service/internal/feeds"
	"github.com/WildanFrananda/kinetix-notary-service/internal/metrics"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/payment"
	"github.com/WildanFrananda/kinetix-notary-service/internal/proving"
	"github.com/WildanFrananda/kinetix-notary-service/internal/registering"
	"github.com/WildanFrananda/kinetix-notary-service/internal/sponsoring"
	"github.com/WildanFrananda/kinetix-notary-service/internal/watching"
	"github.com/WildanFrananda/kinetix-notary-service/internal/wiring"
)

var version = "unknown"

type daemon struct {
	dig.In
	Config    config.Daemon
	Client    *ethclient.Client
	Pool      *pgxpool.Pool
	Payment   *payment.Source
	Ledger    notary.Ledger
	Chain     notary.AnchorChain
	Anchorer  *anchoring.Anchorer
	Watch     *epochwatch.Client
	Runner    *anchoring.Runner
	Prover    *proving.Prover
	Directory *proving.DirectoryCheck
	Sponsors  *sponsoring.Watcher
	Registrar *registering.Registrar
	Answerer  *answering.Answerer
	Answers   notary.AnswerChain
	Watcher   *watching.Watcher
	Evidence  *evidencing.Runner
	Registry  *evidencing.RegistryCheck
	Evidences notary.EvidenceChain
	Lifecycle *feeds.LifecycleFeed
	Delivery  *feeds.DeliveryFeed
	Reviews   *feeds.ReviewFeed
	Handler   http.Handler
}

func main() {
	os.Exit(run())
}

func run() int {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := start(ctx, log); err != nil {
		log.Error("notaryd stopped", "kind", notary.KindOf(err).String(), "err", err.Error())
		return 1
	}

	return 0
}

func start(ctx context.Context, log *slog.Logger) error {
	container, err := wiring.Daemon(
		os.LookupEnv,
		metrics.Build{Service: "kinetix-notary-service", Version: version},
		log,
	)
	if err != nil {
		return notary.Errf(notary.KindUnknown, "notaryd.start", err, "")
	}

	err = container.Invoke(func(d daemon) error {
		defer d.Client.Close()
		defer d.Pool.Close()
		defer func() { _ = d.Payment.Close() }()
		defer func() { _ = d.Watch.Close() }()
		defer func() { _ = d.Lifecycle.Close() }()
		defer func() { _ = d.Delivery.Close() }()
		defer func() { _ = d.Reviews.Close() }()

		return serve(ctx, log, d)
	})
	if err != nil {
		cause := dig.RootCause(err)
		return notary.Errf(notary.KindOf(cause), "notaryd.start", cause, "")
	}

	return nil
}

func serve(ctx context.Context, log *slog.Logger, d daemon) error {
	identity, err := d.Ledger.Identity(ctx)
	if err != nil {
		return err
	}

	if err := d.Anchorer.CheckCouncil(ctx); err != nil {
		return err
	}

	if err := d.Prover.CheckCourt(ctx); err != nil {
		return err
	}

	if err := d.Directory.Check(ctx); err != nil {
		return err
	}

	if err := d.Sponsors.Check(ctx); err != nil {
		return err
	}

	if err := d.Registrar.CheckWithholding(ctx); err != nil {
		return err
	}

	if err := d.Answerer.CheckScan(ctx); err != nil {
		return err
	}

	if err := d.Registry.Check(ctx); err != nil {
		return err
	}

	log.InfoContext(
		ctx,
		"notaryd starting",
		"version",
		version,
		"chain_id",
		identity.ChainID,
		"ledger",
		identity.Address.Hex(),
		"anchorer",
		d.Chain.Anchorer().Hex(),
		"council",
		d.Config.Council.Hex(),
		"watch",
		d.Config.WatchAddress,
		"court",
		d.Config.Chain.Court.Hex(),
		"withholding_court",
		d.Config.Withholding.Hex(),
		"directory",
		d.Config.Directory.Hex(),
		"paymaster",
		d.Config.Paymaster.Hex(),
		"answerer",
		d.Answers.Answerer().Hex(),
		"log_range",
		d.Config.LogRange,
		"currency",
		d.Config.Currency,
		"epoch_interval",
		d.Config.EpochInterval.String(),
		"evidence_registry",
		d.Config.Registry.Hex(),
		"evidence_anchorer",
		d.Evidences.Anchorer().Hex(),
		"evidence_batch_interval",
		d.Config.EvidenceBatch.String(),
		"order",
		d.Config.OrderAddress,
		"matching",
		d.Config.MatchingAddress,
		"review",
		d.Config.ReviewAddress,
	)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	server := &http.Server{
		Addr:              ":" + strconv.Itoa(d.Config.HTTPPort),
		Handler:           d.Handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	served := make(chan error, 1)
	go func() { served <- server.ListenAndServe() }()

	ran := make(chan struct{})

	go func() {
		d.Runner.Run(ctx)
		close(ran)
	}()

	watched := make(chan struct{})

	go func() {
		d.Watcher.Run(ctx)
		close(watched)
	}()

	answered := make(chan struct{})

	go func() {
		d.Answerer.Run(ctx, d.Config.Tick)
		close(answered)
	}()

	evidenced := make(chan struct{})

	go func() {
		d.Evidence.Run(ctx)
		close(evidenced)
	}()

	sponsored := make(chan struct{})

	go func() {
		d.Sponsors.Run(ctx)
		close(sponsored)
	}()

	var failure error

	select {
	case <-ctx.Done():
	case err := <-served:
		if !errors.Is(err, http.ErrServerClosed) {
			failure = notary.Errf(notary.KindUnknown, "notaryd.serve", err, "the HTTP listener stopped")
		}
	}

	cancel()

	shutdown, release := context.WithTimeout(context.WithoutCancel(ctx), 10*time.Second)
	defer release()

	_ = server.Shutdown(shutdown)

	<-ran
	<-watched
	<-answered
	<-sponsored
	<-evidenced
	log.Info("notaryd stopped cleanly")

	return failure
}

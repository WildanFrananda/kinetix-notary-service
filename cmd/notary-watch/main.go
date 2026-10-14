package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/dig"
	"google.golang.org/grpc"

	"github.com/WildanFrananda/kinetix-notary-service/internal/approving"
	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/metrics"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/payment"
	"github.com/WildanFrananda/kinetix-notary-service/internal/wiring"
)

var version = "unknown"

type watch struct {
	dig.In
	Config   config.Watch
	Client   *ethclient.Client
	Payment  *payment.Records
	Approver *approving.Approver
	Signer   notary.ApprovalSigner
	Handler  http.Handler
	GRPC     *grpc.Server
}

func main() {
	os.Exit(run())
}

func run() int {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := start(ctx, log); err != nil {
		log.Error("notary-watch stopped", "kind", notary.KindOf(err).String(), "err", err.Error())
		return 1
	}

	return 0
}

func start(ctx context.Context, log *slog.Logger) error {
	container, err := wiring.Watch(
		os.LookupEnv,
		metrics.Build{Service: "kinetix-notary-watch-service", Version: version},
		log,
	)
	if err != nil {
		return notary.Errf(notary.KindUnknown, "notary-watch.start", err, "")
	}

	err = container.Invoke(func(w watch) error {
		defer w.Client.Close()
		defer func() { _ = w.Payment.Close() }()

		return serve(ctx, log, w)
	})
	if err != nil {
		cause := dig.RootCause(err)
		return notary.Errf(notary.KindOf(cause), "notary-watch.start", cause, "")
	}

	return nil
}

func serve(ctx context.Context, log *slog.Logger, w watch) error {
	if err := w.Approver.CheckCouncil(ctx); err != nil {
		return err
	}

	listener, err := (&net.ListenConfig{}).Listen(ctx, "tcp", ":"+strconv.Itoa(w.Config.GRPCPort))
	if err != nil {
		return notary.Errf(notary.KindUnknown, "notary-watch.serve", err, "the gRPC port cannot be opened")
	}

	log.InfoContext(
		ctx,
		"notary-watch starting",
		"version",
		version,
		"chain_id",
		w.Config.Chain.Ledger.ChainID,
		"ledger",
		w.Config.Chain.Ledger.Address.Hex(),
		"council",
		w.Config.Council.Hex(),
		"member",
		w.Signer.Member().Hex(),
		"grpc_port",
		w.Config.GRPCPort,
	)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	server := &http.Server{
		Addr:              ":" + strconv.Itoa(w.Config.HTTPPort),
		Handler:           w.Handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	served := make(chan error, 1)
	go func() { served <- server.ListenAndServe() }()

	approved := make(chan error, 1)
	go func() { approved <- w.GRPC.Serve(listener) }()

	var failure error

	select {
	case <-ctx.Done():
	case err := <-served:
		if !errors.Is(err, http.ErrServerClosed) {
			failure = notary.Errf(notary.KindUnknown, "notary-watch.serve", err, "the HTTP listener stopped")
		}
	case err := <-approved:
		failure = notary.Errf(notary.KindUnknown, "notary-watch.serve", err, "the gRPC listener stopped")
	}

	cancel()

	shutdown, release := context.WithTimeout(context.WithoutCancel(ctx), 10*time.Second)
	defer release()

	_ = server.Shutdown(shutdown)

	stopped := make(chan struct{})

	go func() {
		w.GRPC.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
	case <-shutdown.Done():
		w.GRPC.Stop()
	}

	log.Info("notary-watch stopped cleanly")

	return failure
}

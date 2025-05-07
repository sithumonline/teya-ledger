package server

import (
	"context"
	"errors"
	"flag"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/alienxp03/teya-ledger/api"
	"github.com/alienxp03/teya-ledger/db"
	"github.com/alienxp03/teya-ledger/handler/transaction"
)

func Start() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	go func() {
		<-ctx.Done()
		stop()
	}()

	addr := flag.String("addr", "0.0.0.0:8080", "HTTP network address")
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		logger.Error("Could not listen", "error", err)
		os.Exit(1)
	}

	db := db.NewMemoryStorage()
	if err := db.Initialize(); err != nil {
		log.Fatalf("Could not initialize database %v", err)
	}

	if err := db.SeedData(); err != nil {
		log.Fatalf("failed to seed data: %v", err)
	}

	storage := db.GetStorage()

	transactioner := transaction.New(storage)
	apiImp := api.New(transactioner, logger)

	srv := &http.Server{
		Handler: apiImp,
	}

	go func() {
		<-ctx.Done()
		logger.Info("Shutting down server")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(ctx)
	}()

	// Start the server
	logger.Info("Ready to accept traffic", "address", *addr)
	if err := srv.Serve(lis); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("Could not start server", "error", err)
		os.Exit(1)
	}
}

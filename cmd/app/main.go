package main

import (
	"context"
	"errors"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	"github.com/ndovnar/family-budget-api/internal/api"
	"github.com/ndovnar/family-budget-api/internal/auth"
	"github.com/ndovnar/family-budget-api/internal/authz"
	"github.com/ndovnar/family-budget-api/internal/config"
	"github.com/ndovnar/family-budget-api/internal/store"
	"github.com/ndovnar/family-budget-api/internal/store/mongo"
)

var version, gitCommit, application string

func main() {
	log.Info().Msgf("%s %s (%s) -- %s", application, version, gitCommit, runtime.Version())

	sigCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCtx.Done()
		log.Info().Msg("shutdown signal received - attempting graceful shutdown")
		cancel()
	}()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	myStore, err := mongo.New(sigCtx, cfg.Mongo, application)
	if err != nil {
		log.Fatal().Err(err).Msg("failed creating store")
	}

	auth := auth.New(cfg.Auth)
	authz := authz.New(auth, myStore)

	group, errCtx := errgroup.WithContext(sigCtx)

	group.Go(func() error {
		return runApi(errCtx, cfg.API, auth, authz, myStore)
	})

	if err := group.Wait(); err != nil {
		if !errors.Is(err, context.Canceled) {
			log.Fatal().Err(err).Msg("main - shutdown completed with error(s)")
		}
	}

	log.Info().Msg("main - shutdown completed without errors")

}

func runApi(ctx context.Context, cfg api.Config, auth *auth.Auth, authz *authz.Authz, myStore store.Store) error {
	cancelCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var serveError error

	api := api.New(cfg, auth, authz, myStore)

	go func() {
		serveError = api.Run(ctx)
		if serveError != nil {
			log.Err(serveError).Msg("api - failed to start listening")
			cancel()
		}
		log.Info().Msg("api - stopped accepting new connections")
	}()

	<-cancelCtx.Done()
	log.Info().Msg("api - initiating graceful shutdown")

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 1*time.Minute)
	defer shutdownRelease()

	err := api.Stop(shutdownCtx)
	if err != nil {
		log.Err(err).Msg("api - shutdown of http server exited with error")
	}

	log.Info().Msg("api - graceful shutdown complete")

	return errors.Join(serveError, err)
}

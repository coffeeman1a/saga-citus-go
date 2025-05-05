package main

import (
	"context"
	"time"

	"github.com/coffeeman1a/saga-citus-go/users-service/internal/config"
	"github.com/coffeeman1a/saga-citus-go/users-service/internal/repository"
	"github.com/coffeeman1a/saga-citus-go/users-service/internal/server"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
)

func main() {
	// init logger and env variables
	config.Init()

	// citus pool
	log.Debug("Starting user-service...")
	log.Debug("Creating DB pool...")
	pool, err := pgxpool.New(context.Background(), config.DatabaseURL)
	if err != nil {
		log.WithError(err).Fatal("unable to create pgxpool")
	}
	defer pool.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Debug("Pinging pgxpool...")
	if err := pool.Ping(ctx); err != nil {
		log.WithError(err).Fatal("pgxpool ping faild, it's likely no any connection to database")
	}

	log.Info("Databse connection pool is up and running")

	log.Debug("Creaying user repository...")
	userRepo := repository.NewPGXRepository(pool)

	log.Debug("Starting server...")
	// start a server
	s := server.StartServer(userRepo)
	if err := s.Run(); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
	log.Info("Server up and running")

}

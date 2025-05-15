package main

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/coffeeman1a/saga-citus-go/orchestrator/internal/config"
	"github.com/coffeeman1a/saga-citus-go/orchestrator/internal/queue"
	"github.com/coffeeman1a/saga-citus-go/orchestrator/internal/saga"
)

func main() {
	config.Init()

	rclient, err := queue.NewRedisClient(config.RedisURL)
	if err != nil {
		log.WithError(err).Fatal("Unable to create redis client")
	}

	log.Debug("Checking redis client...")
	ping, err := rclient.Ping(context.Background()).Result()
	if err != nil {
		log.WithError(err).Fatal("Unbale to ping redis client")
	}

	if ping != "PONG" {
		log.WithError(errors.New("redis client ping failed")).Fatal("Redis client ping failed, it's likely no any connection to redis")
	}

	log.Info("Redis client is up and running")

	log.Info("Orchestrator started. Waiting for events...")

	for {
		event, err := queue.ReadNextEvent(rclient, "saga_queue")
		if err != nil {
			log.WithError(err).Error("Error reading from queue")
			continue
		}

		log.WithFields(log.Fields{
			"step":  event.Step,
			"tx_id": event.TxID,
		}).Info("Processing step")

		err = saga.HandleEvent(event)
		if err != nil {
			log.WithError(err).Warn("Step failed")
		} else {
			log.WithField("step", event.Step).Info("Step success")
		}
	}
}

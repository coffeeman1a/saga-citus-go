package queue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/coffeeman1a/saga-citus-go/orchestrator/internal/models"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func NewRedisClient(url string) (*redis.Client, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	return redis.NewClient(opts), nil
}

func ReadNextEvent(rdb *redis.Client, queueName string) (*models.SagaEvent, error) {
	res, err := rdb.BRPop(ctx, 0*time.Second, queueName).Result()
	if err != nil || len(res) < 2 {
		return nil, err
	}

	var event models.SagaEvent
	if err := json.Unmarshal([]byte(res[1]), &event); err != nil {
		return nil, err
	}
	return &event, nil
}

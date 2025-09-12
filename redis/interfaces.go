package redis

import (
	"time"

	"github.com/redis/go-redis/v9"
)

type Client = redis.Client

type DistributedLock struct {
	client *Client
	key    string
	value  string
	expire time.Duration
}
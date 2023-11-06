package cache

import (
	"github.com/go-redis/redis"
	"time"
)

type Handler struct {
	Client *redis.Client
}

var (
	DefaultExpiration = 5 * time.Minute
	DefaultListenAddr = ":6379"
)

func Initialize(listenAddr string) *Handler {
	return &Handler{Client: redis.NewClient(&redis.Options{
		Addr:     "localhost" + listenAddr,
		Password: "",
		DB:       0,
	})}
}

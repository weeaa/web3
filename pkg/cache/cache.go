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
	DefaultPort       = ":6379"
)

func Initialize(port string) *Handler {
	return &Handler{Client: redis.NewClient(&redis.Options{
		Addr:     "localhost" + port,
		Password: "",
		DB:       0,
	})}
}

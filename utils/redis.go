package utils

import (
	"fmt"
	"os"
	"sync"

	"github.com/redis/go-redis/v9"
)

type redisSingleton struct {
	client *redis.Client
	once   sync.Once
}

var instance *redisSingleton
var once sync.Once

// Retorna a conexão com o redis utilizando uma estratégia de Singleton
func GetRedisClient() *redis.Client {
	once.Do(func() {

		host := os.Getenv("CACHE_HOST")
		port := os.Getenv("CACHE_PORT")

		instance = &redisSingleton{}
		instance.client = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", host, port),
			Password: "",
			DB:       0,
			PoolSize: 10,
		})
	})
	return instance.client
}

package db

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var ctx = context.Background()

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Assure-toi que Redis tourne sur ce port
		Password: "",               // Met un mot de passe si Redis en exige un
		DB:       0,                // Sélectionne la base 0
	})

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		fmt.Println("❌ Erreur de connexion à Redis :", err)
	} else {
		fmt.Println("✅ Connexion Redis réussie !")
	}
}

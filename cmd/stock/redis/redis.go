package redis

import (
	ctx "context"
	"fmt"
	"log"
	t "time"

	"github.com/go-redis/redis/v8"
	conf "github.com/gotway/service-examples/cmd/stock/config"
)

var client *redis.Client

// Init initializes redis client
func Init() {
	opts, err := redis.ParseURL(conf.RedisURL)
	if err != nil {
		log.Fatalf("Error parsing redis URL %s. Error: %s", conf.RedisURL, err)
	}
	client = redis.NewClient(opts)

	if err := client.Ping(ctx.Background()).Err(); err != nil {
		log.Fatalf("Error connecting to redis %s. Error: %s", conf.RedisURL, err)
	}
	log.Print("Connected to redis")
}

// Set stores a key with a TTL expresed in seconds
func Set(key string, value interface{}, ttl int) (int, error) {
	prefixedKey := getPrefixedKey(key)
	narrowedTTL := getNarrowedTTL(ttl)
	duration := t.Duration(narrowedTTL) * t.Second
	_, err := client.Set(ctx.Background(), prefixedKey, value, duration).Result()
	if err != nil {
		return 0, err
	}
	return narrowedTTL, nil
}

// Get gets a key
func Get(key string) (string, error) {
	prefixedKey := getPrefixedKey(key)
	return client.Get(ctx.Background(), prefixedKey).Result()
}

// TTL gets the TTL of a key expressed in seconds
func TTL(key string) (int, error) {
	prefixedKey := getPrefixedKey(key)
	duration, err := client.TTL(ctx.Background(), prefixedKey).Result()
	if err != nil {
		return 0, err
	}
	ttl := int(duration / t.Second)
	return ttl, nil
}

func getPrefixedKey(key string) string {
	return fmt.Sprintf("%s%s", conf.RedisPrefix, key)
}

func getNarrowedTTL(ttl int) int {
	if ttl > 0 && ttl < conf.RedisTTLMax {
		return ttl
	}
	return conf.RedisTTLDefault
}

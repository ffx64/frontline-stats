package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	TTLLeaderboard = 5 * time.Minute
	TTLPlayerStats = 14 * time.Minute
	TTLPlayer      = 30 * time.Minute
	TTLRound       = 5 * time.Minute
	TTLServer      = 60 * time.Minute
	TTLScoreboard  = 5 * time.Minute

	KeyLeaderboardKills     = "leaderboard:kills"
	KeyLeaderboardHeadshots = "leaderboard:headshots"
	KeyLeaderboardVehicles  = "leaderboard:vehicles"
)

func KeyPlayer(guid string) string        { return fmt.Sprintf("player:%s", guid) }
func KeyPlayerStats(guid string) string   { return fmt.Sprintf("player:stats:%s", guid) }
func KeyRound(id string) string           { return fmt.Sprintf("round:%s", id) }
func KeyScoreboard(roundID string) string { return fmt.Sprintf("round:scoreboard:%s", roundID) }
func KeyServer(id string) string          { return fmt.Sprintf("server:%s", id) }

// Get fetches a cached value and deserializes it into T.
// Returns (nil, nil) on cache miss. Returns (nil, err) on Redis error.
func Get[T any](ctx context.Context, rdb *redis.Client, key string) (*T, error) {
	if rdb == nil {
		return nil, nil
	}
	val, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var result T
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Set serializes value and stores it with the given TTL.
// A Redis error is logged but not returned — cache writes are non-fatal.
func Set(ctx context.Context, rdb *redis.Client, key string, value any, ttl time.Duration) {
	if rdb == nil {
		return
	}
	data, err := json.Marshal(value)
	if err != nil {
		log.Printf("[cache] failed to marshal value for key %s: %v", key, err)
		return
	}
	if err := rdb.Set(ctx, key, data, ttl).Err(); err != nil {
		log.Printf("[cache] failed to set key %s: %v", key, err)
	}
}

// Delete removes one or more keys. Errors are logged but not returned.
func Delete(ctx context.Context, rdb *redis.Client, keys ...string) {
	if rdb == nil || len(keys) == 0 {
		return
	}
	if err := rdb.Del(ctx, keys...).Err(); err != nil {
		log.Printf("[cache] failed to delete keys %v: %v", keys, err)
	}
}

// DeletePattern removes all keys matching a glob pattern using SCAN.
// Errors are logged but not returned.
func DeletePattern(ctx context.Context, rdb *redis.Client, pattern string) {
	if rdb == nil {
		return
	}
	var cursor uint64
	for {
		keys, next, err := rdb.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			log.Printf("[cache] failed to scan pattern %s: %v", pattern, err)
			return
		}
		if len(keys) > 0 {
			if err := rdb.Del(ctx, keys...).Err(); err != nil {
				log.Printf("[cache] failed to delete keys for pattern %s: %v", pattern, err)
			}
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
}

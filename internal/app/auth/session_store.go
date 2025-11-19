package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
var sessionTTL = 24 * time.Hour

func hexToUint64(hexStr string) (uint64, error) {
	return strconv.ParseUint(hexStr, 16, 64)
}

// SessionTTLSeconds возвращает время жизни сессии в секундах
func SessionTTLSeconds() int {
	return int(sessionTTL.Seconds())
}

func InitRedis(addr string) error {
	if addr == "" {
		addr = os.Getenv("REDIS_ADDR") // e.g. "localhost:6379"
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr: addr,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		return err
	}
	return nil
}

func newSessionID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// CreateSession stores userID in redis, returns sessionID
func CreateSession(ctx context.Context, userID uint) (string, error) {
	if redisClient == nil {
		return "", errors.New("redis not initialized")
	}
	sid, err := newSessionID()
	if err != nil {
		return "", err
	}
	key := "session:" + sid
	if err := redisClient.Set(ctx, key, userID, sessionTTL).Err(); err != nil {
		return "", err
	}
	return sid, nil
}

// GetUserIDBySession returns userID or 0 if not found
func GetUserIDBySession(ctx context.Context, sid string) (uint, error) {
	if redisClient == nil {
		return 0, errors.New("redis not initialized")
	}
	key := "session:" + sid
	val, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	// parse
	var uid uint64
	uid, err = hexToUint64(val)
	if err == nil {
		return uint(uid), nil
	}
	// try normal parse (db stored as decimal)
	var n uint64
	_, err = fmt.Sscanf(val, "%d", &n)
	if err != nil {
		return 0, err
	}
	return uint(n), nil
}

func DeleteSession(ctx context.Context, sid string) error {
	if redisClient == nil {
		return errors.New("redis not initialized")
	}
	return redisClient.Del(ctx, "session:"+sid).Err()
}

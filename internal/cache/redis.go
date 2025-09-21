package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisCache(addr, password string, db int) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisCache{
		client: rdb,
		ctx:    context.Background(),
	}
}

func (r *RedisCache) Ping() error {
	_, err := r.client.Ping(r.ctx).Result()
	return err
}

// Generic cache operations
func (r *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %v", err)
	}

	return r.client.Set(r.ctx, key, data, expiration).Err()
}

func (r *RedisCache) Get(key string, dest interface{}) error {
	data, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("key not found")
		}
		return err
	}

	return json.Unmarshal([]byte(data), dest)
}

func (r *RedisCache) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

func (r *RedisCache) Exists(key string) bool {
	exists, _ := r.client.Exists(r.ctx, key).Result()
	return exists > 0
}

func (r *RedisCache) Expire(key string, expiration time.Duration) error {
	return r.client.Expire(r.ctx, key, expiration).Err()
}

// Session management
func (r *RedisCache) SetSession(sessionID string, userID uint64, data map[string]interface{}, expiration time.Duration) error {
	sessionKey := fmt.Sprintf("session:%s", sessionID)
	sessionData := map[string]interface{}{
		"user_id":    userID,
		"data":       data,
		"created_at": time.Now(),
	}

	return r.Set(sessionKey, sessionData, expiration)
}

func (r *RedisCache) GetSession(sessionID string) (map[string]interface{}, error) {
	sessionKey := fmt.Sprintf("session:%s", sessionID)
	var sessionData map[string]interface{}

	err := r.Get(sessionKey, &sessionData)
	if err != nil {
		return nil, err
	}

	// Refresh session expiration
	r.Expire(sessionKey, 24*time.Hour)
	return sessionData, nil
}

func (r *RedisCache) DeleteSession(sessionID string) error {
	sessionKey := fmt.Sprintf("session:%s", sessionID)
	return r.Delete(sessionKey)
}

// Cache patterns for common data
func (r *RedisCache) CacheKoperasi(koperasiID uint64, data interface{}) error {
	key := fmt.Sprintf("koperasi:%d", koperasiID)
	return r.Set(key, data, 1*time.Hour)
}

func (r *RedisCache) GetKoperasi(koperasiID uint64, dest interface{}) error {
	key := fmt.Sprintf("koperasi:%d", koperasiID)
	return r.Get(key, dest)
}

func (r *RedisCache) InvalidateKoperasi(koperasiID uint64) error {
	key := fmt.Sprintf("koperasi:%d", koperasiID)
	return r.Delete(key)
}

// Product cache
func (r *RedisCache) CacheProduct(productID uint64, data interface{}) error {
	key := fmt.Sprintf("product:%d", productID)
	return r.Set(key, data, 30*time.Minute)
}

func (r *RedisCache) GetProduct(productID uint64, dest interface{}) error {
	key := fmt.Sprintf("product:%d", productID)
	return r.Get(key, dest)
}

func (r *RedisCache) CacheProductByBarcode(barcode string, data interface{}) error {
	key := fmt.Sprintf("product:barcode:%s", barcode)
	return r.Set(key, data, 30*time.Minute)
}

func (r *RedisCache) GetProductByBarcode(barcode string, dest interface{}) error {
	key := fmt.Sprintf("product:barcode:%s", barcode)
	return r.Get(key, dest)
}

// Wilayah cache (longer TTL for static data)
func (r *RedisCache) CacheWilayah(tipe string, id uint64, data interface{}) error {
	key := fmt.Sprintf("wilayah:%s:%d", tipe, id)
	return r.Set(key, data, 24*time.Hour)
}

func (r *RedisCache) GetWilayah(tipe string, id uint64, dest interface{}) error {
	key := fmt.Sprintf("wilayah:%s:%d", tipe, id)
	return r.Get(key, dest)
}

// Report cache
func (r *RedisCache) CacheReport(reportType string, koperasiID uint64, params string, data interface{}) error {
	key := fmt.Sprintf("report:%s:%d:%s", reportType, koperasiID, params)
	return r.Set(key, data, 5*time.Minute)
}

func (r *RedisCache) GetReport(reportType string, koperasiID uint64, params string, dest interface{}) error {
	key := fmt.Sprintf("report:%s:%d:%s", reportType, koperasiID, params)
	return r.Get(key, dest)
}

// Dashboard analytics cache
func (r *RedisCache) CacheDashboard(koperasiID uint64, data interface{}) error {
	key := fmt.Sprintf("dashboard:%d", koperasiID)
	return r.Set(key, data, 10*time.Minute)
}

func (r *RedisCache) GetDashboard(koperasiID uint64, dest interface{}) error {
	key := fmt.Sprintf("dashboard:%d", koperasiID)
	return r.Get(key, dest)
}

// Rate limiting
func (r *RedisCache) IncrementRateLimit(key string, window time.Duration) (int64, error) {
	pipe := r.client.TxPipeline()
	incr := pipe.Incr(r.ctx, key)
	pipe.Expire(r.ctx, key, window)
	_, err := pipe.Exec(r.ctx)
	if err != nil {
		return 0, err
	}
	return incr.Val(), nil
}

func (r *RedisCache) CheckRateLimit(key string, limit int64) (bool, int64, error) {
	count, err := r.client.Get(r.ctx, key).Int64()
	if err == redis.Nil {
		return true, 0, nil
	}
	if err != nil {
		return false, 0, err
	}
	return count < limit, count, nil
}

// Distributed lock
func (r *RedisCache) AcquireLock(key string, value string, expiration time.Duration) (bool, error) {
	lockKey := fmt.Sprintf("lock:%s", key)
	ok, err := r.client.SetNX(r.ctx, lockKey, value, expiration).Result()
	return ok, err
}

func (r *RedisCache) ReleaseLock(key string, value string) error {
	lockKey := fmt.Sprintf("lock:%s", key)

	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`

	_, err := r.client.Eval(r.ctx, script, []string{lockKey}, value).Result()
	return err
}

// List operations for leaderboards/rankings
func (r *RedisCache) AddToSortedSet(key string, score float64, member interface{}) error {
	data, err := json.Marshal(member)
	if err != nil {
		return err
	}

	return r.client.ZAdd(r.ctx, key, &redis.Z{
		Score:  score,
		Member: string(data),
	}).Err()
}

func (r *RedisCache) GetTopFromSortedSet(key string, count int64) ([]string, error) {
	result, err := r.client.ZRevRange(r.ctx, key, 0, count-1).Result()
	return result, err
}

// Pub/Sub for real-time notifications
func (r *RedisCache) Publish(channel string, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return r.client.Publish(r.ctx, channel, data).Err()
}

func (r *RedisCache) Subscribe(channel string) *redis.PubSub {
	return r.client.Subscribe(r.ctx, channel)
}

// Batch operations
func (r *RedisCache) BatchSet(items map[string]interface{}, expiration time.Duration) error {
	pipe := r.client.TxPipeline()

	for key, value := range items {
		data, err := json.Marshal(value)
		if err != nil {
			return err
		}
		pipe.Set(r.ctx, key, data, expiration)
	}

	_, err := pipe.Exec(r.ctx)
	return err
}

func (r *RedisCache) BatchGet(keys []string) (map[string]string, error) {
	pipe := r.client.TxPipeline()
	cmds := make([]*redis.StringCmd, len(keys))

	for i, key := range keys {
		cmds[i] = pipe.Get(r.ctx, key)
	}

	_, err := pipe.Exec(r.ctx)
	if err != nil && err != redis.Nil {
		return nil, err
	}

	result := make(map[string]string)
	for i, cmd := range cmds {
		if val, err := cmd.Result(); err == nil {
			result[keys[i]] = val
		}
	}

	return result, nil
}

// Close connection
func (r *RedisCache) Close() error {
	return r.client.Close()
}
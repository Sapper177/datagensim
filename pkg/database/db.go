package database

import (
	"fmt"
	"time"
	"context"

    "github.com/redis/go-redis/v9"
)

// RedisClient is a struct that holds the Redis client and its options.
type RedisClient struct {
	client *redis.Client
	options *redis.Options
	ctx   context.Context
}

// NewRedisClient creates a new Redis client with the given options.
func NewRedisClient(addr string, psswd string, db int, readTimeout time.Duration, writeTimeout time.Duration) *RedisClient {
	options := &redis.Options{
		Addr:         addr,
		Password:     psswd,
		DB:           db,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
	client := redis.NewClient(options)
	return &RedisClient{
		client:  client,
		options: options,
	}
}

// Set sets a key-value pair in Redis with an expiration time.
func (r *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	err := r.client.Set(r.ctx, key, value, expiration).Err()
	return HandleDbError(err, key, "set data")
}

// Get retrieves a value from Redis by its key.
func (r *RedisClient) Get(key string) (string, error) {
	val, err := r.client.Get(r.ctx, key).Result()
	return val, HandleDbError(err, key, "retrieve data")
}

// Del deletes a key from Redis.
func (r *RedisClient) Del(key string) error {
	err := r.client.Del(r.ctx, key).Err()
	return HandleDbError(err, key, "delete data")
}

// Close closes the Redis client connection.
func (r *RedisClient) Close() error {
	err := r.client.Close()
	return HandleDbError(err, "", "close connection")
}

// Set up Redis Database to follow this same structure and these functions should work
// 	<bus>:
//   	- <payload_id>
//   	- <payload_id>
//   	- ...
func (r *RedisClient) GetPayloads(bus string) ([]string, error) {
	stringSlice, err := r.client.LRange(r.ctx, bus, 0, -1).Result()
	return stringSlice, HandleDbError(err, bus, "retrieve bus payloads")
}

// 	<payload_id>:
//     	packet_type: <type>
//     	frequency: <frequency>
func (r* RedisClient) GetPayloadInfo(payloadId string) (map[string]string, error) {
	retrievedMap, err := r.client.HGetAll(r.ctx, payloadId).Result()
	return retrievedMap, HandleDbError(err, payloadId, "retrieve bus payloads")
}

//	<payload_id>_data:
//       - <data_id>
//       - <data_id>
//       - <data_id>
//		 ...
func (r *RedisClient) GetPayloadData(payloadId string) ([]string, error) {
	payloadDataKey := payloadId + "_data"
	stringSlice, err := r.client.LRange(r.ctx, payloadDataKey, 0, -1).Result()
	return stringSlice, HandleDbError(err, payloadDataKey, "retrieve payload data")
}

// 	<data_id>:
//		name: <name>
//		value: <value>
//		raw_value: <value>
// 		data_type: <type>
//		offset: <value>
func (r* RedisClient) GetData(data_id string) (map[string]string, error) {
	retrievedMap, err := r.client.HGetAll(r.ctx, data_id).Result()
	return retrievedMap, HandleDbError(err, data_id, "retrieve data")
}

// 	<data_id>-info: <- example analog info
//		min: <value>
//		max: <value>
//		step: <value>
//		frequency: <value>
//		calibration: <calib_id>
//		calibration_type: <type>

// 	<calib_id>:	<- example analog calibration
//		<slope>: <value>
//		<intercept>: <value>
//		<constants>:
//			<c1>: <value>
//			<c2>: <value>
//			<c3>: <value>
//			...
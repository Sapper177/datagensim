package database

import (
	"fmt"
	"github.com/redis/go-redis/v9"
)

func HandleDbError(err error, key string,  msg string) (error){
	if err == redis.Nil {
		// Key doesn't exist or is not a list
		return fmt.Errorf("key '%s' does not exist or is empty", key)
	} else if err != nil {
		// Other Redis error (network, etc.)
		return fmt.Errorf("failed to %s using key '%s' from Redis: %v", msg, key, err)
	}
	return nil
}
package database

import (
	"fmt"
	"github.com/redis/go-redis/v9"
)

func HandleDbError(err error, key string,  msg string) (error){
	if err == redis.Nil {
		// Key doesn't exist or is not a list
		return fmt.Errorf("Key '%s' does not exist or is empty.\n", key)
	} else if err != nil {
		// Other Redis error (network, etc.)
		return fmt.Errorf("Failed to %s using key '%s' from Redis: %v", msg, key, err)
	}
	return nil
}
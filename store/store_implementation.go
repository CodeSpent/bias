package store

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AbstractBaseStoreImplementation struct {
	DB    *gorm.DB
	Redis redis.Client
}

/*func NewAbstractBaseStore(db *gorm.DB, redisClient redis.Client) Store {
	return &AbstractBaseStoreImplementation{
		DB:    db,
		Redis: redisClient,
	}
}
*/

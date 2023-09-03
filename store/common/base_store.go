package common

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AbstractBaseStore struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func NewAbstractBaseStore(db *gorm.DB, redis *redis.Client) *AbstractBaseStore {
	return &AbstractBaseStore{
		DB:    db,
		Redis: redis,
	}
}

func (s *AbstractBaseStore) List() error {
	return s.DB.Error
}

func (s *AbstractBaseStore) Create(model interface{}) error {
	return s.DB.Create(model).Error
}

func (s *AbstractBaseStore) Retrieve(model interface{}, id uint) error {
	return s.DB.First(model, id).Error
}

func (s *AbstractBaseStore) Update(model interface{}) error {
	return s.DB.Save(model).Error
}

func (s *AbstractBaseStore) Delete(model interface{}) error {
	return s.DB.Delete(model).Error
}

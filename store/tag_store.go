package store

import (
	"bias/models"
	"bias/store/common"
	"context"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type TagRepository interface {
	CreateTag(tag *models.TagModel) error
	GetTagByID(id uint) (*models.TagModel, error)
	UpdateTag(tag *models.TagModel) error
	DeleteTag(id uint) error
}

type TagStore struct {
	common.AbstractBaseStore
	DB    *gorm.DB
	Redis *redis.Client
	context.Context
}

func NewTagStore(db *gorm.DB, redisClient *redis.Client, ctx context.Context) *TagStore {
	return &TagStore{
		DB:      db,
		Redis:   redisClient,
		Context: ctx,
	}
}

func (s *TagStore) GetAllTags() ([]*models.TagModel, error) {
	var tags []*models.TagModel
	if err := s.DB.Table("tags").Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (s *TagStore) CreateTag(tag *models.TagModel) error {
	return s.DB.Create(tag).Error
}

func (s *TagStore) GetTagByID(id uint) (*models.TagModel, error) {
	var tag models.TagModel
	if err := s.DB.First(&tag, id).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (s *TagStore) UpdateTag(tag *models.TagModel) error {
	return s.DB.Save(tag).Error
}

func (s *TagStore) DeleteTag(id uint) error {
	tag := models.TagModel{ID: id}
	return s.DB.Delete(&tag).Error
}

package store

import (
	"bias/models"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"strconv"
)

type StreamRepository interface {
	CreateStream(stream *models.StreamModel) error
	GetStreamByID(id uint) (*models.StreamModel, error)
	UpdateStream(stream *models.StreamModel) error
	DeleteStream(stream *models.StreamModel) error
	FindStreamsByGameID(gameID string) ([]models.StreamModel, error)
	GetAllStreams() ([]models.StreamModel, error)
}

type StreamStore struct {
	DB      *gorm.DB
	Redis   *redis.Client
	Context context.Context
}

func NewStreamStore(db *gorm.DB, redisClient *redis.Client, ctx context.Context) *StreamStore {
	return &StreamStore{
		DB:      db,
		Redis:   redisClient,
		Context: ctx,
	}
}

func (s *StreamStore) GetAllStreams() ([]models.StreamModel, error) {
	var streams []models.StreamModel

	streamsMap, err := s.Redis.HGetAll(s.Context, "streams").Result()
	if err != nil {
		return nil, err
	}

	for _, streamJSON := range streamsMap {
		var stream models.StreamModel
		if err := json.Unmarshal([]byte(streamJSON), &stream); err != nil {
			return nil, err
		}

		streams = append(streams, stream)
	}

	return streams, nil
}

func (s *StreamStore) CreateStream(stream *models.StreamModel) error {
	streamJSON, err := json.Marshal(stream)
	if err != nil {
		return err
	}

	stream.ID = uuid.New()

	err = s.Redis.HSet(s.Context, "streams", stream.ID, streamJSON).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *StreamStore) GetStreamByID(id uint) (*models.StreamModel, error) {
	streamJSON, err := s.Redis.HGet(s.Context, "streams", strconv.Itoa(int(id))).Result()
	if err != nil {
		return nil, err
	}

	var stream models.StreamModel
	err = json.Unmarshal([]byte(streamJSON), &stream)
	if err != nil {
		return nil, err
	}

	return &stream, nil
}

func (s *StreamStore) UpdateStream(stream *models.StreamModel) error {
	updatedJSON, err := json.Marshal(stream)
	if err != nil {
		return err
	}

	err = s.Redis.HSet(s.Context, "streams", stream.ID, updatedJSON).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *StreamStore) DeleteStream(stream *models.StreamModel) error {
	err := s.Redis.HDel(s.Context, "streams", stream.ID.String()).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *StreamStore) FindStreamsByGameID(gameID string) ([]models.StreamModel, error) {
	streamsMap, err := s.Redis.HGetAll(s.Context, "streams").Result()
	if err != nil {
		return nil, err
	}

	var streams []models.StreamModel
	for _, streamJSON := range streamsMap {
		var stream models.StreamModel
		err = json.Unmarshal([]byte(streamJSON), &stream)
		if err != nil {
			return nil, err
		}
		streams = append(streams, stream)
	}

	return streams, nil
}

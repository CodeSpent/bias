package models

import (
	_ "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"time"
)

type Stream interface {
	GetID() string
	GetUserID() string
	GetUserName() string
	GetGameID() string
	GetTitle() string
	GetViewerCount() int
	GetStartedAt() time.Time
}

type StreamModel struct {
	ID          uuid.UUID
	UserID      string    `validate:"required"`
	UserName    string    `validate:"required"`
	GameID      string    `validate:"required"`
	Title       string    `validate:"required"`
	ViewerCount int       `validate:"min=0"`
	StartedAt   time.Time `validate:"required"`
	Tags        []TagModel
}

func (s StreamModel) GetID() uuid.UUID {
	return s.ID
}

func (s StreamModel) GetUserID() string {
	return s.UserID
}

func (s StreamModel) GetUserName() string {
	return s.UserName
}

func (s StreamModel) GetGameID() string {
	return s.GameID
}

func (s StreamModel) GetTitle() string {
	return s.Title
}

func (s StreamModel) GetViewerCount() int {
	return s.ViewerCount
}

func (s StreamModel) GetStartedAt() string {
	return s.StartedAt.String()
}

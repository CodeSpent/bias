package models

import (
	"github.com/lib/pq"
	"time"
)

type Tag interface {
	GetID() uint
	GetName() string
	GetType() string
	GetParent() Tag
}

type TagModel struct {
	ID        uint
	Name      string
	Keywords  pq.StringArray `gorm:"type:text[]"`
	Type      string
	ParentID  uint
	Parent    *TagModel `gorm:"foreignkey:ParentID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t TagModel) GetID() uint {
	return t.ID
}

func (t TagModel) GetName() string {
	return t.Name
}

func (t TagModel) GetType() string {
	return t.Type
}

func (t TagModel) GetParent() Tag {
	return t.Parent
}

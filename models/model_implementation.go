package models

import (
	"gorm.io/gorm"
	"time"
)

type AbstractBaseModelImplementation struct {
	ID         uint
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	_modelName string
}

func NewAbstractBaseModel(store StoreInterface) AbstractBaseModel {
	return &AbstractBaseModelImplementation{}
}

func (m *AbstractBaseModelImplementation) GetModelName() string {
	return m._modelName
}

package models

type AbstractBaseModel interface {
}

type BaseModelImplementation struct {
	ID uint
}

type StoreInterface interface {
	List(model interface{}) error
	Create(model interface{}) error
	Retrieve(model interface{}, id uint) error
	Update(model interface{}) error
	Delete(model interface{}) error
}

package models

import "github.com/google/uuid"

type MetaRepository interface {
	AddMeta(data Meta) error
	DeleteMeta(data Meta) error
	UpdateMeta(uid uuid.UUID, data Meta) error
	GetMeta(uid uuid.UUID) (Meta, error)
}

type Meta struct {
	UUID          uuid.UUID
	NumberOfPages *uint32
	Height        *float32
	Width         *float32
	Images        *map[uint32]string
}

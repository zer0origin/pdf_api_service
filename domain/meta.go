package domain

import "github.com/google/uuid"

type MetaRepository interface {
	AddMeta(data MetaData) error
	DeleteMeta(data MetaData) error
	UpdateMeta(data MetaData) error
}

type MetaData struct {
	UUID          uuid.UUID
	NumberOfPages *uint32
	Height        *float32
	Width         *float32
	Images        *map[uint32]string
}

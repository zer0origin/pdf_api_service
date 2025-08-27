package models

import "github.com/google/uuid"

type MetaRepository interface {
	AddMeta(data Meta) error
	DeleteMeta(data Meta) error
	UpdateMeta(uid uuid.UUID, data Meta) error
	GetMeta(uid uuid.UUID) (Meta, error)
}

type Meta struct {
	DocumentUUID  uuid.UUID `json:"documentUUID" example:"ba3ca973-5052-4030-a528-39b49736d8ad"`
	NumberOfPages *uint32   `json:"numberOfPages" example:"31"`
	Width         *float32  `json:"width" example:"1920"`
	Height        *float32  `json:"height" example:"1080"`
	Images        *map[uint32]string
	OwnerUUID     *uuid.UUID `json:"ownerUUID" example:"34906041-2d68-45a2-9671-9f0ba89f31a9"`
	OwnerType     *int       `json:"ownerType" example:"1"`
}

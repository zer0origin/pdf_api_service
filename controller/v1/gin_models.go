package v1

import (
	"github.com/google/uuid"
	"pdf_service_api/models"
)

type GetDocumentRequest struct {
	DocumentUUID *uuid.UUID `json:"documentUUID"`
	OwnerUUID    *uuid.UUID `json:"ownerUUID"`
}

type CreateRequest struct {
	DocumentBase64String string     `json:"documentBase64String"`
	DocumentTitle        *string    `json:"documentTitle"`
	OwnerUUID            *uuid.UUID `json:"ownerUUID"`
	OwnerType            *int       `json:"ownerType"`
}

type AddNewSelectionRequest struct {
	DocumentUUID    *uuid.UUID                        `json:"documentUUID,omitempty"`
	IsComplete      bool                              `json:"isComplete,omitempty"`
	Settings        *string                           `json:"settings,omitempty"`
	SelectionBounds *map[int][]models.SelectionBounds `json:"selectionBounds,omitempty"`
}

type AddMetaRequestGenerated struct {
	DocumentUUID         uuid.UUID `json:"documentUUID" `
	OwnerUUID            uuid.UUID `json:"ownerUUID"`
	DocumentBase64String *string   `json:"documentBase64String"`
}

type UpdateMetaRequest struct {
	UUID          uuid.UUID
	NumberOfPages *uint32
	Height        *float32
	Width         *float32
	Images        *map[uint32]string
}

type DeleteMetaRequest struct {
	UUID uuid.UUID
}

type AddMetaResponse struct {
	MetaUUID uuid.UUID `json:"metaUUID"`
}

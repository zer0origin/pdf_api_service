package v1

import (
	"pdf_service_api/models"

	"github.com/google/uuid"
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

type AddMetaRequest struct {
	DocumentUUID         uuid.UUID `json:"documentUUID" `
	OwnerUUID            uuid.UUID `json:"ownerUUID"`
	OwnerType            int       `json:"ownerType"`
	DocumentBase64String *string   `json:"documentBase64String"`
}

type UpdateMetaRequest struct {
	UUID          uuid.UUID
	NumberOfPages *uint32
	Height        *float32
	Width         *float32
	Images        *map[string]string
}

type DeleteMetaRequest struct {
	UUID uuid.UUID
}

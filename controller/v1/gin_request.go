package v1

import (
	"github.com/google/uuid"
	"pdf_service_api/models"
)

type GetDocumentRequest struct {
	DocumentUuid uuid.UUID `json:"document_uuid"`
}

type UploadRequest struct {
	DocumentBase64String string `json:"documentBase64String"`
}

type AddNewSelectionRequest struct {
	DocumentUUID    *uuid.UUID                        `json:"documentUUID,omitempty"`
	IsComplete      bool                              `json:"isComplete,omitempty"`
	Settings        *string                           `json:"settings,omitempty"`
	SelectionBounds *map[int][]models.SelectionBounds `json:"selectionBounds,omitempty"`
}

type AddMetaRequest struct {
	NumberOfPages *uint32
	Height        *float32
	Width         *float32
	Images        *map[uint32]string
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

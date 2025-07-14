package v1

import "github.com/google/uuid"

type GetDocumentRequest struct {
	DocumentUuid uuid.UUID `json:"document_uuid"`
}

type UploadRequest struct {
	DocumentBase64String *string `json:"documentBase64String"`
}

type AddNewSelectionRequest struct {
	DocumentID      *uuid.UUID `json:"documentID,omitempty"`
	IsComplete      bool       `json:"isComplete,omitempty"`
	Settings        *string    `json:"settings,omitempty"`
	SelectionBounds *string    `json:"selectionBounds,omitempty"`
}

package models

import "github.com/google/uuid"

type GetDocumentRequest struct {
	DocumentUuid uuid.UUID `json:"document_uuid"`
}

type UploadRequest struct {
	DocumentBase64String *string `json:"documentBase64String"`
}

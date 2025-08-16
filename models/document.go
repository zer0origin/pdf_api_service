package models

import (
	"github.com/google/uuid"
	"time"
)

type Document struct {
	Uuid          uuid.UUID    `json:"documentUUID" example:"ba3ca973-5052-4030-a528-39b49736d8ad"`
	DocumentTitle *string      `json:"documentTitle,omitempty"`
	TimeCreated   *time.Time   `json:"timeCreated,omitempty"`
	OwnerUUID     *uuid.UUID   `json:"ownerUUID,omitempty"`
	OwnerType     *int         `json:"ownerType,omitempty"`
	PdfBase64     *string      `json:"pdfBase64,omitempty"`
	SelectionData *[]Selection `json:"selectionData,omitempty"`
}

type DocumentRepository interface {
	UploadDocument(document Document) error
	GetDocumentByDocumentUUID(id uuid.UUID) (Document, error)
	GetDocumentByOwnerUUID(id uuid.UUID) ([]Document, error)
	DeleteDocumentById(id uuid.UUID) error
}

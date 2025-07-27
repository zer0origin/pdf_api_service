package models

import "github.com/google/uuid"

type Document struct {
	Uuid          uuid.UUID    `json:"documentUUID" example:"ba3ca973-5052-4030-a528-39b49736d8ad"`
	PdfBase64     *string      `json:"pdfBase64"`
	SelectionData *[]Selection `json:"selectionData,omitempty"`
}

type DocumentRepository interface {
	UploadDocument(document Document) error
	GetDocumentByDocumentUUID(id uuid.UUID) (Document, error)
	DeleteDocumentById(id uuid.UUID) error
}

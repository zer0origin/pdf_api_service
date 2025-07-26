package domain

import "github.com/google/uuid"

type Document struct {
	Uuid          uuid.UUID    `json:"documentUUID"`
	PdfBase64     *string      `json:"pdfBase64"`
	SelectionData *[]Selection `json:"selectionData,omitempty"`
}

type DocumentRepository interface {
	UploadDocument(document Document) error
	GetDocumentByDocumentUUID(id uuid.UUID) (Document, error)
	DeleteDocumentById(id uuid.UUID) error
}

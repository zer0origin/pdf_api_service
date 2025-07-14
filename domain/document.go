package domain

import "github.com/google/uuid"

type DocumentRepository interface {
	UploadDocument(document Document) error
	GetDocumentById(id uuid.UUID) (Document, error)
	DeleteDocumentById(id uuid.UUID) error
}

type Document struct {
	Uuid          uuid.UUID    `json:"documentUUID"`
	PdfBase64     *string      `json:"pdfBase64"`
	SelectionData *[]Selection `json:"selectionData,omitempty"`
}

type Selection struct {
	Uuid            uuid.UUID  `json:"selectionUUID"`
	DocumentID      *uuid.UUID `json:"documentID,omitempty"`
	IsComplete      bool       `json:"isComplete,omitempty"`
	Settings        string     `json:"settings,omitempty"`
	SelectionBounds string     `json:"selection_bounds,omitempty"`
}

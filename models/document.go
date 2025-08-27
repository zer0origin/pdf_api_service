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
	GetDocumentByDocumentUUID(document, owner uuid.UUID, excludes Exclude) (Document, error)
	GetDocumentByOwnerUUID(owner uuid.UUID, limit int8, offset int8, excludes Exclude) ([]Document, error)
	DeleteDocumentById(documentUuid, ownerUuid uuid.UUID) error
}

type Exclude map[string]bool

func (e Exclude) DocumentTitle(value bool) Exclude {
	e["documentTitle"] = value
	return e
}

func (e Exclude) TimeCreated(value bool) Exclude {
	e["timeCreated"] = value
	return e
}

func (e Exclude) OwnerUUID(value bool) Exclude {
	e["ownerUUID"] = value
	return e
}

func (e Exclude) OwnerType(value bool) Exclude {
	e["ownerType"] = value
	return e
}

func (e Exclude) PdfBase64(value bool) Exclude {
	e["pdfBase64"] = value
	return e
}

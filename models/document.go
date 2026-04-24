package models

import (
	"time"

	"github.com/google/uuid"
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
	GetDocumentByDocumentUUID(document, owner uuid.UUID, excludes DocumentExcludes) (Document, error)
	GetDocumentByOwnerUUID(owner uuid.UUID, limit uint32, offset uint32, excludes DocumentExcludes) ([]Document, error)
	DeleteDocumentById(documentUuid, ownerUuid uuid.UUID) error
}

type DocumentExcludes map[string]bool

func (e DocumentExcludes) DocumentTitle(value bool) DocumentExcludes {
	e["documentTitle"] = value
	return e
}

func (e DocumentExcludes) TimeCreated(value bool) DocumentExcludes {
	e["timeCreated"] = value
	return e
}

func (e DocumentExcludes) OwnerUUID(value bool) DocumentExcludes {
	e["ownerUUID"] = value
	return e
}

func (e DocumentExcludes) OwnerType(value bool) DocumentExcludes {
	e["ownerType"] = value
	return e
}

func (e DocumentExcludes) PdfBase64(value bool) DocumentExcludes {
	e["pdfBase64"] = value
	return e
}

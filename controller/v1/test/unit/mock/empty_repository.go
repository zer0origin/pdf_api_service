package mock

import (
	"github.com/google/uuid"
	"pdf_service_api/domain"
)

type EmptyDocumentRepository struct {
}

func (e EmptyDocumentRepository) DeleteDocumentById(id uuid.UUID) error {
	panic("implement me")
}

func (e EmptyDocumentRepository) GetDocumentById(_ uuid.UUID) (domain.Document, error) {
	panic("implement me")
}

func (e EmptyDocumentRepository) UploadDocument(_ domain.Document) error {
	panic("implement me")
}

type EmptySelectionRepository struct {
}

func (e EmptySelectionRepository) GetSelectionsByDocumentId(uid uuid.UUID) ([]domain.Selection, error) {
	panic("implement me")
}

func (e EmptySelectionRepository) DeleteSelectionBySelectionUUID(uid uuid.UUID) error {
	panic("implement me")
}

func (e EmptySelectionRepository) AddNewSelection(selection domain.Selection) error {
	panic("implement me")
}

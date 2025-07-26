package mock

import (
	"github.com/google/uuid"
	"pdf_service_api/domain"
)

type EmptyDocumentRepository struct {
}

func (e EmptyDocumentRepository) DeleteDocumentById(_ uuid.UUID) error {
	panic("implement me")
}

func (e EmptyDocumentRepository) GetDocumentByDocumentUUID(_ uuid.UUID) (domain.Document, error) {
	panic("implement me")
}

func (e EmptyDocumentRepository) UploadDocument(_ domain.Document) error {
	panic("implement me")
}

type EmptySelectionRepository struct {
}

func (e EmptySelectionRepository) GetSelectionsByDocumentId(_ uuid.UUID) ([]domain.Selection, error) {
	panic("implement me")
}

func (e EmptySelectionRepository) DeleteSelectionBySelectionUUID(_ uuid.UUID) error {
	panic("implement me")
}

func (e EmptySelectionRepository) AddNewSelection(_ domain.Selection) error {
	panic("implement me")
}

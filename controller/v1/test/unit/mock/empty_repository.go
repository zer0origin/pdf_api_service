package mock

import (
	"github.com/google/uuid"
	"pdf_service_api/models"
)

type EmptyDocumentRepository struct {
}

func (e EmptyDocumentRepository) GetDocumentByOwnerUUID(id uuid.UUID, excludes map[string]bool) ([]models.Document, error) {
	panic("implement me")
}

func (e EmptyDocumentRepository) DeleteDocumentById(_ uuid.UUID) error {
	panic("implement me")
}

func (e EmptyDocumentRepository) GetDocumentByDocumentUUID(_ uuid.UUID, _ map[string]bool) (models.Document, error) {
	panic("implement me")
}

func (e EmptyDocumentRepository) UploadDocument(_ models.Document) error {
	panic("implement me")
}

type EmptySelectionRepository struct {
}

func (e EmptySelectionRepository) GetSelectionsByDocumentId(_ uuid.UUID) ([]models.Selection, error) {
	panic("implement me")
}

func (e EmptySelectionRepository) DeleteSelectionBySelectionUUID(_ uuid.UUID) error {
	panic("implement me")
}

func (e EmptySelectionRepository) AddNewSelection(_ models.Selection) error {
	panic("implement me")
}

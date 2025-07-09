package mock

import (
	"github.com/google/uuid"
	"pdf_service_api/models"
)

type EmptyRepository struct {
}

func (e EmptyRepository) GetDocumentById(_ uuid.UUID) (models.Document, error) {
	//TODO implement me
	panic("implement me")
}
func (e EmptyRepository) UploadDocument(_ models.Document) error {
	panic("implement me")
}

type MapRepository struct {
	Repo map[uuid.UUID]models.Document
}

func (m *MapRepository) UploadDocument(document models.Document) error {
	m.Repo[document.Uuid] = document
	return nil
}

func (m *MapRepository) GetDocumentById(uuid uuid.UUID) (models.Document, error) {
	return m.Repo[uuid], nil
}

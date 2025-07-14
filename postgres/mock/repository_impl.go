package mock

import (
	"github.com/google/uuid"
	"pdf_service_api/domain"
)

type EmptyRepository struct {
}

func (e EmptyRepository) DeleteDocumentById(id uuid.UUID) error {
	return nil
}
func (e EmptyRepository) GetDocumentById(_ uuid.UUID) (domain.Document, error) {
	//TODO implement me
	panic("implement me")
}
func (e EmptyRepository) UploadDocument(_ domain.Document) error {
	return nil
}

type MapRepository struct {
	Repo map[uuid.UUID]domain.Document
}

func (m *MapRepository) UploadDocument(document domain.Document) error {
	m.Repo[document.Uuid] = document
	return nil
}

func (m *MapRepository) GetDocumentById(uuid uuid.UUID) (domain.Document, error) {
	return m.Repo[uuid], nil
}

func (m *MapRepository) DeleteDocumentById(uuid uuid.UUID) error {
	delete(m.Repo, uuid)

	return nil
}

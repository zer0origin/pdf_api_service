package mock

import (
	"github.com/google/uuid"
	"pdf_service_api/domain"
)

type MapDocumentRepository struct {
	Repo map[uuid.UUID]domain.Document
}

func (m *MapDocumentRepository) UploadDocument(document domain.Document) error {
	m.Repo[document.Uuid] = document
	return nil
}

func (m *MapDocumentRepository) GetDocumentById(uuid uuid.UUID) (domain.Document, error) {
	return m.Repo[uuid], nil
}

func (m *MapDocumentRepository) DeleteDocumentById(uuid uuid.UUID) error {
	delete(m.Repo, uuid)

	return nil
}

type MapSelectionRepository struct {
	Repo map[uuid.UUID]domain.Selection
}

func (m *MapSelectionRepository) GetSelectionsByDocumentId(uid uuid.UUID) ([]domain.Selection, error) {
	var ss []domain.Selection

	for _, selection := range m.Repo {
		if selection.DocumentID != nil && *selection.DocumentID == uid {
			ss = append(ss, selection)
		}
	}

	return ss, nil
}

func (m *MapSelectionRepository) DeleteSelectionBySelectionUUID(uid uuid.UUID) error {
	delete(m.Repo, uid)
	return nil
}

func (m *MapSelectionRepository) AddNewSelection(selection domain.Selection) error {
	m.Repo[selection.Uuid] = selection
	return nil
}

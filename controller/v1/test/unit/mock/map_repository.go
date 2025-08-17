package mock

import (
	"github.com/google/uuid"
	"pdf_service_api/models"
)

type MapDocumentRepository struct {
	Repo map[uuid.UUID]models.Document
}

func (m *MapDocumentRepository) GetDocumentByOwnerUUID(uid uuid.UUID, excludes map[string]bool) ([]models.Document, error) {
	ss := make([]models.Document, 0)
	for _, selection := range m.Repo {
		if selection.OwnerUUID != nil && *selection.OwnerUUID == uid {
			ss = append(ss, selection)
		}
	}

	return ss, nil
}

func (m *MapDocumentRepository) UploadDocument(document models.Document) error {
	m.Repo[document.Uuid] = document
	return nil
}

func (m *MapDocumentRepository) GetDocumentByDocumentUUID(uuid uuid.UUID, exclude map[string]bool) (models.Document, error) {
	return m.Repo[uuid], nil
}

func (m *MapDocumentRepository) DeleteDocumentById(uuid uuid.UUID) error {
	delete(m.Repo, uuid)

	return nil
}

type MapSelectionRepository struct {
	Repo map[uuid.UUID]models.Selection
}

func (m *MapSelectionRepository) GetSelectionsBySelectionUUID(uid uuid.UUID) ([]models.Selection, error) {
	ss := make([]models.Selection, 0)
	ss = append(ss, m.Repo[uid])
	return ss, nil
}

func (m *MapSelectionRepository) DeleteSelectionByDocumentUUID(uid uuid.UUID) error {
	for selUuid, selection := range m.Repo {
		if selection.DocumentUUID != nil && *selection.DocumentUUID == uid {
			delete(m.Repo, selUuid)
		}
	}

	return nil
}

func (m *MapSelectionRepository) GetSelectionsByDocumentUUID(uid uuid.UUID) ([]models.Selection, error) {
	ss := make([]models.Selection, 0)
	for _, selection := range m.Repo {
		if selection.DocumentUUID != nil && *selection.DocumentUUID == uid {
			ss = append(ss, selection)
		}
	}

	return ss, nil
}

func (m *MapSelectionRepository) DeleteSelectionBySelectionUUID(uid uuid.UUID) error {
	delete(m.Repo, uid)
	return nil
}

func (m *MapSelectionRepository) AddNewSelection(selection models.Selection) error {
	m.Repo[selection.Uuid] = selection
	return nil
}

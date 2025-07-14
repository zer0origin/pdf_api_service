package domain

import "github.com/google/uuid"

type Selection struct {
	Uuid            uuid.UUID  `json:"selectionUUID"`
	DocumentID      *uuid.UUID `json:"documentID,omitempty"`
	IsComplete      bool       `json:"isComplete,omitempty"`
	Settings        *string    `json:"settings,omitempty"`
	SelectionBounds *string    `json:"selectionBounds,omitempty"`
}

type SelectionRepository interface {
	GetSelectionsByDocumentId(uid uuid.UUID) ([]Selection, error)
	DeleteSelectionBySelectionUUID(uid uuid.UUID) error
	AddNewSelection(selection Selection) error
}

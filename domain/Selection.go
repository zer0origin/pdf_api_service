package domain

import "github.com/google/uuid"

type Selection struct {
	Uuid            uuid.UUID  `json:"selectionUUID"`
	DocumentID      *uuid.UUID `json:"documentID,omitempty"`
	IsComplete      bool       `json:"isComplete,omitempty"`
	Settings        string     `json:"settings,omitempty"`
	SelectionBounds string     `json:"selection_bounds,omitempty"`
}

type SelectionRepository interface {
	GetSelectionBySelectionId(uid uuid.UUID) ([]Selection, error)
	DeleteSelectionBySelectionUUIDFunction(uid uuid.UUID) error
}

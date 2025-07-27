package models

import "github.com/google/uuid"

type Selection struct {
	Uuid            uuid.UUID                  `json:"selectionUUID"`
	DocumentUUID    *uuid.UUID                 `json:"documentUUID,omitempty"`
	IsComplete      bool                       `json:"isComplete,omitempty"`
	Settings        *string                    `json:"settings,omitempty"`
	SelectionBounds *map[int][]SelectionBounds `json:"selectionBounds,omitempty"`
}

type SelectionRepository interface {
	GetSelectionsByDocumentUUID(uid uuid.UUID) ([]Selection, error)
	GetSelectionsBySelectionUUID(uid uuid.UUID) ([]Selection, error)
	DeleteSelectionBySelectionUUID(uid uuid.UUID) error
	AddNewSelection(selection Selection) error
	DeleteSelectionByDocumentUUID(uid uuid.UUID) error
}

type SelectionBounds struct {
	SelectionMethod *string `json:"extract_method" example:"None"`
	X1              float64 `json:"x1" example:"43.122"`
	X2              float64 `json:"x2" example:"13"`
	Y1              float64 `json:"y1" example:"52.125"`
	Y2              float64 `json:"y2" example:"27.853"`
}

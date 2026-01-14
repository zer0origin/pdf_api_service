package models

import "github.com/google/uuid"

type Selection struct {
	Uuid         uuid.UUID    `json:"selectionUUID"`
	PageKey      *string      `json:"pageKey,omitempty"`
	DocumentUUID *uuid.UUID   `json:"documentUUID,omitempty"`
	Coordinates  *Coordinates `json:"coordinates,omitempty"`
}

type SelectionRepository interface {
	GetSelectionsByDocumentUUID(uid uuid.UUID) ([]Selection, error)
	GetSelectionsBySelectionUUID(uid uuid.UUID) ([]Selection, error)
	DeleteSelectionBySelectionUUID(uid uuid.UUID) error
	AddNewSelection(selection Selection) error
	DeleteSelectionByDocumentUUID(uid uuid.UUID) error
}

type Coordinates struct {
	X1 float64 `json:"x1" example:"43.122"`
	Y1 float64 `json:"y1" example:"52.125"`
	X2 float64 `json:"x2" example:"13"`
	Y2 float64 `json:"y2" example:"27.853"`
}

package models

import "github.com/google/uuid"

type Document struct {
	Uuid          uuid.UUID    `json:"documentUUID"`
	PdfBase64     *string      `json:"pdfBase64"`
	SelectionData *[]Selection `json:"selectionData,omitempty"`
}

type Selection struct {
	Uuid            uuid.UUID                 `json:"selectionUUID"`
	DocumentID      *uuid.UUID                `json:"documentID,omitempty"`
	IsComplete      bool                      `json:"isComplete,omitempty"`
	Settings        map[string]interface{}    `json:"settings,omitempty"`
	SelectionBounds map[int][]SelectionBounds `json:"selection_bounds,omitempty"`
}

type SelectionBounds struct {
	SelectionMethod *string `json:"extract_method"` //TODO: This functionality isn't even implemented yet.
	X1              float64 `json:"x1"`
	X2              float64 `json:"x2"`
	Y1              float64 `json:"y1"`
	Y2              float64 `json:"y2"`
}

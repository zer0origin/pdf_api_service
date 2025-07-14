package models

import "github.com/google/uuid"

type Document struct {
	Uuid          uuid.UUID    `json:"documentUUID"`
	PdfBase64     *string      `json:"pdfBase64"`
	SelectionData *[]Selection `json:"selectionData,omitempty"`
}

type Selection struct {
	Uuid            uuid.UUID              `json:"selectionUUID"`
	DocumentID      *uuid.UUID             `json:"documentID,omitempty"`
	IsComplete      bool                   `json:"isComplete,omitempty"`
	Settings        map[string]interface{} `json:"settings,omitempty"`
	SelectionBounds string                 `json:"selection_bounds,omitempty"`
}

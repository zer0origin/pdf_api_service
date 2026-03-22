package models

type Extraction = map[string]map[float32]TextData
type TextData struct {
	text                string
	textCoordinate      Coordinates
	SelectionCoordinate Coordinates
}

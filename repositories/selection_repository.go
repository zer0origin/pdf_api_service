package repositories

import (
	"database/sql"
	"github.com/google/uuid"
	"pdf_service_api/database"
	"pdf_service_api/models"
)

type SelectionRepository interface {
	GetSelectionBySelectionId(uid uuid.UUID) ([]models.Selection, error)
}

type selectionRepository struct {
	databaseManager database.ConfigForDatabase
}

func NewSelectionRepository() SelectionRepository {
	return selectionRepository{}
}

func (s selectionRepository) GetSelectionBySelectionId(uid uuid.UUID) ([]models.Selection, error) {
	var dataArr []models.Selection
	getSelection := getSelectionByDocumentUUIDFunction(uid, func(data []models.Selection) {
		dataArr = data
	})

	err := s.databaseManager.WithConnection(getSelection)
	if err != nil {
		return dataArr, err
	}

	return dataArr, nil
}

func getSelectionByDocumentUUIDFunction(uid uuid.UUID, callback func(data []models.Selection)) func(db *sql.DB) error {
	return func(db *sql.DB) error {
		sqlStatement := `SELECT "Selection_UUID", "Document_UUID", "Selection_bounds" FROM selection_table where "Document_UUID" = $1`

		rows, err := db.Query(sqlStatement, uid.String())
		if err != nil {
			return err

		}

		var dataArr []models.Selection
		for rows.Next() {
			data := models.Selection{}
			err := rows.Scan(&data.Uuid, &data.DocumentID, &data.SelectionBounds)
			if err != nil {
				return err
			}

			dataArr = append(dataArr, data)
		}

		callback(dataArr)
		return nil
	}
}

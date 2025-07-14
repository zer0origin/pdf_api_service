package repositories

import (
	"database/sql"
	"github.com/google/uuid"
	"pdf_service_api/database"
	"pdf_service_api/models"
)

type SelectionRepository interface {
	GetSelectionBySelectionId(uid uuid.UUID) ([]models.Selection, error)
	DeleteSelectionBySelectionUUIDFunction(uid uuid.UUID) error
}

type selectionRepository struct {
	databaseManager database.ConfigForDatabase
}

func NewSelectionRepository(db database.ConfigForDatabase) SelectionRepository {
	return selectionRepository{databaseManager: db}
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

func (s selectionRepository) DeleteSelectionBySelectionUUIDFunction(uid uuid.UUID) error {
	err := s.databaseManager.WithConnection(deleteSelectionBySelectionUUIDFunction(uid))

	if err != nil {
		return err
	}

	return nil
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

func deleteSelectionBySelectionUUIDFunction(uid uuid.UUID) func(db *sql.DB) error {
	return func(db *sql.DB) error {
		sqlStatement := `DELETE FROM selection_table WHERE "Selection_UUID" = $1`
		_, err := db.Exec(sqlStatement, uid)

		if err != nil {
			return err
		}

		return nil
	}
}

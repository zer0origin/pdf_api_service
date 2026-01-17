package postgres

import (
	"database/sql"
	"encoding/json"
	"errors"
	"pdf_service_api/models"

	"github.com/google/uuid"
)

type selectionRepository struct {
	databaseManager DatabaseHandler
}

func NewSelectionRepository(db DatabaseHandler) models.SelectionRepository {
	return selectionRepository{databaseManager: db}
}

func (s selectionRepository) AddNewSelection(selection models.Selection) error {
	err := s.databaseManager.WithConnection(AddNewSelectionFunction(selection))
	if err != nil {
		return err
	}

	return nil
}

func (s selectionRepository) GetSelectionsBySelectionUUID(uid uuid.UUID) ([]models.Selection, error) {
	var ss []models.Selection
	getSelection := getSelectionBySelectionUUIDFunction(uid, func(data []models.Selection) {
		ss = data
	})

	err := s.databaseManager.WithConnection(getSelection)
	if err != nil {
		return ss, err
	}

	return ss, nil
}

func (s selectionRepository) GetSelectionsByDocumentUUID(uid uuid.UUID) ([]models.Selection, error) {
	ss := make([]models.Selection, 0)
	getSelection := getSelectionByDocumentUUIDFunction(uid, func(data []models.Selection) {
		ss = data
	})

	err := s.databaseManager.WithConnection(getSelection)
	if err != nil {
		return ss, err
	}

	return ss, nil
}

func (s selectionRepository) DeleteSelectionByDocumentUUID(uid uuid.UUID) error {
	err := s.databaseManager.WithConnection(deleteSelectionByDocumentUUIDFunction(uid))
	if err != nil {
		return err
	}

	return nil
}

func (s selectionRepository) DeleteSelectionBySelectionUUID(uid uuid.UUID) error {
	err := s.databaseManager.WithConnection(deleteSelectionBySelectionUUIDFunction(uid))
	if err != nil {
		return err
	}

	return nil
}

func AddNewSelectionFunction(selection models.Selection) func(db *sql.DB) error {
	return func(db *sql.DB) error {
		sqlStatement := `insert into selection_table ("Selection_UUID", "Document_UUID", "Coordinates", "Page_Key") values ($1, $2, $3, $4);`

		pageKey := selection.PageKey
		selUid := selection.Uuid
		if selUid == uuid.Nil {
			return errors.New("selection uuid cannot be nil")
		}

		docUid := selection.DocumentUUID
		if *docUid == uuid.Nil {
			return errors.New("selection uuid cannot be nil")
		}

		bytes, err := json.Marshal(selection.Coordinates)
		if err != nil {
			return err
		}

		_, err = db.Exec(sqlStatement, selUid, docUid, bytes, pageKey)
		if err != nil {
			return err
		}

		return nil
	}
}

func getSelectionByDocumentUUIDFunction(uid uuid.UUID, callback func(data []models.Selection)) func(db *sql.DB) error {
	return func(db *sql.DB) error {
		sqlStatement := `SELECT "Selection_UUID", "Document_UUID", "Coordinates", "Page_Key" FROM selection_table where "Document_UUID" = $1`

		rows, err := db.Query(sqlStatement, uid.String())
		if err != nil {
			return err

		}

		//var ss []models.Selection
		ss := make([]models.Selection, 0)
		for rows.Next() {
			data := models.Selection{}
			err := rows.Scan(&data.Uuid, &data.DocumentUUID, &data.Coordinates, &data.PageKey)
			if err != nil {
				return err
			}

			ss = append(ss, data)
		}

		callback(ss)
		return nil
	}
}

func getSelectionBySelectionUUIDFunction(uid uuid.UUID, callback func(data []models.Selection)) func(db *sql.DB) error {
	return func(db *sql.DB) error {
		sqlStatement := `SELECT "Selection_UUID", "Document_UUID", "Coordinates", "Page_Key" FROM selection_table where "Selection_UUID" = $1`

		rows, err := db.Query(sqlStatement, uid.String())
		if err != nil {
			return err

		}

		var ss []models.Selection
		for rows.Next() {
			data := models.Selection{}
			err := rows.Scan(&data.Uuid, &data.DocumentUUID, &data.Coordinates, &data.PageKey)
			if err != nil {
				return err
			}

			ss = append(ss, data)
		}

		callback(ss)
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

func deleteSelectionByDocumentUUIDFunction(uid uuid.UUID) func(db *sql.DB) error {
	return func(db *sql.DB) error {
		sqlStatement := `DELETE FROM selection_table WHERE "Document_UUID" = $1`
		_, err := db.Exec(sqlStatement, uid)
		if err != nil {
			return err
		}

		return nil
	}
}

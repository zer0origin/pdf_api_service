package postgres

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"pdf_service_api/models"
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
		sqlStatement := `insert into selection_table ("Selection_UUID", "Document_UUID", "isCompleted", "Settings", "Selection_bounds") values ($1, $2, $3, $4, $5);`

		selUid := selection.Uuid
		if selUid == uuid.Nil {
			return errors.New("selection uuid cannot be nil")
		}

		docUid := selection.DocumentUUID
		if *docUid == uuid.Nil {
			return errors.New("selection uuid cannot be nil")
		}

		isComplete := selection.IsComplete
		settings := selection.Settings
		if settings == nil || *settings == "" {
			settings = func() *string { v := "{}"; return &v }()
		}

		selBounds := selection.SelectionBounds

		_, err := db.Exec(sqlStatement, selUid, docUid, isComplete, settings, selBounds)

		if err != nil {
			return err
		}

		return nil
	}
}

func getSelectionByDocumentUUIDFunction(uid uuid.UUID, callback func(data []models.Selection)) func(db *sql.DB) error {
	return func(db *sql.DB) error {
		sqlStatement := `SELECT "Selection_UUID", "Document_UUID", "Selection_bounds" FROM selection_table where "Document_UUID" = $1`

		rows, err := db.Query(sqlStatement, uid.String())
		if err != nil {
			return err

		}

		//var ss []models.Selection
		ss := make([]models.Selection, 0)
		for rows.Next() {
			data := models.Selection{}
			err := rows.Scan(&data.Uuid, &data.DocumentUUID, &data.SelectionBounds)
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
		sqlStatement := `SELECT "Selection_UUID", "Document_UUID", "Selection_bounds" FROM selection_table where "Selection_UUID" = $1`

		rows, err := db.Query(sqlStatement, uid.String())
		if err != nil {
			return err

		}

		var ss []models.Selection
		for rows.Next() {
			data := models.Selection{}
			err := rows.Scan(&data.Uuid, &data.DocumentUUID, &data.SelectionBounds)
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

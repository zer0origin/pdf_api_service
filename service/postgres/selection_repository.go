package postgres

import (
	"database/sql"
	"encoding/json"
	"errors"
	"pdf_service_api/models"

	"github.com/google/uuid"
	pg "github.com/lib/pq"
)

type selectionRepository struct {
	databaseManager DatabaseHandler
}

func (s selectionRepository) GetMapOfSelectionsBySelectionUUID(uid []uuid.UUID) (map[uuid.UUID]models.Selection, error) {
	ss := make(map[uuid.UUID]models.Selection)
	err := s.databaseManager.WithConnection(GetMapOfSelectionsBySelectionUUIDFunction(uid, func(data []models.Selection) {
		for _, sel := range data {
			ss[sel.Uuid] = sel
		}
	}))

	return ss, err
}

func NewSelectionRepository(db DatabaseHandler) models.SelectionRepository {
	return selectionRepository{databaseManager: db}
}

func (s selectionRepository) AddNewSelection(selection models.Selection) error {
	err := s.databaseManager.WithConnection(AddNewSelectionFunction(selection))
	return err
}

func (s selectionRepository) GetSelectionBySelectionUUID(uid uuid.UUID) (models.Selection, error) {
	var ss models.Selection
	getSelection := getSelectionBySelectionUUIDFunction(uid, func(data models.Selection) {
		ss = data
	})

	err := s.databaseManager.WithConnection(getSelection)
	return ss, err
}

func (s selectionRepository) GetSelectionListByDocumentUUID(uid uuid.UUID) ([]models.Selection, error) {
	ss := make([]models.Selection, 0)
	getSelection := getSelectionListByDocumentUUIDFunction(uid, func(data []models.Selection) {
		ss = data
	})

	err := s.databaseManager.WithConnection(getSelection)
	return ss, err
}

func (s selectionRepository) DeleteSelectionByDocumentUUID(uid uuid.UUID) error {
	err := s.databaseManager.WithConnection(deleteSelectionByDocumentUUIDFunction(uid))
	return err
}

func (s selectionRepository) DeleteSelectionBySelectionUUID(uid uuid.UUID) error {
	err := s.databaseManager.WithConnection(deleteSelectionBySelectionUUIDFunction(uid))
	return err
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

func getSelectionListByDocumentUUIDFunction(uid uuid.UUID, callback func(data []models.Selection)) func(db *sql.DB) error {
	return func(db *sql.DB) error {
		sqlStatement := `SELECT "Selection_UUID", "Document_UUID", "Coordinates", "Page_Key" FROM selection_table where "Document_UUID" = $1`

		rows, err := db.Query(sqlStatement, uid.String())
		if err != nil {
			return err

		}

		results := make([]models.Selection, 0)
		for rows.Next() {
			data := models.Selection{}

			var coordinateStr sql.NullString
			err := rows.Scan(&data.Uuid, &data.DocumentUUID, &coordinateStr, &data.PageKey)
			if err != nil {
				return err
			}

			if coordinateStr.Valid {
				coordinate := models.Coordinates{}
				err = json.Unmarshal([]byte(coordinateStr.String), &coordinate)
				if err != nil {
					return err
				}

				data.Coordinates = &coordinate
			}
			results = append(results, data)
		}

		callback(results)
		return nil
	}
}

func getSelectionBySelectionUUIDFunction(uid uuid.UUID, callback func(data models.Selection)) func(db *sql.DB) error {
	return func(db *sql.DB) error {
		sqlStatement := `SELECT "Selection_UUID", "Document_UUID", "Coordinates", "Page_Key" FROM selection_table where "Selection_UUID" = $1`

		rows := db.QueryRow(sqlStatement, uid.String())

		var results models.Selection
		data := models.Selection{}
		var coordinateStr sql.NullString
		err := rows.Scan(&data.Uuid, &data.DocumentUUID, &coordinateStr, &data.PageKey)
		if err != nil {
			return err
		}

		if coordinateStr.Valid {
			coordinate := models.Coordinates{}
			err = json.Unmarshal([]byte(coordinateStr.String), &coordinate)
			if err != nil {
				return err
			}

			data.Coordinates = &coordinate
		}

		results = data
		callback(results)
		return nil
	}
}

func deleteSelectionBySelectionUUIDFunction(uid uuid.UUID) func(db *sql.DB) error {
	return func(db *sql.DB) error {
		sqlStatement := `DELETE FROM selection_table WHERE "Selection_UUID" = $1`
		_, err := db.Exec(sqlStatement, uid)
		return err
	}
}

func deleteSelectionByDocumentUUIDFunction(uid uuid.UUID) func(db *sql.DB) error {
	return func(db *sql.DB) error {
		sqlStatement := `DELETE FROM selection_table WHERE "Document_UUID" = $1`
		_, err := db.Exec(sqlStatement, uid)
		return err
	}
}

func GetMapOfSelectionsBySelectionUUIDFunction(uids []uuid.UUID, callback func(data []models.Selection)) func(db *sql.DB) error {
	return func(db *sql.DB) error {
		sqlStatement := `SELECT "Selection_UUID", "Document_UUID", "Coordinates", "Page_Key" FROM selection_table WHERE "Selection_UUID" = any ($1)`

		rows, err := db.Query(sqlStatement, pg.Array(uids))
		if err != nil {
			return err
		}

		results := make([]models.Selection, 0)
		for rows.Next() {
			data := models.Selection{}
			var coordinateStr sql.NullString
			err := rows.Scan(&data.Uuid, &data.DocumentUUID, &coordinateStr, &data.PageKey)
			if err != nil {
				return err
			}

			if coordinateStr.Valid {
				coordinate := models.Coordinates{}
				err = json.Unmarshal([]byte(coordinateStr.String), &coordinate)
				if err != nil {
					return err
				}

				data.Coordinates = &coordinate
			}

			results = append(results, data)
		}

		callback(results)
		return err
	}
}

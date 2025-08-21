package postgres

import (
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"pdf_service_api/models"
)

type metaRepository struct {
	DatabaseHandler DatabaseHandler
}

func NewMetaRepository(db DatabaseHandler) models.MetaRepository {
	return metaRepository{DatabaseHandler: db}
}

func (m metaRepository) AddMeta(data models.Meta) error {
	if err := m.DatabaseHandler.WithConnection(addMetaDataFunction(data)); err != nil {
		return err
	}

	return nil
}

func (m metaRepository) DeleteMeta(data models.Meta) error {
	if err := m.DatabaseHandler.WithConnection(removeMetaDataFunction(data)); err != nil {
		return err
	}

	return nil
}

func (m metaRepository) UpdateMeta(uid uuid.UUID, data models.Meta) error {
	if err := m.DatabaseHandler.WithConnection(updateMetaDataFunction(uid, data)); err != nil {
		return err
	}

	return nil
}

func (m metaRepository) GetMeta(uid uuid.UUID) (models.Meta, error) {
	returnedData := &models.Meta{}
	callbackFunction := func(data models.Meta) error {
		*returnedData = data
		return nil
	}

	if err := m.DatabaseHandler.WithConnection(getMetaDataFunction(uid, callbackFunction)); err != nil {
		return models.Meta{}, err
	}

	return *returnedData, nil
}

func addMetaDataFunction(data models.Meta) func(db *sql.DB) error {
	return func(db *sql.DB) error {
		SqlStatement := `INSERT INTO documentmeta_table ("Document_UUID", "Number_Of_Pages", "Height", "Width", "Images") values ($1, $2, $3, $4, $5)`
		if _, err := db.Exec(SqlStatement, data.DocumentUUID, data.NumberOfPages, data.Height, data.Width, data.Images); err != nil {
			return err
		}

		return nil
	}
}

func removeMetaDataFunction(data models.Meta) func(db *sql.DB) error {
	return func(db *sql.DB) error {
		SqlStatement := `DELETE FROM documentmeta_table WHERE "Document_UUID" = $1`
		if _, err := db.Exec(SqlStatement, data.DocumentUUID); err != nil {
			return err
		}

		return nil
	}
}

func updateMetaDataFunction(uid uuid.UUID, data models.Meta) func(db *sql.DB) error {
	return func(db *sql.DB) error {
		SqlStatement := `UPDATE documentmeta_table SET "Number_Of_Pages" = COALESCE($1, "Number_Of_Pages"), "Height" = COALESCE($2, "Height"), "Width" = COALESCE($3, "Width"), "Images" = COALESCE($4, "Images") where "Document_UUID" = $5`
		bytes, err := json.Marshal(data.Images)
		if err != nil {
			return err
		}

		if _, err := db.Exec(SqlStatement, data.NumberOfPages, data.Height, data.Width, string(bytes), uid); err != nil {
			return err
		}

		return nil
	}
}

func getMetaDataFunction(uid uuid.UUID, callback func(data models.Meta) error) func(db *sql.DB) error {
	return func(db *sql.DB) error {
		meta := &models.Meta{}
		SqlStatement := `SELECT "Document_UUID", "Number_Of_Pages", "Height", "Width", "Images" FROM documentmeta_table where "Document_UUID" = $1`

		row := db.QueryRow(SqlStatement, uid)
		err := row.Scan(&meta.DocumentUUID, &meta.NumberOfPages, &meta.Height, &meta.Width, &meta.Images)
		if err != nil {
			return err
		}

		return callback(*meta)
	}
}

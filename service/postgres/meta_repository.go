package postgres

import (
	"database/sql"
	"encoding/json"
	"pdf_service_api/models"

	"github.com/google/uuid"
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

func (m metaRepository) GetMeta(documentUid, ownerUid uuid.UUID) (models.Meta, error) {
	returnedData := &models.Meta{}
	callbackFunction := func(data models.Meta) error {
		*returnedData = data
		return nil
	}

	if err := m.DatabaseHandler.WithConnection(getMetaDataFunction(documentUid, ownerUid, callbackFunction)); err != nil {
		return models.Meta{}, err
	}

	return *returnedData, nil
}

func (m metaRepository) GetMetaPagination(documentUid, ownerUid uuid.UUID, start, end uint16) (models.Meta, error) {
	returnedData := &models.Meta{}
	callbackFunction := func(data models.Meta) error {
		*returnedData = data
		return nil
	}

	if err := m.DatabaseHandler.WithConnection(getMetaDataPaginationFunction(documentUid, ownerUid, start, end, callbackFunction)); err != nil {
		return models.Meta{}, err
	}

	return *returnedData, nil
}

func addMetaDataFunction(data models.Meta) func(db *sql.DB) error {
	return func(db *sql.DB) error {
		SqlStatement := `INSERT INTO documentmeta_table ("Document_UUID", "Number_Of_Pages", "Height", "Width", "Images") values ($1, $2, $3, $4, $5)`

		imagesJson, err := json.Marshal(data.Images)
		if err != nil {
			return err
		}

		if _, err := db.Exec(SqlStatement, data.DocumentUUID, data.NumberOfPages, data.Height, data.Width, imagesJson); err != nil {
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

func getMetaDataFunction(documentUid, ownerUid uuid.UUID, callback func(data models.Meta) error) func(db *sql.DB) error {
	return func(db *sql.DB) error {
		meta := &models.Meta{}
		SqlStatement := `SELECT mt."Document_UUID", mt."Number_Of_Pages", mt."Height", mt."Width", mt."Images" FROM documentmeta_table as mt join public.document_table dt on dt."Document_UUID" = mt."Document_UUID" where mt."Document_UUID" = $1 and dt."Owner_UUID" = $2`

		row := db.QueryRow(SqlStatement, documentUid, ownerUid)

		var imageStr string
		err := row.Scan(&meta.DocumentUUID, &meta.NumberOfPages, &meta.Height, &meta.Width, &imageStr)
		if err != nil {
			return err
		}

		imageMap := make(map[uint32]string)
		err = json.Unmarshal([]byte(imageStr), &imageMap)
		if err != nil {
			return err
		}
		meta.Images = &imageMap

		return callback(*meta)
	}
}

func getMetaDataPaginationFunction(documentUid, ownerUid uuid.UUID, start, end uint16, callback func(data models.Meta) error) func(db *sql.DB) error {
	return func(db *sql.DB) error {
		meta := &models.Meta{}
		SqlStatement := `select mt."Document_UUID",
       mt."Number_Of_Pages",
       mt."Height",
       mt."Width",
       Pagination_Images
from documentmeta_table as mt
    inner join public.document_table dt on dt."Document_UUID" = mt."Document_UUID"
         cross join (select coalesce(jsonb_object_agg(j.key, j.value), '{}') AS Pagination_Images
                     from documentmeta_table mt,
                          jsonb_each(mt."Images"::jsonb) j
                     where key::int BETWEEN $3 and $4
                       and mt."Document_UUID" = $1)
where mt."Document_UUID" = $1 and dt."Owner_UUID" = $2;`
		row := db.QueryRow(SqlStatement, documentUid, ownerUid, start, end)

		var imageStr string
		err := row.Scan(&meta.DocumentUUID, &meta.NumberOfPages, &meta.Height, &meta.Width, &imageStr)
		if err != nil {
			return err
		}

		imageMap := make(map[uint32]string)
		err = json.Unmarshal([]byte(imageStr), &imageMap)
		if err != nil {
			return err
		}
		meta.Images = &imageMap

		return callback(*meta)
	}
}

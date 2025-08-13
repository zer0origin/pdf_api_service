package postgres

import (
	"database/sql"
	"github.com/google/uuid"
	"pdf_service_api/models"
)

type documentRepository struct {
	databaseManager DatabaseHandler
}

func NewDocumentRepository(databaseManager DatabaseHandler) models.DocumentRepository {
	return documentRepository{databaseManager: databaseManager}
}

func (d documentRepository) DeleteDocumentById(uuid uuid.UUID) error {
	err := d.databaseManager.WithConnection(deleteDocumentSqlDatabase(uuid))
	if err != nil {
		return err
	}

	return nil
}

func (d documentRepository) GetDocumentByOwnerUUID(uuid uuid.UUID) ([]models.Document, error) {
	ss := make([]models.Document, 0)
	err := d.databaseManager.WithConnection(getDocumentByOwnerUUIDFunction(uuid, func(data []models.Document) {
		ss = data
	}))
	if err != nil {
		return ss, err
	}

	return ss, nil
}

func (d documentRepository) GetDocumentByDocumentUUID(uuid uuid.UUID) (models.Document, error) {
	document := &models.Document{}
	err := d.databaseManager.WithConnection(getDocumentByDocumentUUIDFunction(uuid, func(data models.Document) {
		*document = data
	}))

	if err != nil {
		return models.Document{}, err
	}

	return *document, nil
}

func (d documentRepository) UploadDocument(document models.Document) error {
	uploadDocumentSQL := createUploadDocumentSqlDatabase(&document) //create callback
	err := d.databaseManager.WithConnection(uploadDocumentSQL)
	if err != nil {
		return err
	}

	return nil
}

func getDocumentByDocumentUUIDFunction(uid uuid.UUID, callback func(data models.Document)) func(db *sql.DB) error {
	return func(db *sql.DB) error {
		sqlStatement := `SELECT "Document_UUID", "Document_Title", "Document_Base64", "Time_Created", "Owner_UUID", "Owner_Type" FROM document_table WHERE "Document_UUID" = $1`
		rows := db.QueryRow(sqlStatement, uid)
		if rows.Err() != nil {
			return rows.Err()
		}

		document := models.Document{}
		err := rows.Scan(&document.Uuid, &document.DocumentTitle, &document.PdfBase64, &document.TimeCreated, &document.OwnerUUID, &document.OwnerType)
		if err != nil {
			return err
		}

		callback(document)
		return nil
	}
}

func getDocumentByOwnerUUIDFunction(uid uuid.UUID, callback func(data []models.Document)) func(db *sql.DB) error {
	return func(db *sql.DB) error {
		sqlStatement := `SELECT "Document_UUID", "Document_Title", "Document_Base64", "Time_Created", "Owner_UUID", "Owner_Type" FROM document_table WHERE "Owner_UUID" = $1 order by "Time_Created"`
		rows, err := db.Query(sqlStatement, uid)
		if err != nil {
			return rows.Err()
		}

		dd := make([]models.Document, 0)

		for rows.Next() {
			document := models.Document{}
			err := rows.Scan(&document.Uuid, &document.DocumentTitle, &document.PdfBase64, &document.TimeCreated, &document.OwnerUUID, &document.OwnerType)
			if err != nil {
				return err
			}

			dd = append(dd, document)
		}

		callback(dd)
		return nil
	}
}

func createUploadDocumentSqlDatabase(document *models.Document) func(db *sql.DB) error {
	return func(db *sql.DB) error {
		sqlStatement := `insert into document_table("Document_UUID", "Document_Title", "Document_Base64") values ($1, $2, $3) returning "Document_UUID"`
		_, err := db.Exec(sqlStatement, document.Uuid, document.DocumentTitle, document.PdfBase64)

		if err != nil {
			return err
		}

		return nil
	}
}

func deleteDocumentSqlDatabase(uuid uuid.UUID) func(db *sql.DB) error {
	return func(db *sql.DB) error {
		sqlStatement := `DELETE FROM document_table where "Document_UUID" = $1`
		_, err := db.Exec(sqlStatement, uuid)

		if err != nil {
			return err
		}

		return nil
	}
}

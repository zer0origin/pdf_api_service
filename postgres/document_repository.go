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

func (d documentRepository) GetDocumentByDocumentUUID(uid uuid.UUID) (models.Document, error) {
	document := &models.Document{}
	err := d.databaseManager.WithConnection(getDocumentByUUIDFunction(uid, func(data models.Document) {
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

func getDocumentByUUIDFunction(uid uuid.UUID, callback func(data models.Document)) func(db *sql.DB) error {
	return func(db *sql.DB) error {
		sqlStatement := `SELECT "Document_UUID", "Document_Base64" FROM document_table WHERE "Document_UUID" = $1`
		rows := db.QueryRow(sqlStatement, uid)
		if rows.Err() != nil {
			return rows.Err()
		}

		document := &models.Document{}
		err := rows.Scan(&document.Uuid, &document.PdfBase64)
		if err != nil {
			return err
		}

		callback(*document)
		return nil
	}
}

func createUploadDocumentSqlDatabase(document *models.Document) func(db *sql.DB) error {
	return func(db *sql.DB) error {
		sqlStatement := `insert into document_table values ($1, $2) returning "Document_UUID"`
		_, err := db.Exec(sqlStatement, document.Uuid, document.PdfBase64)

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

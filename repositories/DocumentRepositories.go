package repositories

import (
	"database/sql"
	"github.com/google/uuid"
	"pdf_service_api/database"
	"pdf_service_api/models"
)

type DocumentRepository interface {
	UploadDocument(document models.Document) error
	GetDocumentById(id uuid.UUID) (models.Document, error)
}

type documentRepository struct {
	databaseManager database.ConfigForDatabase
}

func NewDocumentRepository(databaseManager database.ConfigForDatabase) DocumentRepository {
	return documentRepository{databaseManager: databaseManager}
}

func (d documentRepository) GetDocumentById(uid uuid.UUID) (models.Document, error) {
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

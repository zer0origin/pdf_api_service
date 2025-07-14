package postgres

import (
	"database/sql"
	"github.com/google/uuid"
	"pdf_service_api/domain"
)

type documentRepository struct {
	databaseManager ConfigForDatabase
}

func NewDocumentRepository(databaseManager ConfigForDatabase) domain.DocumentRepository {
	return documentRepository{databaseManager: databaseManager}
}

func (d documentRepository) DeleteDocumentById(uuid uuid.UUID) error {
	err := d.databaseManager.WithConnection(deleteDocumentSqlDatabase(uuid))
	if err != nil {
		return err
	}

	return nil
}

func (d documentRepository) GetDocumentById(uid uuid.UUID) (domain.Document, error) {
	document := &domain.Document{}
	err := d.databaseManager.WithConnection(getDocumentByUUIDFunction(uid, func(data domain.Document) {
		*document = data
	}))

	if err != nil {
		return domain.Document{}, err
	}

	return *document, nil
}

func (d documentRepository) UploadDocument(document domain.Document) error {
	uploadDocumentSQL := createUploadDocumentSqlDatabase(&document) //create callback
	err := d.databaseManager.WithConnection(uploadDocumentSQL)
	if err != nil {
		return err
	}

	return nil
}

func getDocumentByUUIDFunction(uid uuid.UUID, callback func(data domain.Document)) func(db *sql.DB) error {
	return func(db *sql.DB) error {
		sqlStatement := `SELECT "Document_UUID", "Document_Base64" FROM document_table WHERE "Document_UUID" = $1`
		rows := db.QueryRow(sqlStatement, uid)
		if rows.Err() != nil {
			return rows.Err()
		}

		document := &domain.Document{}
		err := rows.Scan(&document.Uuid, &document.PdfBase64)
		if err != nil {
			return err
		}

		callback(*document)
		return nil
	}
}

func createUploadDocumentSqlDatabase(document *domain.Document) func(db *sql.DB) error {
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

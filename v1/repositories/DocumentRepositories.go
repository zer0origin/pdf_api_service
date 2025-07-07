package repositories

import (
	"database/sql"
	"github.com/google/uuid"
	"pdf_service_api/database"
	"pdf_service_api/v1/models"
)

type DocumentRepository interface {
	UploadDocument(document models.Document) (uuid.UUID, error)
	GetDocumentById(id uuid.UUID) (models.Document, error)
}

type documentRepository struct {
}

func NewDocumentRepository() DocumentRepository {
	return documentRepository{}
}

func (d documentRepository) GetDocumentById(uid uuid.UUID) (models.Document, error) {
	document := &models.Document{}
	err := database.WithConnection(getDocumentByUUIDFunction(uid, func(data models.Document) {
		*document = data
	}))

	if err != nil {
		return models.Document{}, err
	}

	return *document, nil
}

func (d documentRepository) UploadDocument(document models.Document) (uuid.UUID, error) {
	u := uuid.New()
	uploadDocumentSQL := createUploadDocumentSqlDatabase(&document, u) //create callback
	err := database.WithConnection(uploadDocumentSQL)
	if err != nil {
		return uuid.Nil, err
	}

	return u, nil
}

// SQL Query to database | TODO: MOVE!
func getDocumentByUUIDFunction(uid uuid.UUID, callback func(data models.Document)) func(db *sql.DB) error {
	return func(db *sql.DB) error {
		sqlStatement := "SELECT document_id FROM documents WHERE document_id = $1"

		rows := db.QueryRow(sqlStatement, uid)

		if rows.Err() != nil {
			return rows.Err()
		}

		document := &models.Document{}
		err := rows.Scan(document)
		if err != nil {
			return err
		}

		callback(*document)
		return nil
	}
}

// SQL Query to database | TODO: MOVE!
func createUploadDocumentSqlDatabase(document *models.Document, insertedUUID uuid.UUID) func(db *sql.DB) error {
	return func(db *sql.DB) error {
		sqlStatement := `insert into document_table values ($1, $2) returning "Document_UUID"`
		_, err := db.Exec(sqlStatement, insertedUUID, document.PdfBase64, insertedUUID)

		if err != nil {
			return err
		}

		return nil
	}
}

package postgres

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"pdf_service_api/models"
	"text/template"
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

func (d documentRepository) GetDocumentByDocumentUUID(uuid uuid.UUID, excludes map[string]bool) (models.Document, error) {
	document := &models.Document{}
	err := d.databaseManager.WithConnection(getDocumentByDocumentUUIDFunction(uuid, excludes, func(data models.Document) {
		*document = data
	}))

	if err != nil {
		return models.Document{}, err
	}

	return *document, nil
}

func (d documentRepository) UploadDocument(document models.Document) error {
	uploadDocumentSQL := createDocumentFunction(&document) //create callback
	err := d.databaseManager.WithConnection(uploadDocumentSQL)
	if err != nil {
		return err
	}

	return nil
}

func getDocumentByDocumentUUIDFunction(uid uuid.UUID, excludes map[string]bool, callback func(data models.Document)) func(db *sql.DB) error {
	return func(db *sql.DB) error {
		sqlStatement := `SELECT {{if .documentTitle }}{{else}}"Document_Title", {{end}}{{if .pdfBase64 }}{{else}}"Document_Base64", {{end}}{{if .timeCreated }}{{else}}"Time_Created", {{end}}{{if .ownerUUID }}{{else}}"Owner_UUID", {{end}}{{if .ownerType }}{{else}}"Owner_Type",{{end}} "Document_UUID" FROM document_table WHERE "Document_UUID" = $1`

		templ, err := template.New("documentQuery").Parse(sqlStatement)
		var buffer bytes.Buffer
		err = templ.Execute(&buffer, excludes)
		if err != nil {
			fmt.Println(err)
		}

		generatedSQL := buffer.String()
		fmt.Println("Generated SQL:")
		fmt.Println(generatedSQL)

		rows := db.QueryRow(generatedSQL, uid)
		if rows.Err() != nil {
			return rows.Err()
		}

		document := models.Document{}

		scanDestinations := make([]any, 0)
		if !excludes["documentTitle"] {
			scanDestinations = append(scanDestinations, &document.DocumentTitle)
		}

		if !excludes["pdfBase64"] {
			scanDestinations = append(scanDestinations, &document.PdfBase64)
		}

		if !excludes["timeCreated"] {
			scanDestinations = append(scanDestinations, &document.TimeCreated)
		}

		if !excludes["ownerUUID"] {
			scanDestinations = append(scanDestinations, &document.OwnerUUID)
		}

		if !excludes["ownerType"] {
			scanDestinations = append(scanDestinations, &document.OwnerType)
		}
		scanDestinations = append(scanDestinations, &document.Uuid)

		err = rows.Scan(scanDestinations...)
		if err != nil {
			return err
		}

		callback(document)
		return nil
	}
}

func getDocumentByOwnerUUIDFunction(uid uuid.UUID, callback func(data []models.Document)) func(db *sql.DB) error {
	return func(db *sql.DB) error {
		sqlStatement := `SELECT "Document_UUID", "Document_Title", "Document_Base64", "Time_Created", "Owner_UUID", "Owner_Type" FROM document_table WHERE "Owner_UUID" = $1 order by "Time_Created" DESC`
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

func createDocumentFunction(document *models.Document) func(db *sql.DB) error {
	return func(db *sql.DB) error {
		sqlStatement := `insert into document_table("Document_UUID", "Document_Title", "Document_Base64", "Owner_UUID", "Owner_Type") values ($1, $2, $3, $4, $5) returning "Document_UUID"`
		_, err := db.Exec(sqlStatement, document.Uuid, document.DocumentTitle, document.PdfBase64, document.OwnerUUID, document.OwnerType)

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

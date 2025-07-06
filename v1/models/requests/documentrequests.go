package requests

import "pdf_service_api/v1/models"

type GetDocumentRequest struct {
	Document models.Document `json:"document"`
}

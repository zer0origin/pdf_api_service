package requests

import "github.com/google/uuid"

type GetDocumentRequest struct {
	DocumentUuid uuid.UUID `json:"document_uuid"`
}

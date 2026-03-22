package v1

type ExtractionController struct {
}

// extraction handles the HTTP POST request to take selections and retrieve the text.
//
// @Summary
// @Description
// @Tags extraction
// @Accept
// @Produce
// @Param   request body v1.DeleteMetaRequest true "Metadata deletion request"
// @Success 200 {object} map[string]bool "Successful deletion"
// @Failure 400 "Bad request, typically due to invalid input"
// @Failure 500 "Internal server error, typically due to database issues"
// @Router /meta [delete]
func (t SelectionController) extractAsText() {

}

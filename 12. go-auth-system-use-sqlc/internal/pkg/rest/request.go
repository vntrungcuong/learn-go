package rest

import (
	"encoding/json"
	"net/http"

	"go-auth-system/internal/pkg/validator"
)

func ParseAndValidate[T any](w http.ResponseWriter, r *http.Request, req *T) bool {
	// 1. Decode JSON
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		Error(w, r, http.StatusBadRequest, "Invalid JSON format", nil)
		return false
	}

	// 2. Validate tags
	if err := validator.ValidateStruct(req); err != nil {
		rawErrs := validator.MapValidationError(err)

		// Map from validator.ValidationError to response.ErrorItem
		var apiErrs []ErrorItem
		for _, e := range rawErrs {
			apiErrs = append(apiErrs, ErrorItem{
				Field:   e.Field,
				Message: e.Message,
				Code:    e.Tag,
			})
		}

		Error(w, r, http.StatusUnprocessableEntity, "Validation failed", apiErrs)
		return false
	}

	return true
}

package validation

import (
	"event-broker-document-api/helper"

	validation "github.com/go-ozzo/ozzo-validation"
)

func Required(value string, field string) {

	err := validation.Validate(value,
		validation.Required,
	)

	if err != nil {
		helper.ResponseError(400, "field '"+field+"' error: "+err.Error(), 1025)
	}
}

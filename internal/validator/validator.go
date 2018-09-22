package validator

import (
	"github.com/codetaming/indy-ingest/internal/model"
)

type Validator interface {
	Validate(schemaUrl string, bodyJson string) (model.ValidationResult, error)
}

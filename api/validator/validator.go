package validator

import "github.com/xeipuuv/gojsonschema"

type ValidationResult struct {
	Valid   bool
	Message string
	Errors  []string
}

func Validate(schemaUrl string, bodyJson string) (ValidationResult) {

	schemaLoader := gojsonschema.NewReferenceLoader(schemaUrl)
	documentLoader := gojsonschema.NewStringLoader(bodyJson)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)

	if err != nil {
		panic(err.Error())
	}

	var message string
	var errors []string

	if result.Valid() {
		message = "The document is valid"
	} else {
		message = "The document is not valid"
		for _, desc := range result.Errors() {
			errors = append(errors, desc.Description())
		}
	}

	vr := ValidationResult{
		Valid:   result.Valid(),
		Message: message,
		Errors:  errors,
	}

	return vr
}

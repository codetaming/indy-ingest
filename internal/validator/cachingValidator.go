package validator

import (
	"errors"
	"github.com/codetaming/indy-ingest/internal/model"
	"github.com/codetaming/indy-ingest/internal/validator/cache_provider"
	"github.com/xeipuuv/gojsonschema"
	"io/ioutil"
	"log"
	"net/http"
)

type CachingValidator struct {
	logger        *log.Logger
	cacheProvider cache_provider.CacheProvider
}

func (v *CachingValidator) SaveCache() {
	v.cacheProvider.SaveCache()
}

func (v *CachingValidator) loadFromHTTP(address string) (string, error) {
	resp, err := http.Get(address)
	if err != nil {
		v.logger.Printf("error retreiving %s: %s", address, err.Error())
		return "", err
	}

	// must return HTTP Status 200 OK
	if resp.StatusCode != http.StatusOK {
		v.logger.Printf("error retreiving %s: %d", address, resp.StatusCode)
		return "", err
	}

	bodyBuff, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		v.logger.Printf("error reading body of %s: %s", address, err.Error())
		return "", err
	}

	return string(bodyBuff), nil
}

func (v *CachingValidator) Validate(schemaUrl string, bodyJson string) (model.ValidationResult, error) {
	if schemaUrl == "" {
		return model.ValidationResult{}, errors.New("schema URL cannot be empty")
	}

	if bodyJson == "" {
		return model.ValidationResult{}, errors.New("JSON cannot be empty")
	}

	if !v.cacheProvider.ExistsInCache(schemaUrl) {
		v.logger.Printf("loading %s from internet", schemaUrl)
		schemaJson, err := v.loadFromHTTP(schemaUrl)
		if err != nil {
			v.logger.Printf("error retreiving schema: %s", err.Error())
		}
		v.cacheProvider.AddToCache(schemaUrl, schemaJson)
	}
	sl := gojsonschema.NewStringLoader(v.cacheProvider.RetrieveFromCache(schemaUrl))
	dl := gojsonschema.NewStringLoader(bodyJson)
	result, err := gojsonschema.Validate(sl, dl)

	if err != nil {
		return model.ValidationResult{}, err
	}

	var message string
	var validationErrors []string

	if result.Valid() {
		message = "The document is valid"
	} else {
		message = "The document is not valid"
		for _, desc := range result.Errors() {
			validationErrors = append(validationErrors, desc.Description())
		}
	}

	vr := model.ValidationResult{
		Valid:   result.Valid(),
		Message: message,
		Errors:  validationErrors,
	}

	return vr, nil
}

func NewCachingValidator(logger *log.Logger, schemaCacheFile string) (Validator, error) {
	cp, err := cache_provider.NewLocalCacheProvider(logger, schemaCacheFile)
	c := &CachingValidator{
		logger:        logger,
		cacheProvider: cp,
	}
	return c, err
}

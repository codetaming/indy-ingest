package validator

import (
	"encoding/json"
	"errors"
	"github.com/codetaming/indy-ingest/internal/model"
	"github.com/xeipuuv/gojsonschema"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type CachingValidator struct {
	logger          *log.Logger
	schemaCacheFile string
	schemaCache     map[string]string
	loadedCacheSize int
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

	_, exists := v.schemaCache[schemaUrl]
	if exists {
		v.logger.Printf("loading %s from cache", schemaUrl)
	} else {
		v.logger.Printf("loading %s from internet", schemaUrl)
		schemaJson, err := v.loadFromHTTP(schemaUrl)
		if err != nil {
			v.logger.Printf("error retreiving schema: %s", err.Error())
		}
		v.schemaCache[schemaUrl] = schemaJson
		v.logger.Printf("added %s to cache", schemaUrl)
		v.saveCache()
	}
	sl := gojsonschema.NewStringLoader(string(v.schemaCache[schemaUrl]))
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

type Schemas struct {
	Schemas []Schema `json:"schemas"`
}

type Schema struct {
	Id      string          `json:"id"`
	Content json.RawMessage `json:"content"`
}

func NewCachingValidator(logger *log.Logger, schemaCacheFile string) (Validator, error) {
	c := &CachingValidator{
		logger:          logger,
		schemaCacheFile: schemaCacheFile,
		schemaCache:     make(map[string]string),
	}
	c.loadCache()
	return c, nil
}

func (v *CachingValidator) saveCache() {
	//only save cache if it has gained new entries
	if len(v.schemaCache) > v.loadedCacheSize {
		var schemas Schemas
		for k, v := range v.schemaCache {
			{
				schema := Schema{
					Id:      k,
					Content: json.RawMessage(v),
				}
				schemas.Schemas = append(schemas.Schemas, schema)
			}
		}
		j, err := json.MarshalIndent(schemas, "", "  ")
		if err != nil {
			v.logger.Printf("error marshalling cache: %s", err.Error())
		}
		v.logger.Printf("saving cache")
		err = ioutil.WriteFile(v.schemaCacheFile, j, 0644)
		if err != nil {
			v.logger.Printf("error saving cache: %s", err.Error())
		}
	} else {
		v.logger.Printf("cache unchanged, not saving")
	}
}

func (v *CachingValidator) loadCache() {
	v.schemaCache = make(map[string]string)
	jsonFile, err := os.Open(v.schemaCacheFile)
	if err != nil {
		v.logger.Printf("error loading cache file %s", v.schemaCacheFile)
	} else {
		var schemas Schemas
		byteValue, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			v.logger.Printf("error reading schema file: %s", err.Error())
		}
		err = json.Unmarshal(byteValue, &schemas)
		if err != nil {
			v.logger.Printf("error unmarshalling json: %s", err.Error())
		}
		for _, s := range schemas.Schemas {
			if s.Id != "" {
				content, err := json.Marshal(s.Content)
				if err != nil {
					v.logger.Printf("error marshalling content: %s", err.Error())
				}
				v.schemaCache[s.Id] = string(content)
				v.loadedCacheSize++
				v.logger.Printf("loaded schema %s into cache", s.Id)
			} else {
				v.logger.Printf("failed to map from %s", v.schemaCacheFile)
			}
		}
	}
	defer jsonFile.Close()
}

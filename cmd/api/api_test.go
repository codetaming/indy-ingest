package api

import (
	"github.com/codetaming/indy-ingest/internal/persistence/local"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var a *API
var logger *log.Logger

type testDef struct {
	name                   string
	in                     *http.Request
	out                    *httptest.ResponseRecorder
	expectedLocationHeader string
	expectedStatus         int
	expectedBody           string
}

func TestHandlers_Validate(t *testing.T) {
	tests := []testDef{
		{
			name:           "Validate null body",
			in:             requestWithValidationHeaders("../../data/examples/null.json"),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "JSON cannot be empty\n",
		},
		{
			name:           "Validate empty body",
			in:             requestWithValidationHeaders("../../data/examples/empty.json"),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusOK,
			expectedBody:   "{\"valid\":false,\"message\":\"The document is not valid\",\"errors\":[\"describedBy is required\",\"schema_type is required\",\"biomaterial_core is required\",\"organ is required\"]}\n",
		},
		{
			name:           "Validate with valid JSON",
			in:             requestWithValidationHeaders("../../data/examples/valid.json"),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusOK,
			expectedBody:   "{\"valid\":true,\"message\":\"The document is valid\",\"errors\":null}\n",
		},
		{
			name:           "Validate with invalid JSON",
			in:             requestWithValidationHeaders("../../data/examples/invalid.json"),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusOK,
			expectedBody:   "{\"valid\":false,\"message\":\"The document is not valid\",\"errors\":[\"biomaterial_id is required\",\"Additional property k is not allowed\"]}\n",
		},
		{
			name:           "Validate no header",
			in:             baseRequest("../../data/examples/invalid.json"),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Link header must be provided\n",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			a.Validate(test.out, test.in)
			assert.Equal(t, test.expectedStatus, test.out.Code)
			assert.Equal(t, test.expectedBody, test.out.Body.String())
		})
	}
}

func requestWithValidationHeaders(bodyFile string) *http.Request {
	request := baseRequest(bodyFile)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Link", `<https://schema.humancellatlas.org/type/biomaterial/5.1.0/specimen_from_organism>; rel="describedby"`)
	return request
}

func baseRequest(bodyFile string) *http.Request {
	body, err := os.Open(bodyFile)
	if err != nil {
		logger.Fatalf("failed to open test file: %s: %v", bodyFile, err)
	}
	return httptest.NewRequest("POST", "/validate", body)
}

func TestHandlers_CreateDataset(t *testing.T) {
	tests := []testDef{
		{
			name:                   "Create Dataset",
			in:                     httptest.NewRequest("POST", "/dataset", nil),
			out:                    httptest.NewRecorder(),
			expectedLocationHeader: "/dataset/.+",
			expectedStatus:         http.StatusCreated,
			expectedBody:           "{\"owner\":\".+\",\"dataset_id\":\".+\",\"created\":\".+\"}",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			a.CreateDataset(test.out, test.in)
			assert.Regexp(t, test.expectedLocationHeader, test.out.Header()["Location"])
			assert.Equal(t, test.expectedStatus, test.out.Code)
			assert.Regexp(t, test.expectedBody, test.out.Body)
		})
	}
}

func init() {
	logger = log.New(os.Stdout, "ingest-test ", log.LstdFlags|log.Lshortfile)
	dataStore, err := local.NewInMemoryDataStore(logger)
	if err != nil {
		logger.Fatalf("failed to create data store: %v", err)
	}

	fileStore, err := local.NewLocalFileStore(logger, "/tmp")
	if err != nil {
		logger.Fatalf("failed to create file store: %v", err)
	}
	a = NewAPI(logger, dataStore, fileStore)
}

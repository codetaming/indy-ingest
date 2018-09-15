package api

import (
	"github.com/codetaming/indy-ingest/internal/persistence"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var a *API

func TestHandlers_CreateDataset(t *testing.T) {
	tests := []struct {
		name                   string
		in                     *http.Request
		out                    *httptest.ResponseRecorder
		expectedLocationHeader string
		expectedStatus         int
		expectedBody           string
	}{
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
	logger := log.New(os.Stdout, "ingest-test ", log.LstdFlags|log.Lshortfile)
	dataStore, err := persistence.NewInMemoryDataStore(logger)
	if err != nil {
		logger.Fatalf("failed to create data store: %v", err)
	}

	fileStore, err := persistence.NewLocalFileStore(logger, "/tmp")
	if err != nil {
		logger.Fatalf("failed to create file store: %v", err)
	}
	a = NewAPI(logger, dataStore, fileStore)
}

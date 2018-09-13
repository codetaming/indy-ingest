package persistence

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestLocalFileStore(t *testing.T) {
	logger := log.New(os.Stdout, "ingest-test ", log.LstdFlags|log.Lshortfile)
	fileStore, err := NewLocalFileStore(logger, "/tmp")
	assert.Nil(t, err)
	key := "testKey"
	content := "test"
	location, err := fileStore.StoreMetadata(key, content)
	assert.Nil(t, err)
	println(location)
	retrievedContent, err := fileStore.RetrieveMetadata(key)
	assert.Equal(t, content, retrievedContent)
}

package persistence

import (
	"fmt"
	"io/ioutil"
	"log"
)

type LocalFileStore struct {
	logger        *log.Logger
	fileStoreRoot string
}

func (l *LocalFileStore) StoreMetadata(key string, content string) (location string, err error) {
	d1 := []byte(content)
	filename := fmt.Sprintf("%s/%s", l.fileStoreRoot, key)
	err = ioutil.WriteFile(filename, d1, 0644)
	if err != nil {
		l.logger.Printf("Error writing file: %s", filename)
	}
	return filename, err
}

func (l *LocalFileStore) RetrieveMetadata(key string) (content string, err error) {
	filename := fmt.Sprintf("%s/%s", l.fileStoreRoot, key)
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		l.logger.Printf("Error reading file: %s", filename)
	}
	return string(dat), err
}

func NewLocalFileStore(logger *log.Logger, fileStoreRoot string) (FileStore, error) {
	return &LocalFileStore{
		logger:        logger,
		fileStoreRoot: fileStoreRoot,
	}, nil
}

package persistence

import "log"

type LocalFileStore struct {
	logger *log.Logger
}

func (LocalFileStore) StoreMetadata(key string, bodyJson string) (string, error) {
	panic("implement me")
}

func (LocalFileStore) RetrieveMetadata(key string) (string, error) {
	panic("implement me")
}

func NewLocalFileStore(logger *log.Logger) (FileStore, error) {
	return LocalFileStore{
		logger: logger,
	}, nil
}

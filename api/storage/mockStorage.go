package storage

type MockStorage struct{}

func (MockStorage) StoreMetadata(key string, bodyJson string) (string, error) {
	return key, nil
}

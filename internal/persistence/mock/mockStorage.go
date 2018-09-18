package mock

type MockStorage struct{}

func (MockStorage) StoreMetadata(key string, bodyJson string) (string, error) {
	return key, nil
}

func (MockStorage) RetrieveMetadata(key string) (string, error) {
	return "{}", nil
}

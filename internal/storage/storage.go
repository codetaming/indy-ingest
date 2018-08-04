package storage

type MetadataStorer interface {
	StoreMetadata(key string, bodyJson string) (string, error)
}

type MetadataRetriever interface {
	RetrieveMetadata(key string) (string, error)
}

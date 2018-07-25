package storage

type MetadataStorer interface {
	StoreMetadata(key string, bodyJson string) (string, error)
}

type MetadataRetriever interface {
	RetrieveMetadata(datasetId string, metadataId string) (string, error)
}

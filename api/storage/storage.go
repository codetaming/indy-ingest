package storage

type MetadataStorer interface {
	StoreMetadata(key string, bodyJson string) (string, error)
}

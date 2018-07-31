package publication

type MetadataCreatedPublisher interface {
	PublishMetadataCreated(metadataUrl string) (err error)
}

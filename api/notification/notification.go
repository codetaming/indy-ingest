package notification

type MetadataCreatedNotifier interface {
	NotifyMetadataCreated(metadataUrl string) (err error)
}

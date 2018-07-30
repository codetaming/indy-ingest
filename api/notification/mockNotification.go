package notification

import "log"

type MockNotification struct{}

func (MockNotification) NotifyMetadataCreated(metadataUrl string) (err error) {
	log.Print("metadata created: " + metadataUrl)
	return nil
}

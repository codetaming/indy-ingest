package mock

import "log"

type MockPublication struct{}

func (MockPublication) PublishMetadataCreated(metadataUrl string) (err error) {
	log.Print("metadata created: " + metadataUrl)
	return nil
}

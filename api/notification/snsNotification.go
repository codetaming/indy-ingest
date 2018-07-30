package notification

type SNSNotification struct{}

func (SNSNotification) NotifyMetadataCreated(metadataUrl string) (err error) {
	return nil
}

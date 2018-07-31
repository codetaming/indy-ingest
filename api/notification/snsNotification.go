package notification

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"os"
)

type SNSNotification struct{}

var snsSvc *sns.SNS

func init() {
	region := os.Getenv("AWS_REGION")
	if ses, err := session.NewSession(&aws.Config{
		Region: &region,
	}); err != nil {
		fmt.Println(fmt.Sprintf("Failed to connect to AWS: %s", err.Error()))
	} else {
		snsSvc = sns.New(ses)
	}
}

func (SNSNotification) NotifyMetadataCreated(metadataUrl string) (err error) {
	return nil
}

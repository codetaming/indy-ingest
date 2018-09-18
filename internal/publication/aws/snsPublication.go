package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"log"
	"os"
)

type SnsPublication struct{}

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

func (SnsPublication) PublishMetadataCreated(metadataUrl string) (err error) {
	topic := os.Getenv("SNS_METADATA_CREATED")
	params := &sns.PublishInput{
		Message:  aws.String(metadataUrl),
		TopicArn: aws.String(topic),
	}
	resp, err := snsSvc.Publish(params)
	if err != nil {
		log.Printf(err.Error())
	}
	log.Println(resp.MessageId)
	return err
}

package storage

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"log"
	"os"
	"strings"
)

type S3Storage struct{}

var s3Uploader *s3manager.Uploader

func init() {
	region := os.Getenv("AWS_REGION")
	if ses, err := session.NewSession(&aws.Config{
		Region: &region,
	}); err != nil {
		fmt.Println(fmt.Sprintf("Failed to connect to AWS: %s", err.Error()))
	} else {
		s3Uploader = s3manager.NewUploader(ses)
	}
}

func (S3Storage) StoreMetadata(key string, bodyJson string) (string, error) {
	contentType := "application/json"
	log.Printf("Uploading file to S3: " + key)
	upParams := &s3manager.UploadInput{
		Bucket:      aws.String(os.Getenv("METADATA_BUCKET")),
		ContentType: &contentType,
		Key:         &key,
		Body:        strings.NewReader(bodyJson),
	}
	result, err := s3Uploader.Upload(upParams)
	if err != nil {
		log.Panic(fmt.Sprintf("failed to create S3 file, %v", err))
		return "", err
	}
	return result.Location, nil
}

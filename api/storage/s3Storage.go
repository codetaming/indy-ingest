package storage

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type S3Storage struct{}

var s3Uploader *s3manager.Uploader
var s3Svc *s3.S3

func init() {
	region := os.Getenv("AWS_REGION")
	if ses, err := session.NewSession(&aws.Config{
		Region: &region,
	}); err != nil {
		fmt.Println(fmt.Sprintf("Failed to connect to AWS: %s", err.Error()))
	} else {
		s3Uploader = s3manager.NewUploader(ses)
		s3Svc = s3.New(ses)
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

func (S3Storage) RetrieveMetadata(key string) (string, error) {
	result, err := s3Svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("METADATA_BUCKET")),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Panic(fmt.Sprintf("failed to retrieve S3 file, %v", err))
		return "", err
	}
	if b, err := ioutil.ReadAll(result.Body); err == nil {
		return string(b), nil
	}
	return "", nil
}

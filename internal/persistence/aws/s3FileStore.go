package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/codetaming/indy-ingest/internal/persistence"
	"io/ioutil"
	"log"
	"strings"
)

type S3FileStore struct {
	logger         *log.Logger
	s3Uploader     *s3manager.Uploader
	s3Svc          *s3.S3
	metadataBucket string
}

func NewS3FileStore(logger *log.Logger, region string, metadataBucket string) (persistence.FileStore, error) {
	if ses, err := session.NewSession(&aws.Config{
		Region: &region,
	}); err != nil {
		fmt.Println(fmt.Sprintf("Failed to connect to AWS: %s", err.Error()))
		return nil, err
	} else {
		return S3FileStore{
			logger:         logger,
			s3Uploader:     s3manager.NewUploader(ses),
			s3Svc:          s3.New(ses),
			metadataBucket: metadataBucket,
		}, nil
	}
}

func (f S3FileStore) StoreMetadata(key string, bodyJson string) (string, error) {
	contentType := "application/json"
	f.logger.Printf("Uploading file to S3: " + key)
	upParams := &s3manager.UploadInput{
		Bucket:      &f.metadataBucket,
		ContentType: &contentType,
		Key:         &key,
		Body:        strings.NewReader(bodyJson),
	}
	result, err := f.s3Uploader.Upload(upParams)
	if err != nil {
		f.logger.Panic(fmt.Sprintf("failed to create S3 file, %v", err))
		return "", err
	}
	return result.Location, nil
}

func (f S3FileStore) RetrieveMetadata(key string) (string, error) {
	result, err := f.s3Svc.GetObject(&s3.GetObjectInput{
		Bucket: &f.metadataBucket,
		Key:    aws.String(key),
	})
	if err != nil {
		f.logger.Panic(fmt.Sprintf("failed to retrieve S3 file, %v", err))
		return "", err
	}
	if b, err := ioutil.ReadAll(result.Body); err == nil {
		return string(b), nil
	}
	return "", err
}

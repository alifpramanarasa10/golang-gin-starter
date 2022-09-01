package aws

import (
	"fmt"
	"gin-starter/config"
	"gin-starter/utils"
	"mime/multipart"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
)

// S3Bucket define struct for aws s3 bucket integration
type S3Bucket struct {
	cfg     config.Config
	session *session.Session
}

// NewS3Bucket initiate S3 bucket SDK
func NewS3Bucket(cfg config.Config, session *session.Session) *S3Bucket {
	return &S3Bucket{
		cfg:     cfg,
		session: session,
	}
}

// Upload uploads file to bucket
func (s *S3Bucket) Upload(f *multipart.FileHeader, folder string) (string, error) {
	bucket := s.cfg.AWS.BucketName

	var err error

	src, err := f.Open()
	if err != nil {
		return "", errors.Wrap(err, "[CloudStorageService-Upload] error open file")
	}

	defer func() {
		if err := src.Close(); err != nil {
			fmt.Println("error while closing gin context form file :", err)
		}
	}()

	splitFilename := strings.Split(f.Filename, ".")
	fileMime := ""

	if len(splitFilename) > 0 {
		fileMime = splitFilename[len(splitFilename)-1]
	}

	fileName := fmt.Sprintf("%s.%s", utils.SHAEncrypt(f.Filename), fileMime)

	fileStored := fmt.Sprintf("%s/%s", folder, fileName)

	uploader := s3manager.NewUploader(s.session)

	// Upload to the s3 bucket
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		ACL:    aws.String("public-read"),
		Key:    aws.String(fileName),
		Body:   src,
	})

	if err != nil {
		return "", fmt.Errorf("[S3CloudStorage] error while uploading file : %s", err)
	}

	return fileStored, nil
}

// Delete delete file from bucket
func (s *S3Bucket) Delete(path string) error {
	svc := s3.New(s.session)

	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: &s.cfg.AWS.BucketName,
		Key:    &path,
	})
	if err != nil {
		return err
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: &s.cfg.AWS.BucketName,
		Key:    &path,
	})

	if err != nil {
		return err
	}

	return nil
}

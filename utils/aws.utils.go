package utils

import (
	"gin-starter/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// InitAWSS3 initiate aws session based on defined credentials
func InitAWSS3(cfg config.Config) *session.Session {
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(cfg.AWS.Region),
			Credentials: credentials.NewStaticCredentials(
				cfg.AWS.AccessKeyID,
				cfg.AWS.SecretAccessKey,
				"", // a token will be created when the session is used.
			),
		})

	if err != nil {
		panic(err)
	}

	return sess
}

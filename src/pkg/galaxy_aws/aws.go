package galaxy_aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jctanner/galaxygo/pkg/galaxy_logger"
	"github.com/jctanner/galaxygo/pkg/galaxy_settings"
)

var settings = galaxy_settings.NewGalaxySettings()
var logger = galaxy_logger.Logger{}

func GetS3ObjectByFilepath(filepath string) *s3.GetObjectOutput {

	// get the access key id
	aws_access_key := settings.Aws_access_key
	logger.Debug(fmt.Sprintf("s3 access ley %v", aws_access_key))

	// get the secret key
	aws_secret_key := settings.Aws_secret_key
	logger.Debug(fmt.Sprintf("s3 secret key %v", aws_secret_key))

	// get the s3 region
	aws_region := settings.Aws_s3_region
	logger.Debug(fmt.Sprintf("s3 region %v", aws_region))

	// get the s3 url
	s3_endpoint_url := settings.Aws_s3_endpoint_url
	logger.Debug(fmt.Sprintf("s3 endpoint url %v", s3_endpoint_url))

	// get the s3 bucket
	s3_bucket_name := settings.Aws_s3_bucket_name
	logger.Debug(fmt.Sprintf("s3 bucket name %v", s3_bucket_name))

	// set the creds
	creds := credentials.NewStaticCredentials(aws_access_key, aws_secret_key, "")
	logger.Debug(fmt.Sprintf("aws credentials %v", creds))

	// Create a new aws session
	sess, err := session.NewSession(&aws.Config{
		Endpoint:    aws.String(s3_endpoint_url),
		Region:      aws.String(aws_region),
		Credentials: creds,
	})
	if err != nil {
		fmt.Println("Failed to create session", err)
		return nil
	}

	// enable http request and response logging
	if settings.Debug {
		sess.Config.WithLogLevel(aws.LogDebugWithHTTPBody)
	}

	filekey := s3_bucket_name + "/" + filepath
	logger.Debug(fmt.Sprintf("s3 filekey %v", filekey))

	// Create a new S3 service client
	svc := s3.New(sess)

	resp, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s3_bucket_name),
		Key:    aws.String(filekey),
	})
	defer resp.Body.Close()

	return resp
}

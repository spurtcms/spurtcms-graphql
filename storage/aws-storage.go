package storage

import (
	"fmt"
	"log"

	"github.com/99designs/gqlgen/graphql"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// create session
func CreateS3Session(AwsCredentials map[string]interface{}) (ses *s3.S3, err error) {

	var awsId, awsKey, awsRegion string

	if AwsCredentials != nil {

		awsId = AwsCredentials["AccessId"].(string)

		awsKey = AwsCredentials["AccessKey"].(string)

		awsRegion = AwsCredentials["Region"].(string)

		// awsBucket =  AwsCredentials["BucketName"].(string)

	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsId, awsKey, ""),
	})

	if err != nil {

		log.Println("Error creating session: ", err)

		return nil, err

	}

	svc := s3.New(sess)

	return svc, nil

}

func CreateS3Sess(AwsCredentials map[string]interface{}) *session.Session {

	var awsId, awsKey, awsRegion string

	if AwsCredentials != nil {

		awsId = AwsCredentials["AccessId"].(string)

		awsKey = AwsCredentials["AccessKey"].(string)

		awsRegion = AwsCredentials["Region"].(string)

	}

	// The session the S3 Uploader will use
	sess := session.Must(session.NewSession(
		&aws.Config{
			Region:      &awsRegion,
			Credentials: credentials.NewStaticCredentials(awsId, awsKey, ""),
		},
	))

	return sess
}

/*upload files to s3 */
func UploadFileS3(AwsCredentials map[string]interface{}, upload *graphql.Upload, filePath string) error {

	log.Println("enter",AwsCredentials)

	session := CreateS3Sess(AwsCredentials)

	awsBucket := AwsCredentials["BucketName"].(string)

	fmt.Println("upload filename :==", filePath)

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(session)

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(awsBucket),
		Key:    aws.String(filePath),
		Body:   upload.File,
		ACL:    aws.String("public-read"),
	})

	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}

	fmt.Printf("file uploaded to, %s\n", aws.StringValue(&result.Location))

	return nil
}

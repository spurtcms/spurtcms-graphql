package storage

import (
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// create session
func CreateS3Session(awsSession *session.Session) (ses *s3.S3) {

	svc := s3.New(awsSession)

	return svc

}

func CreateAwsSession(AwsCredentials map[string]interface{}) *session.Session {

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
func UploadFileS3(AwsCredentials map[string]interface{}, upload *graphql.Upload, base64String string, filePath string) error {

	session := CreateAwsSession(AwsCredentials)

	awsBucket := AwsCredentials["BucketName"].(string)

	fmt.Println("upload filename :==", filePath)

	var s3Body io.Reader
	
	var err error

	if base64String != "" {

		base64data := base64String[strings.IndexByte(base64String, ',')+1:]

		imageData, err := base64.StdEncoding.DecodeString(base64data)

		if err != nil {
			
			return err
		}

		s3Body = strings.NewReader(string(imageData))

	} else if upload != nil {

		s3Body = upload.File

	}

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(session)

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(awsBucket),
		Key:    aws.String(filePath),
		Body:   s3Body,
		ACL:    aws.String("public-read"),
	})

	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}

	fmt.Printf("file uploaded to, %s\n", aws.StringValue(&result.Location))

	return nil
}

func CheckS3FileExistence(AwsCredentials map[string]interface{}, fileName string) (bool, error) {

	session := CreateAwsSession(AwsCredentials)

	s3Svc := CreateS3Session(session)

	awsBucket := AwsCredentials["BucketName"].(string)

	// HeadObject to check if the file exists
	obj, err := s3Svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(awsBucket),
		Key:    aws.String(fileName),
	})

	fmt.Println("checking", obj)

	if err != nil {

		fmt.Printf("s3 storage file exist error: %v\n", err)

		return false, err
	}

	return true, nil

}

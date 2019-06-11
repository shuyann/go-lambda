package main

import (
	"archive/zip"
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func Handler(ctx context.Context, s3Event events.S3Event) error {
	log.Printf("INFO: start s3 function.")
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := s3.New(sess)
	for _, record := range s3Event.Records {
		s3Record := record.S3
		// get s3 object from s3
		obj, err := svc.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(s3Record.Bucket.Name),
			Key:    aws.String(s3Record.Object.Key),
		})
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case s3.ErrCodeNoSuchBucket:
					log.Fatalf("bucket %v does not exist", s3Record.Bucket.Name)
				case s3.ErrCodeNoSuchKey:
					log.Fatalf("key %v does not exist", s3Record.Object.Key)
				default:
					log.Fatalf("aws error %v", aerr.Error())
				}
			}
		}
		// create zip file.
		zipName, zipPath, err := createZipFile(obj, s3Record.Object.Key)
		if err != nil {
			log.Fatal(err)
		}
		// upload zip file to s3.
		zipFile, err := os.Open(zipPath)
		if err != nil {
			log.Fatal(err)
		}
		svc.PutObject(&s3.PutObjectInput{
			Body:   zipFile,
			Bucket: aws.String(os.Getenv("bucket")),
			Key:    aws.String(zipName),
		})
	}
	log.Printf("INFO: end s3 function.")

	return nil
}

func createZipFile(obj *s3.GetObjectOutput, key string) (zipFileName, zipFilePath string, err error) {
	log.Printf("INFO: start create zip file key: %v", key)
	tmpPath := "/tmp/"
	fileName := tmpPath + key
	zipFileName = fileNameWithoutExt(key) + ".zip"
	zipFilePath = tmpPath + zipFileName

	content, err := ioutil.ReadAll(obj.Body)
	if err != nil {
		return
	}
	file, err := os.Create(fileName)
	if err != nil {
		return
	}
	defer file.Close()

	if _, err = file.Write(content); err != nil {
		return
	}

	dest, err := os.Create(zipFilePath)
	if err != nil {
		return
	}
	defer dest.Close()

	zipWriter := zip.NewWriter(dest)
	defer zipWriter.Close()

	src, err := os.Open(fileName)
	if err != nil {
		return
	}

	writer, err := zipWriter.Create(fileName)
	if _, err = io.Copy(writer, src); err != nil {
		return
	}
	log.Printf("INFO: end create zip file key: %v", key)
	return
}

func fileNameWithoutExt(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}

func main() {
	lambda.Start(Handler)
}

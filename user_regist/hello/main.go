package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/google/uuid"
	"github.com/guregu/dynamo"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"
)

type Request events.APIGatewayProxyRequest
type Response events.APIGatewayProxyResponse
type UserItem struct {
	ID         string `dynamo:"id"`
	Name       string `dynamo:"name"`
	Email      string `dynamo:"email"`
	Url        string `dynamo:"url"`
	AcceptedAt int    `dynamo:"accepted_at"`
}

func Handler(ctx context.Context, req Request) (Response, error) {
	// Create session using the default region and credentials
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	db := dynamo.New(sess)
	table := db.Table("user")

	bucket := os.Getenv("bucket")
	key := os.Getenv("key")
	email := os.Getenv("email")

	values, err := url.ParseQuery(req.Body)
	if err != nil {
		log.Fatal(err)
	}

	username, email := values["username"][0], values["email"][0]

	now, err := strconv.Atoi(time.Now().Format("20060102150405"))
	if err != nil {
		log.Fatal(err)
	}

	s3Cli := s3.New(sess)
	request, _ := s3Cli.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	preSignedUrl, err := request.Presign(15 * time.Minute)
	if err != nil {
		log.Fatal(err)
	}

	userItem := UserItem{
		ID:         uuid.New().String(),
		Name:       username,
		Email:      email,
		Url:        preSignedUrl,
		AcceptedAt: now,
	}

	if err = table.Put(userItem).Run(); err != nil {
		log.Fatal(err)
	}

	// Please change region to us-east for ses use
	sesCli := ses.New(sess)

	sesInput := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(email),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(preSignedUrl),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String("Data"),
			},
		},
		Source: aws.String(email),
	}
	_, err = sesCli.SendEmail(sesInput)
	if err != nil {
		log.Fatal(err)
	}

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            "<html><body>Thank you register</body></html>",
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}

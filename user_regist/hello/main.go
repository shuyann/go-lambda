package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/google/uuid"
	"github.com/guregu/dynamo"
	"log"
	"net/url"
	"strconv"
	"time"
)

type Request events.APIGatewayProxyRequest
type Response events.APIGatewayProxyResponse
type UserItem struct {
	ID         string `dynamo:"id"`
	Name       string `dynamo:"name"`
	Email      string `dynamo:"email"`
	AcceptedAt int    `dynamo:"accepted_at"`
}

func Handler(ctx context.Context, req Request) (Response, error) {
	// Create session using the default region and credentials
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	db := dynamo.New(sess)
	table := db.Table("user")

	// Parse form data.
	values, err := url.ParseQuery(req.Body)
	if err != nil {
		log.Fatal(err)
	}

	username, email := values["username"][0], values["email"][0]

	now, err := strconv.Atoi(time.Now().Format("20060102150405"))
	if err != nil {
		log.Fatal(err)
	}

	userItem := UserItem{
		ID:         uuid.New().String(),
		Name:       username,
		Email:      email,
		AcceptedAt: now,
	}
	fmt.Println("INFO userItem: ", userItem)

	if err = table.Put(userItem).Run(); err != nil {
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

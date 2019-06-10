package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
)

type SimpleRequest struct {
	X string `json:"x"`
	Y string `json:"y"`
}

type SimpleResponse struct {
	Result string `json:"result"`
}

func Handler(ctx context.Context, request SimpleRequest) (SimpleResponse, error) {

	fmt.Println("x = ", request.X, "y = ", request.Y)
	x, err := strconv.Atoi(request.X)
	if err != nil {
		return SimpleResponse{}, err
	}
	y, err := strconv.Atoi(request.Y)
	if err != nil {
		return SimpleResponse{}, err
	}
	result := x / y

	return SimpleResponse{Result: strconv.Itoa(result)}, nil
}

func main() {
	lambda.Start(Handler)
}

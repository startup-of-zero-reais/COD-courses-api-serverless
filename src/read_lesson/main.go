package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
)

func main() {
	lambda.Start(handler)
}

type MyResponse struct {
	Message string `json:"message"`
}

func (m *MyResponse) String() string {
	bytes, _ := json.Marshal(m)
	return string(bytes)
}

func handler(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	
	return events.APIGatewayProxyResponse{
		Body:	   (&MyResponse{Message: "Hello World!"}).String(),
		StatusCode: 200,
		IsBase64Encoded: false,
	}, nil
}
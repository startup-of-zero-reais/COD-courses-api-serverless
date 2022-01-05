package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"log"
	"write_lesson/common"
)

type (
	Handler struct {
		dynamoClient *dynamodb.Client
		table        *string
	}
)

func main() {
	dynamoClient, err := common.DynamoConnection()
	if err != nil {
		log.Fatalf("failed to load config, %v", err)
	}

	h := Handler{
		dynamoClient: dynamoClient,
		table:        aws.String("code-craft-courses-table"),
	}

	lambda.Start(h.HandleLambda)
}

func (receiver Handler) HandleLambda(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return common.ResponseProxy(200, common.NewDataResponse(event.Body), nil)
}

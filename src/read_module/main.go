package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"log"
	"read_module/common"
)

type Handler struct {
	dynamoClient *dynamodb.Client
	table        *string
}

func main() {
	dynamoClient, err := common.DynamoConnection()
	if err != nil {
		log.Fatalf("failed to load connection: %v", err)
	}

	h := Handler{
		dynamoClient: dynamoClient,
		table:        aws.String("code-craft-courses-table"),
	}

	lambda.Start(h.HandleLambda)
}

func (h Handler) HandleLambda(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.QueryStringParameters["action"] != "search" {
		return h.Get(request)
	}

	return h.Search(request)
}

func (h Handler) Get(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return common.ResponseProxy(200, common.NewMessage("OK"), nil)
}

func (h Handler) Search(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return common.ResponseProxy(200, common.NewMessage("OK"), nil)
}

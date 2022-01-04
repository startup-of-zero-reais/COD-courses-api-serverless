package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
	"read_lesson/common"
)

type Handler struct {
	dynamoClient *dynamodb.Client
	table        *string
}

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

func (h Handler) HandleLambda(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if event.QueryStringParameters["action"] != "search" {
		return h.HandleSingleScan(event)
	}

	return h.HandleSearch(event)
}

func (h Handler) HandleSingleScan(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var lessonId string
	if lessonId = request.QueryStringParameters["id"]; lessonId == "" {
		return common.ResponseProxy(400, common.NewMessage("voce deve fornecer o id"), nil)
	}

	var sectionId string
	if sectionId = request.QueryStringParameters["section_id"]; sectionId == "" {
		return common.ResponseProxy(400, common.NewMessage("voce deve fornecer o identificador de seção"), nil)
	}

	input := &dynamodb.GetItemInput{
		TableName: h.table,
		Key: map[string]types.AttributeValue{
			"ID":       &types.AttributeValueMemberS{Value: fmt.Sprintf("lesson_%s", lessonId)},
			"ParentID": &types.AttributeValueMemberS{Value: fmt.Sprintf("section_%s", sectionId)},
		},
	}

	output, err := h.dynamoClient.GetItem(context.TODO(), input)
	if err != nil {
		return common.ResponseProxy(500, common.NewMessage(err.Error()), err)
	}

	var item common.Lesson
	err = attributevalue.UnmarshalMap(output.Item, &item)
	if err != nil {
		return common.ResponseProxy(500, common.NewMessage(err.Error()), err)
	}

	if &item == nil || item.ID == "" {
		return common.ResponseProxy(404, common.NewMessage("nenhuma aula encontrada com esses parâmetros de busca"), nil)
	}

	return common.ResponseProxy(200, common.NewDataResponse(&item), nil)
}

func (h Handler) HandleSearch(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return common.ResponseProxy(200, common.NewDataResponse("Ok"), nil)
}

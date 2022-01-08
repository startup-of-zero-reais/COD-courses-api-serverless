package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
	"read_section/common"
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

func (h Handler) HandleLambda(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.QueryStringParameters["action"] != "search" {
		return h.Get(request)
	}

	return h.Search(request)
}

func (h Handler) Get(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var sectionId string
	if sectionId = request.QueryStringParameters["section"]; sectionId == "" {
		return common.ResponseProxy(400, common.NewMessage("voce deve fornecer a section"), nil)
	}

	var moduleId string
	if moduleId = request.QueryStringParameters["module"]; moduleId == "" {
		return common.ResponseProxy(400, common.NewMessage("voce deve fornecer o module"), nil)
	}

	input := &dynamodb.GetItemInput{
		TableName: h.table,
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("MODULE#%s", moduleId)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("SECTION#%s", sectionId)},
		},
	}

	output, err := h.dynamoClient.GetItem(context.TODO(), input)
	if err != nil {
		return common.ResponseProxy(500, common.NewMessage(err.Error()), nil)
	}

	var item common.Section
	err = attributevalue.UnmarshalMap(output.Item, &item)
	if err != nil {
		return common.ResponseProxy(500, common.NewMessage(err.Error()), nil)
	}

	if &item == nil || item.SK == "" {
		return common.ResponseProxy(404, common.NewMessage("Seção não encontrada"), nil)
	}

	return common.ResponseProxy(200, common.NewDataResponse(&item), nil)
}

func (h Handler) Search(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	moduleId := request.QueryStringParameters["module"]

	if moduleId == "" {
		return common.ResponseProxy(400, common.NewMessage("Forneça um modulo para busca"), nil)
	}

	hashValue := fmt.Sprintf("MODULE#%s", moduleId)
	input := new(dynamodb.QueryInput)
	input.TableName = h.table

	expr, err := expression.NewBuilder().WithKeyCondition(
		expression.KeyAnd(
			expression.Key("PK").Equal(expression.Value(hashValue)),
			expression.KeyBeginsWith(expression.Key("SK"), "SECTION#"),
		),
	).Build()
	if err != nil {
		return common.ResponseProxy(500, common.NewMessage(err.Error()), err)
	}

	input.KeyConditionExpression = expr.KeyCondition()
	input.ExpressionAttributeValues = expr.Values()
	input.ExpressionAttributeNames = expr.Names()

	output, err := h.dynamoClient.Query(context.TODO(), input)
	if err != nil {
		return common.ResponseProxy(500, common.NewMessage("Falha ao executar query"), err)
	}

	var sections []common.Section
	err = attributevalue.UnmarshalListOfMaps(output.Items, &sections)
	if err != nil {
		return common.ResponseProxy(500, common.NewMessage("falha ao executar unmarshal"), err)
	}

	return common.ResponseProxy(200, common.NewDataResponse(sections), nil)
}

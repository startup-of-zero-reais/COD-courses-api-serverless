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
	"read_course/common"
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
	if userId := request.Headers["user"]; userId == "" {
		return common.ResponseProxy(406, common.NewMessage("user é um header obrigatório"), nil)
	}

	if request.QueryStringParameters["action"] != "search" {
		return h.Get(request)
	}

	return h.Search(request)
}

func (h Handler) Get(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var userId string
	if userId = request.Headers["user"]; userId == "" {
		return common.ResponseProxy(406, common.NewMessage("user é um header obrigatório"), nil)
	}

	var courseId string
	if courseId = request.QueryStringParameters["course"]; courseId == "" {
		return common.ResponseProxy(400, common.NewMessage("course é um parâmetro de busca obrigatório"), nil)
	}

	hashValue := fmt.Sprintf("USER#%s", userId)
	rangeValue := fmt.Sprintf("COURSE#%s", courseId)

	input := &dynamodb.GetItemInput{
		TableName: h.table,
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: hashValue},
			"SK": &types.AttributeValueMemberS{Value: rangeValue},
		},
	}

	output, err := h.dynamoClient.GetItem(context.TODO(), input)
	if err != nil {
		return common.ResponseProxy(500, common.NewMessage(fmt.Sprintf("erro ao recuperar itens: %v", err)), nil)
	}

	var item common.Course
	err = attributevalue.UnmarshalMap(output.Item, &item)
	if err != nil {
		return common.ResponseProxy(500, common.NewMessage(fmt.Sprintf("erro ao executar unmarshal: %v", err)), nil)
	}

	if &item == nil || item.SK == "" {
		return common.ResponseProxy(404, common.NewMessage("nenhum resultado encontrado"), err)
	}

	return common.ResponseProxy(200, common.NewDataResponse(item), nil)
}

func (h Handler) Search(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userId := request.Headers["user"]

	hashValue := fmt.Sprintf("USER#%s", userId)
	input := new(dynamodb.QueryInput)
	input.TableName = h.table

	expr, err := expression.NewBuilder().WithKeyCondition(
		expression.KeyAnd(
			expression.Key("PK").Equal(expression.Value(hashValue)),
			expression.KeyBeginsWith(expression.Key("SK"), "COURSE#"),
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

	if output.Count <= 0 {
		return common.ResponseProxy(404, common.NewMessage("Nenhum resultado encontrado para a busca"), nil)
	}

	var modules []common.Course
	err = attributevalue.UnmarshalListOfMaps(output.Items, &modules)

	if err != nil {
		return common.ResponseProxy(500, common.NewMessage("falha ao executar unmarshal"), err)
	}

	return common.ResponseProxy(200, common.NewDataResponse(modules), nil)
}

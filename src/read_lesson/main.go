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
	if lessonId = request.QueryStringParameters["lesson"]; lessonId == "" {
		return common.ResponseProxy(400, common.NewMessage("voce deve fornecer a lesson"), nil)
	}

	var sectionId string
	if sectionId = request.QueryStringParameters["section"]; sectionId == "" {
		return common.ResponseProxy(400, common.NewMessage("voce deve fornecer o identificador de seção"), nil)
	}

	input := &dynamodb.GetItemInput{
		TableName: h.table,
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("SECTION#%s", sectionId)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("LESSON#%s", lessonId)},
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

	if &item == nil || item.SK == "" {
		return common.ResponseProxy(404, common.NewMessage("nenhuma aula encontrada com esses parâmetros de busca"), nil)
	}

	return common.ResponseProxy(200, common.NewDataResponse(&item), nil)
}

func (h Handler) HandleSearch(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sectionId := request.QueryStringParameters["section"]
	moduleId := request.QueryStringParameters["module"]
	courseId := request.QueryStringParameters["course"]

	if sectionId == "" && courseId == "" && moduleId == "" {
		return common.ResponseProxy(400, common.NewMessage("Busque por curso, módulo ou seção de curso"), nil)
	}

	input := new(dynamodb.QueryInput)
	input.TableName = h.table

	var hashKey string
	var hashValue string

	if courseId != "" {
		input.IndexName = aws.String("CourseLessonsIndex")
		hashKey = "ParentCourse"
		hashValue = fmt.Sprintf("COURSE#%s", courseId)
	}

	if moduleId != "" {
		input.IndexName = aws.String("ModuleLessonsIndex")
		hashKey = "ParentModule"
		hashValue = fmt.Sprintf("MODULE#%s", moduleId)
	}

	if sectionId != "" {
		input.IndexName = nil
		hashKey = "PK"
		hashValue = fmt.Sprintf("SECTION#%s", sectionId)
	}

	expr, _ := expression.NewBuilder().WithKeyCondition(
		expression.KeyAnd(
			expression.Key(hashKey).Equal(expression.Value(hashValue)),
			expression.KeyBeginsWith(expression.Key("SK"), "LESSON#"),
		),
	).Build()

	input.KeyConditionExpression = expr.KeyCondition()
	input.ExpressionAttributeValues = expr.Values()
	input.ExpressionAttributeNames = expr.Names()

	output, err := h.dynamoClient.Query(context.TODO(), input)
	if err != nil {
		return common.ResponseProxy(500, common.NewMessage("Falha ao executar scan"), err)
	}

	var lessons []common.Lesson
	err = attributevalue.UnmarshalListOfMaps(output.Items, &lessons)
	if err != nil {
		return common.ResponseProxy(500, common.NewMessage("Falha ao executar unmarshal"), err)
	}

	return common.ResponseProxy(200, common.NewDataResponse(lessons), nil)
}

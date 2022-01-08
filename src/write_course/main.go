package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
	"reflect"
	"write_course/common"
)

type Handler struct {
	dynamoClient *dynamodb.Client
	table        *string
}

func main() {
	dynamoClient, err := common.DynamoConnection()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	h := Handler{
		dynamoClient: dynamoClient,
		table:        aws.String("code-craft-courses-table"),
	}

	lambda.Start(h.HandleLambda)
}

func (h Handler) HandleLambda(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var userId string
	if userId = event.Headers["user"]; userId == "" {
		return common.ResponseProxy(406, common.NewMessage("você deve fornecer header do user"), nil)
	}

	if event.Body == "" {
		return common.ResponseProxy(400, common.NewMessage("você deve fornecer os dados para serem adicionados"), nil)
	}

	course, err := h.ValidateBody(event.Body)
	if err != nil {
		return common.ResponseProxy(400, common.NewMessage(fmt.Sprintf("erro de validação: %s", err.Error())), nil)
	}

	item := map[string]types.AttributeValue{
		"PK":          &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userId)},
		"SK":          &types.AttributeValueMemberS{Value: fmt.Sprintf("COURSE#%s", course.SK)},
		"Thumb":       &types.AttributeValueMemberS{Value: course.Thumb},
		"Title":       &types.AttributeValueMemberS{Value: course.Title},
		"Description": &types.AttributeValueMemberS{Value: course.Description},
		"Owner":       &types.AttributeValueMemberS{Value: fmt.Sprintf("OWNER#%s", userId)},
		"CartOpen":    &types.AttributeValueMemberBOOL{Value: course.CartOpen},
		"CreatedAt":   &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", course.CreatedAt)},
		"UpdatedAt":   &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", course.UpdatedAt)},
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: h.table,
	}

	_, err = h.dynamoClient.PutItem(context.TODO(), input)
	if err != nil {
		return common.ResponseProxy(500, common.NewMessage(
			fmt.Sprintf("erro ao criar novo curso: %s", err.Error()),
		), nil)
	}

	return common.ResponseProxy(200, common.NewDataResponse(course), nil)
}

func (h Handler) ValidateBody(message string) (*common.Course, error) {
	var course common.Course
	err := json.Unmarshal([]byte(message), &course)
	if err != nil {
		return nil, err
	}

	nonRequiredFields := []string{
		"created_at",
		"updated_at",
		"cart_open",
	}

	inSlice := func(fieldName string) bool {
		for _, nonRequired := range nonRequiredFields {
			if fieldName == nonRequired {
				return true
			}
		}

		return false
	}

	reflected := reflect.ValueOf(course)
	refType := reflect.TypeOf(course)

	for i := 0; i < reflected.NumField(); i++ {
		field := reflected.Field(i)
		fieldName := refType.Field(i).Tag.Get("json")
		if field.IsZero() && !inSlice(fieldName) {
			return nil, errors.New(fmt.Sprintf("O campo %s é obrigatório", fieldName))
		}
	}

	(&course).BeforeCreate()

	return &course, nil
}

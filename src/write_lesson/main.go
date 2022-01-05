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

func (h Handler) HandleLambda(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if event.Body == "" {
		return common.ResponseProxy(400, common.NewMessage("você deve fornecer os dados para serem adicionados"), nil)
	}

	lesson, err := h.ValidateBody(event.Body)
	if err != nil {
		return common.ResponseProxy(400, common.NewMessage(fmt.Sprintf("erro de validação: %s", err.Error())), nil)
	}

	artifacts := map[string]types.AttributeValue{}
	for artifact, value := range lesson.Artifacts {
		artifacts[artifact] = &types.AttributeValueMemberS{Value: value}
	}

	item := map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("SECTION#%s", lesson.PK)},
		"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("LESSON#%s", lesson.SK)},
		"Artifacts": &types.AttributeValueMemberM{
			Value: artifacts,
		},
		"ParentCourse": &types.AttributeValueMemberS{Value: fmt.Sprintf("COURSE#%s", lesson.ParentCourse)},
		"ParentModule": &types.AttributeValueMemberS{Value: fmt.Sprintf("MODULE#%s", lesson.ParentModule)},
		"Thumb":        &types.AttributeValueMemberS{Value: lesson.Thumb},
		"Title":        &types.AttributeValueMemberS{Value: lesson.Title},
		"VideoSource":  &types.AttributeValueMemberS{Value: lesson.VideoSource},
		"CreatedAt":    &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", lesson.CreatedAt)},
		"UpdatedAt":    &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", lesson.UpdatedAt)},
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: h.table,
	}

	_, err = h.dynamoClient.PutItem(context.TODO(), input)
	if err != nil {
		return common.ResponseProxy(500, common.NewMessage(fmt.Sprintf("erro ao criar nova aula: %s", err.Error())), nil)
	}

	return common.ResponseProxy(200, common.NewDataResponse(lesson), nil)
}

func (h Handler) ValidateBody(message string) (*common.Lesson, error) {
	var lesson common.Lesson
	err := json.Unmarshal([]byte(message), &lesson)
	if err != nil {
		return nil, err
	}

	nonRequiredFields := []string{
		"created_at",
		"updated_at",
		"lesson_id",
		"artifacts",
	}

	inSlice := func(fieldName string) bool {
		for _, nonRequired := range nonRequiredFields {
			if fieldName == nonRequired {
				return true
			}
		}

		return false
	}

	reflected := reflect.ValueOf(lesson)
	refType := reflect.TypeOf(lesson)

	for i := 0; i < reflected.NumField(); i++ {
		field := reflected.Field(i)
		fieldName := refType.Field(i).Tag.Get("json")
		if field.IsZero() && !inSlice(fieldName) {
			return nil, errors.New(fmt.Sprintf("O campo %s é obrigatório", fieldName))
		}
	}

	(&lesson).BeforeCreate()

	return &lesson, nil
}

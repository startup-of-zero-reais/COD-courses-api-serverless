package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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

	return common.ResponseProxy(200, common.NewDataResponse(lesson), nil)
}

func (h Handler) ValidateBody(message string) (*common.Lesson, error) {
	var lesson common.Lesson
	err := json.Unmarshal([]byte(message), &lesson)
	if err != nil {
		return nil, err
	}

	reflected := reflect.ValueOf(lesson)
	refType := reflect.TypeOf(lesson)

	for i := 0; i < reflected.NumField(); i++ {
		field := reflected.Field(i)
		fieldName := refType.Field(i).Tag.Get("json")
		if field.IsZero() && fieldName != "artifacts" {
			return nil, errors.New(fmt.Sprintf("O campo %s é obrigatório", fieldName))
		}
	}

	return &lesson, nil
}

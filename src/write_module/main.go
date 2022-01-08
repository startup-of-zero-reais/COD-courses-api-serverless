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
	"write_module/common"
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
	if event.Body == "" {
		return common.ResponseProxy(400, common.NewMessage("você deve fornecer os dados para serem adicionados"), nil)
	}

	module, err := h.ValidateBody(event.Body)
	if err != nil {
		return common.ResponseProxy(400, common.NewMessage(fmt.Sprintf("erro de validação: %s", err.Error())), nil)
	}

	item := map[string]types.AttributeValue{
		"PK":            &types.AttributeValueMemberS{Value: fmt.Sprintf("COURSE#%s", module.PK)},
		"SK":            &types.AttributeValueMemberS{Value: fmt.Sprintf("MODULE#%s", module.SK)},
		"Title":         &types.AttributeValueMemberS{Value: module.Title},
		"SectionsOrder": &types.AttributeValueMemberM{Value: ParseMap(module.SectionsOrder)},
		"UnlockAfter":   &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", module.UnlockAfter)},
		"CreatedAt":     &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", module.CreatedAt)},
		"UpdatedAt":     &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", module.UpdatedAt)},
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: h.table,
	}

	_, err = h.dynamoClient.PutItem(context.TODO(), input)
	if err != nil {
		return common.ResponseProxy(500, common.NewMessage(
			fmt.Sprintf("erro ao criar novo módulo: %s", err.Error()),
		), nil)
	}

	return common.ResponseProxy(200, common.NewDataResponse(module), nil)
}

func (h Handler) ValidateBody(message string) (*common.Module, error) {
	var module common.Module
	err := json.Unmarshal([]byte(message), &module)
	if err != nil {
		return nil, err
	}

	nonRequiredFields := []string{
		"sections_order",
		"created_at",
		"updated_at",
	}

	inSlice := func(fieldName string) bool {
		for _, nonRequired := range nonRequiredFields {
			if fieldName == nonRequired {
				return true
			}
		}

		return false
	}

	reflected := reflect.ValueOf(module)
	refType := reflect.TypeOf(module)

	for i := 0; i < reflected.NumField(); i++ {
		field := reflected.Field(i)
		fieldName := refType.Field(i).Tag.Get("json")
		if field.IsZero() && !inSlice(fieldName) {
			return nil, errors.New(fmt.Sprintf("O campo %s é obrigatório", fieldName))
		}
	}

	(&module).BeforeCreate()

	return &module, nil
}

func ParseMap(toParse map[string]string) map[string]types.AttributeValue {
	result := map[string]types.AttributeValue{}

	for key, value := range toParse {
		result[key] = &types.AttributeValueMemberS{Value: value}
	}

	return result
}

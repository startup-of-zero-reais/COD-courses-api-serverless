package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"log"
)

type (
	Item struct {
		ID       string `json:"course_id"`
		ParentID string `json:"module_id"`
	}
)

func main() {
	lambda.Start(handler)
}

var table = aws.String("code-craft-courses-table")

func handler(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	client, err := DynamoConnection()
	if err != nil {
		return ResponseProxy(500, &MessageResponse{}, err)
	}

	exprBuilder := expression.NewBuilder()

	if courseId := event.QueryStringParameters["course_id"]; courseId != "" {
		exprBuilder = exprBuilder.WithFilter(
			expression.Name("ID").Equal(expression.Value(courseId)),
		)
	} else {
		return ResponseProxy(400, &MessageResponse{Message: "voce deve fornecer o course_id"}, nil)
	}

	if moduleId := event.QueryStringParameters["module_id"]; moduleId != "" {
		exprBuilder = exprBuilder.WithFilter(
			expression.Name("ParentID").Equal(expression.Value(moduleId)),
		)
	}

	expr, err := exprBuilder.WithProjection(
		expression.NamesList(expression.Name("ID"), expression.Name("ParentID")),
	).Build()
	if err != nil {
		return ResponseProxy(500, &MessageResponse{}, err)
	}

	log.Printf("Expr: %+v\n", expr)

	input := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 table,
	}

	output, err := client.Scan(context.TODO(), input)
	if err != nil {
		return ResponseProxy(500, &MessageResponse{}, err)
	}

	var items []Item

	err = attributevalue.UnmarshalListOfMaps(output.Items, &items)
	if err != nil {
		return ResponseProxy(500, &MessageResponse{}, err)
	}

	log.Printf("Request: %s\n", event.QueryStringParameters)
	if err != nil {
		return ResponseProxy(500, &MessageResponse{}, err)
	}

	return ResponseProxy(200, &DataResponse{Data: items}, nil)
}

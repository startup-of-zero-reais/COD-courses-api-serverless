package main

import "github.com/aws/aws-lambda-go/events"

type (
	Stringify interface {
		String() string
	}
)

func Json(s Stringify) string {
	if s == nil {
		return ""
	}

	return s.String()
}

func ResponseProxy(statusCode int, body Stringify, err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode:      statusCode,
		Body:            Json(body),
		IsBase64Encoded: false,
	}, err
}

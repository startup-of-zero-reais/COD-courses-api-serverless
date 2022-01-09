package common

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
)

type (
	Stringify interface {
		String() string
	}
)

func ToString(s interface{}) string {
	bytes, _ := json.Marshal(s)
	return string(bytes)
}

func Json(s Stringify) string {
	if s == nil {
		return ""
	}

	return s.String()
}

func ResponseProxy(statusCode int, body Stringify, err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       Json(body),
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
		IsBase64Encoded: false,
	}, err
}

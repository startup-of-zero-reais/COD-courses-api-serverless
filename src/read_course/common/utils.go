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
	response := events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       Json(body),
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "Origin, Content-Type, user, X-Api-key, Authorization",
		},
		IsBase64Encoded: false,
	}

	if err != nil {
		response.Body = err.Error()
	}

	return response, nil
}

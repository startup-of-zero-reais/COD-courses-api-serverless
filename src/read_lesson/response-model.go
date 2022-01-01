package main

import "encoding/json"

type (
	MessageResponse struct {
		Message string `json:"message"`
	}

	DataResponse struct {
		Data interface{} `json:"data"`
	}
)

func (m *MessageResponse) String() string {
	bytes, _ := json.Marshal(m)
	return string(bytes)
}

func (d *DataResponse) String() string {
	bytes, _ := json.Marshal(d)
	return string(bytes)
}

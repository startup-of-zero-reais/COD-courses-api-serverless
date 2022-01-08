package common

type (
	MessageResponse struct {
		Message string `json:"message"`
	}

	DataResponse struct {
		Data interface{} `json:"data,omitempty"`
	}
)

func NewMessage(m string) *MessageResponse {
	return &MessageResponse{
		Message: m,
	}
}

func NewDataResponse(data interface{}) *DataResponse {
	return &DataResponse{
		Data: data,
	}
}

func (m *MessageResponse) String() string {
	return ToString(m)
}

func (d *DataResponse) String() string {
	return ToString(d)
}

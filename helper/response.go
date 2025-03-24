package helper

type Response struct {
	// Code    int    `json:"code"`
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

func CustomResponse(status bool, message string) Response {
	return Response{
		Status:  status,
		Message: message,
	}
}

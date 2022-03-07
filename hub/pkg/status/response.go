package status

type Response struct {
	Message string `json:"message"`
}

func Error(err error) Response {
	return Response{
		Message: err.Error(),
	}
}

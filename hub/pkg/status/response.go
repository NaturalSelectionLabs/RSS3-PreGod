package status

const (
	Version = "rss3.io/version/v0.4.0"
)

type Response struct {
	Version string      `json:"version"`
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func Data(data interface{}) Response {
	return Response{
		Version: Version,
		Data:    data,
	}
}

func Error(err error) Response {
	return Response{
		Version: Version,
		Message: err.Error(),
	}
}

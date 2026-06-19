package errors

type AppError interface {
	error
	Code() string
	Status() int
	Message() string
	ToJSON() []byte
}

type errorInfo struct {
	message string
	code    string
	status  int
}

package web

type ErrorTypeErr string

const (
	MESSAGE_BAD_REQUEST ErrorTypeErr = "Failed to send message, please try again later"
)

func (e ErrorTypeErr) Error() map[string]any {
	return map[string]any{
		"error": e,
	}
}

package emailerrors

import "fmt"

const (
	errCodeEmailSendingFailed = "ERR_EMAIL_SENDING_FAILED"
)

type EmailError struct {
	Code string
	Err  error
}

func (e *EmailError) Error() string {
	return fmt.Sprintf("Code: %s, error: %s", e.Code, e.Err)
}

func (e *EmailError) Wrap(err error) *EmailError {
	return &EmailError{
		Code: e.Code,
		Err:  err,
	}
}

func NewEmailSendingFailedError() *EmailError {
	return &EmailError{
		Code: errCodeEmailSendingFailed,
	}
}

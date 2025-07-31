package currency

import (
	"fmt"
)

type errCode int

const (
	CodeNotFound errCode = iota
	CodeInvalid
	CodeUnknown
)

type errOrigin string

const (
	cmcOrigin          errOrigin = "CmcClient"
	currencyRepoOrigin errOrigin = "CurrencyStorage"
	serviceOrigin      errOrigin = "Service"
)

type ServiceError struct {
	Code       errCode
	Origin     errOrigin
	CausingErr error
	Message    string
}

func NewServiceError(
	causingErr error,
	origin errOrigin,
	message string,
	code errCode,
) *ServiceError {
	return &ServiceError{
		CausingErr: causingErr,
		Origin:     origin,
		Message:    message,
		Code:       code,
	}
}

func (e *ServiceError) Error() string {
	return fmt.Sprintf(
		"error in service: %s, origin: %s, cause: %v",
		e.Message,
		e.Origin,
		e.CausingErr.Error(),
	)
}

func (e *ServiceError) Unwrap() error {
	return e.CausingErr
}

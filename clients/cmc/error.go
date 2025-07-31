package cmc

import "fmt"

type errCode int

type CmcClientError struct {
	Message    string
	StatusCode int
}

func NewCmcClientError(message string, httpCode int) *CmcClientError {
	return &CmcClientError{
		Message:    message,
		StatusCode: httpCode,
	}
}

func (e *CmcClientError) Error() string {
	return fmt.Sprintf("error performing cmc request: %v, status code: %v", e.Message, e.StatusCode)
}

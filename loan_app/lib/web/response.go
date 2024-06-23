package web

import "net/http"

const (
	UnauthorizedRequest = "unauthorized"
	BadRequest          = "bad_request"
	InternalServerError = "internal_server_error"
	NotFoundError       = "not_found"
	APIVersion1         = "1.0"
)

var (
	ErrUnauthorizedRequest = func(desc string) *APIResponse {
		return NewAPIErrorResponse(UnauthorizedRequest, desc, http.StatusUnauthorized)
	}
	ErrNotFound = func(desc string) *APIResponse {
		return NewAPIErrorResponse(NotFoundError, desc, http.StatusNotFound)
	}
	ErrBadRequest = func(desc string) *APIResponse {
		return NewAPIErrorResponse(BadRequest, desc, http.StatusBadRequest)
	}
	// internal error message should not be exposed to client. We can log such cases as error logs.
	ErrInternalServerError = NewAPIErrorResponse(InternalServerError, "something went wrong", http.StatusInternalServerError)
)

type APIResponse struct {
	APIVersion string       `json:"api_version"`
	Data       interface{}  `json:"data,omitempty"`
	Error      *errorDetail `json:"error,omitempty"`
	Success    bool         `json:"success"`
	httpStatus int
}

type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (a *APIResponse) HTTPStatusCode() int {
	return a.httpStatus
}

func NewAPISuccessResponse(data interface{}, httpCode int) *APIResponse {
	return &APIResponse{
		APIVersion: APIVersion1,
		Success:    true,
		httpStatus: httpCode,
		Data:       data,
	}
}

func NewAPIErrorResponse(errCode string, desc string, httpCode int) *APIResponse {
	return &APIResponse{
		APIVersion: APIVersion1,
		Success:    false,
		httpStatus: httpCode,
		Error: &errorDetail{
			Code:    errCode,
			Message: desc,
		},
	}
}

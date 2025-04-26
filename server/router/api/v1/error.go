package v1

type ErrorCode int32

const (
	Unimplemented        ErrorCode = 0
	InvalidRequest       ErrorCode = 1
	Internal             ErrorCode = 2
	Unauthenticated      ErrorCode = 3
	Unknown              ErrorCode = 4
)
type ErrorResponse struct {
	Code    ErrorCode    `json:"code"`
    Message string       `json:"message"`
}

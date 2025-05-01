package v1

type ErrorCode int32

const (
	Unimplemented        ErrorCode = 0
	InvalidRequest       ErrorCode = 1
	Internal             ErrorCode = 2
	Unauthenticated      ErrorCode = 3
	NotFound             ErrorCode = 4
	PermissionDenied     ErrorCode = 5
	Unknown              ErrorCode = 6
)
type ErrorResponse struct {
	Code    ErrorCode    `json:"code"`
    Message string       `json:"message"`
}

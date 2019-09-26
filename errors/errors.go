package errors

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

// Error
type Error struct {
	Op     Op
	Kind   Kind
	Level  Level
	Fields []Field
	Err    error
}

// Op operation string
type Op string

// Kind
type Kind int

// Level
type Level int

const (
	KindUnknown Kind = iota
	KindInvalid
	KindPermission
	KindExist
	KindNotExist
	KindInternal
)

const (
	LevelError Level = iota
	LevelWarn
	LevelInfo
	LevelDebug
)

// Field
type Field struct {
	Key   string
	Value string
}

// Error
func (e *Error) Error() string {
	return ""
}

// String
func (k Kind) String() string {
	switch k {
	case KindUnknown:
		return "unknown error"
	case KindInvalid:
		return "invalid error"
	case KindPermission:
		return "permission error"
	case KindExist:
		return "exit error"
	case KindNotExist:
		return "not exit error"
	case KindInternal:
		return "internal error"
	}

	return "unknown error kind"
}

// GrpcCode
func (k Kind) GrpcCode() codes.Code {
	switch k {
	case KindUnknown:
		return codes.Unknown
	case KindInvalid:
		return codes.InvalidArgument
	case KindPermission:
		return codes.PermissionDenied
	case KindExist:
		return codes.AlreadyExists
	case KindNotExist:
		return codes.NotFound
	case KindInternal:
		return codes.Internal
	}

	return codes.Unknown
}

// HttpCode
func (k Kind) HttpStatus() int {
	switch k {
	case KindUnknown:
		return http.StatusInternalServerError
	case KindInvalid:
		return http.StatusBadRequest
	case KindPermission:
		return http.StatusUnauthorized
	case KindExist:
		return http.StatusConflict
	case KindNotExist:
		return http.StatusNotFound
	case KindInternal:
		return http.StatusInternalServerError
	}

	return http.StatusInternalServerError
}

// E
func E(args ...interface{}) error {
	if len(args) == 0 {
		panic("call to errors.E with no arguments")
	}

	e := &Error{Fields: []Field{}}
	for _, arg := range args {
		switch arg := arg.(type) {
		case Op:
			e.Op = arg
		case Kind:
			e.Kind = arg
		case Level:
			e.Level = arg
		case Field:
			e.Fields = append(e.Fields, arg)
		case error:
			e.Err = arg
		default:
			panic("call with unknown arg to errors.E")
		}
	}

	return e
}

// Ops
func Ops(e *Error) []Op {
	res := []Op{e.Op}

	subErr, ok := e.Err.(*Error)
	if !ok {
		return res
	}

	res = append(res, Ops(subErr)...)

	return res
}

// Fields
func Fields(e *Error) []Field {
	var res []Field
	res = append(res, e.Fields...)

	subErr, ok := e.Err.(*Error)
	if !ok {
		return res
	}

	res = append(res, Fields(subErr)...)

	return res
}

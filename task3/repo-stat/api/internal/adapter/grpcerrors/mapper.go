package grpcerrors

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ToHTTP(err error) (int, string) {
	if err == nil {
		return 500, "unknown error"
	}

	st, ok := status.FromError(err)
	if ok {
		code, _ := mapCode(st.Code())
		return code, st.Message()
	}

	return 500, err.Error()
}

func mapCode(code codes.Code) (int, string) {
	switch code {
	case codes.OK:
		return 200, "ok"
	case codes.InvalidArgument, codes.FailedPrecondition:
		return 400, "bad request"
	case codes.NotFound:
		return 404, "not found"
	case codes.AlreadyExists:
		return 409, "conflict"
	case codes.PermissionDenied:
		return 403, "forbidden"
	case codes.Unauthenticated:
		return 401, "unauthorized"
	case codes.Unavailable:
		return 503, "service unavailable"
	case codes.DeadlineExceeded:
		return 504, "gateway timeout"
	default:
		return 500, fmt.Sprintf("internal error: %v", code)
	}
}

package interceptor

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/itsLeonB/cocoon/internal/logging"
	"github.com/itsLeonB/ezutil"
	"github.com/rotisserie/eris"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewErrorInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			return resp, nil
		}

		// Already a gRPC status → just return
		if _, ok := status.FromError(err); ok {
			return resp, err
		}

		appErr := constructAppError(err)
		code := httpStatusToGRPCCode(appErr.HttpStatusCode)
		return resp, status.Error(code, appErr.Message)
	}
}

func httpStatusToGRPCCode(httpStatus int) codes.Code {
	switch httpStatus {
	case http.StatusOK:
		return codes.OK
	case http.StatusBadRequest:
		return codes.InvalidArgument
	case http.StatusUnauthorized:
		return codes.Unauthenticated
	case http.StatusForbidden:
		return codes.PermissionDenied
	case http.StatusNotFound:
		return codes.NotFound
	case http.StatusConflict:
		return codes.AlreadyExists
	case http.StatusTooManyRequests:
		return codes.ResourceExhausted
	case http.StatusNotImplemented:
		return codes.Unimplemented
	case http.StatusServiceUnavailable:
		return codes.Unavailable
	case http.StatusGatewayTimeout:
		return codes.DeadlineExceeded
	case http.StatusUnprocessableEntity:
		return codes.InvalidArgument // often mapped like this
	default:
		return codes.Internal
	}
}

func constructAppError(err error) ezutil.AppError {
	if appErr, ok := err.(ezutil.AppError); ok {
		return appErr
	}

	originalErr := eris.Unwrap(err)

	switch t := originalErr.(type) {
	case validator.ValidationErrors:
		var errors []string
		for _, e := range t {
			errors = append(errors, e.Error())
		}
		return ezutil.ValidationError(errors)

	case *json.SyntaxError:
		// replace config.MsgInvalidJson with whatever message you want
		return ezutil.BadRequestError("invalid json")

	default:
		// EOF error from json package is unexported; check both equality and string
		if originalErr == io.EOF || (originalErr != nil && originalErr.Error() == "EOF") {
			// replace config.MsgMissingBody with your message
			return ezutil.BadRequestError("missing request body")
		}

		logging.Logger.Errorf("unhandled error of type: %T\n", originalErr)
		logging.Logger.Errorf(eris.ToString(err, true))
		return ezutil.InternalServerError()
	}
}

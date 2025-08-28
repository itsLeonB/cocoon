package interceptor

import (
	"context"
	"time"

	"github.com/itsLeonB/ezutil/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// NewLoggerInterceptor logs incoming requests, responses, durations, and errors.
func NewLoggerInterceptor(logger ezutil.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()

		// Call handler
		resp, err := handler(ctx, req)

		// Duration
		elapsed := time.Since(start)

		// Extract gRPC status code (if error)
		st, _ := status.FromError(err)

		if err != nil {
			logger.Errorf(
				"[gRPC] method=%s duration=%s status=%s error=%v",
				info.FullMethod,
				elapsed,
				st.Code(),
				st.Message(),
			)
		} else {
			logger.Infof(
				"[gRPC] method=%s duration=%s status=OK",
				info.FullMethod,
				elapsed,
			)
		}

		return resp, err
	}
}

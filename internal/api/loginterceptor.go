package api

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func loggingUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	startingTime := time.Now()
	resp, err := handler(ctx, req)

	rpcStatus, _ := status.FromError(err)
	logGRPCRequest(ctx, info.FullMethod, rpcStatus, time.Since(startingTime))

	return resp, err
}

func logGRPCRequest(ctx context.Context, method string, rpcStatus *status.Status, elapsedTime time.Duration) {
	lvl := levelByGRPCCode(rpcStatus.Code())
	slog.Log(ctx, lvl, "Request handled", slog.String("handler_type", "gRPC"),
		slog.String("method", method), slog.String("response_code", rpcStatus.Code().String()), slog.String("error_message", rpcStatus.Message()), slog.Duration("elapsed_time", elapsedTime))
}

func levelByGRPCCode(code codes.Code) slog.Level {
	if code == http.StatusInternalServerError {
		return slog.LevelError
	}

	return slog.LevelInfo
}

package interceptors

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func isValidToken(authHeader string) bool {
	// Format: "Bearer secret123"
	parts := strings.Split(authHeader, " ")
	return len(parts) == 2 && parts[0] == "Bearer" && parts[1] == "secret123"
}

func UnaryAuthInterceptors(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Missing Metadata")
	}

	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 || !isValidToken(authHeaders[0]) {
		return nil, status.Error(codes.Unauthenticated, "Invalid or missing token")
	}

	return handler(ctx, req)
}


func StreamAuthInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	md, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		return status.Errorf(codes.Unauthenticated, "Missing metadata")
	}

	authHeader := md["authorization"]
	if len(authHeader) == 0 || authHeader[0] != "Bearer super-secret-token" {
		return status.Errorf(codes.Unauthenticated, "Invalid or missing token")
	}

	// ✅ Token is valid, continue with stream
	return handler(srv, ss)
}

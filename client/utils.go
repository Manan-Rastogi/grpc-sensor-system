package main

import (
	"context"
	"google.golang.org/grpc/metadata"
)

// Helper function
func withAuthCtx() context.Context {
	return metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{
		"authorization": "Bearer secret123",
	}))
}

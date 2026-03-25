package grpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"log"
	"os"
	"service/internal/pkg/grpc/example"
	"time"
)

var (
	TestingRPCClient example.TestingServiceClient

	TestingRPCActive bool
)

// TODO: Hanya contoh. nanti langsung hapus saja
func InitTestingRPCClient() func() {
	addr := os.Getenv("GRPC_TESTING_HOST")

	if addr != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		keepaliveParam := keepalive.ClientParameters{
			Time:                60 * time.Second,
			Timeout:             20 * time.Second,
			PermitWithoutStream: true,
		}

		conn, err := grpc.DialContext(ctx, addr,
			grpc.WithInsecure(),
			grpc.WithBlock(),
			grpc.WithKeepaliveParams(keepaliveParam),
		)
		if err != nil {
			log.Panicf("Did not connect to %s: %v", addr, err)
		}

		TestingRPCClient = example.NewTestingServiceClient(conn)

		TestingRPCActive = true

		cleanup := func() {
			cancel()
			conn.Close()
		}

		return cleanup
	}

	return func() {}
}

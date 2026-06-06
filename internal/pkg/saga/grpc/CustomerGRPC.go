package grpc

import (
	"context"
	"log"
	"os"
	"service/internal/pkg/grpc/customer"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var (
	CustomerRPCCClient customer.CustomerServiceClient

	CustomerRPCActive bool
)

func InitCustomerClient() func() {
	addr := os.Getenv("GRPC_CUSTOMER_HOST")

	if addr != "" {
		ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)

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
			log.Panic("Failed to connect to Customer gRPC: %w", err)
		}

		CustomerRPCCClient = customer.NewCustomerServiceClient(conn)
		CustomerRPCActive = true

		cleanup := func() {
			cancle()
			conn.Close()
		}

		return cleanup
	}
	return func() {}
}

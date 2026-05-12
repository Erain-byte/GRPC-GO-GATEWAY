//go:generate protoc --go_out=. --go-grpc_out=. proto/user.proto
//go:generate protoc --go_out=. --go-grpc_out=. proto/admin.proto
//go:generate protoc --go_out=. --go-grpc_out=. proto/ai.proto

package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("监听失败: %v", err)
	}

	s := grpc.NewServer()

	log.Println("gRPC Gateway 启动 :8080")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}

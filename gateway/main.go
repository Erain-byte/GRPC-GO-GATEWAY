//go:generate protoc --go_out=. --go-grpc_out=. proto/user.proto
//go:generate protoc --go_out=. --go-grpc_out=. proto/admin.proto
//go:generate protoc --go_out=. --go-grpc_out=. proto/ai.proto

package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Erain-byte/GRPC-GO-GATEWAY/gateway/config"
	"google.golang.org/grpc"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("配置加载失败: %v", err)
	}

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("监听失败: %v", err)
	}

	s := grpc.NewServer()

	log.Printf("gRPC Gateway 启动 %s", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}

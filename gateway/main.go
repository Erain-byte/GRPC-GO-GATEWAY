//go:generate protoc --go_out=. --go-grpc_out=. proto/user.proto
//go:generate protoc --go_out=. --go-grpc_out=. proto/admin.proto
//go:generate protoc --go_out=. --go-grpc_out=. proto/ai.proto

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Erain-byte/GRPC-GO-GATEWAY/gateway/config"
	"github.com/Erain-byte/GRPC-GO-GATEWAY/gateway/internal/forwarder"
	"google.golang.org/grpc"
)

func main() {
	// 1. 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("配置加载失败: %v", err)
	}

	// 3. 初始化服务注册中心（自动完成所有转发器的创建）
	registry, err := forwarder.NewServiceRegistry(cfg)
	if err != nil {
		log.Fatalf("初始化服务注册中心失败: %v", err)
	}

	// 确保退出时关闭所有连接
	if registry.AdminFwd != nil {
		defer registry.AdminFwd.Close()
	}

	// 4. 创建 Context 用于优雅退出
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 3. 启动 gRPC 服务
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	lis, err := net.Listen("tcp", addr) // 监听端口
	if err != nil {
		log.Fatalf("监听失败: %v", err)
	}

	s := grpc.NewServer()

	// 自动注册所有已配置的微服务转发器
	registry.RegisterAll(s)

	log.Printf("gRPC Gateway 启动 %s", addr)

	// 在协程中启动服务
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Printf("gRPC 服务异常: %v", err)
		}
	}()

	// 4. 监听系统退出信号 (Ctrl+C)
	quit := make(chan os.Signal, 1)                      // 创建一个信号通道接受退出信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 监听退出信号
	<-quit

	log.Println("收到退出信号，正在关闭网关服务...")
	cancel()

	// 5. 优雅停止 gRPC 服务
	s.GracefulStop()
	log.Println("网关服务已安全退出")
}

package main

import (
	"fmt"
	"log"
	"net"
	"sync"

	"admin/config"
	"admin/internal/database"
	"admin/internal/redis"
	"admin/internal/service"
	"admin/router"

	"github.com/Erain-byte/GRPC-GO-GATEWAY/gateway/pb/admin"

	"google.golang.org/grpc"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 初始化数据库
	database.InitDB(cfg)

	// 初始化 Redis
	redis.InitRedis(cfg)

	// 创建 WaitGroup 用于等待所有服务结束
	var wg sync.WaitGroup

	// 启动 HTTP 服务
	wg.Add(1)
	go func() {
		defer wg.Done()
		startHTTPServer(cfg)
	}()

	// 启动 gRPC 服务
	wg.Add(1)
	go func() {
		defer wg.Done()
		startGRPCServer(cfg)
	}()

	// 等待所有服务结束
	wg.Wait()
}

// startHTTPServer 启动 HTTP 服务
func startHTTPServer(cfg *config.Config) {
	r := router.SetupRouter()
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("HTTP 服务启动 %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("HTTP 服务启动失败: %v", err)
	}
}

// startGRPCServer 启动 gRPC 服务
func startGRPCServer(cfg *config.Config) {
	// 监听端口
	addr := fmt.Sprintf(":%s", cfg.Server.GrpcPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("gRPC 监听失败: %v", err)
	}

	// 创建 gRPC 服务器
	s := grpc.NewServer()

	// 注册 Admin 服务
	admin.RegisterAdminServiceServer(s, &service.AdminServer{})

	log.Printf("gRPC 服务启动 %s", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("gRPC 服务启动失败: %v", err)
	}
}

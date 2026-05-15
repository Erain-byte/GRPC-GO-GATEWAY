package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"admin/config"
	"admin/internal/consul"
	"admin/internal/database"
	"admin/internal/redis"
	"admin/internal/service"
	"admin/router"

	"github.com/Erain-byte/GRPC-GO-GATEWAY/gateway/pb/admin"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	// 1. 加载配置
	cfg := config.LoadConfig()

	// 2. 初始化核心依赖（失败则直接终止）
	database.InitDB(cfg)
	redis.InitRedis(cfg)

	// 3. 注册服务到 Consul
	if err := consul.RegisterService(cfg); err != nil {
		log.Printf("警告: Consul 注册失败，服务将继续运行但无法被发现: %v", err)
	} else {
		// 确保退出时注销服务
		defer consul.DeregisterService(cfg)
	}

	// 3. 创建 Context 用于优雅退出
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 4. 启动 gRPC 服务
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Server.GrpcPort))
	if err != nil {
		log.Fatalf("gRPC 监听失败: %v", err)
	}
	grpcServer := grpc.NewServer()
	admin.RegisterAdminServiceServer(grpcServer, &service.AdminServer{})

	// 注册 gRPC 标准健康检查服务（Consul 需要这个来判定服务是否存活）
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	// 设置整体服务状态为 SERVING
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	log.Printf("gRPC 服务启动 :%s", cfg.Server.GrpcPort)

	// 在协程中启动 gRPC 服务
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("gRPC 服务异常: %v", err)
		}
	}()

	// 5. 启动 HTTP 服务
	go startHTTPServer(ctx, cfg)

	// 6. 监听系统退出信号 (Ctrl+C)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // 阻塞等待退出信号，如 Ctrl+C

	log.Println("收到退出信号，正在关闭服务...")
	cancel() // 通知所有监听 ctx 的协程（如 HTTP）退出

	// 7. 优雅停止 gRPC 服务
	// GracefulStop 会等待所有正在处理的 RPC 请求完成，然后停止服务
	grpcServer.GracefulStop()
	log.Println("gRPC 服务已停止")

	// 8. 清理资源
	if database.DB != nil {
		sqlDB, _ := database.DB.DB()
		sqlDB.Close()
	}
	if redis.Rdb != nil {
		redis.Rdb.Close()
	}
	log.Println("服务已安全退出")
}

// startHTTPServer 启动 HTTP 服务
func startHTTPServer(ctx context.Context, cfg *config.Config) {
	r := router.SetupRouter()
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("HTTP 服务启动 %s", addr)

	// 使用 Shutdown 方法实现优雅退出
	server := &http.Server{Addr: addr, Handler: r}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP 服务启动失败: %v", err)
		}
	}()

	// 等待 Context 取消
	<-ctx.Done()
	log.Println("正在关闭 HTTP 服务...")
	// 给 HTTP 服务 5 秒时间处理完剩余请求
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP 服务强制关闭: %v", err)
	}
}

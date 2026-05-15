package forwarder

import (
	"context"
	"fmt"
	"log"

	adminPb "github.com/Erain-byte/GRPC-GO-GATEWAY/gateway/pb/admin"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// AdminForwarder Admin 服务转发器
type AdminForwarder struct {
	adminPb.UnimplementedAdminServiceServer
	client adminPb.AdminServiceClient
	conn   *grpc.ClientConn
}

// NewAdminForwarder 创建并初始化 Admin 转发器（通过 Consul 发现）
func NewAdminForwarder(consulAddr, serviceName string) (*AdminForwarder, error) {
	// 1. 连接 Consul
	consulConfig := api.DefaultConfig()
	consulConfig.Address = consulAddr
	client, err := api.NewClient(consulConfig)
	if err != nil {
		return nil, fmt.Errorf("创建 Consul 客户端失败: %v", err)
	}

	// 2. 从 Consul 获取健康的服务实例
	entries, _, err := client.Health().Service(serviceName, "grpc", true, nil)
	if err != nil || len(entries) == 0 {
		return nil, fmt.Errorf("在 Consul 中未找到可用的 %s 服务", serviceName)
	}

	// 3. 选取第一个健康的实例（生产环境可做负载均衡）
	service := entries[0].Service
	addr := fmt.Sprintf("%s:%d", service.Address, service.Port)

	log.Printf("Gateway 通过 Consul 发现 Admin 服务地址: %s", addr)

	// 4. 建立 gRPC 连接
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &AdminForwarder{
		client: adminPb.NewAdminServiceClient(conn),
		conn:   conn,
	}, nil
}

// Close 关闭连接
func (f *AdminForwarder) Close() {
	if f.conn != nil {
		f.conn.Close()
	}
}

// --- 实现 AdminServiceServer 接口的所有方法 ---

func (f *AdminForwarder) HealthCheck(ctx context.Context, req *adminPb.HealthCheckRequest) (*adminPb.HealthCheckResponse, error) {
	return f.client.HealthCheck(ctx, req)
}

func (f *AdminForwarder) Login(ctx context.Context, req *adminPb.LoginRequest) (*adminPb.LoginResponse, error) {
	return f.client.Login(ctx, req)
}

func (f *AdminForwarder) GetAdmin(ctx context.Context, req *adminPb.GetAdminRequest) (*adminPb.Admin, error) {
	return f.client.GetAdmin(ctx, req)
}

func (f *AdminForwarder) UpdateAdmin(ctx context.Context, req *adminPb.UpdateAdminRequest) (*adminPb.Admin, error) {
	return f.client.UpdateAdmin(ctx, req)
}

func (f *AdminForwarder) UploadAvatar(ctx context.Context, req *adminPb.UploadAvatarRequest) (*adminPb.UploadResponse, error) {
	return f.client.UploadAvatar(ctx, req)
}

func (f *AdminForwarder) ListAdmins(ctx context.Context, req *adminPb.ListAdminsRequest) (*adminPb.ListAdminsResponse, error) {
	return f.client.ListAdmins(ctx, req)
}

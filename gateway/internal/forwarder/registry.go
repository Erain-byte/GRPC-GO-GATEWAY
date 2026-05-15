package forwarder

import (
	"log"

	"github.com/Erain-byte/GRPC-GO-GATEWAY/gateway/config"
	adminPb "github.com/Erain-byte/GRPC-GO-GATEWAY/gateway/pb/admin"
	"google.golang.org/grpc"
)

// ServiceRegistry 微服务转发器注册中心
type ServiceRegistry struct {
	AdminFwd *AdminForwarder
	// UserFwd  *UserForwarder // 未来扩展
}

// NewServiceRegistry 创建并初始化所有转发器
func NewServiceRegistry(cfg *config.Config) (*ServiceRegistry, error) {
	r := &ServiceRegistry{}

	// 1. 初始化 Admin 转发器（通过 Consul 发现）
	if cfg.Consul.Addr != "" && cfg.Services.AdminServiceName != "" {
		adminFwd, err := NewAdminForwarder(cfg.Consul.Addr, cfg.Services.AdminServiceName)
		if err != nil {
			return nil, err
		}
		r.AdminFwd = adminFwd
	}

	return r, nil
}

// RegisterAll 统一注册所有已初始化的转发器到 gRPC 服务器
func (r *ServiceRegistry) RegisterAll(s *grpc.Server) {
	if r.AdminFwd != nil {
		adminPb.RegisterAdminServiceServer(s, r.AdminFwd)
		log.Println("[Registry] 已注册 Admin 服务转发器")
	}

	// 未来扩展示例：
	// if r.UserFwd != nil {
	//     userPb.RegisterUserServiceServer(s, r.UserFwd)
	//     log.Println("[Registry] 已注册 User 服务转发器")
	// }
}

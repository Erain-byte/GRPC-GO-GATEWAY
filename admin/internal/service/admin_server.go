package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Erain-byte/GRPC-GO-GATEWAY/gateway/pb/admin"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// AdminServer 实现 AdminService 接口
type AdminServer struct {
	admin.UnimplementedAdminServiceServer
}

// HealthCheck 健康检查
func (s *AdminServer) HealthCheck(ctx context.Context, req *admin.HealthCheckRequest) (*admin.HealthCheckResponse, error) {
	return &admin.HealthCheckResponse{
		Status:    "serving",
		Timestamp: time.Now().Unix(),
		Version:   "1.0.0",
	}, nil
}

// Login 管理员登录
func (s *AdminServer) Login(ctx context.Context, req *admin.LoginRequest) (*admin.LoginResponse, error) {
	// TODO: 验证用户名密码（从数据库查询）
	if req.Username != "admin" || req.Password != "123456" {
		return nil, fmt.Errorf("用户名或密码错误")
	}

	// TODO: 生成 JWT Token
	token := "mock_jwt_token_" + time.Now().Format("20060102150405")

	// 构造管理员信息
	now := timestamppb.Now()
	adminInfo := &admin.Admin{
		Id:          1,
		Username:    req.Username,
		Email:       "admin@example.com",
		Phone:       "13800138000",
		AvatarUrl:   "",
		Role:        admin.AdminRole_ROLE_SUPER_ADMIN,
		Status:      admin.AdminStatus_STATUS_ACTIVE,
		Permissions: []string{"user.manage", "order.manage", "system.config"},
		LastLoginAt: now,
		LastLoginIp: "127.0.0.1", // TODO: 从请求中获取真实IP
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return &admin.LoginResponse{
		Token: token,
		Admin: adminInfo,
	}, nil
}

// GetAdmin 获取管理员信息
func (s *AdminServer) GetAdmin(ctx context.Context, req *admin.GetAdminRequest) (*admin.Admin, error) {
	// TODO: 从数据库查询
	now := timestamppb.Now()
	return &admin.Admin{
		Id:          req.Id,
		Username:    "admin",
		Email:       "admin@example.com",
		Phone:       "13800138000",
		AvatarUrl:   "",
		Role:        admin.AdminRole_ROLE_SUPER_ADMIN,
		Status:      admin.AdminStatus_STATUS_ACTIVE,
		Permissions: []string{"user.manage", "order.manage", "system.config"},
		LastLoginAt: now,
		LastLoginIp: "127.0.0.1",
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// UpdateAdmin 更新管理员信息
func (s *AdminServer) UpdateAdmin(ctx context.Context, req *admin.UpdateAdminRequest) (*admin.Admin, error) {
	// TODO: 更新数据库
	now := timestamppb.Now()
	return &admin.Admin{
		Id:          req.Id,
		Username:    req.Username,
		Email:       req.Email,
		Phone:       "",
		AvatarUrl:   "",
		Role:        admin.AdminRole_ROLE_ADMIN,
		Status:      admin.AdminStatus_STATUS_ACTIVE,
		Permissions: []string{"user.manage"},
		LastLoginAt: now,
		LastLoginIp: "127.0.0.1",
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// UploadAvatar 上传头像
func (s *AdminServer) UploadAvatar(ctx context.Context, req *admin.UploadAvatarRequest) (*admin.UploadResponse, error) {
	// TODO: 验证文件、保存到OSS
	return &admin.UploadResponse{
		FileUrl:  req.AvatarUrl,
		FileId:   "file_123",
		FileType: "avatar",
		FileSize: 102400,
	}, nil
}

// ListAdmins 列出管理员
func (s *AdminServer) ListAdmins(ctx context.Context, req *admin.ListAdminsRequest) (*admin.ListAdminsResponse, error) {
	// TODO: 从数据库分页查询
	now := timestamppb.Now()
	admins := []*admin.Admin{
		{
			Id:          1,
			Username:    "admin",
			Email:       "admin@example.com",
			Role:        admin.AdminRole_ROLE_SUPER_ADMIN,
			Status:      admin.AdminStatus_STATUS_ACTIVE,
			Permissions: []string{"user.manage", "order.manage"},
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			Id:          2,
			Username:    "operator",
			Email:       "operator@example.com",
			Role:        admin.AdminRole_ROLE_OPERATOR,
			Status:      admin.AdminStatus_STATUS_ACTIVE,
			Permissions: []string{"order.manage"},
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	return &admin.ListAdminsResponse{
		Admins: admins,
		Total:  int32(len(admins)),
	}, nil
}

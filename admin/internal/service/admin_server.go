package service

import (
	"context"
	"time"

	"admin/internal/database"
	"admin/internal/models"
	"admin/internal/utils"

	"github.com/Erain-byte/GRPC-GO-GATEWAY/gateway/pb/admin"
	pb "github.com/Erain-byte/GRPC-GO-GATEWAY/gateway/pb/admin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AdminServer 实现 AdminService 接口
type AdminServer struct {
	pb.UnimplementedAdminServiceServer
}

// HealthCheck 健康检查
func (s *AdminServer) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &admin.HealthCheckResponse{
		Status:    "serving",
		Timestamp: time.Now().Unix(),
		Version:   "1.0.0",
	}, nil
}

// Login 管理员登录
func (s *AdminServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// 1. 参数校验（由 protoc-gen-validate 生成）
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "参数校验失败: %v", err)
	}

	// 2. 从数据库查询用户
	var adminModel models.Admin
	result := database.DB.Where("username = ?", req.Username).First(&adminModel)
	if result.Error != nil {
		return nil, status.Errorf(codes.Unauthenticated, "用户名或密码错误")
	}

	// 3. 校验密码（bcrypt 比对）
	if !utils.CheckPassword(req.Password, adminModel.Password) {
		return nil, status.Errorf(codes.Unauthenticated, "用户名或密码错误")
	}

	// 4. 检查账号状态
	if adminModel.Status != 1 { // 1 代表正常
		return nil, status.Errorf(codes.PermissionDenied, "账号已被禁用")
	}

	// 5. 生成 JWT Token
	token, err := utils.GenerateToken(uint(adminModel.ID), adminModel.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Token 生成失败: %v", err)
	}

	// 6. 更新最后登录信息
	now := time.Now()
	database.DB.Model(&adminModel).Updates(map[string]interface{}{
		"last_login_at": now,
		"last_login_ip": "127.0.0.1", // TODO: 从 gRPC metadata 获取真实 IP
	})

	// 7. 构造返回结果
	loginTime := timestamppb.New(now)
	createTime := timestamppb.New(adminModel.CreatedAt)
	updateTime := timestamppb.New(adminModel.UpdatedAt)

	return &pb.LoginResponse{
		Token: token,
		Admin: &pb.Admin{
			Id:          int64(adminModel.ID),
			Username:    adminModel.Username,
			Email:       adminModel.Email,
			Phone:       adminModel.Phone,
			AvatarUrl:   adminModel.AvatarURL,
			Role:        pb.AdminRole(adminModel.Role),
			Status:      pb.AdminStatus(adminModel.Status),
			Permissions: []string{"all"}, // TODO: 根据角色加载权限
			LastLoginAt: loginTime,
			LastLoginIp: adminModel.LastLoginIP,
			CreatedAt:   createTime,
			UpdatedAt:   updateTime,
		},
	}, nil
}

// GetAdmin 获取管理员信息
func (s *AdminServer) GetAdmin(ctx context.Context, req *pb.GetAdminRequest) (*pb.Admin, error) {
	// TODO: 从数据库查询
	now := timestamppb.Now()
	return &pb.Admin{
		Id:          req.Id,
		Username:    "admin",
		Email:       "admin@example.com",
		Phone:       "13800138000",
		AvatarUrl:   "",
		Role:        pb.AdminRole_ROLE_SUPER_ADMIN,
		Status:      pb.AdminStatus_STATUS_ACTIVE,
		Permissions: []string{"user.manage", "order.manage", "system.config"},
		LastLoginAt: now,
		LastLoginIp: "127.0.0.1",
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// UpdateAdmin 更新管理员信息
func (s *AdminServer) UpdateAdmin(ctx context.Context, req *pb.UpdateAdminRequest) (*pb.Admin, error) {
	// TODO: 更新数据库
	now := timestamppb.Now()
	return &pb.Admin{
		Id:          req.Id,
		Username:    req.Username,
		Email:       req.Email,
		Phone:       "",
		AvatarUrl:   "",
		Role:        pb.AdminRole_ROLE_ADMIN,
		Status:      pb.AdminStatus_STATUS_ACTIVE,
		Permissions: []string{"user.manage"},
		LastLoginAt: now,
		LastLoginIp: "127.0.0.1",
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// UploadAvatar 上传头像
func (s *AdminServer) UploadAvatar(ctx context.Context, req *pb.UploadAvatarRequest) (*pb.UploadResponse, error) {
	// TODO: 验证文件、保存到OSS
	return &pb.UploadResponse{
		FileUrl:  req.AvatarUrl,
		FileId:   "file_123",
		FileType: "avatar",
		FileSize: 102400,
	}, nil
}

// ListAdmins 列出管理员
func (s *AdminServer) ListAdmins(ctx context.Context, req *pb.ListAdminsRequest) (*pb.ListAdminsResponse, error) {
	// TODO: 从数据库分页查询
	now := timestamppb.Now()
	admins := []*pb.Admin{
		{
			Id:          1,
			Username:    "admin",
			Email:       "admin@example.com",
			Role:        pb.AdminRole_ROLE_SUPER_ADMIN,
			Status:      pb.AdminStatus_STATUS_ACTIVE,
			Permissions: []string{"user.manage", "order.manage"},
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			Id:          2,
			Username:    "operator",
			Email:       "operator@example.com",
			Role:        pb.AdminRole_ROLE_OPERATOR,
			Status:      pb.AdminStatus_STATUS_ACTIVE,
			Permissions: []string{"order.manage"},
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	return &pb.ListAdminsResponse{
		Admins: admins,
		Total:  int32(len(admins)),
	}, nil
}

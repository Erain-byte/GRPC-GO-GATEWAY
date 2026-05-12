# Go Gateway 项目

Go 语言实现的 API 网关，用于迭代替换 Python gateway 项目。

## 项目架构

```
┌─────────────────────────────────────────────┐
│         Go Gateway (API 网关)                │
│  - HTTP/WebSocket 入口                       │
│  - 服务发现 (Consul)                         │
│  - 负载均衡                                  │
│  - gRPC 客户端                               │
│  - 认证/授权                                 │
└──────────────┬──────────────────────────────┘
               │ gRPC
     ┌─────────┼──────────┐
     ▼         ▼          ▼
  User      Admin       AI
  Service   Service    Service
  (Go)      (Go)      (Python)
     │         │          │
     └─────────┴──────────┘
               │
               ▼
         Consul (服务注册中心)
```

## 目录结构

```
gateway/
├── api/                    # API 层
│   ├── http/              # HTTP 路由和 Handler
│   └── websocket/         # WebSocket 接口
│
├── internal/              # 内部模块
│   ├── discovery/         # 服务发现（Consul 客户端）
│   ├── client/            # gRPC 客户端管理
│   └── middleware/        # 中间件（认证、跨域、日志）
│
├── pkg/                   # 公共工具包
├── config/                # 配置管理
├── proto/                 # Proto 文件（接口定义）
│   ├── user.proto
│   ├── admin.proto
│   └── ai.proto
├── pb/                    # 生成的 gRPC 代码
│   ├── user/              # User 服务代码
│   │   ├── user.pb.go
│   │   └── user_grpc.pb.go
│   ├── admin/             # Admin 服务代码
│   │   ├── admin.pb.go
│   │   └── admin_grpc.pb.go
│   └── ai/                # AI 服务代码
│       ├── ai.pb.go
│       └── ai_grpc.pb.go
├── docs/                  # 文档
├── main.go                # 入口文件
├── go.mod
├── go.sum
└── README.md
```

## 已完成的 Proto 接口

### UserService (user.proto)
- ✅ HealthCheck - 健康检查
- ✅ Login - 用户登录
- ✅ Logout - 用户登出
- ✅ Register - 用户注册
- ✅ GetUser - 获取用户信息
- ✅ CreateUser - 创建用户
- ✅ UpdateUser - 更新用户
- ✅ DeleteUser - 删除用户
- ✅ ListUsers - 列出用户

### AdminService (admin.proto)
- ✅ HealthCheck - 健康检查
- ✅ Login - 管理员登录
- ✅ GetAdmin - 获取管理员信息
- ✅ UpdateAdmin - 更新管理员信息
- ✅ ListAdmins - 列出管理员

### AIService (ai.proto)
- ✅ HealthCheck - 健康检查
- ✅ Chat - 聊天对话（一元 RPC）
- ✅ Embedding - 文本嵌入向量
- ✅ StreamChat - 流式聊天（服务端流式 RPC）

## 下一步开发计划

### 0. 实现 gRPC 服务端（最优先）⭐⭐⭐⭐⭐
- [ ] 创建 Admin Service 项目结构
- [ ] 实现 HealthCheck 方法
- [ ] 实现 Login 方法（JWT Token 生成）
- [ ] 实现 GetAdmin/UpdateAdmin/ListAdmins 方法
- [ ] 实现 UploadAvatar 方法
- [ ] 启动 Admin gRPC 服务（端口 50051）
- [ ] 测试连通性（grpcurl）
- [ ] 创建 User Service 项目结构
- [ ] 实现 User 服务所有方法
- [ ] 集成数据库（GORM + MySQL）

### 1. 实现服务发现模块
- [ ] 安装 Consul 依赖
- [ ] 实现 Consul 客户端封装
- [ ] 服务注册功能（启动时自动注册）
- [ ] 服务发现功能（Gateway 查询服务地址）
- [ ] 健康检查集成
- [ ] 服务下线注销

### 2. 实现 gRPC 客户端管理
- [ ] 连接池管理
- [ ] User 服务客户端封装
- [ ] Admin 服务客户端封装
- [ ] AI 服务客户端封装
- [ ] 自动重连机制
- [ ] 负载均衡（随机/轮询）

### 3. 实现拦截器（Interceptors）
- [ ] 日志拦截器（记录请求/响应）
- [ ] 错误处理拦截器（统一错误格式）
- [ ] 认证拦截器（验证 JWT Token）
- [ ] 限流拦截器
- [ ] 超时控制拦截器

### 4. 实现 HTTP 路由（Gateway）
- [ ] Gin 框架集成
- [ ] 用户认证路由（登录/注册）
- [ ] 聊天接口路由（一元/流式）
- [ ] 管理员路由
- [ ] 文件上传路由（获取 OSS 签名）
- [ ] 中间件（认证、跨域、日志）

### 5. 配置管理
- [ ] .env 文件支持
- [ ] 配置加载模块
- [ ] 环境变量读取
- [ ] 多环境配置（dev/prod）

### 6. 实现 WebSocket 支持
- [ ] WebSocket 握手
- [ ] 消息转发（Gateway <-> AI Service）
- [ ] 心跳检测
- [ ] 断线重连

### 7. 数据库集成
- [ ] 选择数据库（MySQL/PostgreSQL）
- [ ] 安装 GORM
- [ ] 定义数据模型（Admin, User）
- [ ] 实现 CRUD 操作
- [ ] 数据库迁移

### 8. 认证与授权
- [ ] JWT Token 生成
- [ ] JWT Token 验证
- [ ] 权限检查（基于 AdminRole）
- [ ] Token 刷新机制

### 9. 日志系统
- [ ] 集成 zap/logrus
- [ ] 结构化日志输出
- [ ] 日志级别控制
- [ ] 日志文件滚动

### 10. 测试与优化
- [ ] 单元测试
- [ ] 集成测试
- [ ] 性能测试
- [ ] 错误处理完善


## 技术栈

- **语言**: Go 1.26.2
- **Web 框架**: Gin（计划）
- **RPC**: gRPC + Protobuf
- **服务发现**: Consul（计划）
- **配置管理**: godotenv（计划）
- **认证**: JWT（计划）

## 服务说明

### Gateway (API 网关)
- **端口**: 8080 (gRPC)
- **功能**:
  - HTTP/WebSocket 入口
  - 服务发现和负载均衡
  - 请求转发到后端服务
  - 统一认证和授权

### Admin Service (管理服务)
- **端口**: 8081 (HTTP)
- **功能**:
  - 管理员登录/登出
  - 用户管理
  - 系统配置

### User Service (用户服务)
- **端口**: 50051 (gRPC)
- **功能**:
  - 用户注册/登录
  - 用户信息管理
  - JWT Token 生成

### AI Service (AI 服务)
- **端口**: 50052 (gRPC)
- **语言**: Python
- **功能**:
  - 聊天对话（一元/流式）
  - 文本嵌入向量生成

## 快速开始

### 1. 安装依赖

```bash
# 安装 protoc 编译器
# 下载: https://github.com/protocolbuffers/protobuf/releases
# 解压到 D:\protoc，添加 D:\protoc\bin 到 PATH

# 安装 Go 插件
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### 2. 生成 Proto 代码

```bash
cd gateway
go generate
```

### 3. 启动 Consul

```bash
# 下载 Consul: https://www.consul.io/downloads
consul agent -dev
```

### 4. 启动服务

```bash
# 启动 Admin 服务
cd ../admin
go run main.go

# 启动 Gateway
go run main.go
```

### 5. 访问测试

```bash
# 健康检查
curl http://localhost:8081/health

# 用户登录
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123456"}'
```

## 开发流程

### 修改 Proto 文件后重新生成代码

```bash
# 生成所有服务的代码
protoc --go_out=. --go-grpc_out=. proto/user.proto proto/admin.proto proto/ai.proto

# 或单独生成某个服务
protoc --go_out=. --go-grpc_out=. proto/user.proto
```

### 添加新依赖

```bash
# 例如安装 Consul 客户端
go get github.com/hashicorp/consul/api

# 安装 Gin
go get github.com/gin-gonic/gin
```

### 服务注册

每个服务启动时自动向 Consul 注册：
```go
discovery.RegisterService("user-service", "localhost:50051")
```

### 服务调用

Gateway 通过 Consul 发现服务地址：
```go
address := discovery.GetServiceAddress("user-service")
conn := grpc.Dial(address)
client := user.NewUserServiceClient(conn)
```

## 注意事项

- `pb/` 目录下的代码是自动生成的，不要手动修改
- 修改接口需要编辑 `proto/*.proto` 文件，然后运行 `go generate`
- `.pb.go` 文件应该提交到 Git
- 保持 proto 文件和生成代码的同步
- 导入时使用：`import "gateway/pb/user"`

package consul

import (
	"fmt"
	"log"
	"net"
	"time"

	"admin/config"

	"github.com/hashicorp/consul/api"
)

// RegisterService 将当前服务注册到 Consul
func RegisterService(cfg *config.Config) error {
	// 1. 创建 Consul 客户端配置
	consulConfig := api.DefaultConfig()
	consulConfig.Address = cfg.Consul.Addr

	client, err := api.NewClient(consulConfig) // 创建 Consul 客户端
	if err != nil {
		return fmt.Errorf("创建 Consul 客户端失败: %v", err)
	}

	// 2. 确定服务注册地址（配置优先）
	addr := cfg.Consul.ServiceAddress
	if addr == "" {
		// 如果未配置，则尝试自动获取本机非回环 IP
		ip, err := getLocalIP()
		if err != nil {
			return fmt.Errorf("无法获取本地 IP 且未配置 CONSUL_SERVICE_ADDRESS: %v", err)
		}
		addr = ip
	}

	// 3. 定义服务注册信息（使用 TTL 健康检查）
	registration := &api.AgentServiceRegistration{
		ID:      cfg.Consul.ServiceID,
		Name:    cfg.Consul.ServiceName,
		Port:    cfg.Consul.ServicePort,
		Address: addr,
		Tags:    []string{"grpc"},
		Check: &api.AgentServiceCheck{
			CheckID:                        "service:" + cfg.Consul.ServiceID,
			TTL:                            "15s", // 设置 TTL 为 15 秒
			DeregisterCriticalServiceAfter: "30s", // 30 秒未收到心跳则注销
		},
	}

	// 4. 执行注册
	log.Printf("正在向 Consul (%s) 注册服务: %s [Addr: %s, Port: %d]",
		cfg.Consul.Addr, cfg.Consul.ServiceName, addr, cfg.Consul.ServicePort)

	if err := client.Agent().ServiceRegister(registration); err != nil {
		return fmt.Errorf("服务注册失败: %v", err)
	}

	// 5. 启动心跳协程
	go startHeartbeat(client, cfg.Consul.ServiceID)

	log.Println("[SUCCESS] Admin 服务已成功注册到 Consul，心跳监测已启动")
	return nil
}

// startHeartbeat 定期向 Consul 发送心跳
func startHeartbeat(client *api.Client, serviceID string) {
	checkID := "service:" + serviceID
	ticker := time.NewTicker(10 * time.Second) // 每 10 秒发送一次心跳
	defer ticker.Stop()

	for range ticker.C {
		if err := client.Agent().UpdateTTL(checkID, "", api.HealthPassing); err != nil {
			log.Printf("[WARN] Consul 心跳发送失败: %v", err)
		}
	}
}

// DeregisterService 从 Consul 注销服务（用于优雅退出）
func DeregisterService(cfg *config.Config) error {
	consulConfig := api.DefaultConfig()
	consulConfig.Address = cfg.Consul.Addr

	client, err := api.NewClient(consulConfig)
	if err != nil {
		return err
	}

	log.Printf("正在从 Consul 注销服务: %s", cfg.Consul.ServiceID)
	return client.Agent().ServiceDeregister(cfg.Consul.ServiceID)
}

// getLocalIP 获取本机非回环 IP 地址
func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("未找到合适的 IP 地址")
}

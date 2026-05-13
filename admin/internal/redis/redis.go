package redis

import (
	"context"
	"fmt"
	"log"

	"admin/config"

	"github.com/go-redis/redis/v8"
)

var Rdb *redis.Client

// InitRedis 初始化 Redis 连接（包含连接池配置）
func InitRedis(cfg *config.Config) {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		// 连接池配置
		PoolSize:     100, // 最大连接数
		MinIdleConns: 10,  // 最小空闲连接数
	})

	// 测试连接
	ctx := context.Background()
	if _, err := Rdb.Ping(ctx).Result(); err != nil {
		log.Fatalf("Redis 连接失败: %v", err)
	}

	log.Println("Redis 连接成功")
}

// GetShardedKey 根据 ID 生成分片 Key（简单取模示例，生产环境建议使用一致性哈希）
func GetShardedKey(prefix string, id interface{}, shardCount int) string {
	// 这里仅作演示，实际项目中应使用 crc32 或 murmur3 算法计算哈希
	hash := fmt.Sprintf("%v", id)
	shard := 0
	// 简单的哈希逻辑
	for _, c := range hash {
		shard += int(c)
	}
	shard = shard % shardCount
	return fmt.Sprintf("%s:shard%d:%v", prefix, shard, id)
}

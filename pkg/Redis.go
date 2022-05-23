package pkg

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	Rdb             *redis.Client
	redisDir        string
	redisDbFilename string
)

// RedisClient 连接 Redis
func RedisClient(pwd string) (err error) {

	Rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", Rhost, Rport),
		Password: pwd, // 密码认证
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	pong, err := Rdb.Ping(ctx).Result()
	if err != nil {
		return err
	}
	if strings.Contains(pong, "PONG") {
		redisVersion()
	}
	return nil

}

// RedisCmd 执行 Redis 命令
func RedisCmd(cmd string) interface{} {

	ctx := context.Background()

	var argsInterface []interface{}
	args := strings.Fields(cmd)
	for _, arg := range args {
		argsInterface = append(argsInterface, arg)
	}

	info, err := Rdb.Do(ctx, argsInterface...).Result()
	if err != nil {
		Err(err)
		return ""
	}
	return info
}

// 获取 Redis 基本信息
func redisVersion() bool {
	info := RedisCmd("info")
	if strings.Contains(info.(string), "redis_version") {
		Info("获取 Redis 基本信息")
		os := ReString(info, "os:.*")
		version := ReString(info, "redis_version:.*")
		Success(os)
		Success(version)
		dir := RedisCmd("config get dir")
		redisDir = redisString(dir)[4:]
		Success(redisDir)

		file := RedisCmd("config get dbfilename")
		redisDbFilename = redisString(file)[11:]
		Success(redisDbFilename)
		return true
	}
	return false
}

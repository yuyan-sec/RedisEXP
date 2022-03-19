package pkg

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	Rdb                       *redis.Client
	Lhost                     string
	Lport                     string
	Rhost                     string
	Rport                     string
	PWD                       string
	redisDir, redisDbfilename string
)

// RedisClient 连接 Redis
func RedisClient() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", Rhost, Rport),
		Password: PWD, // 密码认证
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	_, err := Rdb.Ping(ctx).Result()
	if err != nil {
		if strings.Contains(err.Error(), "context deadline exceeded") {
			Info("Redis 连接超时")
			os.Exit(0)
		}

		Err(err)

		if strings.Contains(err.Error(), "NOAUTH Authentication required.") {
			Info("Redis 需要密码认证")
			os.Exit(0)
		}
		if strings.Contains(err.Error(), "ERR invalid password") {
			log.Println("Redis 认证密码错误!")
			os.Exit(0)
		}

		return
	}

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
func redisVersion() {
	Info("获取 Redis 基本信息")
	info := RedisCmd("info")
	os := redisRe(info, "os:.*")
	version := redisRe(info, "redis_version:.*")
	Success(os)
	Success(version)
	dir := RedisCmd("config get dir")
	redisDir = redisString(dir)[4:]
	Success(redisDir)

	file := RedisCmd("config get dbfilename")
	redisDbfilename = redisString(file)[11:]
	Success(redisDbfilename)

}

// 正则
func redisRe(info interface{}, s string) string {
	reg := regexp.MustCompile(s)
	list := reg.FindAllStringSubmatch(info.(string), -1)
	return list[0][0]
}

func redisString(i interface{}) string {
	switch v := i.(type) {
	case []interface{}:
		s := ""
		for _, i := range v {
			s += i.(string) + " "
		}
		return s
	}
	return ""

}

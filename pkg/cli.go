package pkg

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// 循环执行 Redis 命令
func LoopRedis(rhost, rport string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s:%s> ", rhost, rport)
		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimRight(cmd, "\r\n")

		if strings.EqualFold(cmd, "exit") {
			break
		}
		// 执行命令
		i, err := RedisCmd(cmd)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(i.(string))
	}

}

// 循环执行命令
func LoopCmd(dll string) {

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimRight(cmd, "\r\n")

		if strings.EqualFold(cmd, "exit") {
			CloseSlave(dll)
			break
		}
		RunCmd(cmd)

	}
}

// 循环执行CVE
func LoopCVE() {

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimRight(cmd, "\r\n")

		if strings.EqualFold(cmd, "exit") {
			break
		}
		RedisLua(cmd)

	}
}

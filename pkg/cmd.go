package pkg

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// 循环执行命令
func loopCmd(s string) {
	Info("执行命令")
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimRight(cmd, "\r\n")
		if cmd == "exit" || cmd == "q" || cmd == "quit" {
			if strings.Contains(s, "exec") {
				CloseSlave("exec")
			}
			break
		}
		// 执行命令
		if strings.Contains(s, "exec") {
			RunCmd(cmd)
		} else if strings.Contains(s, "lua") {
			RedisLua(cmd)
		}

	}
}

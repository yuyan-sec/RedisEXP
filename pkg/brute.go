package pkg

import (
	"os"
	"strings"
	"sync"
)

var wg sync.WaitGroup

// 爆破密码
func brutePWD() {

	ch := make(chan struct{}, 1)
	for _, value := range data {
		wg.Add(1)
		ch <- struct{}{}
		go func() {
			defer wg.Done()
			err := RedisClient(value)
			if err == nil {
				Success("成功爆破到 Redis 密码：" + value)
				os.Exit(0)
			} else if strings.Contains(err.Error(), "ERR Client sent AUTH, but no password is set") {
				Success("存在未授权 Redis , 不需要输入密码")
				os.Exit(0)
			} else {
				Err(err)
			}
			<-ch
		}()
	}

	wg.Wait()
	Info("未发现 Redis 密码")

}

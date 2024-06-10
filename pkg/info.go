package pkg

import (
	"fmt"
	"strings"
)

var (
	Redis_dir        string
	Redis_dbfilename string
)

func RedisVersion(show bool) {
	info, _ := RedisCmd("info server")
	Redis_dir = getRedisValue("config get dir")
	Redis_dbfilename = getRedisValue("config get dbfilename")

	if show {
		for _, s := range strings.Split(info.(string), "\r\n") {
			switch {
			case strings.Contains(s, "redis_version"), strings.Contains(s, "os"), strings.Contains(s, "arch_bits"):
				fmt.Printf("[%v]\n", s)
			}
		}

		fmt.Printf("[dir: %s]\n", Redis_dir)

		fmt.Printf("[dbfilename: %s]\n", Redis_dbfilename)
	}

}

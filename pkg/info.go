package pkg

import (
	"fmt"
	"strings"
)

func RedisVersion() {
	info, _ := RedisCmd("info server")

	for _, s := range strings.Split(info.(string), "\r\n") {
		switch {
		case strings.Contains(s, "redis_version"), strings.Contains(s, "os"), strings.Contains(s, "arch_bits"):
			fmt.Printf("[%v]\n", s)
		}
	}

	dir, _ := RedisCmd("config get dir")
	fmt.Println(dir)
	dbfilename, _ := RedisCmd("config get dbfilename")
	fmt.Println(dbfilename)

}

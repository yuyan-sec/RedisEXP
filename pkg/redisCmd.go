package pkg

import (
	"context"
	"regexp"
	"strings"
)

// redisCmd 执行 Redis 命令
func RedisCmd(cmd string) (interface{}, error) {
	ctx := context.Background()

	var argsInterface []interface{}

	// 处理输入字符串有空格的问题
	if strings.Contains(cmd, "\"") || strings.Contains(cmd, "'") {
		oldString := reString(cmd, `(['"])(.*?)(['"])`)

		//newString := strings.ReplaceAll(oldString, " ", "$_$_$_$_$_$")
		//cmd = strings.ReplaceAll(cmd, oldString, newString)
		//cmd = strings.ReplaceAll(cmd, "\"", "")
		//cmd = strings.ReplaceAll(cmd, "'", "")

		cmd = strings.NewReplacer(oldString, strings.ReplaceAll(oldString, " ", "$_$_$_$_$_$"), "\"", "", "'", "").Replace(cmd)

	}

	args := strings.Fields(cmd)
	for _, arg := range args {
		arg = strings.ReplaceAll(arg, "$_$_$_$_$_$", " ")
		argsInterface = append(argsInterface, arg)
	}

	info, err := Rdb.Do(ctx, argsInterface...).Result()
	if err != nil {
		return nil, err
	}
	return info, nil
}

// 正则匹配
func reString(info interface{}, s string) string {
	reg := regexp.MustCompile(s)
	list := reg.FindAllStringSubmatch(info.(string), -1)
	return list[0][0]
}

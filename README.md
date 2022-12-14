# Eval 2 Term !
> 连接一句话`webshell`，并获取到可交互的虚拟终端

![](https://txc.gtimg.com/data/383193/2022/0908/322a397ed3b47ce68efefc52ec676ae5.png)

## 环境需求
1. `PHP` 支持`proc_open`函数（默认支持）    
2. 目标支持可写目录（默认在当前目录，可修改代码为其他可写目录）

## 程序说明
`Eval2Term` 是一个能直接连接`eval`一句话，并得到一个可交互终端的小工具。    
比如在一些环境下：
1. 目标机器不出网
2. 没有权限开启`sock`端口
3. 目标没有反弹程序或脚本环境
4. 本机无外网IP可供监听

## 使用说明
安装`golang`，下载源码，编译：
``` bash
$ cd eval2term
$ go build
```
编译得出文件：`eval2term`

用法：
``` bash
$ ./eval2term -url http://target/shell.php -pwd ant
```
按 `Ctrl+D` 结束程序

> 说明：服务端代码仅在 `Linux` 机器上测试，客户端代码仅在 `macOS` 系统上测试
> 如无法顺利使用，请参考下方原理说明，用你聪慧的大脑和勤劳的双手，改造之！


## 原理说明
由于我们没法采用`socket`进行长连接通信，所以采用了`HTTP`长轮询的思路：    
首先，注入代码，`php`后台运行`bash`进程    
然后，客户端监听用户按键，把按键数据发送到目标    
接着，后台获取按键数据，发送到进程句柄，把结果输出到文件    
最后，客户端再轮询获取结果，输出    

## 缺点
这种方式在以前，应该是非常好用的。   
随着科技发展，未来应该采用 `HTTP2`、`websocket`等方式，以实现长连接操作。

这种方式的缺点有如下：
1. 有延迟，因为我们发送按键数据、获取命令结果，都需要发起`http`请求    
2. 同上，多次发起请求，可能会导致防火墙拦截    
3. 暂不支持中文等语言（要解决按键问题）    

## 待做
目前程序主要功能和逻辑已完成，你可以根据自己需求改动代码使用。    
我觉得，后续应该会有如下功能：
1. 多脚本支持
2. 多用户操作
3. 流量加密等常规操作
4. 还没想好

## 后记
疫情，封城，宅家，想起微博密码    
五年，奔三，回望，我们都已长大

- - -

我是蚁逅，奔三前留点作品的少年。    
微博：https://weibo.com/antoor    
微信：`Shell_Way`

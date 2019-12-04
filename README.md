# 项目介绍

**go-netbus**【网络直通车】是为了解决内网穿透问题建立的项目。

## 功能列表

- 基于 TCP 协议
- 支持多端口穿透
- 支持断线重连
- 支持指定访问端口
- 增加接入身份校验，解决服务端能被任何人使用的问题

## 启动方式

支持两种方式启动:

- 纯命令行启动（不需要配置文件）
- 配置文件启动（需要配置文件）**【推荐】**

### 命令行启动

```bash
# 启动服务端

$ netbus -server <port> <custom-port-key> <random-port-key>

# 注释
# port                服务端端口，不要使用保留端口，必填
# key                 6-16 个字符，用于身份校验

```

```bash
# 启动客户端

$ netbus -client <key> <server:port> <local:port> [access-port] [max-redial-times]

# 注释
# key                与服务端保持一致
# server:port        服务端地址，格式如：45.32.78.129:6666
# local:port         被代理服务地址，多个以逗号隔开，比如：127.0.0.1:8080,127.0.0.1:9200
# access-port        访问端口，与 local:port 一一对应，多个以逗号隔开，比如：9090,10200， 可选，若未填访问端口保持与 local:port 一致
# max-redial-times   最大重连次数，-1 表示无限重连，可选参数
```

### 配置文件启动

**服务端配置**
```ini
# 服务端配置
[server]
# 代理端口
port = 6666
# Key 6-16 个字符，用于身份校验
key = winshu
```

**客户端配置**
```ini
# 客户端配置
[client]
# Key 与服务端保持一致
key = winshu
# 服务端地址，格式 ip:port
server-host = 127.0.0.1:6666
# 内网被代理服务地址(多个用逗号隔开，端口不能相同)，格式 ip:port,ip:port
local-host = 127.0.0.1:3306
# 访问端口(可选，未设置时访问端口与代理端口相同)
access-port = 13306
# 最大重连次数，-1表示一直重连
max-redial-times = 20
```

**启动命令**
```bash
# 启动服务端
$ netbus -server

# 启动客户端
$ netbus -client
```
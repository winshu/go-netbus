# 项目介绍

**go-netbus**【网络直通车】是为了解决内网穿透问题建立的项目。

## 功能列表

- 基于 TCP 协议
- 支持多端口穿透
- 支持断线重连
- 支持指定访问端口
- 增加接入身份校验，解决服务端能被任何人使用的问题
- 单服务支持多个连接

## 工作原理

![netbus_architecture](doc/netbus_architecture.png)

通常我们需要访问一个服务，只需要直接通过地址连接就可以了，但是如果服务处在一个与访问者隔离的内网环境里，我们就不能访问了。
不过也不是没有办法，如果内网服务可以访问外网的话，我们就可以找到突破点。

### 实现步骤

**服务端部分**
1. 启动`tunnel_bridge`， 监听`tunnel_client`的连接。
2. 启动`tunnel_server`， 监听`用户`的访问请求。

**客户端部分**
1. 启动`tunnel_client`， 向`tunnel_bridge`拨号，建立连接。
2. 如果`用户`发起请求，则向`内网服务`拨号，建立连接。
3. 处理连接双方的通讯，实现穿透。

## 启动方式

支持两种方式启动:

- 纯命令行启动（不需要配置文件）
- 配置文件启动（需要配置文件）**【推荐】**

### 命令行启动

```bash
# 启动服务端

$ netbus -server <key> <port>

# 注释
# key                 6-16 个字符，用于身份校验
# port                服务端端口，不要使用保留端口，必填

```

```bash
# 启动客户端

$ netbus -client <key> <server:port> <local:port> [access-port] [max-redial-times]

# 注释
# key                6-16 个字符，用于身份校验，与服务端保持一致
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
# 6-16 个字符，用于身份校验，与服务端保持一致
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
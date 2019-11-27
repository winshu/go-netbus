# 项目介绍

**go-netbus**【网络直通车】是为了解决内网穿透问题而建立的项目。

## 功能列表

- 基于 TCP 协议
- 支持多端口穿透
- 支持无限重连机制、有限次数重连机制
- 支持两种访问级别：指定访问端口、随机访问端口
- 服务端生成代理端口时回写给客户端，解决随机端口的问题
- 增加代理接入身份校验，解决服务端能被任何人使用的问题

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
# custom-port-key     指定端口访问 key， 建议长度 6~8 比如：custom
# random-port-key     随机商品访问 key， 建议长度 6~8 比如：random

```

```bash
# 启动客户端

$ netbus -client <key> <server:port> <local:port> [access-port] [max-redial-times]

# 注释
# key                保持与服务端的 custom-port-key 或 random-port-key 一致，决定了权限
# server:port        服务端地址，格式如：45.32.78.129:6666
# local:port         被代理服务地址，多个以逗号隔开，比如：127.0.0.1:8080,127.0.0.1:9200
# access-port        访问端口，与 local:port 一一对应，多个以逗号隔开，比如：9090,10200， 可选，若未填访问端口保持与 local:port 一致
# max-redial-times   最大重连次数，-1 表示无限重连，可选参数
```

### 配置文件启动

配置文件`config.ini`需与启动文件置于同一目录。

通常情况下，服务端配置与客户端配置是分开的。（单机测试时可以用一个文件）

**服务端配置**
```ini
# 服务端配置
[server]
# 代理端口
port = 6666
# 自定义端口 Key 不要太长，6-8 个字符
custom-port-key = custom
# 随机端口 Key 不要太长，6-8 个字符
random-port-key = random
```

**客户端配置**
```ini
# 客户端配置
[client]
# Key 与服务端保持一致
key = custom
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
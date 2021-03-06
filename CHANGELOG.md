# 版本变动日志

## Version 1.0.0

- 实现内网穿透
- 支持多端口穿透
- 服务端生成代理端口时回写给客户端，解决随机端口的问题

## Version 1.0.1
- 增加代理接入身份校验，解决服务端能被任何人使用的问题
- 支持客户端指定代理端口模式

## Version 1.0.2
- 使用延迟创建内网服务连接的方式，解决不停连接内网服务的问题
- 解决内网服务断开，访问服务仍未关闭的问题
- 去除随机端口访问模式
- 简化通讯协议

## Version 1.0.3
- 单服务支持多个连接，服务更稳定
- 通讯协议加入原端口

## Version 1.0.4
- 解决隧道桥中断，连接没有及时处理的问题
- 解决断线不能真正重连的问题

## Version 1.0.5
- 引入基于有效期的身份失效机制
- 增加版本校验机制，避免客户端与服务端版本不一致的问题

## TODO

- 增加通信加密
- 客户端引入连接池
- 增加心跳机制检测网络是否通畅
- 支持连接复用

通讯协议

- 通讯结果    1个字节(1: 成功，其他：失败)
- 版本号      4个字节
- 原端口      4个字节
- 访问端口    4个字节
- Key        建议长度 6-16 字符串

协议最大长度不能超过 255

举例
1|1|3306|13306|winshu

## 简易编译打包脚本

```shell script
# linux
git pull
go build -o netbus main.go
mkdir netbus_linux_amd64
mv netbus netbus_linux_amd64
cp config.ini netbus_linux_amd64
tar -zvcf netbus_linux_amd64.tar.gz netbus_linux_amd64
rm -rf netbus_linux_amd64
mv netbus_linux_amd64.tar.gz /mnt/d/

```

```bash
# windows
@echo off
git pull
go build -o netbus.exe main.go
mkdir netbus_windows_amd64
move /Y netbus.exe netbus_windows_amd64
xcopy /Y config.ini netbus_windows_amd64
xcopy /Y doc\*.bat netbus_windows_amd64
del /Q netbus_windows_amd64.zip
"C:/Program Files (x86)/WinRAR/WinRAR.exe" a netbus_windows_amd64.zip netbus_windows_amd64
rmdir /S/Q netbus_windows_amd64

```



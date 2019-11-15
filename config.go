package main

type ProxyConfig struct {
	proxyPort  int
	randomPort bool
}

type ClientConfig struct {
	proxyHost string
	localHost string
}

var (
	proxyConfig  ProxyConfig
	clientConfig ClientConfig
)


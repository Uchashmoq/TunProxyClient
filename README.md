# TunProxyClient
一个简单的采用aes加密的socks5代理

[服务器仓库](https://github.com/Uchashmoq/TunProxyServer) 

#### 编译

go build

#### 运行

go run main.go 或前往release下载可执行文件

#### 配置文件config.json

```json
{
    "StaticKey":"HrGZo2uaSgccL4Ke", 用于aes CBC模式加密的初始向量，16字节，与服务器上保持一致
    "ServerAddr":"127.0.0.1:14445",	服务器ip
    "ClientAddr":"127.0.0.1:9050"	本地ip
}
```

#### 浏览器插件配置

![p1](https://github.com/Uchashmoq/TunProxyClient/blob/main/%E5%9B%BE%E7%89%871.png)

**浏览器插件搜索 switchyomega 并安装，或者其他支持socks5代理的工具**

![p2](https://github.com/Uchashmoq/TunProxyClient/blob/main/%E5%9B%BE%E7%89%872.png)

![p3](https://github.com/Uchashmoq/TunProxyClient/blob/main/%E5%9B%BE%E7%89%873.png)

**创建一个情景模式，自己起名（这里起名为tun），配置：网址协议sock5，代理服务器127.0.0.1,端口9050,无需密码,应用并保存**

![p4](https://github.com/Uchashmoq/TunProxyClient/blob/main/%E5%9B%BE%E7%89%874.png)

**打开插件，选择刚才的情景模式（这里起名叫maharsock）打开，启动代理。不用了的话就选择系统代理** 

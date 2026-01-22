# LiteSpeed

LiteSpeed is a simple tool for batch test ss/ssr/v2ray/trojan/clash servers.  
Feature
- 支持ss/ssr/v2ray/trojan/clash订阅链接
- 支持ss/ssr/v2ray/trojan/clash节点链接
- 支持ss/ssr/v2ray/trojan/clash订阅或节点文件
- support ss/ssr/v2ray/trojan/clash subscription url,
- support ss/ssr/v2ray/trojan/clash profile links
- support ss/ssr/v2ray/trojan/clash subscription or profile file, 

感谢@xxf098 大佬，本工具是在LiteSpeedTest的基础上修改的，去掉了部分打印，去掉了web相关的内容，只保留了测速相关功能，新增检测节点国家的功能
关于测速，只保证了节点的联通性，不展示具体的节点速度。

 ![build]() 

### Build
```bash
$env:GOOS = "linux"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"

go build -trimpath -ldflags "-s -w" -o lite-linux-with-ip
```
### Start

 ```
 go run . -input "https://xxxx"
 ```

 ```
 ./lite-linux-with-ip -input "https://xxxx"
 ```

## Credits

- [clash](https://github.com/Dreamacro/clash)
- [stairspeedtest-reborn](https://github.com/tindy2013/stairspeedtest-reborn)
- [gg](https://github.com/fogleman/gg)
- [LiteSpeedTest](https://github.com/xxf098/LiteSpeedTest)


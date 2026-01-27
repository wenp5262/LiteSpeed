# LiteSpeed

LiteSpeed is a simple tool for batch test ss/ssr/v2ray/trojan/clash servers.  
Feature
- 支持ss/ssr/v2ray/trojan/clash订阅链接
- 支持ss/ssr/v2ray/trojan/clash节点链接
- 支持ss/ssr/v2ray/trojan/clash订阅或节点文件
- support ss/ssr/v2ray/trojan/clash subscription url,
- support ss/ssr/v2ray/trojan/clash profile links
- support ss/ssr/v2ray/trojan/clash subscription or profile file, 

本项目基于 LiteSpeedTest 修改而来，感谢 @xxf098 大佬的原始贡献。

在原项目基础上：

移除了冗余的日志打印

去除了 Web 相关功能

仅保留核心的测速与连通性检测能力

新增 节点国家/地区检测 功能

需要说明的是：
本工具的测速结果仅用于判断节点的连通性状态，不用于展示或评估节点的实际带宽或速度性能。



This project is a modified version of LiteSpeedTest, with special thanks to @xxf098 for the original work.

Changes and features include:

Removed unnecessary log outputs

Stripped all web-related components

Retained only the core connectivity and testing logic

Added support for node country/region detection

Note:
The speed test in this project is designed only to verify node connectivity. It does not measure or display actual node bandwidth or performance.


### Build
build.ps1 打包脚本
```bash
$env:GOOS = "linux"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"

go build -trimpath -ldflags "-s -w" -o LiteSpeed
```

### Start

 ```
 go run . -input "https://xxxx"
 ```

 ```
 ./LiteSpeed -input "https://xxxx"
 ```
```
./LiteSpeed_windows_386.exe -input "https://xxxx"
```

### Output
```azure
[
  {
    "id": "",
    "remarks": "",
    "link": "trojan://869e9086806483ca4744a4cb0f3d6e16@160.16.89.176:1933/?type=tcp\u0026security=tls\u0026sni=www.nintendogames.net#%E6%97%A5%E6%9C%AC3%7C%40ripaojiedian",
    "protocol": "trojan",
    "server": "160.16.89.176:1933",
    "speed": "",
    "remoteip": "116.80.60.36",
    "country": "日本"
  },
  {
    "id": "",
    "remarks": "",
    "link": "trojan://869e9086806483ca4744a4cb0f3d6e16@160.16.89.176:3076/?type=tcp\u0026security=tls\u0026sni=www.nintendogames.net#%E6%96%B0%E5%8A%A0%E5%9D%A1%7C%40ripaojiedian",
    "protocol": "trojan",
    "server": "160.16.89.176:3076",
    "speed": "",
    "remoteip": "178.128.114.70",
    "country": "新加坡"
  },
  {
    "id": "",
    "remarks": "",
    "link": "trojan://869e9086806483ca4744a4cb0f3d6e16@58.152.18.124:443/?type=tcp\u0026security=tls\u0026sni=www.nintendogames.net#%E9%A6%99%E6%B8%AF3%7C%40ripaojiedian",
    "protocol": "trojan",
    "server": "58.152.18.124:443",
    "speed": "",
    "remoteip": "58.152.18.124",
    "country": "香港"
  },
  {
    "id": "",
    "remarks": "",
    "link": "ss://YWVzLTI1Ni1jZmI6WG44aktkbURNMDBJZU8lIyQjZkpBTXRzRUFFVU9wSC9ZV1l0WXFERm5UMFNW@103.186.155.51:38388#%E8%B6%8A%E5%8D%97%7C%40ripaojiedian",
    "protocol": "ss",
    "server": "103.186.155.51:38388",
    "speed": "",
    "remoteip": "103.186.155.51",
    "country": "越南"
  },
  {
    "id": "",
    "remarks": "",
    "link": "vmess://eyJ2IjogIjIiLCAicHMiOiAiREU0Nu+9nOacuuaIv++9nOS9jumjjumZqSIsICJhZGQiOiAidjMzLmhkYWNkLmNvbSIsICJwb3J0IjogIjMwODMzIiwgImFpZCI6IDIsICJzY3kiOiAiYXV0byIsICJuZXQiOiAidGNwIiwgInR5cGUiOiAibm9uZSIsICJ0bHMiOiAiIiwgImlkIjogImNiYjNmODc3LWQxZmItMzQ0Yy04N2E5LWQxNTNiZmZkNTQ4NCJ9",
    "protocol": "vmess",
    "server": "v33.hdacd.com:30833",
    "speed": "",
    "remoteip": "164.92.247.180",
    "country": "德国"
  },
  {
    "id": "",
    "remarks": "",
    "link": "vmess://eyJhZGQiOiJ2MTAuaGRhY2QuY29tIiwiYWlkIjoiMiIsImFscG4iOiIiLCJob3N0IjoiIiwiaWQiOiJjYmIzZjg3Ny1kMWZiLTM0NGMtODdhOS1kMTUzYmZmZDU0ODQiLCJuZXQiOiJ0Y3AiLCJwYXRoIjoiLyIsInBvcnQiOiIzMDgwNyIsInBzIjoi6aaZ5rivfEByaXBhb2ppZWRpYW4iLCJzY3kiOiJhdXRvIiwic25pIjoiIiwidGxzIjoiIiwidHlwZSI6Im5vbmUiLCJ2IjoiMiJ9",
    "protocol": "vmess",
    "server": "v10.hdacd.com:30807",
    "speed": "",
    "remoteip": "112.118.97.96",
    "country": "香港"
  },
  {
    "id": "",
    "remarks": "",
    "link": "vmess://eyJ2IjogIjIiLCAicHMiOiAiQ04z772c5a625a69772c6Z2e5bi45a6J5YWoIiwgImFkZCI6ICIxODMuMjM2LjUxLjM2IiwgInBvcnQiOiAiNTkwMDMiLCAiYWlkIjogMCwgInNjeSI6ICJhdXRvIiwgIm5ldCI6ICJ0Y3AiLCAidHlwZSI6ICJub25lIiwgInRscyI6ICIiLCAiaWQiOiAiNDE4MDQ4YWYtYTI5My00Yjk5LTliMGMtOThjYTM1ODBkZDI0In0=",
    "protocol": "vmess",
    "server": "183.236.51.36:59003",
    "speed": "",
    "remoteip": "183.236.51.36",
    "country": "中国"
  },
  {
    "id": "",
    "remarks": "",
    "link": "vmess://eyJ2IjogIjIiLCAicHMiOiAiVVMxNTTvvZzmnLrmiL/vvZzkvY7po47pmakiLCAiYWRkIjogInY0LmhkYWNkLmNvbSIsICJwb3J0IjogIjMwODA0IiwgImFpZCI6IDIsICJzY3kiOiAiYXV0byIsICJuZXQiOiAidGNwIiwgInR5cGUiOiAibm9uZSIsICJ0bHMiOiAiIiwgImlkIjogImNiYjNmODc3LWQxZmItMzQ0Yy04N2E5LWQxNTNiZmZkNTQ4NCJ9",
    "protocol": "vmess",
    "server": "v4.hdacd.com:30804",
    "speed": "",
    "remoteip": "216.144.235.210",
    "country": "美国"
  },
  {
    "id": "",
    "remarks": "",
    "link": "trojan://869e9086806483ca4744a4cb0f3d6e16@58.152.18.124:443?sni=www.nintendogames.net#%E9%A6%99%E6%B8%AF3%7C%40ripaojiedian",
    "protocol": "trojan",
    "server": "58.152.18.124:443",
    "speed": "",
    "remoteip": "58.152.18.124",
    "country": "香港"
  },
  {
    "id": "",
    "remarks": "",
    "link": "trojan://869e9086806483ca4744a4cb0f3d6e16@160.16.89.176:1933?allowInsecure=0\u0026sni=www.nintendogames.net#JP_speednode_0009",
    "protocol": "trojan",
    "server": "160.16.89.176:1933",
    "speed": "",
    "remoteip": "116.80.60.36",
    "country": "日本"
  },
  {
    "id": "",
    "remarks": "",
    "link": "trojan://869e9086806483ca4744a4cb0f3d6e16@160.16.89.176:3076?allowInsecure=0\u0026sni=www.nintendogames.net#JP_speednode_0010",
    "protocol": "trojan",
    "server": "160.16.89.176:3076",
    "speed": "",
    "remoteip": "178.128.114.70",
    "country": "新加坡"
  },
  {
    "id": "",
    "remarks": "",
    "link": "vmess://eyJ2IjogIjIiLCAicHMiOiAiU0cxNe+9nOacuuaIv++9nOS9jumjjumZqSIsICJhZGQiOiAidjEyLmhkYWNkLmNvbSIsICJwb3J0IjogIjMwODEyIiwgImFpZCI6IDAsICJzY3kiOiAiYXV0byIsICJuZXQiOiAidGNwIiwgInR5cGUiOiAibm9uZSIsICJ0bHMiOiAiIiwgImlkIjogImNiYjNmODc3LWQxZmItMzQ0Yy04N2E5LWQxNTNiZmZkNTQ4NCJ9",
    "protocol": "vmess",
    "server": "v12.hdacd.com:30812",
    "speed": "",
    "remoteip": "172.104.180.238",
    "country": "新加坡"
  }
]

```

## Credits

- [clash](https://github.com/Dreamacro/clash)
- [stairspeedtest-reborn](https://github.com/tindy2013/stairspeedtest-reborn)
- [gg](https://github.com/fogleman/gg)
- [LiteSpeedTest](https://github.com/xxf098/LiteSpeedTest)


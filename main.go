package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/wenp5262/LiteSpeed/cli"
	"github.com/wenp5262/LiteSpeed/config"
	"github.com/wenp5262/LiteSpeed/request"
)

type Server struct {
	Id       string `json:"id"`
	Remarks  string `json:"remarks"`
	Link     string `json:"link"`
	Protocol string `json:"protocol"`
	Server   string `json:"server"`
	Speed    string `json:"speed"`
	RemoteIP string `json:"remoteip"`
	Country  string `json:"country"`
}

// 只要能拿到 country 就算联通；否则视为过滤
func main() {
	input := flag.String("input", "", "subscription URL / file path / base64 / raw profile text")
	concurrency := flag.Int("c", 20, "concurrency")
	showFiltered := flag.Bool("show-filtered", true, "print nodes filtered during test (country empty / errors) to stderr")
	flag.Parse()

	if strings.TrimSpace(*input) == "" {
		fmt.Fprintln(os.Stderr, "missing -input")
		os.Exit(2)
	}

	// 解析输入（这里就可能发生“测试前过滤”）
	rep, err := cli.ParseLinksReport(*input)
	if err != nil {
		fmt.Fprintln(os.Stderr, "parse input failed:", err)
		os.Exit(1)
	}

	links := rep.Links
	if len(links) == 0 {
		fmt.Fprintln(os.Stderr, "no nodes found in input")
		os.Exit(1)
	}

	sem := make(chan struct{}, *concurrency)
	var wg sync.WaitGroup

	// 成功结果通道（并发安全收集）
	okCh := make(chan Server, len(links))

	for i, link := range links {
		wg.Add(1)
		sem <- struct{}{}
		go func(id int, l string) {
			defer wg.Done()
			defer func() { <-sem }()

			// 默认留空（按你要求：remarks/id/speed 不用管）
			var (
				proto  string
				server string
			)

			// 尽量从 Link2Config 提取 protocol/server（提取不到就留空）
			if cfg, e := config.Link2Config(l); e == nil && cfg != nil {
				// cfg.Protocol 是否存在取决于你的实现；若没有就用 Link 前缀判断
				// 这里尽可能兼容：优先 cfg.Server/cfg.Port，再 fallback
				if cfg.Server != "" && cfg.Port != 0 {
					server = fmt.Sprintf("%s:%d", cfg.Server, cfg.Port)
				}
				// 你的 cfg 结构如果有 Protocol 字段，可在此赋值；否则用前缀猜
				// proto = cfg.Protocol
			}

			// 前缀猜协议（避免 cfg 不含 Protocol 字段）
			if proto == "" {
				proto = guessProtocol(l)
			}

			// 你的判定：能获取到国家就算联通
			ip, country := request.GetIPThroughLink(id, l)
			if strings.TrimSpace(country) == "" {
				if *showFiltered {
					// 这里打印“测试阶段过滤”的节点
					// 注意：不使用 log，避免时间戳
					//fmt.Fprintf(os.Stderr, "FILTERED(during test) protocol=%s server=%s reason=country empty link=%s\n",
					//	proto, server, l)
				}
				return
			}

			okCh <- Server{
				Id:       "",
				Remarks:  "",
				Link:     l,
				Protocol: proto,
				Server:   server,
				Speed:    "",
				RemoteIP: ip,
				Country:  country,
			}
		}(i, link)
	}

	wg.Wait()
	close(okCh)

	// 汇总结果
	results := make([]Server, 0, len(links))
	for s := range okCh {
		results = append(results, s)
	}

	// 最终：只输出 JSON 到 stdout
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "json marshal failed:", err)
		os.Exit(1)
	}
	fmt.Println(string(data))
}

// 通过 link 前缀猜协议，避免依赖 cfg.Protocol 字段
func guessProtocol(link string) string {
	l := strings.ToLower(strings.TrimSpace(link))
	switch {
	case strings.HasPrefix(l, "vmess://"):
		return "vmess"
	case strings.HasPrefix(l, "vless://"):
		return "vless"
	case strings.HasPrefix(l, "trojan://"):
		return "trojan"
	case strings.HasPrefix(l, "ss://"):
		return "ss"
	case strings.HasPrefix(l, "ssr://"):
		return "ssr"
	case strings.HasPrefix(l, "http://"):
		return "http"
	case strings.HasPrefix(l, "https://"):
		return "https"
	default:
		return ""
	}
}

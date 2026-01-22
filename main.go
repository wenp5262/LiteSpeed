package main

import (
	"flag"
	"fmt"
	"log"
	"sync"

	"github.com/wenp5262/LiteSpeed/cli"
	"github.com/wenp5262/LiteSpeed/config"
	"github.com/wenp5262/LiteSpeed/request"
)

type ConnectableNode struct {
	Remarks string `json:"remarks"`
	Country string `json:"country"`
	IP      string `json:"ip"`
	Link    string `json:"link"`
}

// 判定“能联通”的标准：能够成功获取到节点出口国家（country != ""）
func main() {
	input := flag.String("input", "", "subscription URL / file path / base64 / raw profile text")
	concurrency := flag.Int("c", 20, "concurrency")
	flag.Parse()

	if *input == "" {
		log.Fatal("missing -input")
	}

	links, err := cli.ParseLinks(*input)
	if err != nil {
		log.Fatalf("parse input failed: %v", err)
	}
	if len(links) == 0 {
		log.Fatal("no nodes found in input")
	}

	sem := make(chan struct{}, *concurrency)
	var wg sync.WaitGroup
	out := make(chan *ConnectableNode, len(links))

	for i, link := range links {
		wg.Add(1)
		sem <- struct{}{}
		go func(id int, l string) {
			defer wg.Done()
			defer func() { <-sem }()

			cfg, err := config.Link2Config(l)
			remarks := ""
			if err == nil && cfg != nil {
				remarks = cfg.Remarks
			}

			ip, country := request.GetIPThroughLink(id, l)
			if country == "" {
				return
			}
			out <- &ConnectableNode{
				Remarks: remarks,
				Country: country,
				IP:      ip,
				Link:    l,
			}
		}(i, link)
	}

	wg.Wait()
	close(out)

	// 输出：仅能联通节点
	fmt.Println("remarks\tcountry\tip\tlink")
	for n := range out {
		fmt.Printf("%s\t%s\t%s\t%s\n", n.Remarks, n.Country, n.IP, n.Link)
	}
}

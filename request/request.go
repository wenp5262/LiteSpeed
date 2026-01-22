package request

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/wenp5262/LiteSpeed/common"
	"github.com/wenp5262/LiteSpeed/config"
	C "github.com/wenp5262/LiteSpeed/constant"
	"github.com/wenp5262/LiteSpeed/dns"
	"github.com/wenp5262/LiteSpeed/outbound"
	"github.com/wenp5262/LiteSpeed/transport/resolver"
	"github.com/wenp5262/LiteSpeed/utils"
)

const (
	remoteHost   = "clients3.google.com"
	generate_204 = "http://clients3.google.com/generate_204"
)

var (
	httpRequest = []byte("GET /generate_204 HTTP/1.1\r\nHost: clients3.google.com\r\nUser-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Safari/537.36\r\n\r\n")
	tcpTimeout  = 2200 * time.Millisecond
)

type PingOption struct {
	Attempts int
	TimeOut  time.Duration
}

func PingVmess(vmessOption *outbound.VmessOption) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), tcpTimeout)
	defer cancel()
	vmess, err := outbound.NewVmess(vmessOption)
	if err != nil {
		return 0, err
	}
	meta := &C.Metadata{
		NetWork: 0,
		Type:    0,
		SrcPort: "",
		DstPort: "80",
		Host:    remoteHost,
	}
	remoteConn, err := vmess.DialContext(ctx, meta)
	if err != nil {
		return 0, err
	}
	return pingInternal(remoteConn)
}

func parseFirstLine(buf []byte) (int, error) {
	bNext := buf
	var b []byte
	var err error
	for len(b) == 0 {
		if b, bNext, err = nextLine(bNext); err != nil {
			return 0, err
		}
	}

	// parse protocol
	n := bytes.IndexByte(b, ' ')
	if n < 0 {
		return 0, fmt.Errorf("cannot find whitespace in the first line of response %q", buf)
	}
	b = b[n+1:]

	// parse status code
	statusCode, n, err := parseUintBuf(b)
	if err != nil {
		return 0, fmt.Errorf("cannot parse response status code: %s. Response %q", err, buf)
	}
	if len(b) > n && b[n] != ' ' {
		return 0, fmt.Errorf("unexpected char at the end of status code. Response %q", buf)
	}

	if statusCode == 204 || statusCode == 200 {
		return len(buf) - len(bNext), nil
	}
	return 0, errors.New("Wrong Status Code")
}

func nextLine(b []byte) ([]byte, []byte, error) {
	nNext := bytes.IndexByte(b, '\n')
	if nNext < 0 {
		return nil, nil, errors.New("need more data: cannot find trailing lf")
	}
	n := nNext
	if n > 0 && b[n-1] == '\r' {
		n--
	}
	return b[:n], b[nNext+1:], nil
}

func parseUintBuf(b []byte) (int, int, error) {
	n := len(b)
	if n == 0 {
		return -1, 0, errors.New("empty integer")
	}
	v := 0
	for i := 0; i < n; i++ {
		c := b[i]
		k := c - '0'
		if k > 9 {
			if i == 0 {
				return -1, i, errors.New("unexpected first char found. Expecting 0-9")
			}
			return v, i, nil
		}
		vNew := 10*v + int(k)
		// Test for overflow.
		if vNew < v {
			return -1, i, errors.New("too long int")
		}
		v = vNew
	}
	return v, n, nil
}

func PingTrojan(trojanOption *outbound.TrojanOption) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), tcpTimeout)
	defer cancel()
	trojan, err := outbound.NewTrojan(trojanOption)
	if err != nil {
		return 0, err
	}
	meta := &C.Metadata{
		NetWork: 0,
		Type:    0,
		SrcPort: "",
		DstPort: "80",
		Host:    remoteHost,
	}
	remoteConn, err := trojan.DialContext(ctx, meta)
	if err != nil {
		return 0, err
	}
	return pingInternal(remoteConn)
}

func PingSS(ssOption *outbound.ShadowSocksOption) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), tcpTimeout)
	defer cancel()
	ss, err := outbound.NewShadowSocks(ssOption)
	if err != nil {
		return 0, err
	}
	meta := &C.Metadata{
		NetWork: 0,
		Type:    0,
		SrcPort: "",
		DstPort: "80",
		Host:    remoteHost,
	}
	remoteConn, err := ss.DialContext(ctx, meta)
	if err != nil {
		return 0, err
	}
	return pingInternal(remoteConn)
}

func PingSSR(ssrOption *outbound.ShadowSocksROption) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), tcpTimeout)
	defer cancel()
	ssr, err := outbound.NewShadowSocksR(ssrOption)
	if err != nil {
		return 0, err
	}
	meta := &C.Metadata{
		NetWork: 0,
		Type:    0,
		SrcPort: "",
		DstPort: "80",
		Host:    remoteHost,
	}
	remoteConn, err := ssr.DialContext(ctx, meta)
	if err != nil {
		return 0, err
	}
	return pingInternal(remoteConn)
}

type PingResult struct {
	elapse int64
	err    error
}

func PingLink(link string, attempts int) (int64, error) {
	opt := PingOption{
		Attempts: attempts,
		TimeOut:  tcpTimeout,
	}
	return PingLinkInternal(link, opt)
}

func PingLinkInternal(link string, pingOption PingOption) (int64, error) {
	matches, err := utils.CheckLink(link)
	if err != nil {
		return 0, err
	}
	var option interface{}
	switch strings.ToLower(matches[1]) {
	case "vmess":
		option, err = config.VmessLinkToVmessOption(link)
	case "trojan":
		option, err = config.TrojanLinkToTrojanOption(link)
	case "http":
		option, err = config.HttpLinkToHttpOption(link)
	case "ss":
		option, err = config.SSLinkToSSOption(link)
	case "ssr":
		option, err = config.SSRLinkToSSROption(link)
	default:
		return 0, common.NewError("Not Suported Link")
	}
	if err != nil {
		return 0, err
	}
	var elapse int64
	var ip string
	var country string
	if pingOption.TimeOut > 0 {
		tcpTimeout = pingOption.TimeOut
	}
	log.Println("进入 PingLinkInternal")
	err = utils.ExponentialBackoff(pingOption.Attempts, 100).On(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), tcpTimeout)
		defer cancel()
		pingChan := make(chan PingResult, 1)
		go func(pingChan chan<- PingResult) {
			start := time.Now()
			elp, err := Ping(option)
			elapse := time.Since(start)
			if elapse > 2000*time.Second {
				elp = 0
			}
			pingResult := PingResult{elapse: elp, err: err}
			pingChan <- pingResult
		}(pingChan)
		for {
			select {
			case pingResult := <-pingChan:
				{
					elapse = pingResult.elapse
					log.Printf("节点延迟：%dms", elapse)
					if pingResult.err != nil {
						return pingResult.err
					}
					// ✅ 节点测试成功后，获取出口 IP 和国家
					info, err := getIPThroughOption(option)
					if err == nil {
						ip = info.Query
						country = info.Country
						log.Printf("节点出口IP：%s, 出口国家：%s", ip, country)
					} else {
						// 出口信息获取失败不影响整体成功
						ip = ""
						country = ""
					}
					return pingResult.err
				}
			case <-ctx.Done():
				return fmt.Errorf("ping time out")
			}
		}

	})
	return elapse, err
}

func PingContext(ctx context.Context, option interface{}) (int64, error) {
	var d outbound.ContextDialer
	var err error
	meta := &C.Metadata{
		NetWork: 0,
		Type:    C.TEST,
		SrcPort: "",
		DstPort: "80",
		Host:    remoteHost,
		Timeout: tcpTimeout,
	}
	if ssOption, ok := option.(*outbound.ShadowSocksOption); ok {
		d, err = outbound.NewShadowSocks(ssOption)
		if err != nil {
			return 0, err
		}
	}
	if ssrOption, ok := option.(*outbound.ShadowSocksROption); ok {
		d, err = outbound.NewShadowSocksR(ssrOption)
		if err != nil {
			return 0, err
		}
		dialerCtx := func(ctx context.Context, network, addr string) (net.Conn, error) {
			return d.DialContext(ctx, meta)
		}
		return pingHTTPClient(ctx, generate_204, tcpTimeout, dialerCtx)
	}
	if vmessOption, ok := option.(*outbound.VmessOption); ok {
		d, err = outbound.NewVmess(vmessOption)
		if err != nil {
			return 0, err
		}
	}
	if trojanOption, ok := option.(*outbound.TrojanOption); ok {
		d, err = outbound.NewTrojan(trojanOption)
		if err != nil {
			return 0, err
		}
	}
	if httpOption, ok := option.(*outbound.HttpOption); ok {
		d = outbound.NewHttp(*httpOption)
	}
	if d == nil {
		return 0, errors.New("not support config")
	}
	remoteConn, err := d.DialContext(ctx, meta)
	if err != nil {
		return 0, err
	}
	return pingInternal(remoteConn)
}

func Ping(option interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), tcpTimeout)
	defer cancel()
	return PingContext(ctx, option)
}

func pingHTTPClient(ctx context.Context, url string, timeout time.Duration, dialCtx func(ctx context.Context, network, addr string) (net.Conn, error)) (int64, error) {
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport, Timeout: timeout}
	if dialCtx != nil {
		httpTransport.DialContext = dialCtx
	}
	defer httpClient.CloseIdleConnections()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	start := time.Now()
	response, err := httpClient.Do(req)
	now := time.Now()
	if err != nil {
		return 0, err
	}
	elapse := now.Sub(start).Milliseconds()
	defer response.Body.Close()
	if response.StatusCode != 204 && response.StatusCode != 200 {
		return 0, fmt.Errorf("wrong status code %d", response.StatusCode)
	}
	return elapse, nil
}

func pingInternal(remoteConn net.Conn) (int64, error) {
	defer remoteConn.Close()
	remoteConn.SetDeadline(time.Now().Add(tcpTimeout))
	start := time.Now()
	// httpRequest := "GET /generate_204 HTTP/1.1\r\nHost: %s\r\nUser-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Safari/537.36\r\n\r\n"
	if _, err := remoteConn.Write(httpRequest); err != nil {
		return 0, err
	}
	buf := make([]byte, 128)
	_, err := remoteConn.Read(buf)
	if err != nil && err != io.EOF {
		return 0, err
	}
	_, err = parseFirstLine(buf)
	if err != nil {
		return 0, err
	}
	elapsed := time.Since(start).Milliseconds()
	// fmt.Print(string(buf))
	// fmt.Printf("server: %s port: %d elapsed: %d\n", vmessOption.Server, vmessOption.Port, elapsed)
	return elapsed, nil
}

func init() {
	resolver.DefaultResolver = dns.DefaultResolver()
}

type IPInfo struct {
	Query   string `json:"query"`   // IP
	Country string `json:"country"` // 国家
}

func GetIPThroughLink(id int, link string) (string, string) {

	var ip string = ""
	var country string = ""
	matches, err := utils.CheckLink(link)
	if err != nil {
		return ip, country
	}
	var option interface{}
	switch strings.ToLower(matches[1]) {
	case "vmess":
		option, err = config.VmessLinkToVmessOption(link)
	case "trojan":
		option, err = config.TrojanLinkToTrojanOption(link)
	case "http":
		option, err = config.HttpLinkToHttpOption(link)
	case "ss":
		option, err = config.SSLinkToSSOption(link)
	case "ssr":
		option, err = config.SSRLinkToSSROption(link)
	default:
		return ip, country
	}
	if err != nil {
		return ip, country
	}

	info, err := getIPThroughOption(option)
	if err == nil {
		ip = info.Query
		country = info.Country
		//log.Printf("%d 节点出口IP：%s, 出口国家：%s", id, ip, country)
	}
	return ip, country
}

func getIPThroughOption(option interface{}) (*IPInfo, error) {
	var d outbound.ContextDialer
	var err error
	meta := &C.Metadata{
		NetWork: 0,
		Type:    C.TEST,
		Host:    "ip-api.com",
		DstPort: "80",
		Timeout: tcpTimeout,
	}

	switch o := option.(type) {
	case *outbound.VmessOption:
		d, err = outbound.NewVmess(o)
	case *outbound.TrojanOption:
		d, err = outbound.NewTrojan(o)
	case *outbound.ShadowSocksOption:
		d, err = outbound.NewShadowSocks(o)
	case *outbound.ShadowSocksROption:
		d, err = outbound.NewShadowSocksR(o)
	case *outbound.HttpOption:
		d = outbound.NewHttp(*o)
	default:
		return nil, fmt.Errorf("unsupported option type")
	}
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return d.DialContext(ctx, meta)
			},
		},
		Timeout: 8 * time.Second,
	}

	req, err := http.NewRequest("GET", "http://ip-api.com/json?lang=zh-CN", nil)

	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var info IPInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

package cli

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/wenp5262/LiteSpeed/config"
	"github.com/wenp5262/LiteSpeed/utils"
)

var (
	ErrInvalidData = errors.New("invalid data")
)

// ParseReport describes parsed links and skipped lines (filtered before any test).
type ParseReport struct {
	Links      []string
	Skipped    []SkippedLine
	TotalLines int
}

type SkippedLine struct {
	Line   string `json:"line"`
	Reason string `json:"reason"`
}

// ParseLinksReport is like ParseLinks but also returns which lines were filtered out during parsing.
func ParseLinksReport(message string) (*ParseReport, error) {
	opt := ParseOption{Type: PARSE_ANY}
	return ParseLinksReportWithOption(message, opt)
}

func ParseLinksReportWithOption(message string, opt ParseOption) (*ParseReport, error) {
	message = strings.TrimSpace(message)

	// URL
	if opt.Type == PARSE_URL || utils.IsUrl(message) {
		links, skipped, total, err := getSubscriptionLinksReport(message)
		if err != nil {
			return nil, err
		}
		return &ParseReport{Links: links, Skipped: skipped, TotalLines: total}, nil
	}
	// File
	if opt.Type == PARSE_FILE || utils.IsFilePath(message) {
		links, skipped, total, err := parseFileReport(message)
		if err != nil {
			return nil, err
		}
		return &ParseReport{Links: links, Skipped: skipped, TotalLines: total}, nil
	}
	// Base64 explicit
	if opt.Type == PARSE_BASE64 {
		links, skipped, total, err := parseBase64Report(message)
		if err != nil {
			return nil, err
		}
		return &ParseReport{Links: links, Skipped: skipped, TotalLines: total}, nil
	}

	// Try decode as base64 subscription
	if decoded, err := utils.DecodeB64(message); err == nil {
		links, skipped, total := parseProfilesReport(decoded)
		return &ParseReport{Links: links, Skipped: skipped, TotalLines: total}, nil
	}

	// Clash yaml / raw profile text
	if strings.Contains(message, "proxies:") {
		links, err := parseClash(message)
		if err != nil {
			return nil, err
		}
		return &ParseReport{Links: links, Skipped: nil, TotalLines: len(links)}, nil
	}
	if strings.Contains(message, "vmess://") ||
		strings.Contains(message, "trojan://") ||
		strings.Contains(message, "ssr://") ||
		strings.Contains(message, "ss://") {
		links, skipped, total := parseProfilesReport(message)
		return &ParseReport{Links: links, Skipped: skipped, TotalLines: total}, nil
	}
	return nil, ErrInvalidData
}

// ParseLinks parses subscription input which can be: URL, file path, base64, clash yaml, or raw profiles text.
func ParseLinks(message string) ([]string, error) {
	opt := ParseOption{Type: PARSE_ANY}
	return ParseLinksWithOption(message, opt)
}

type PAESE_TYPE int

const (
	PARSE_ANY PAESE_TYPE = iota
	PARSE_URL
	PARSE_FILE
	PARSE_BASE64
	PARSE_CLASH
	PARSE_PROFILE
)

type ParseOption struct {
	Type PAESE_TYPE
}

func ParseLinksWithOption(message string, opt ParseOption) ([]string, error) {
	message = strings.TrimSpace(message)

	// URL
	if opt.Type == PARSE_URL || utils.IsUrl(message) {
		return getSubscriptionLinks(message)
	}
	// File
	if opt.Type == PARSE_FILE || utils.IsFilePath(message) {
		return parseFile(message)
	}
	// Base64 explicit
	if opt.Type == PARSE_BASE64 {
		return parseBase64(message)
	}

	// Try decode as base64 subscription
	if decoded, err := utils.DecodeB64(message); err == nil {
		return parseProfiles(decoded)
	}

	// Clash yaml / raw profile text
	if strings.Contains(message, "proxies:") {
		return parseClash(message)
	}
	if strings.Contains(message, "vmess://") ||
		strings.Contains(message, "trojan://") ||
		strings.Contains(message, "ssr://") ||
		strings.Contains(message, "ss://") {
		return parseProfiles(message)
	}
	return nil, ErrInvalidData
}

func getSubscriptionLinks(link string) ([]string, error) {
	c := http.Client{Timeout: 20 * time.Second}
	resp, err := c.Get(link)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if isYamlFile(link) {
		return scanClashProxies(resp.Body)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	dataStr := string(data)

	decoded, err := utils.DecodeB64(dataStr)
	if err != nil {
		// Not base64: try treat as clash / profiles
		if strings.Contains(dataStr, "proxies:") {
			return parseClash(dataStr)
		}
		if strings.Contains(dataStr, "vmess://") ||
			strings.Contains(dataStr, "trojan://") ||
			strings.Contains(dataStr, "ssr://") ||
			strings.Contains(dataStr, "ss://") {
			return parseProfiles(dataStr)
		}
		return []string{}, err
	}
	return parseProfiles(decoded)
}

func parseBase64(message string) ([]string, error) {
	msg, err := utils.DecodeB64(message)
	if err != nil {
		return nil, err
	}
	return parseProfiles(msg)
}

func parseProfilesReport(message string) ([]string, []SkippedLine, int) {
	message = strings.ReplaceAll(message, "\r\n", "\n")
	lines := strings.Split(message, "\n")
	links := make([]string, 0, len(lines))
	skipped := make([]SkippedLine, 0)
	total := 0
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		total++
		ll := strings.ToLower(l)
		if strings.HasPrefix(ll, "vmess://") ||
			strings.HasPrefix(ll, "trojan://") ||
			strings.HasPrefix(ll, "ssr://") ||
			strings.HasPrefix(ll, "ss://") ||
			strings.HasPrefix(ll, "vless://") ||
			strings.HasPrefix(ll, "http://") ||
			strings.HasPrefix(ll, "https://") {
			links = append(links, l)
			continue
		}
		reason := "unsupported scheme"
		if strings.Contains(ll, "://") {
			reason = "unsupported scheme: " + strings.SplitN(ll, "://", 2)[0]
		}
		skipped = append(skipped, SkippedLine{Line: l, Reason: reason})
	}
	return links, skipped, total
}

func parseBase64Report(message string) ([]string, []SkippedLine, int, error) {
	msg, err := utils.DecodeB64(message)
	if err != nil {
		return nil, nil, 0, err
	}
	links, skipped, total := parseProfilesReport(msg)
	return links, skipped, total, nil
}

func parseFileReport(filepath string) ([]string, []SkippedLine, int, error) {
	filepath = strings.TrimSpace(filepath)
	if _, err := os.Stat(filepath); err != nil {
		return nil, nil, 0, err
	}
	// clash yaml
	if isYamlFile(filepath) {
		links, err := parseClashFileByLine(filepath)
		if err != nil {
			return nil, nil, 0, err
		}
		return links, nil, len(links), nil
	}
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, nil, 0, err
	}
	// try base64 first
	links, skipped, total, err := parseBase64Report(string(data))
	if err == nil {
		return links, skipped, total, nil
	}
	// fallback to raw profiles
	links2, skipped2, total2 := parseProfilesReport(string(data))
	return links2, skipped2, total2, nil
}

func getSubscriptionLinksReport(link string) ([]string, []SkippedLine, int, error) {
	c := http.Client{Timeout: 20 * time.Second}
	resp, err := c.Get(link)
	if err != nil {
		return nil, nil, 0, err
	}
	defer resp.Body.Close()

	if isYamlFile(link) {
		links, err := scanClashProxies(resp.Body)
		if err != nil {
			return nil, nil, 0, err
		}
		return links, nil, len(links), nil
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, 0, err
	}
	dataStr := string(data)

	decoded, err := utils.DecodeB64(dataStr)
	if err != nil {
		// Not base64: try treat as clash / profiles
		if strings.Contains(dataStr, "proxies:") {
			links, err := parseClash(dataStr)
			if err != nil {
				return nil, nil, 0, err
			}
			return links, nil, len(links), nil
		}
		links, skipped, total := parseProfilesReport(dataStr)
		return links, skipped, total, nil
	}
	links, skipped, total := parseProfilesReport(decoded)
	return links, skipped, total, nil
}

func parseProfiles(message string) ([]string, error) {
	// split by lines, keep supported schemes
	message = strings.ReplaceAll(message, "\r\n", "\n")
	lines := strings.Split(message, "\n")
	links := make([]string, 0, len(lines))
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		if strings.HasPrefix(strings.ToLower(l), "vmess://") ||
			strings.HasPrefix(strings.ToLower(l), "trojan://") ||
			strings.HasPrefix(strings.ToLower(l), "ssr://") ||
			strings.HasPrefix(strings.ToLower(l), "ss://") ||
			strings.HasPrefix(strings.ToLower(l), "vless://") ||
			strings.HasPrefix(strings.ToLower(l), "http://") ||
			strings.HasPrefix(strings.ToLower(l), "https://") {
			links = append(links, l)
		}
	}
	return links, nil
}

func parseClash(input string) ([]string, error) {
	return parseClashByte([]byte(input))
}

func parseClashByte(data []byte) ([]string, error) {
	cc, err := config.ParseClash(data)
	if err != nil {
		return nil, err
	}
	return cc.Proxies, nil
}

func parseFile(filepath string) ([]string, error) {
	filepath = strings.TrimSpace(filepath)
	if _, err := os.Stat(filepath); err != nil {
		return nil, err
	}
	// clash yaml
	if isYamlFile(filepath) {
		return parseClashFileByLine(filepath)
	}
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	links, err := parseBase64(string(data))
	if err != nil && len(data) > 2048 {
		preview := string(data[:2048])
		if strings.Contains(preview, "proxies:") {
			return scanClashProxies(bytes.NewReader(data))
		}
		if strings.Contains(preview, "vmess://") ||
			strings.Contains(preview, "trojan://") ||
			strings.Contains(preview, "ssr://") ||
			strings.Contains(preview, "ss://") {
			return parseProfiles(string(data))
		}
	}
	if err != nil {
		return parseProfiles(string(data))
	}
	return links, nil
}

func scanClashProxies(r io.Reader) ([]string, error) {
	proxiesStart := false
	var data []byte
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		b := scanner.Bytes()
		trimLine := strings.TrimSpace(string(b))
		if trimLine == "proxy-groups:" || trimLine == "rules:" || trimLine == "Proxy Group:" {
			break
		}
		if !proxiesStart && (trimLine == "proxies:" || trimLine == "Proxy:") {
			proxiesStart = true
			b = []byte("proxies:")
		}
		if proxiesStart {
			if _, err := config.ParseBaseProxy(trimLine); err != nil {
				continue
			}
			data = append(data, b...)
			data = append(data, byte('\n'))
		}
	}
	return parseClashByte(data)
}

func parseClashFileByLine(filepath string) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return scanClashProxies(file)
}

func isYamlFile(path string) bool {
	return strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml")
}

// Some subscriptions are base64-encoded but missing padding. Provide a fallback for raw base64 block.
func DecodeBase64Std(s string) (string, error) {
	s = strings.TrimSpace(s)
	// normalize URL-safe base64
	s = strings.ReplaceAll(s, "-", "+")
	s = strings.ReplaceAll(s, "_", "/")
	if m := len(s) % 4; m != 0 {
		s += strings.Repeat("=", 4-m)
	}
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

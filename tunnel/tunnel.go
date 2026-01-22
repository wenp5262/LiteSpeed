package tunnel

import (
    "fmt"
    "net"
    "strconv"
    "strings"
)

type AddrType int

const (
    IPv4 AddrType = iota
    IPv6
    DomainName
)

type Address struct {
    Type AddrType
    Host string
    Port int
}

func (a *Address) String() string {
    return net.JoinHostPort(a.Host, strconv.Itoa(a.Port))
}

// NewAddressFromAddr parses a typical "host:port" address.
// This is a minimal implementation to satisfy compilation in trimmed mode.
func NewAddressFromAddr(network, address string) (*Address, error) {
    host, portStr, err := net.SplitHostPort(address)
    if err != nil {
        // handle raw host without port
        if strings.Contains(address, ":") {
            return nil, err
        }
        return nil, fmt.Errorf("missing port in address: %s", address)
    }
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return nil, err
    }
    t := DomainName
    ip := net.ParseIP(host)
    if ip != nil {
        if ip.To4() != nil {
            t = IPv4
        } else {
            t = IPv6
        }
    }
    return &Address{Type: t, Host: host, Port: port}, nil
}

// The following types are intentionally minimal. They exist to keep package compilation successful
// for the trimmed, country-only connectivity check build.

type Tunnel interface{}

type Conn interface{}

type Server interface {
    Close() error
}

type Client interface {
    Close() error
}

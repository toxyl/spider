package services

import (
	"encoding/base64"
	"fmt"
	"net"
	"strconv"

	gostringgenerator "github.com/toxyl/go-string-generator"
	"github.com/toxyl/go-string-generator/tokens"
	"github.com/toxyl/spider/random"
)

var (
	stringGenerator *tokens.RandomStringGenerator
	DataDir         = ""
)

func Spider2name(port int) string {
	for name, service := range *services {
		for _, p := range service.Ports {
			if p == port {
				return fmt.Sprintf("%s (%d)", name, port)
			}
		}
	}
	return fmt.Sprintf("Port %d", port)
}

func Banner(port int) string {
	for _, service := range *services {
		for _, p := range service.Ports {
			if p == port {
				if stringGenerator == nil {
					stringGenerator = gostringgenerator.NewGenerator(DataDir, func(err error) {})
				}
				return stringGenerator.Generate(service.Banner)
			}
		}
	}
	return base64.RawStdEncoding.EncodeToString([]byte(random.Taunt()))
}

// extractPort takes a "IP:Port" formatted string and returns the port as integer.
// If extraction fails the function returns 0.
func extractPort(ipv4Addr string) int {
	_, port, err := net.SplitHostPort(ipv4Addr)
	if err != nil {
		return 0
	}
	p, err := strconv.Atoi(port)
	if err != nil {
		p = 0
	}
	return p
}

// extractHost takes a "IP:Port" formatted string and returns the ip portion.
// If extraction fails the function returns an empty string.
func extractHost(ipv4Addr string) string {
	ip, _, err := net.SplitHostPort(ipv4Addr)
	if err != nil {
		return ""
	}
	return ip
}

func Host(addr net.Addr) string {
	return extractHost(addr.String())
}

func Conn2spider(conn net.Conn) int {
	return extractPort(conn.LocalAddr().String())
}

func Conn2prey(conn net.Conn) string {
	return Host(conn.RemoteAddr())
}

func Conn2banner(conn net.Conn) string {
	return Banner(Conn2spider(conn))
}

func Conn2name(conn net.Conn) string {
	return Spider2name(Conn2spider(conn))
}

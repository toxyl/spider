package main

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/toxyl/glog"
	gostringgenerator "github.com/toxyl/go-string-generator"
)

var (
	gen = gostringgenerator.NewGenerator("/etc/spider/data/", func(err error) {})
)

func getSpiderMetricName(spider int, section string) string {
	return fmt.Sprintf("spider_%d_%s", spider, section)
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

func port(addr net.Addr) int {
	return extractPort(addr.String())
}

func host(addr net.Addr) string {
	return extractHost(addr.String())
}

func randomStringFromList(strings ...string) string {
	if len(strings) <= 0 {
		return ""
	}
	var i int = randomInt(0, len(strings)-1)
	return strings[i]
}

func randomLinebreak() string {
	return randomStringFromList("\r\n", "\r", "\n")
}

func fileDelete(path string) error {
	return os.Remove(path)
}

func fileExists(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()
	_, err = file.Stat()
	return err == nil
}

func fileRead(path string) (string, error) {
	if !fileExists(path) {
		return "", fmt.Errorf("file %s does not exist", path)
	}
	bytes, err := os.ReadFile(path)
	return string(bytes), err
}

func fileWrite(path, content string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("could not create %s file", path)
	}
	return os.WriteFile(f.Name(), []byte(content), 0644)
}

func randomInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	n := max - min + 1
	if n <= 0 {
		return min
	}
	return rand.Intn(n) + min
}

// garbageString produces a string (length is randomly chosen between 1 and `n`)
// consisting of random (non)-printable characters.
func garbageString(n int) string {
	garbage := make([]byte, randomInt(1, n))
	rand.Read(garbage)
	return string(garbage)
}

func spider2name(port int) string {
	for name, service := range *services {
		for _, p := range service.Ports {
			if p == port {
				return fmt.Sprintf("%s (%d)", name, port)
			}
		}
	}
	return fmt.Sprintf("Port %d", port)
}

func randomTaunt() string {
	return gen.Generate(randomStringFromList(config.Taunts...))
}

func banner(port int) string {
	for _, service := range *services {
		for _, p := range service.Ports {
			if p == port {
				return gen.Generate(service.Banner)
			}
		}
	}
	return base64.RawStdEncoding.EncodeToString([]byte(randomTaunt()))
}

func conn2spider(conn net.Conn) int {
	return port(conn.LocalAddr())
}

func conn2prey(conn net.Conn) string {
	return host(conn.RemoteAddr())
}

func conn2banner(conn net.Conn) string {
	return banner(conn2spider(conn))
}

func conn2name(conn net.Conn) string {
	return spider2name(conn2spider(conn))
}

func connWrite(conn net.Conn, msg string) error {
	_, err := conn.Write([]byte(msg))
	return err
}

func colorizeAction(action string) string {
	if action == "kills" {
		action = glog.WrapRed(action)
	} else if action == "attacks" {
		action = glog.WrapOrange(action)
	} else if action == "pokes" {
		action = glog.WrapYellow(action)
	} else {
		action = glog.Auto(action)
	}
	return action
}

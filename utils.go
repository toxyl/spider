package main

import (
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"
)

// ExtractPort takes a "IP:Port" formatted string and returns the port as integer.
// If extraction fails the function returns 0.
func ExtractPort(ipv4Addr string) int {
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

// ExtractHost takes a "IP:Port" formatted string and returns the ip portion.
// If extraction fails the function returns an empty string.
func ExtractHost(ipv4Addr string) string {
	ip, _, err := net.SplitHostPort(ipv4Addr)
	if err != nil {
		return ""
	}
	return ip
}

func ExtractPortFromAddr(addr net.Addr) int {
	return ExtractPort(addr.String())
}

func ExtractHostFromAddr(addr net.Addr) string {
	return ExtractHost(addr.String())
}

func GetRandomStringFromList(strings ...string) string {
	if len(strings) <= 0 {
		return ""
	}
	var i int = randomInt(0, len(strings)-1)
	return strings[i]
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

func StringToInt(s string, defaultValue int) int {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		i = int64(defaultValue)
	}
	return int(i)
}

func randomInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	n := max - min + 1
	if n <= 0 {
		return min
	}
	return rand.Intn(n) + min
}

// GenerateGarbageString produces a string (length is randomly chosen between 1 and `n`)
// consisting of random (non)-printable characters.
func GenerateGarbageString(n int) string {
	garbage := make([]byte, randomInt(1, n))
	rand.Read(garbage)
	return string(garbage)
}

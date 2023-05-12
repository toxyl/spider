package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/toxyl/glog"
)

var (
	log    = glog.NewLoggerSimple("Spider")
	config *Config
)

func info(conn net.Conn, action string) {
	log.Info("Spider at %s %s prey %s", glog.Int(ExtractPortFromAddr(conn.LocalAddr())), glog.Auto(action), glog.ConnRemote(conn, false))
}

func randomTaunt() string {
	return GetRandomStringFromList(config.Taunts...)
}

func attackPrey(conn net.Conn) {
	if conn == nil {
		return
	}
	info(conn, "attacks")
	n := 0
	attack := true
	t := time.Now().Add(time.Duration(config.AttackLength) * time.Second)
	for attack {
		if conn == nil {
			break
		}
		n = randomInt(100, 1000)
		_, _ = conn.Write([]byte(GenerateGarbageString(n)))
		attack = t.After(time.Now()) || n%7 != 0 || n%5 != 0
	}
}

func killPrey(conn net.Conn) {
	if conn == nil {
		return
	}

	info(conn, "kills")
	_, _ = conn.Write([]byte("\n" + randomTaunt() + "\n"))
	time.Sleep(5 * time.Second)
	conn.Close()
}

func catchPrey(conn net.Conn) {
	attackPrey(conn)
	killPrey(conn)
}

func buildWebs() {
	for _, port := range config.Ports {
		srv := fmt.Sprintf("%s:%d", config.Host, port)

		listener, err := net.Listen("tcp", srv)
		if err != nil {
			log.NotOK("Spider backs off, someone is already at %s...", glog.Auto(port))
			continue
		}
		log.Default("Spider at %s builds web...", glog.Auto(port))
		go func() {
			for {
				conn, err := listener.Accept()
				if err != nil {
					info(conn, "fails to catch")
					log.Error("%s", glog.Error(err))
					conn.Close()
					continue
				}
				host := ExtractHostFromAddr(conn.RemoteAddr())
				isWhitelisted := false
				for _, wl := range config.Whitelist {
					if host == wl {
						isWhitelisted = true
						break
					}
				}
				if isWhitelisted {
					info(conn, "avoids")
					conn.Close()
					continue
				}
				go catchPrey(conn)
			}
		}()
	}
}

func init() {
	glog.LoggerConfig.ShowRuntimeMilliseconds = false
	glog.LoggerConfig.ShowSubsystem = false
	glog.LoggerConfig.ShowRuntimeSeconds = false
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: spider [config file]")
		return
	}

	c, err := LoadConfig(os.Args[1])
	if err != nil {
		log.Error("Failed to load config: %s", glog.Error(err))
	}
	config = c

	buildWebs()

	for {
		time.Sleep(15 * time.Second) // to keep the application open
	}
}

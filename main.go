package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/toxyl/glog"
)

var (
	log      = glog.NewLoggerSimple("Spider")
	config   *Config
	services *Services
)

func info(conn net.Conn, action string) {
	log.Info("%s spider %s prey %s", glog.Auto(GetSpiderName(ExtractPortFromAddr(conn.LocalAddr()))), glog.Auto(action), glog.ConnRemote(conn, false))
}

func randomTaunt() string {
	return GetRandomStringFromList(config.Taunts...)
}

func attackPrey(conn net.Conn) {
	if conn == nil {
		return
	}
	info(conn, "attacks")
	t := time.Now().Add(time.Duration(config.AttackLength) * time.Second)
	_, _ = conn.Write([]byte(GetSpiderBanner(ExtractPortFromAddr(conn.LocalAddr())) + "\n")) // we first send a "proper" banner
	time.Sleep(5 * time.Second)                                                              // and then sleep a bit so our target has some time to process the banner
	for t.After(time.Now()) {
		if conn == nil {
			break
		}
		_, _ = conn.Write([]byte(GenerateGarbageString(randomInt(100, 10000))))
		time.Sleep(time.Duration(randomInt(1, 5)) * time.Second)
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
			log.NotOK("%s spider backs off, someone is already there...", glog.Auto(GetSpiderName(port)))
			continue
		}
		log.Default("%s spider builds web...", glog.Auto(GetSpiderName(port)))
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

	s, err := LoadServices()
	if err != nil {
		log.Error("Failed to load services: %s", glog.Error(err))
	}
	services = s

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

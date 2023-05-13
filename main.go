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
	stats    = NewStats()
)

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

func connWrite(conn net.Conn, msg string) {
	_, _ = conn.Write([]byte(msg))
}

func connInfo(conn net.Conn, action string) {
	if action == "kills" {
		action = glog.WrapRed(action)
	} else if action == "attacks" {
		action = glog.WrapOrange(action)
	} else {
		action = glog.Auto(action)
	}
	log.Info("%s spider %s prey %s", glog.Auto(conn2name(conn)), action, glog.Auto(conn2prey(conn)))
}

func colorizeAction(action string) string {
	if action == "kills" {
		action = glog.WrapRed(action)
	} else if action == "attacks" {
		action = glog.WrapOrange(action)
	} else {
		action = glog.Auto(action)
	}
	return action
}

func spiderInfo(spider int, action string, suffix string) {
	log.Info("%s spider %s%s", glog.Auto(spider2name(spider)), colorizeAction(action), suffix)
}

func spiderNotOK(spider int, action string, suffix string) {
	log.NotOK("%s spider %s%s", glog.Auto(spider2name(spider)), colorizeAction(action), suffix)
}

func spiderOK(spider int, action string, suffix string) {
	log.OK("%s spider %s%s", glog.Auto(spider2name(spider)), colorizeAction(action), suffix)
}

func randomTaunt() string {
	return randomStringFromList(config.Taunts...)
}

func attackPrey(conn net.Conn) {
	if conn == nil {
		return
	}
	connInfo(conn, "attacks")
	stats.AddPrey(conn2spider(conn))
	defer stats.AddKill(conn2spider(conn))

	t := time.Now().Add(time.Duration(config.AttackLength) * time.Second)
	connWrite(conn, conn2banner(conn)+"\n") // we first send a "proper" banner
	time.Sleep(5 * time.Second)             // and then sleep a bit so our target has some time to process the banner
	for t.After(time.Now()) {
		if conn == nil {
			break
		}
		connWrite(conn, garbageString(10000))
		time.Sleep(time.Duration(randomInt(1, 5)) * time.Second)
	}
}

func killPrey(conn net.Conn) {
	if conn == nil {
		return
	}

	connInfo(conn, "kills")
	connWrite(conn, "\n"+randomTaunt()+"\n")
	time.Sleep(5 * time.Second)
	conn.Close()
}

func catchPrey(conn net.Conn) {
	attackPrey(conn)
	killPrey(conn)
}

func buildWebs() {
	for _, spider := range config.Spiders {
		srv := fmt.Sprintf("%s:%d", config.Host, spider)

		listener, err := net.Listen("tcp", srv)
		if err != nil {
			spiderNotOK(spider, "backs off", "someone is already there...")
			continue
		}
		spiderInfo(spider, "builds", " web...")
		stats.AddSpider(spider)
		go func() {
			for {
				conn, err := listener.Accept()
				if err != nil {
					connInfo(conn, "fails to catch")
					log.Error("%s", glog.Error(err))
					conn.Close()
					continue
				}
				host := host(conn.RemoteAddr())
				isWhitelisted := false
				for _, wl := range config.Whitelist {
					if host == wl {
						isWhitelisted = true
						break
					}
				}
				if isWhitelisted {
					connInfo(conn, "avoids")
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
		stats.Print()
		time.Sleep(1 * time.Minute)
	}
}

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

func connInfo(conn net.Conn, action string) {
	log.Info("%s spider %s prey %s", glog.Auto(conn2name(conn)), colorizeAction(action), glog.Auto(conn2prey(conn)))
}

func spiderInfo(spider int, action string, suffix string) {
	log.Info("%s spider %s%s", glog.Auto(spider2name(spider)), colorizeAction(action), suffix)
}

func spiderNotOK(spider int, action string, suffix string) {
	log.NotOK("%s spider %s%s", glog.Auto(spider2name(spider)), colorizeAction(action), suffix)
}

func randomTaunt() string {
	return gen.Generate(randomStringFromList(config.Taunts...))
}

func attackPrey(conn net.Conn) {
	t := time.Now().Add(time.Duration(config.AttackLength) * time.Second)
	err := connWrite(conn, conn2banner(conn)+"\n") // we first send a "proper" banner
	if err != nil {
		return
	}
	time.Sleep(5 * time.Second) // and then sleep a bit so our target has some time to process the banner
	for t.After(time.Now()) {
		if conn == nil {
			break
		}
		err = connWrite(conn, garbageString(10000))
		if err != nil {
			return
		}
		time.Sleep(time.Duration(randomInt(1, 5)) * time.Second)
	}
}

func killPrey(conn net.Conn) {
	err := connWrite(conn, "\n"+randomTaunt()+"\n")
	if err != nil {
		return
	}
	time.Sleep(5 * time.Second)
	conn.Close()
}

func catchPrey(conn net.Conn) {
	if conn == nil {
		return
	}
	spider := conn2spider(conn)
	connInfo(conn, "attacks")
	stats.AddPrey(spider)
	defer func() {
		connInfo(conn, "kills")
		killPrey(conn)
		stats.AddKill(spider)
	}()
	attackPrey(conn)
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

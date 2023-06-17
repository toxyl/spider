package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/toxyl/glog"
)

var (
	log        = glog.NewLoggerSimple("Spider")
	config     *Config
	services   *Services
	stats      *Stats
	pathConfig string
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
	spider := conn2spider(conn)
	timeAttackStart := time.Now()
	timeAttackEnd := timeAttackStart.Add(time.Duration(config.AttackLength) * time.Second)
	defer func() {
		connInfo(conn, "kills")
		connWrite(conn, randomLinebreak()+randomTaunt()+randomLinebreak())
		time.Sleep(5 * time.Second)
		conn.Close()
		stats.AddKill(spider, float64(time.Now().Unix()-timeAttackStart.Unix()))
	}()

	connInfo(conn, "pokes")
	stats.AddPrey(spider)
	err := connWrite(conn, conn2banner(conn)+randomLinebreak()) // we first send a "proper" banner
	if err != nil {
		return
	}
	time.Sleep(5 * time.Second) // and then sleep a bit so our target has some time to process the banner

	connInfo(conn, "attacks")
	for timeAttackEnd.After(time.Now()) {
		err = connWrite(conn, garbageString(10000))
		if err != nil {
			return
		}
		time.Sleep(time.Duration(randomInt(1, 5)) * time.Second)
	}
}

func catchPrey(conn net.Conn) {
	if conn == nil {
		return
	}
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

	err := LoadServices()
	if err != nil {
		log.Error("Failed to load services: %s", glog.Error(err))
	}

	pathConfig, err = filepath.Abs(os.Args[1])
	if err != nil {
		log.Error("Failed to generate absolute path for config: %s", glog.Error(err))
	}

	err = LoadConfig(pathConfig)
	if err != nil {
		log.Error("Failed to load config: %s", glog.Error(err))
	}

	stats = NewStats()

	buildWebs()

	for {
		stats.Print()
		time.Sleep(1 * time.Minute)
	}
}

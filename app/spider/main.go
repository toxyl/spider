package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/toxyl/glog"
	metrics "github.com/toxyl/metric-nexus"
	"github.com/toxyl/spider/log"
	"github.com/toxyl/spider/random"
	"github.com/toxyl/spider/services"
	"github.com/toxyl/spider/utils"
)

var (
	statistics *Stats
)

func attackPrey(conn net.Conn) {
	spider := services.Conn2spider(conn)
	timeAttackStart := time.Now()
	timeAttackEnd := timeAttackStart.Add(time.Duration(config.AttackLength) * time.Second)

	statistics.AddPrey(spider)
	defer func(conn net.Conn) {
		if conn != nil {
			SpiderPreyInfo(conn, "kills")
			_ = utils.ConnWrite(conn, random.Linebreak()+random.Taunt()+random.Linebreak())
			time.Sleep(5 * time.Second)
			_ = conn.Close()
		}
		statistics.AddKill(spider, float64(time.Now().Unix()-timeAttackStart.Unix()))
	}(conn)
	SpiderPreyInfo(conn, "pokes")
	err := utils.ConnWrite(conn, services.Conn2banner(conn)+random.Linebreak()) // we first send a "proper" banner
	if err != nil {
		return
	}
	time.Sleep(5 * time.Second) // and then sleep a bit so our target has some time to process the banner

	SpiderPreyInfo(conn, "attacks")
	for timeAttackEnd.After(time.Now()) {
		err = utils.ConnWrite(conn, random.Garbage(10000))
		if err != nil {
			return
		}
		t := random.Int(1, 5)
		if t == 3 { // i.e. 1 in 5 chance to insert a linebreak
			_ = utils.ConnWrite(conn, random.Linebreak())
		}
		time.Sleep(time.Duration(t) * time.Second)
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
			SpiderFailed(spider, "backs off", "someone is already there...")
			continue
		}
		SpiderInfo(spider, "builds", " web...")
		statistics.AddSpider(spider)
		go func() {
			for {
				conn, err := listener.Accept()
				if err != nil {
					SpiderPreyInfo(conn, "fails to catch")
					log.Error("%s", glog.Error(err))
					_ = conn.Close()
					continue
				}
				host := services.Host(conn.RemoteAddr())
				isWhitelisted := false
				for _, wl := range config.Whitelist {
					if host == wl {
						isWhitelisted = true
						break
					}
				}
				if isWhitelisted {
					SpiderPreyInfo(conn, "avoids")
					_ = conn.Close()
					continue
				}
				go catchPrey(conn)
			}
		}()
	}
}

func init() {
	glog.LoggerConfig.ShowIndicator = true
	glog.LoggerConfig.ShowSubsystem = false
	glog.LoggerConfig.ShowDateTime = true
	glog.LoggerConfig.ShowRuntimeMilliseconds = false
	glog.LoggerConfig.ShowRuntimeSeconds = false
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Failed to retrieve the executable path:", err)
		return
	}

	random.DataDir = filepath.Join(filepath.Dir(exePath), "data")
	if !utils.FileExists(random.DataDir) {
		_ = os.MkdirAll(random.DataDir, 0755)
	}
	services.DataDir = random.DataDir
	random.Taunts = config.Taunts
}

func getMemoryUsage() (alloc, sys, heapAlloc, heapSys int) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return int(m.Alloc), int(m.Sys), int(m.HeapAlloc), int(m.HeapSys)
}

func trackUptimeAndMemoryUsage() {
	client := metrics.NewClient(config.MetricNexus.Host, config.MetricNexus.Port, config.MetricNexus.Key, true)
	client.Create(utils.GetMetricName(0, "uptime"), "How long the cluster has been online. Note that this is the sum of all nodes.")
	client.Create(utils.GetMetricName(0, "mem_alloc"), "How much memory the cluster actively uses. Note that this is the sum of all nodes.")
	client.Create(utils.GetMetricName(0, "mem_sys"), "How much memory the cluster has requested from the OS. Note that this is the sum of all nodes.")
	client.Create(utils.GetMetricName(0, "mem_heap_alloc"), "How much heap memory the cluster actively uses. Note that this is the sum of all nodes.")
	client.Create(utils.GetMetricName(0, "mem_heap_sys"), "How much heap memory the cluster has requested from the OS. Note that this is the sum of all nodes.")

	t := time.Now()
	lastAlloc, lastSys, lastHeapAlloc, lastHeapSys := getMemoryUsage()
	for {
		time.Sleep(10 * time.Second)
		alloc, sys, heapAlloc, heapSys := getMemoryUsage()
		client.Add(utils.GetMetricName(0, "mem_alloc"), alloc-lastAlloc)
		client.Add(utils.GetMetricName(0, "mem_sys"), sys-lastSys)
		client.Add(utils.GetMetricName(0, "mem_heap_alloc"), heapAlloc-lastHeapAlloc)
		client.Add(utils.GetMetricName(0, "mem_heap_sys"), heapSys-lastHeapSys)
		lastAlloc = alloc
		lastSys = sys
		lastHeapAlloc = heapAlloc
		lastHeapSys = heapSys
		client.Add(utils.GetMetricName(0, "uptime"), time.Since(t).Seconds())
		t = time.Now()
	}
}

func main() {
	statistics = NewStats()
	statistics.AddHost()
	buildWebs()
	trackUptimeAndMemoryUsage()
}

package main

import (
	"os"
	"strings"
	"sync"

	"github.com/toxyl/glog"
	metrics "github.com/toxyl/metric-nexus"
	"github.com/toxyl/spider/log"
	stats "github.com/toxyl/spider/stats"
	"github.com/toxyl/spider/utils"
)

type Stats struct {
	lock    *sync.Mutex
	spiders map[int]*stats.SpiderStats
	client  *metrics.Client
}

func (s *Stats) AddHost() {
	s.lock.Lock()
	defer s.lock.Unlock()
	file := utils.GetMetricFileName(0, "hosts")
	hdir, err := os.UserHomeDir()
	if err != nil {
		log.Error("Failed to get user home dir for %s: %s", glog.File(file), glog.Error(err))
	}
	file = strings.ReplaceAll(file, "~/", hdir+"/")

	if !utils.FileExists(file) {
		met := utils.GetMetricName(0, "hosts")
		s.client.Create(met, "How many hosts handle this port.")
		s.client.Add(met, 1)
		err := os.WriteFile(file, []byte{'1'}, 0644)
		if err != nil {
			log.Error("Failed to create %s: %s", glog.File(file), glog.Error(err))
		}
	}
}

func (s *Stats) AddSpider(spider int) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.spiders[spider]; ok {
		return // we already have this spider
	}
	s.spiders[spider] = stats.NewSpiderStats(spider, config.MetricNexus.Host, config.MetricNexus.Port, config.MetricNexus.Key)
}

func (s *Stats) AddKill(spider int, timeWasted float64) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.spiders[spider]; !ok {
		return // we can't add kills to imaginary spiders
	}
	s.spiders[spider].RemovePrey(timeWasted)
}

func (s *Stats) AddPrey(spider int) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.spiders[spider]; !ok {
		return // we can't add prey to imaginary spiders
	}
	s.spiders[spider].AddPrey()
}

func NewStats() *Stats {
	s := &Stats{
		lock:    &sync.Mutex{},
		spiders: map[int]*stats.SpiderStats{},
		client:  metrics.NewClient(config.MetricNexus.Host, config.MetricNexus.Port, config.MetricNexus.Key, true),
	}
	return s
}

package main

import (
	"os"
	"sync"

	metrics "github.com/toxyl/metric-nexus"
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
	fHostRegistered := utils.GetMetricFileName(0, "hosts")

	met := utils.GetMetricName(0, "hosts")
	s.client.Create(met, "How many hosts handle this port.")

	if !utils.FileExists(fHostRegistered) {
		s.client.Add(met, 1)
		_ = os.WriteFile(fHostRegistered, []byte{'1'}, 0644)
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

package main

import (
	"sync"

	metrics "github.com/toxyl/metric-nexus"
	stats "github.com/toxyl/spider/stats"
)

type Stats struct {
	lock    *sync.Mutex
	spiders map[int]*stats.SpiderStats
	client  *metrics.Client
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

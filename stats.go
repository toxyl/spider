package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/toxyl/glog"
)

type Stats struct {
	lock     *sync.Mutex
	LastKill map[int]time.Time
	Kills    map[int]int
	Prey     map[int]int
}

func (s *Stats) AddSpider(spider int) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.Kills[spider] = 0
	s.Prey[spider] = 0
}

func (s *Stats) AddKill(spider int) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.Kills[spider]; !ok {
		s.Kills[spider] = 0
	}
	s.Kills[spider]++
	s.LastKill[spider] = time.Now()

	if _, ok := s.Prey[spider]; !ok {
		s.Prey[spider] = 1
	}
	s.Prey[spider]--
}

func (s *Stats) AddPrey(spider int) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.Prey[spider]; !ok {
		s.Prey[spider] = 0
	}
	s.Prey[spider]++
}

func (s *Stats) IsStarving(spider int) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	if v, ok := s.Prey[spider]; ok {
		if v > 0 {
			return false // is busy with prey, should be ok
		}
	}

	if _, ok := s.LastKill[spider]; !ok {
		return true // has no kills and no prey, must be hungry
	}

	if s.LastKill[spider].Add(1 * time.Hour).Before(time.Now()) {
		return true // the spider hasn't eaten for an hour, so yeah, it's hungy
	}

	return false // the spider ate within the last hour, it's alright
}

func (s *Stats) Print() {
	s.lock.Lock()
	defer s.lock.Unlock()

	busy := []string{}
	starving := []string{}
	waiting := []string{}
	for spider, kills := range s.Kills {
		s.lock.Unlock()
		isStarving := s.IsStarving(spider)
		s.lock.Lock()
		if isStarving {
			starving = append(starving, glog.Highlight(fmt.Sprint(spider))+" ("+glog.IntAmount(kills, "kill", "kills")+")")
		} else if prey, ok := s.Prey[spider]; ok {
			busy = append(busy, glog.Highlight(fmt.Sprint(spider))+" ("+glog.IntAmount(prey, "prey", "prey")+")")
		} else if kills > 0 {
			waiting = append(waiting, glog.Highlight(fmt.Sprint(spider))+" ("+glog.IntAmount(kills, "kill", "kills")+")")
		}
	}
	if len(waiting) > 0 {
		log.OK("Spiders waiting: %s", glog.Auto(waiting))
	}
	if len(busy) > 0 {
		log.Success("Spiders busy: %s", glog.Auto(busy))
	}
	if len(starving) > 0 {
		log.NotOK("Spiders starving: %s", glog.Auto(starving))
	}
}

func NewStats() *Stats {
	s := &Stats{
		lock:     &sync.Mutex{},
		LastKill: map[int]time.Time{},
		Kills:    map[int]int{},
		Prey:     map[int]int{},
	}
	return s
}

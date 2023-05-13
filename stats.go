package main

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/toxyl/glog"
)

type Stats struct {
	lock                  *sync.Mutex
	LastKill              map[int]time.Time
	Kills                 map[int]int
	Prey                  map[int]int
	changedSinceLastPrint bool
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
	s.changedSinceLastPrint = true

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
	s.changedSinceLastPrint = true

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
	if !s.changedSinceLastPrint {
		return // no need to print again
	}

	spiders := []int{}
	for spider := range s.Kills {
		spiders = append(spiders, spider)
	}
	sort.Ints(spiders)
	var (
		totalKills    = 0
		totalPrey     = 0
		totalBusy     = 0
		totalWaiting  = 0
		totalStarving = 0
		colSpider     = glog.NewTableColumnRight("Spider")
		colBusy       = glog.NewTableColumnCenter("Busy")
		colWaiting    = glog.NewTableColumnCenter("Waiting")
		colStarving   = glog.NewTableColumnCenter("Starving")
		colPrey       = glog.NewTableColumnRight("Prey 🪰")
		colKills      = glog.NewTableColumnRight("Kills 💀")
	)

	for _, spider := range spiders {
		s.lock.Unlock()
		isStarving := s.IsStarving(spider)
		s.lock.Lock()
		kills := s.Kills[spider]
		prey := s.Prey[spider]
		isBusy := prey > 0
		isWaiting := prey == 0 && !isStarving

		colSpider.Push(glog.Highlight(fmt.Sprint(spider)))
		colBusy.Push(isBusy)
		colWaiting.Push(isWaiting)
		colStarving.Push(isStarving)
		colPrey.Push(prey)
		colKills.Push(kills)

		totalPrey += prey
		totalKills += kills
		if isBusy {
			totalBusy++
		}
		if isWaiting {
			totalWaiting++
		}
		if isStarving {
			totalStarving++
		}
	}

	colSpider.Push(glog.Highlight("Total"))
	colBusy.Push(totalBusy)
	colWaiting.Push(totalWaiting)
	colStarving.Push(totalStarving)
	colPrey.Push(totalPrey)
	colKills.Push(totalKills)
	glog.LoggerConfig.ShowIndicator = false
	log.Table(colSpider, colBusy, colWaiting, colStarving, colPrey, colKills)
	glog.LoggerConfig.ShowIndicator = true
	s.changedSinceLastPrint = false
}

func NewStats() *Stats {
	s := &Stats{
		lock:                  &sync.Mutex{},
		LastKill:              map[int]time.Time{},
		Kills:                 map[int]int{},
		Prey:                  map[int]int{},
		changedSinceLastPrint: false,
	}
	return s
}

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
	TimeWasted            map[int]float64
	Kills                 map[int]int
	Prey                  map[int]int
	changedSinceLastPrint bool
}

func (s *Stats) AddSpider(spider int) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.Kills[spider] = 0
	s.Prey[spider] = 0
	s.LastKill[spider] = time.Now()
	s.TimeWasted[spider] = 0.0
}

func (s *Stats) AddKill(spider int, timeWasted float64) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.changedSinceLastPrint = true

	if _, ok := s.Kills[spider]; !ok {
		s.Kills[spider] = 0
	}
	s.TimeWasted[spider] += timeWasted
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
		totalKills     = 0
		totalPrey      = 0
		totalAttacking = 0
		totalWaiting   = 0
		totalStarving  = 0
		timeWasted     = 0.0
		colSpider      = glog.NewTableColumnCenterCustom("🤖", ' ', fmt.Sprint)
		colAttacking   = glog.NewTableColumnCenterCustom("🪖", ' ', fmt.Sprint)
		colWaiting     = glog.NewTableColumnCenterCustom("🚬", ' ', fmt.Sprint)
		colStarving    = glog.NewTableColumnCenterCustom("🍴", ' ', fmt.Sprint)
		colPrey        = glog.NewTableColumnCenterCustom("🪰", ' ', fmt.Sprint)
		colKills       = glog.NewTableColumnCenterCustom("💀", ' ', fmt.Sprint)
		colTime        = glog.NewTableColumnCenterCustom("⌛", ' ', fmt.Sprint)
		colTimeAvg     = glog.NewTableColumnCenterCustom("⌛/💀", ' ', fmt.Sprint)
	)
	emptyCell := glog.WrapDarkGray("*")
	for _, spider := range spiders {
		s.lock.Unlock()
		isStarving := s.IsStarving(spider)
		s.lock.Lock()
		kills := s.Kills[spider]
		prey := s.Prey[spider]
		wasted := s.TimeWasted[spider]
		killAvg := wasted / float64(kills)
		isAttacking := prey > 0
		isWaiting := prey == 0 && !isStarving

		colSpider.Push(glog.Highlight(fmt.Sprint(spider)))
		if isAttacking {
			totalAttacking++
			colAttacking.Push("🪖")
		} else {
			colAttacking.Push(emptyCell)
		}
		if isWaiting {
			totalWaiting++
			colWaiting.Push("🚬")
		} else {
			colWaiting.Push(emptyCell)
		}
		if isStarving {
			totalStarving++
			colStarving.Push("🍴")
		} else {
			colStarving.Push(emptyCell)
		}
		if prey > 0 {
			totalPrey += prey
			colPrey.Push(glog.Int(prey))
		} else {
			colPrey.Push(emptyCell)
		}
		if kills > 0 {
			totalKills += kills
			colKills.Push(glog.Int(kills))
		} else {
			colKills.Push(emptyCell)
		}
		if wasted > 0 {
			timeWasted += wasted
			colTime.Push(glog.DurationShort(wasted, glog.DURATION_SCALE_AVERAGE))
			colTimeAvg.Push(glog.DurationShort(killAvg, glog.DURATION_SCALE_AVERAGE))
		} else {
			colTime.Push(emptyCell)
			colTimeAvg.Push(emptyCell)
		}
	}

	glog.LoggerConfig.ShowDateTime = false
	glog.LoggerConfig.ShowIndicator = false
	log.Default(" ")
	log.Default(glog.Bold()+glog.Underline()+"Status Update"+glog.Reset()+" %s", glog.Time(time.Now()))
	log.Table(colSpider, colAttacking, colWaiting, colStarving, colPrey, colKills, colTime, colTimeAvg)
	log.Default(
		"🪖 %s %s 🚬 %s %s 🍴 %s %s",
		glog.PadLeft(glog.Int(totalAttacking), 4, ' '),
		glog.PadRight(glog.Auto("attacking"), 10, ' '),
		glog.PadLeft(glog.Int(totalWaiting), 4, ' '),
		glog.PadRight(glog.Auto("waiting"), 10, ' '),
		glog.PadLeft(glog.Int(totalStarving), 4, ' '),
		glog.PadRight(glog.Auto("starving"), 10, ' '),
	)
	log.Default(
		"🤖 %s %s 🪰 %s %s 💀 %s %s",
		glog.PadLeft(glog.Int(len(spiders)), 4, ' '),
		glog.PadRight(glog.Auto("spiders"), 10, ' '),
		glog.PadLeft(glog.Int(totalPrey), 4, ' '),
		glog.PadRight(glog.Auto("prey"), 10, ' '),
		glog.PadLeft(glog.Int(totalKills), 4, ' '),
		glog.PadRight(glog.Auto("kills"), 10, ' '),
	)
	log.Default(
		"              ⌛ %s %s",
		glog.PadLeft(glog.DurationShort(timeWasted, glog.DURATION_SCALE_AVERAGE), 9, ' '),
		glog.PadRight(glog.Auto("wasted"), 10, ' '),
	)
	log.Default(" ")
	glog.LoggerConfig.ShowIndicator = true
	glog.LoggerConfig.ShowDateTime = true
	s.changedSinceLastPrint = false
}

func NewStats() *Stats {
	s := &Stats{
		lock:                  &sync.Mutex{},
		LastKill:              map[int]time.Time{},
		TimeWasted:            map[int]float64{},
		Kills:                 map[int]int{},
		Prey:                  map[int]int{},
		changedSinceLastPrint: false,
	}
	return s
}

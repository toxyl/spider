package stats

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/toxyl/glog"
	metrics "github.com/toxyl/metric-nexus"
	"github.com/toxyl/spider/log"
	"github.com/toxyl/spider/utils"
)

func PrintStatsTable(prey, kills map[int]int, wasted, active map[int]float64) {
	spiderIDs := []int{}
	for spider := range prey {
		spiderIDs = append(spiderIDs, spider)
	}
	sort.Ints(spiderIDs)
	var (
		totalKills   = 0
		totalPrey    = 0
		totalActive  = 0.0
		totalWasted  = 0.0
		colSpider    = glog.NewTableColumnRightCustom("Spider", ' ', fmt.Sprint)
		colPrey      = glog.NewTableColumnRightCustom("Prey", ' ', fmt.Sprint)
		colKills     = glog.NewTableColumnRightCustom("Kills", ' ', fmt.Sprint)
		colActive    = glog.NewTableColumnCenterCustom("Active", ' ', fmt.Sprint)
		colWasted    = glog.NewTableColumnCenterCustom("Wasted", ' ', fmt.Sprint)
		colWastedAvg = glog.NewTableColumnCenterCustom("Avg. Wasted", ' ', fmt.Sprint)
	)
	emptyCell := glog.WrapDarkGray("*")
	for _, spider := range spiderIDs {
		kills := kills[spider]
		prey := prey[spider]
		active := active[spider]
		wasted := wasted[spider]
		killAvg := wasted / float64(kills)

		colSpider.Push(glog.Highlight(fmt.Sprint(spider)))
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
			totalWasted += wasted
			colWasted.Push(glog.DurationShort(wasted, glog.DURATION_SCALE_AVERAGE))
			colWastedAvg.Push(glog.DurationShort(killAvg, glog.DURATION_SCALE_AVERAGE))
		} else {
			colWasted.Push(emptyCell)
			colWastedAvg.Push(emptyCell)
		}
		if active > 0 {
			totalActive += active
			colActive.Push(glog.DurationShort(active, glog.DURATION_SCALE_AVERAGE))
		} else {
			colActive.Push(emptyCell)
		}
	}
	colSpider.Push(glog.Auto("All"))
	colPrey.Push(glog.Int(totalPrey))
	colKills.Push(glog.Int(totalKills))
	colActive.Push(glog.DurationShort(totalActive, glog.DURATION_SCALE_AVERAGE))
	colWasted.Push(glog.DurationShort(totalWasted, glog.DURATION_SCALE_AVERAGE))
	colWastedAvg.Push(glog.DurationShort(totalWasted/float64(totalKills), glog.DURATION_SCALE_AVERAGE))
	log.Normal(" ")
	log.Table(colSpider, colPrey, colKills, colActive, colWasted, colWastedAvg)
	log.Normal(" ")
}

type SpiderStats struct {
	client         *metrics.Client
	lock           *sync.Mutex
	prey           int
	kills          int
	wasted         float64
	metricPrey     string
	metricKills    string
	metricWasted   string
	metricActive   string
	metricFilePrey string
}

func (ss *SpiderStats) AddPrey() {
	ss.lock.Lock()
	defer ss.lock.Unlock()
	ss.prey++
}

func (ss *SpiderStats) RemovePrey(timeWasted float64) {
	ss.lock.Lock()
	defer ss.lock.Unlock()
	ss.prey--
	ss.kills++
	ss.wasted += timeWasted
}

func (ss *SpiderStats) Update(timeWaited time.Duration) {
	ss.lock.Lock()
	defer ss.lock.Unlock()

	k := ss.kills
	ss.kills = 0
	p := ss.prey
	ss.prey = 0
	w := ss.wasted
	ss.wasted = 0

	ss.client.Add(ss.metricPrey, p)
	ss.client.Add(ss.metricKills, k)
	ss.client.Add(ss.metricWasted, w)
	ss.client.Add(ss.metricActive, timeWaited.Seconds())
	utils.FileWrite(ss.metricFilePrey, fmt.Sprint(p))
}

func NewSpiderStats(spider int, nexusHost string, nexusPort int, nexusKey string) *SpiderStats {
	ss := &SpiderStats{
		client:         metrics.NewClient(nexusHost, nexusPort, nexusKey, true),
		lock:           &sync.Mutex{},
		prey:           0,
		kills:          0,
		wasted:         0,
		metricPrey:     utils.GetMetricName(spider, "prey"),
		metricKills:    utils.GetMetricName(spider, "kills"),
		metricWasted:   utils.GetMetricName(spider, "wasted"),
		metricActive:   utils.GetMetricName(spider, "active"),
		metricFilePrey: utils.GetMetricFileName(spider, "prey"),
	}

	ss.client.Create(ss.metricKills, "Kills made by the spider.")
	ss.client.Create(ss.metricPrey, "Prey being attacked by the spider.")
	ss.client.Create(ss.metricWasted, "Time wasted by the spider.")
	ss.client.Create(ss.metricActive, "How long the spider has been active.")

	if v, err := utils.FileRead(ss.metricFilePrey); err == nil {
		// we might have crashed, let's remove the prey since it has disappeared with us
		ss.client.Subtract(ss.metricPrey, v)
		if err := utils.FileDelete(ss.metricFilePrey); err != nil {
			panic("could not remove " + ss.metricFilePrey + " because " + err.Error())
		}
	}

	go func() {
		for {
			t := time.Now()
			time.Sleep(10 * time.Second)
			ss.Update(time.Since(t))
		}
	}()

	return ss
}

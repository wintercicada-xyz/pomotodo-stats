package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
	heatmap "github.com/blurfx/calendar-heatmap"
)

var colors = []string{
	"#ffffff",
	"#fef7f6",
	"#fdefed",
	"#fce7e4",
	"#fbdfdb",
	"#fad7d2",
	"#f9cfc9",
	"#f8c7c0",
	"#f7bfb7",
	"#f6b7ae",
	"#f5afa5",
	"#f4a79c",
	"#f39f93",
	"#f2978a",
	"#f18f81",
	"#f08778",
	"#ef7f6f",
	"#ee7766",
	"#ed6f5d",
	"#ec6754",
}

var myConfig = heatmap.CalendarHeatmapConfig{
	Colors:           colors,
	BlockSize:        11,
	BlockMargin:      2,
	BlockRoundness:   2,
	MonthLabels:      []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"},
	MonthLabelHeight: 15,
	WeekdayLabels:    []string{"", "Mon", "", "Wed", "", "Fri", ""},
}

var myHeatmap = heatmap.New(&myConfig)

func main() {
	addr := "0.0.0.0:8080"
	if len(os.Args) > 1 {
		addr = os.Args[1]
	}
	fmt.Println(http.ListenAndServe(addr, &getStatsHandler{
		pomotodoStatsMap: make(map[string]pomotodoStats),
		refreshPeriod:    time.Hour,
	}))
}

type pomotodoStats struct {
	SVG          []byte
	latestUpdate time.Time
}

func (stats *pomotodoStats) refresh(icsURL string) error {
	calendar, err := getCalendar(icsURL)
	if err != nil {
		fmt.Println(err)
		return err
	}
	stats.latestUpdate = time.Now()
	pomotodoCount, topTags, err := analyzeCalendar(calendar)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(pomotodoCount)
	fmt.Println(topTags)
	end := stats.latestUpdate
	start := end.AddDate(0, -1, 0)

	svg := myHeatmap.Generate(
		heatmap.Date{Year: start.Year(), Month: start.Month(), Day: start.Day()},
		heatmap.Date{Year: end.Year(), Month: end.Month(), Day: end.Day()},
		pomotodoCount,
	)
	stats.SVG = svg.Bytes()
	return nil

}

type getStatsHandler struct {
	pomotodoStatsMap map[string]pomotodoStats
	refreshPeriod    time.Duration
}

var _ http.Handler = (*getStatsHandler)(nil)

func (handler *getStatsHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	icsURL := request.URL.Query().Get("ics")
	URL, err := url.Parse(icsURL)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}
	if URL.Hostname() != "ical.pomotodo.com" {
		responseWriter.WriteHeader(http.StatusForbidden)
		return
	}

	stats, ok := handler.pomotodoStatsMap[icsURL]
	if ok && time.Since(stats.latestUpdate) > handler.refreshPeriod {
		err := stats.refresh(icsURL)
		if err != nil {
			fmt.Println(err)
			responseWriter.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	if !ok {
		newStats := *new(pomotodoStats)
		err := newStats.refresh(icsURL)
		if err != nil {
			fmt.Println(err)
			responseWriter.WriteHeader(http.StatusInternalServerError)
			return
		}
		handler.pomotodoStatsMap[icsURL] = newStats

	}
	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Header().Set("Content-Type", "image/svg+xml")
	responseWriter.Write(handler.pomotodoStatsMap[icsURL].SVG)
}

func getCalendar(url string) (calendar *ics.Calendar, err error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	contentType := response.Header.Get("content-type")

	if !strings.Contains(contentType, "text/calendar") {
		return nil, fmt.Errorf("Expected ics file, get %s", contentType)
	}
	calendar, err = ics.ParseCalendar(response.Body)
	if err != nil {
		return nil, err
	}
	response.Body.Close()
	return calendar, nil
}

func analyzeCalendar(calendar *ics.Calendar) (pomotodoCount map[heatmap.Date]int, topTags []string, err error) {
	pomotodoCount = make(map[heatmap.Date]int)
	tagsCount := make(map[string]int)
	for _, event := range calendar.Events() {
		summary := event.GetProperty("SUMMARY").Value
		strs := strings.Split(summary, " ")
		for _, str := range strs {
			if strings.HasPrefix(str, "#") {
				tagsCount[str]++
			}
		}
		t, err := event.GetStartAt()
		if err != nil {
			return nil, nil, err
		}
		pomotodoCount[heatmap.Date{Year: t.Year(), Month: t.Month(), Day: t.Day()}]++
	}
	//sort tagsCount by its value but use the key as the value
	var keys []string
	for k := range tagsCount {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return tagsCount[keys[i]] > tagsCount[keys[j]]
	})
	topTags = keys
	return pomotodoCount, topTags, nil
}

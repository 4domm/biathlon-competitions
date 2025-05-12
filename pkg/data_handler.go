package pkg

import (
	"fmt"
	"sort"
	"strconv"
	"time"
)

type LapInfo struct {
	Time  time.Duration
	Speed float64
}

type Result struct {
	CompetitorID string
	Status       string
	StartTime    time.Time
	FinishTime   time.Time
	Laps         []LapInfo
	PenaltyTime  time.Duration
	PenaltySpeed float64
	PenaltyStart time.Time
	Hits         int
	Shots        int
	TotalTime    time.Duration
}

type DataHandler struct {
	config         *Config
	resultWriter   *OutputWrapper
	competitorData map[string]*Result
}

func NewDataHandler(cfg *Config, writer *OutputWrapper) *DataHandler {
	return &DataHandler{config: cfg, resultWriter: writer, competitorData: make(map[string]*Result)}
}

func (dh *DataHandler) ProcessEvents(events []Event) {
	for _, event := range events {
		timestamp := event.Time.Format("[15:04:05.000]")
		cid := event.CompetitorID
		if _, ok := dh.competitorData[cid]; !ok {
			dh.competitorData[cid] = &Result{
				CompetitorID: cid,
				Laps:         make([]LapInfo, dh.config.Laps),
			}
		}
		res := dh.competitorData[cid]

		switch event.ID {
		case 1:
			fmt.Printf("%s The competitor(%s) registered\n", timestamp, event.CompetitorID)

		case 2:
			if len(event.ExtraParams) > 0 {
				scheduled, _ := time.Parse("15:04:05.000", event.ExtraParams)
				res.StartTime = scheduled
				fmt.Printf("%s The start time for the competitor(%s) was set by a draw to %s\n", timestamp, event.CompetitorID, event.ExtraParams)
			}

		case 3:
			fmt.Printf("%s The competitor(%s) is on the start line\n", timestamp, event.CompetitorID)

		case 4:
			fmt.Printf("%s The competitor(%s) has started\n", timestamp, event.CompetitorID)

		case 5:
			if len(event.ExtraParams) > 0 {
				fmt.Printf("%s The competitor(%s) is on the firing range(%s)\n", timestamp, event.CompetitorID, event.ExtraParams)
			}

		case 6:
			if len(event.ExtraParams) > 0 {
				n, err := strconv.Atoi(event.ExtraParams)
				if err == nil {
					res.Shots = max(res.Shots, n)
					res.Hits++
					fmt.Printf("%s The target(%s) has been hit by competitor(%s)\n", timestamp, event.ExtraParams, event.CompetitorID)
				}
			}

		case 7:
			fmt.Printf("%s The competitor(%s) left the firing range\n", timestamp, event.CompetitorID)

		case 8:
			res.PenaltyStart = event.Time
			fmt.Printf("%s The competitor(%s) entered the penalty laps\n", timestamp, event.CompetitorID)

		case 9:
			if !res.PenaltyStart.IsZero() {
				penaltyDuration := event.Time.Sub(res.PenaltyStart)
				res.PenaltyTime += penaltyDuration
				if penaltyDuration.Seconds() > 0 {
					res.PenaltySpeed = float64(dh.config.PenaltyLen) / penaltyDuration.Seconds()
				}
			}
			fmt.Printf("%s The competitor(%s) left the penalty laps\n", timestamp, event.CompetitorID)

		case 10:
			lapIndex := -1
			for i := range res.Laps {
				if res.Laps[i].Time == 0 {
					lapIndex = i
					break
				}
			}
			if lapIndex >= 0 {
				start := res.StartTime
				for i := 0; i < lapIndex; i++ {
					start = start.Add(res.Laps[i].Time)
				}
				dur := event.Time.Sub(start)
				speed := float64(dh.config.LapLen) / dur.Seconds()
				res.Laps[lapIndex] = LapInfo{Time: dur, Speed: speed}
				res.FinishTime = event.Time
			}
			fmt.Printf("%s The competitor(%s) ended the main lap\n", timestamp, event.CompetitorID)

		case 11:
			if len(event.ExtraParams) > 0 {
				res.Status = "NotFinished"
				fmt.Printf("%s The competitor(%s) can`t continue: %s\n", timestamp, event.CompetitorID, event.ExtraParams)
			}

		default:
			fmt.Printf("%s Unknown event ID(%d) for competitor(%s)\n", timestamp, event.ID, event.CompetitorID)
		}
	}
}

func (dh *DataHandler) ComputeReport() {
	var results []*Result
	for _, res := range dh.competitorData {
		if res.Status == "" && !res.FinishTime.IsZero() {
			res.TotalTime = res.FinishTime.Sub(res.StartTime)
		} else if res.Status == "" && res.StartTime.IsZero() {
			res.Status = "NotStarted"
		}
		results = append(results, res)
	}
	sort.Slice(results, func(i, j int) bool {
		if results[i].Status != "" && results[j].Status == "" {
			return false
		}
		if results[i].Status == "" && results[j].Status != "" {
			return true
		}
		return results[i].TotalTime < results[j].TotalTime
	})
	for _, res := range results {
		dh.resultWriter.WriteData(res.FormatReport())
	}
}

func (res *Result) FormatReport() string {
	status := fmt.Sprintf("[%s]", res.Status)
	if res.Status == "" {
		status = fmt.Sprintf("[%s]", formatDuration(res.TotalTime))
	}
	report := fmt.Sprintf("%s %s [", status, res.CompetitorID)
	for i, lap := range res.Laps {
		if i > 0 {
			report += ", "
		}
		if lap.Time == 0 {
			report += "{,}"
		} else {
			report += fmt.Sprintf("{%s, %.3f}", formatDuration(lap.Time), lap.Speed)
		}
	}
	report += fmt.Sprintf("] {%s, %.3f} %d/%d", formatDuration(res.PenaltyTime), res.PenaltySpeed, res.Hits, res.Shots)
	return report
}

func formatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	ms := int(d.Milliseconds()) % 1000
	return fmt.Sprintf("%02d:%02d:%02d.%03d", h, m, s, ms)
}

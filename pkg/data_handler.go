package pkg

import (
	"fmt"
)

var eventMessages = map[int]string{
	1:  "The competitor(%s) registered",
	2:  "The start time for the competitor(%s) was set by a draw to %s",
	3:  "The competitor(%s) is on the start line",
	4:  "The competitor(%s) has started",
	5:  "The competitor(%s) is on the firing range(%s)",
	6:  "The target(%s) has been hit by competitor(%s)",
	7:  "The competitor(%s) left the firing range",
	8:  "The competitor(%s) entered the penalty laps",
	9:  "The competitor(%s) left the penalty laps",
	10: "The competitor(%s) ended the main lap",
	11: "The competitor(%s) can`t continue: %s",
}

type DataHandler struct {
	config       *Config
	resultWriter *OutputWrapper
	

}

func NewDataHandler(cfg *Config, writer *OutputWrapper) *DataHandler {
	return &DataHandler{config: cfg, resultWriter: writer}
}

func (dh *DataHandler) ProcessEvents(events []Event) {
	for _, event := range events {
		competitorID := event.CompetitorID
		timestamp := event.Time.Format("[15:04:05.000]")

		template, ok := eventMessages[event.ID]
		if !ok {
			fmt.Printf("%s Unknown event ID(%d) for competitor(%s)", timestamp, event.ID, competitorID)
			continue
		}

		msg := ""
		switch event.ID {
		case 2, 5, 6, 11:
			if len(event.ExtraParams) > 0 {
				if event.ID == 6 {
					msg = fmt.Sprintf(template, event.ExtraParams, competitorID)
				} else {
					msg = fmt.Sprintf(template, competitorID, event.ExtraParams)
				}
			}
		default:
			msg = fmt.Sprintf(template, competitorID)
		}

		fmt.Printf("%s %s \n", timestamp, msg)
	}
}

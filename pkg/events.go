package pkg

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Event struct {
	Time         time.Time
	ID           int
	CompetitorID string
	ExtraParams  string
}

func LoadEvents(path string) ([]Event, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var events []Event
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		event, err := parseEvent(scanner.Text())
		if err != nil {
			fmt.Println("cant parse line")
			continue
		}
		events = append(events, event)
	}
	return events, scanner.Err()
}

func parseEvent(input string) (Event, error) {
	parts := strings.Fields(input)
	timeStr := parts[0][1 : len(parts[0])-1]
	t, err := time.Parse("15:04:05.000", timeStr)
	if err != nil {
		return Event{}, err
	}
	eventID, _ := strconv.Atoi(parts[1])
	return Event{
		Time:         t,
		ID:           eventID,
		CompetitorID: parts[2],
		ExtraParams:  strings.Join(parts[3:], " "),
	}, nil
}

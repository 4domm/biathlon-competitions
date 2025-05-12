package main

import (
	"log"
	"os"

	"github.com/4domm/biathlon-competitions/pkg"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalln("Usage: ./bin config.json events")
		return
	}
	configPath := os.Args[1]
	eventsPath := os.Args[2]
	config, err := pkg.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Error with loadconfig: %v", err)
		return
	}
	events, err := pkg.LoadEvents(eventsPath)
	if err != nil {
		log.Fatal(err)
		return
	}
	dataHandler := pkg.NewDataHandler(config, pkg.NewOutputWrapper(os.Stdout))

	dataHandler.ProcessEvents(events)
	dataHandler.ComputeReport()
}

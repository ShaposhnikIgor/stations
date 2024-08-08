package main

import (
	"fmt"
	"os"
	train "stations/pkg"
	"strconv"
)

func main() {
	if len(os.Args) < 5 {
		fmt.Fprintln(os.Stderr, "Error: Too few command line arguments")
		os.Exit(1)
	}

	filePath := os.Args[1]
	startStation := os.Args[2]
	endStation := os.Args[3]
	numTrains, err := strconv.Atoi(os.Args[4])
	if err != nil || numTrains <= 0 {
		fmt.Fprintln(os.Stderr, "Error: Invalid number of trains")
		os.Exit(1)
	}

	if len(os.Args) > 5 {
		train.ValidateExtraArgs(os.Args[5:])
	}

	stations, connections, err := train.ParseNetworkMap(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	graph := train.NewGraph(connections, stations)
	train.MoveTrains(graph, startStation, endStation, numTrains)
}

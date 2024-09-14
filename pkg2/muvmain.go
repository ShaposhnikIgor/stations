package train2

import (
	"fmt"
	"os"
	"strconv"
)

func Muvmain() {
	if len(os.Args) < 5 {
		fmt.Fprintln(os.Stderr, "Error: Too few command line arguments")
		os.Exit(1)
	}

	CheckArguments(os.Args)

	filePath := os.Args[1]
	startStation := os.Args[2]
	endStation := os.Args[3]
	numTrains, err := strconv.Atoi(os.Args[4])
	if err != nil || numTrains <= 0 {
		fmt.Fprintln(os.Stderr, "Error: Invalid number of trains")
		os.Exit(1)
	}

	if len(os.Args) > 5 {
		ValidateExtraArgs(os.Args[5:])
	}

	// Parse the network map first
	stations, connections, err := ParseNetworkMap(filePath)
	if err != nil {
		Error(err.Error())
		os.Exit(0)
	}

	// Validate the existence of the start and end stations before doing anything else
	ValidateStationExistence(stations, startStation, endStation)

	// Check for invalid connections before creating the graph
	err = CheckInvalidConnections(stations, connections)
	if err != nil {
		Error(err.Error())
		os.Exit(1)
	}

	// Continue with the rest of the program only if the stations exist
	graph := NewGraph(connections, stations)
	MoveTrains(graph, startStation, endStation, numTrains)

	// Check that the start and end stations are different
	ValidateDifferentStations(startStation, endStation)

	// Perform the search for the path between the stations
	path, err := HybridSearch(graph, startStation, endStation, numTrains)
	if err != nil {
		Error(err.Error())
	}

	// Validate the existence of the path found
	ValidatePathExistence(path, startStation, endStation)

	// Validate the number of trains
	ValidateTrainCount(numTrains, err)
}

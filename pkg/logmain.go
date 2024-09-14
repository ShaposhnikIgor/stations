package train

import (
	"fmt"

	"os"

	"strconv"
)

func Logmain() {
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

	ValidateDifferentStations(startStation, endStation)

	// Validate the existence of the start and end stations before doing anything else
	ValidateStationExistence(stations, startStation, endStation)
	err2 := CheckConnectionsExist(stations, connections)
	if err2 != nil {
		Error(err2.Error())
		os.Exit(0)
	}

	stationConnections := BuildConnectionMap(stations, connections)
	//fmt.Println(stationConnections)

	allRoutes, err := FindAllPossibleRoutes(stationConnections, startStation, endStation)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}

	combinationRoutes := FindAllRouteCombinations(allRoutes)

	bestRoute, bestRouteInfo := FindOptimalRoute(numTrains, combinationRoutes)

	DisplayTrainMovements(bestRoute, bestRouteInfo, numTrains)
}

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

	train.CheckArguments(os.Args)

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

	// Parse the network map first
	stations, connections, err := train.ParseNetworkMap(filePath)
	if err != nil {
		train.Error(err.Error())
		os.Exit(0)
	}

	train.ValidateDifferentStations(startStation, endStation)

	// Validate the existence of the start and end stations before doing anything else
	train.ValidateStationExistence(stations, startStation, endStation)
	err2 := train.CheckConnectionsExist(stations, connections)
	if err2 != nil {
		train.Error(err2.Error())
		os.Exit(0)
	}

	stationConnections := train.BuildConnectionMap(stations, connections)
	//fmt.Println(stationConnections)

	allRoutes, err := train.FindAllPossibleRoutes(stationConnections, startStation, endStation)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}

	combinationRoutes := train.FindAllRouteCombinations(allRoutes)

	bestRoute, bestRouteInfo := train.FindOptimalRoute(numTrains, combinationRoutes)

	train.DisplayTrainMovements(bestRoute, bestRouteInfo, numTrains)
}

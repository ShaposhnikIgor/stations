package train

import (
	"fmt"
	"math"
	"os"
	"slices"
	"sort"
)

// BuildConnectionMap creates a map where each station is mapped to a slice of stations it is connected to.
func BuildConnectionMap(stations []Station, connections [][]string) map[string][]string {
	stationConnections := make(map[string][]string)
	var connected []string

	for _, stn := range stations {
		for _, connection := range connections {
			if connection[0] == stn.Name {
				connected = append(connected, connection[1])
			} else if connection[1] == stn.Name {
				connected = append(connected, connection[0])
			}
		}
		if connected != nil {
			stationConnections[stn.Name] = connected
			connected = nil
		}
	}

	return stationConnections
}

// FindAllRoutes finds all possible routes from the start station to the end station.
func FindAllPossibleRoutes(connections map[string][]string, startStation, endStation string) ([][]string, error) {
	var allRoutes [][]string

	var findPaths func(current, destination string, path []string)
	findPaths = func(current, destination string, path []string) {
		path = append(path, current)

		if current == destination {
			route := make([]string, len(path))
			copy(route, path)
			allRoutes = append(allRoutes, route[1:]) // Exclude the start station
			return
		}

		for _, neighbor := range connections[current] {
			if !slices.Contains(path, neighbor) {
				findPaths(neighbor, destination, path)
			}
		}
	}

	findPaths(startStation, endStation, []string{})

	ValidatePathExistence(allRoutes, startStation, endStation)

	sort.Slice(allRoutes, func(i, j int) bool {
		return len(allRoutes[i]) < len(allRoutes[j])
	})

	return allRoutes, nil
}

// FindAllRouteCombinations generates all possible combinations of non-redundant routes.
func FindAllRouteCombinations(allRoutes [][]string) [][][]string {
	var routeCombinations [][][]string

	for startIndex := 0; startIndex < len(allRoutes); startIndex++ {
		currentCombination := [][]string{allRoutes[startIndex]}
		generateCombinations(allRoutes, currentCombination, len(allRoutes), startIndex+1, &routeCombinations)
	}

	return routeCombinations
}

// generateCombinations recursively generates combinations of routes and checks for redundancy.
func generateCombinations(allRoutes [][]string, currentCombination [][]string, totalRoutes int, currentIndex int, routeCombinations *[][][]string) {
	if currentIndex == totalRoutes {
		if isUniqueCombination(currentCombination, routeCombinations) {
			*routeCombinations = append(*routeCombinations, currentCombination)
		}
	} else {
		if canAddRoute(allRoutes[currentIndex], currentCombination) {
			newCombination := append(currentCombination, allRoutes[currentIndex])
			generateCombinations(allRoutes, newCombination, totalRoutes, currentIndex+1, routeCombinations)
		}
		generateCombinations(allRoutes, currentCombination, totalRoutes, currentIndex+1, routeCombinations)
	}
}

// isUniqueCombination checks if a given combination of routes is not redundant compared to existing combinations.
func isUniqueCombination(newCombination [][]string, existingCombinations *[][][]string) bool {
	if len(newCombination) == 1 {
		return isUniqueSingleRoute(newCombination[0], existingCombinations)
	}

	newCombinationLengths := getRouteLengths(newCombination)
	for _, existingCombination := range *existingCombinations {
		existingCombinationLengths := getRouteLengths(existingCombination)
		if areCombinationsEquivalent(newCombinationLengths, existingCombinationLengths) {
			return false
		}
	}

	return true
}

// isUniqueSingleRoute checks if a single route is unique within existing combinations.
func isUniqueSingleRoute(route []string, existingCombinations *[][][]string) bool {
	for _, combination := range *existingCombinations {
		for _, existingRoute := range combination {
			if len(route) == len(existingRoute) {
				return false
			}
		}
	}
	return true
}

// getRouteLengths returns a slice of lengths for each route in the combination.
func getRouteLengths(routes [][]string) []int {
	lengths := make([]int, len(routes))
	for i, route := range routes {
		lengths[i] = len(route)
	}
	return lengths
}

// areCombinationsEquivalent checks if two combinations of routes have the same lengths, after sorting.
func areCombinationsEquivalent(combinationA, combinationB []int) bool {
	if len(combinationA) != len(combinationB) {
		return false
	}
	sort.Ints(combinationA)
	sort.Ints(combinationB)

	for i := range combinationA {
		if combinationA[i] != combinationB[i] {
			return false
		}
	}
	return true
}

// canAddRoute checks if a route can be added to a combination without violating uniqueness.
func canAddRoute(newRoute []string, currentCombination [][]string) bool {
	endStation := os.Args[3]

	for _, route := range currentCombination {
		for _, stationA := range route {
			for _, stationB := range newRoute {
				if stationA == endStation || stationB == endStation {
					continue
				} else if stationA == stationB {
					return false
				}
			}
		}
	}
	return true
}

// FindOptimalRoute determines the best route combination based on the train number and route lengths.
func FindOptimalRoute(trainNumber int, routeCombinations [][][]string) (optimalRoute [][]string, optimalRouteInfo []int) {
	routeLengths := calculateRouteLengths(routeCombinations)
	shortestTurns := math.MaxInt

	for index, lengths := range routeLengths {
		turnsRequired := calculateTurnsForTrains(lengths, trainNumber)
		if shortestTurns > turnsRequired {
			optimalRoute = routeCombinations[index]
			optimalRouteInfo = lengths
			shortestTurns = turnsRequired
		}
	}
	return optimalRoute, optimalRouteInfo
}

// calculateRouteLengths processes the route combinations and returns a slice of route lengths.
func calculateRouteLengths(routeCombinations [][][]string) [][]int {
	var routeLengths [][]int

	for _, routes := range routeCombinations {
		var lengths []int
		for _, route := range routes {
			lengths = append(lengths, len(route))
		}
		routeLengths = append(routeLengths, lengths)
	}
	return routeLengths
}

// calculateTurnsForTrains determines the number of turns required to utilize all trains on the routes.
func calculateTurnsForTrains(routeLengths []int, trainNumber int) int {
	trainCount := 0
	turns := 0

	for trainCount < trainNumber {
		trainCount, turns = incrementTurn(trainCount, turns, routeLengths)
	}
	return turns
}

// incrementTurn increases the turn count and the number of trains placed on routes for the current turn.
func incrementTurn(trainCount int, turns int, routeLengths []int) (int, int) {
	placedTrains := 0
	for _, length := range routeLengths {
		if length <= turns {
			placedTrains++
		}
	}

	trainCount += placedTrains
	return trainCount, turns + 1
}

func DisplayTrainMovements(routePlans [][]string, routeDurations []int, numTrains int) {
	trainAllocation := allocateTrains(routeDurations, numTrains)
	simulateAndDisplayMovements(trainAllocation, routeDurations, routePlans, numTrains)
}

func allocateTrains(routeDurations []int, numTrains int) map[int][]int {
	turn := 1
	trainsAllocated := 0
	trainAllocation := make(map[int][]int, len(routeDurations))

	for trainsAllocated < numTrains {
		for routeIdx, duration := range routeDurations {
			if duration <= turn {
				trainsAllocated++
				if trainsAllocated <= numTrains {
					trainAllocation[routeIdx] = append(trainAllocation[routeIdx], trainsAllocated)
				}
			}
		}
		turn++
	}
	return trainAllocation
}

func simulateAndDisplayMovements(trainAllocation map[int][]int, routeDurations []int, routePlans [][]string, numTrains int) {
	stationStatus := initializeStationStatus(routePlans)
	trainsStatusMap := initializeTrainStatusMap(trainAllocation, routeDurations, routePlans, numTrains)
	performTrainMovements(stationStatus, trainsStatusMap, routePlans, numTrains)
}

func initializeStationStatus(routePlans [][]string) map[string]int {
	stationStatus := make(map[string]int)
	for _, stations := range routePlans {
		for _, station := range stations {
			stationStatus[station] = 0 // All stations are initially available
		}
	}
	return stationStatus
}

type trainStatus struct {
	status               string
	pathNumber           int
	currentStationNumber int
}

func initializeTrainStatusMap(trainAllocation map[int][]int, routeDurations []int, routePlans [][]string, numTrains int) map[int]*trainStatus {
	trainsStatusMap := make(map[int]*trainStatus)

	for train := 1; train <= numTrains; train++ {
		pathIdx := findPathForTrain(train, trainAllocation, routeDurations)
		trainsStatusMap[train] = &trainStatus{
			status:               "starting",
			pathNumber:           pathIdx,
			currentStationNumber: 0,
		}
	}
	return trainsStatusMap
}

func findPathForTrain(train int, trainAllocation map[int][]int, routeDurations []int) int {
	for pathIdx := 0; pathIdx < len(routeDurations); pathIdx++ {
		if contains(trainAllocation[pathIdx], train) {
			return pathIdx
		}
	}
	return -1
}

func contains(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func performTrainMovements(stationStatus map[string]int, trainsStatusMap map[int]*trainStatus, routePlans [][]string, numTrains int) {
	var trainLog string
	var oneLengthPathUsed bool
	endStation := os.Args[3]

	for trainsStatusMap[numTrains].status != "finished" {
		for trainIdx := 1; trainIdx <= numTrains; trainIdx++ {
			if trainStatus, exists := trainsStatusMap[trainIdx]; exists {
				trainLog, oneLengthPathUsed = processTrainMovement(trainStatus, stationStatus, routePlans, trainIdx, endStation, trainLog, oneLengthPathUsed)
			}
		}
		fmt.Println(trainLog)
		trainLog = ""
		oneLengthPathUsed = false
	}
}

func processTrainMovement(trainStatus *trainStatus, stationStatus map[string]int, routePlans [][]string, trainIdx int, endStation, trainLog string, oneLengthPathUsed bool) (string, bool) {
	if trainStatus.status == "moving" {
		nextStation := routePlans[trainStatus.pathNumber][trainStatus.currentStationNumber+1]
		if stationStatus[nextStation] == 0 {
			currentStation := routePlans[trainStatus.pathNumber][trainStatus.currentStationNumber]
			nextStation := routePlans[trainStatus.pathNumber][trainStatus.currentStationNumber+1]
			stationStatus[currentStation] = 0
			trainStatus.currentStationNumber++
			trainLog += fmt.Sprintf("T%d-%v ", trainIdx, nextStation)
			if nextStation == endStation {
				trainStatus.status = "finished"
				stationStatus[nextStation] = 0
			} else {
				stationStatus[nextStation] = 1
			}
		}
	} else if trainStatus.status == "starting" {
		startingStation := routePlans[trainStatus.pathNumber][0]
		if stationStatus[startingStation] == 0 {
			if startingStation == endStation {
				if !oneLengthPathUsed {
					trainStatus.status = "finished"
					oneLengthPathUsed = true
					trainLog += fmt.Sprintf("T%d-%v ", trainIdx, startingStation)
				}
				return trainLog, oneLengthPathUsed
			} else {
				stationStatus[startingStation] = 1
				trainStatus.status = "moving"
			}
			trainLog += fmt.Sprintf("T%d-%v ", trainIdx, startingStation)
		}
	}
	return trainLog, oneLengthPathUsed
}

func ValidateExtraArgs(args []string) {
	for _, arg := range args {
		if arg == "extra" || arg == "bonus" {
			fmt.Println("Handling extra argument:", arg)
		} else {
			fmt.Fprintf(os.Stderr, "Error: Invalid extra argument: %s\n", arg)
			os.Exit(1)
		}
	}
}

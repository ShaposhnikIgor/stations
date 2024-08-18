package train

import (
	"errors"
	"fmt"
	"os"
)

// Error function prints the error message and exits the program.
func Error(message string) {
	fmt.Fprintln(os.Stderr, "Error:", message)
	os.Exit(1)
}

// CheckArguments validates the number of command line arguments.
func CheckArguments(args []string) {
	if len(args) < 5 {
		Error("Too few command line arguments")
	} else if len(args) > 6 {
		Error("Too many command line arguments")
	}
}

// ValidateStationExistence checks if the start and end stations exist in the map.
func ValidateStationExistence(stations map[string]Station, startStation, endStation string) {
	if _, ok := stations[startStation]; !ok {
		Error(fmt.Sprintf("Start station does not exist: %s", startStation))
	}
	if _, ok := stations[endStation]; !ok {
		Error(fmt.Sprintf("End station does not exist: %s", endStation))
	}
}

// ValidateDifferentStations checks that start and end stations are not the same.
func ValidateDifferentStations(startStation, endStation string) {
	if startStation == endStation {
		Error("Start and end station are the same")
	}
}

// ValidatePathExistence checks if there is a valid path between start and end stations.
func ValidatePathExistence(path []string, startStation, endStation string) {
	if len(path) == 0 {
		Error(fmt.Sprintf("No path found from %s to %s", startStation, endStation))
	}
}

// ValidateTrainCount checks if the number of trains is a valid positive integer.
func ValidateTrainCount(numTrains int, err error) {
	if err != nil || numTrains <= 0 {
		Error("Invalid number of trains")
	}
}

// CheckDuplicateRoutes validates that no duplicate connections exist, including reversed ones.
func CheckDuplicateRoutes(connections []Connection) error {
	seen := make(map[string]bool)
	for _, conn := range connections {
		forward := fmt.Sprintf("%s-%s", conn.From, conn.To)
		reverse := fmt.Sprintf("%s-%s", conn.To, conn.From)
		if seen[forward] || seen[reverse] {
			return fmt.Errorf("duplicate connection between %s and %s", conn.From, conn.To)
		}
		seen[forward] = true
	}
	return nil
}

// CheckSections validates that both "stations:" and "connections:" sections are present in the map.
func CheckSections(stationsSectionFound, connectionsSectionFound bool) error {
	if !stationsSectionFound {
		return fmt.Errorf("map does not contain a 'stations:' section")
	}
	if !connectionsSectionFound {
		return fmt.Errorf("map does not contain a 'connections:' section")
	}
	return nil
}

func CheckDuplicateStationNames(stations map[string]Station) error {
	seen := make(map[string]bool)
	for name := range stations {
		if seen[name] {
			// Вывод ошибки в стандартный поток ошибок
			fmt.Fprintf(os.Stderr, "Error: duplicate station name found: %s\n", name)
			return errors.New("duplicate station name detected")
		}
		seen[name] = true
	}
	return nil
}

// CheckDuplicateCoordinates checks if any two stations have the same coordinates.
func CheckDuplicateCoordinates(stations map[string]Station) error {
	coordsMap := make(map[string]bool)
	for _, station := range stations {
		coords := fmt.Sprintf("%d,%d", station.X, station.Y)
		if coordsMap[coords] {
			return fmt.Errorf("two stations exist at the same coordinates: %d, %d", station.X, station.Y)
		}
		coordsMap[coords] = true
	}
	return nil
}

// CheckInvalidConnections validates that all connections refer to existing stations.
func CheckInvalidConnections(stations map[string]Station, connections []Connection) error {
	for _, conn := range connections {
		if _, ok := stations[conn.From]; !ok {
			return fmt.Errorf("connection from non-existent station: %s", conn.From)
		}
		if _, ok := stations[conn.To]; !ok {
			return fmt.Errorf("connection to non-existent station: %s", conn.To)
		}
	}
	return nil
}

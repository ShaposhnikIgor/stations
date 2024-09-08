package train

import (
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
func ValidateStationExistence(stations []Station, startStation, endStation string) {
	startExists := false
	endExists := false

	// Iterate over the slice to find the stations
	for _, station := range stations {
		if station.Name == startStation {
			startExists = true
		}
		if station.Name == endStation {
			endExists = true
		}
	}

	// Check if start station exists
	if !startExists {
		Error(fmt.Sprintf("Start station does not exist: %s", startStation))
	}

	// Check if end station exists
	if !endExists {
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
func ValidatePathExistence(path [][]string, startStation, endStation string) {
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

// CheckDuplicateRoutes checks for duplicate connections, including reversed ones, in a slice of connection pairs.
func CheckDuplicateRoutes(connections [][]string) error {
	seen := make(map[string]bool)
	for _, conn := range connections {
		if len(conn) != 2 {
			return fmt.Errorf("invalid connection format: %v", conn)
		}
		from := conn[0]
		to := conn[1]
		forward := fmt.Sprintf("%s-%s", from, to)
		reverse := fmt.Sprintf("%s-%s", to, from)
		if seen[forward] || seen[reverse] {
			return fmt.Errorf("duplicate connection between %s and %s", from, to)
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

// CheckDuplicateStationNames checks for duplicate station names in a slice of Station.
func CheckDuplicateStationNames(stations []Station) error {
	seen := make(map[string]bool)
	for _, station := range stations {
		name := station.Name
		if seen[name] {
			// Output error message to standard error output
			return fmt.Errorf("duplicate station name detected: %s", name)
		}
		seen[name] = true
	}
	return nil
}

// CheckDuplicateCoordinates checks if any two stations have the same coordinates.
func CheckDuplicateCoordinates(stations []Station) error {
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

func CheckConnectionsExist(stations []Station, connections [][]string) error {
	// Helper function to check if a station exists in the slice
	exists := func(name string) bool {
		for _, station := range stations {
			if station.Name == name {
				return true
			}
		}
		return false
	}

	// Iterate through the connections and validate station existence
	for _, conn := range connections {
		if !exists(conn[0]) {
			return fmt.Errorf("Error: Connection from unknown station: %s", conn[0])
		}
		if !exists(conn[1]) {
			return fmt.Errorf("Error: Connection to unknown station: %s", conn[1])
		}
	}

	return nil
}

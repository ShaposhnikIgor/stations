package train

import (
	"fmt"
	"os"
)

// Error handling for duplicate routes
func CheckDuplicateRoutes(connections []Connection) error {
	routeMap := make(map[string]bool)
	for _, conn := range connections {
		route1 := conn.From + "-" + conn.To
		route2 := conn.To + "-" + conn.From
		if routeMap[route1] || routeMap[route2] {
			return fmt.Errorf("error: duplicate connection between %s and %s", conn.From, conn.To)
		}
		routeMap[route1] = true
	}
	return nil
}

// Error handling for invalid station names
func CheckInvalidStationNames(stations []Station) error {
	for _, station := range stations {
		if station.Name == "" {
			return fmt.Errorf("error: invalid station name")
		}
	}
	return nil
}

// Error handling for missing stations or connections sections
func CheckSections(stationsSectionFound, connectionsSectionFound bool) error {
	if !stationsSectionFound {
		return fmt.Errorf("error: map does not contain a 'stations:' section")
	}
	if !connectionsSectionFound {
		return fmt.Errorf("error: map does not contain a 'connections:' section")
	}
	return nil
}

// Error handling for duplicate station coordinates
func CheckDuplicateCoordinates(stations []Station) error {
	coords := make(map[string]struct{})
	for _, station := range stations {
		coordKey := fmt.Sprintf("%d,%d", station.X, station.Y)
		if _, exists := coords[coordKey]; exists {
			return fmt.Errorf("error: duplicate station coordinates (%d,%d)", station.X, station.Y)
		}
		coords[coordKey] = struct{}{}
	}
	return nil
}

// Error handling for non-positive integer coordinates
func CheckValidCoordinates(stations []Station) error {
	for _, station := range stations {
		if station.X < 0 || station.Y < 0 {
			return fmt.Errorf("error: invalid station coordinates for %s", station.Name)
		}
	}
	return nil
}

// Handle errors with appropriate messages
func HandleError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

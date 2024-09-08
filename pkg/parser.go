package train

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Station struct {
	Name string
	X    int
	Y    int
}

// ParseNetworkMap reads the network map file and returns a list of stations and a 2D slice of connections.
func ParseNetworkMap(filePath string) ([]Station, [][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	var stations []Station
	var connections [][]string
	scanner := bufio.NewScanner(file)
	mode := ""
	stationsSectionFound := false
	connectionsSectionFound := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if hashIndex := strings.Index(line, "#"); hashIndex != -1 {
			line = strings.TrimSpace(line[:hashIndex])
		}

		if len(line) == 0 {
			continue
		}

		if line == "stations:" {
			stationsSectionFound = true
		}
		if line == "connections:" {
			connectionsSectionFound = true
		}
	}

	err = CheckSections(stationsSectionFound, connectionsSectionFound)
	if err != nil {
		return nil, nil, err
	}

	file.Seek(0, 0) // Reset file pointer to the beginning
	scanner = bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Ignore comments
		if hashIndex := strings.Index(line, "#"); hashIndex != -1 {
			line = strings.TrimSpace(line[:hashIndex])
		}

		if len(line) == 0 {
			continue
		}

		if len(stations) > 10000 {
			return nil, nil,

				fmt.Errorf("Error: Map contains more than 10000 stations")

		}

		if line == "stations:" {
			mode = "stations"
			stationsSectionFound = true
			continue
		}
		if line == "connections:" {
			mode = "connections"
			connectionsSectionFound = true
			continue
		}

		if !stationsSectionFound {
			return nil, nil, fmt.Errorf("Error: Map does not contain 'stations:' section")
		}
		if !connectionsSectionFound {
			return nil, nil, fmt.Errorf("Error: Map does not contain 'connections:' section")
		}

		switch mode {
		case "stations":
			parts := strings.Split(line, ",")
			if len(parts) != 3 {
				return nil, nil, fmt.Errorf("error: invalid station format on line: %s", line)
			}
			name := strings.TrimSpace(parts[0])
			x, err1 := strconv.Atoi(strings.TrimSpace(parts[1]))
			y, err2 := strconv.Atoi(strings.TrimSpace(parts[2]))
			if err1 != nil || err2 != nil || x < 0 || y < 0 {
				return nil, nil, fmt.Errorf("error: invalid station coordinates on line: %s", line)
			}
			stations = append(stations, Station{Name: name, X: x, Y: y})
		case "connections":
			parts := strings.Split(line, "-")
			if len(parts) != 2 {
				return nil, nil, fmt.Errorf("error: invalid connection format on line: %s", line)
			}
			from := strings.TrimSpace(parts[0])
			to := strings.TrimSpace(parts[1])
			connections = append(connections, []string{from, to})
		default:
			return nil, nil, fmt.Errorf("error: invalid format on line: %s", line)
		}
	}

	err = CheckDuplicateStationNames(stations)
	if err != nil {
		return nil, nil, err
	}

	err = CheckDuplicateCoordinates(stations)
	if err != nil {
		return nil, nil, err
	}

	err = CheckDuplicateRoutes(connections)
	if err != nil {
		return nil, nil, err
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}
	return stations, connections, nil

}

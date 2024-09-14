package train2

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ParseNetworkMap(filePath string) (map[string]Station, []Connection, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	stations := make(map[string]Station)
	var connections []Connection
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

	// if err := scanner.Err(); err != nil {
	// 	return nil, nil, err
	// }

	err = CheckSections(stationsSectionFound, connectionsSectionFound)
	if err != nil {
		return nil, nil, err
	}

	file.Seek(0, 0) // Reset file pointer to the beginning
	scanner = bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if hashIndex := strings.Index(line, "#"); hashIndex != -1 {
			line = strings.TrimSpace(line[:hashIndex])
		}

		if len(line) == 0 {
			continue
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

		switch mode {
		case "stations":
			parts := strings.Split(line, ",")
			if len(parts) != 3 {
				return nil, nil, fmt.Errorf("error: invalid station format on line: %s", line)
			}
			name := strings.TrimSpace(parts[0])
			// Check for duplicate station names
			if _, exists := stations[name]; exists {
				//fmt.Fprintf(os.Stderr, "Error: duplicate station name found: %s\n", name)
				return nil, nil, fmt.Errorf("duplicate station name found: %s", name)
			}
			x, err1 := strconv.Atoi(strings.TrimSpace(parts[1]))
			y, err2 := strconv.Atoi(strings.TrimSpace(parts[2]))
			if err1 != nil || err2 != nil || x < 0 || y < 0 {
				return nil, nil, fmt.Errorf("error: invalid station coordinates on line: %s", line)
			}
			stations[name] = Station{Name: name, X: x, Y: y}
		case "connections":
			parts := strings.Split(line, "-")
			if len(parts) != 2 {
				return nil, nil, fmt.Errorf("error: invalid connection format on line: %s", line)
			}
			from := strings.TrimSpace(parts[0])
			to := strings.TrimSpace(parts[1])
			connections = append(connections, Connection{From: from, To: to})
		default:
			return nil, nil, fmt.Errorf("error: invalid format on line: %s", line)
		}
	}

	// err = CheckSections(stationsSectionFound, connectionsSectionFound)
	// if err != nil {
	// 	return nil, nil, err
	// }

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

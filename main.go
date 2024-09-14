package main

import (
	"bufio"
	"fmt"
	"os"
	train "stations/pkg"
	train2 "stations/pkg2"
)

func countStations(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		count++
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return count, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error: Please provide the file path to the station map.")
		os.Exit(1)
	}

	filePath := os.Args[1]

	// Counting the number of stations in the file
	stationCount, err := countStations(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to count stations: %v\n", err)
		os.Exit(1)
	}

	// Checking the number of stations and running the appropriate code
	if stationCount > 5000 {
		//fmt.Println("Launching Muvmain because the number of stations is greater than 5000")
		train2.Muvmain()
	} else {
		//fmt.Println("Launching Logmain because the number of stations is less than or equal to 5000")
		train.Logmain()
	}
}

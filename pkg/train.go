package train

import (
	"container/list"
	"errors"
	"fmt"
	"os"
	"strings"
)

func BFS(g *Graph, start, end string) ([]string, error) {
	if _, ok := g.AdjList[start]; !ok {
		return nil, fmt.Errorf("Error: Start station %s does not exist", start)
	}
	if _, ok := g.AdjList[end]; !ok {
		return nil, fmt.Errorf("Error: End station %s does not exist", end)
	}
	if start == end {
		return nil, errors.New("Error: Start and end station are the same")
	}

	visited := make(map[string]bool)
	prev := make(map[string]string)
	queue := list.New()
	queue.PushBack(start)
	visited[start] = true

	for queue.Len() > 0 {
		node := queue.Remove(queue.Front()).(string)
		for _, neighbor := range g.AdjList[node] {
			if !visited[neighbor] {
				visited[neighbor] = true
				prev[neighbor] = node
				queue.PushBack(neighbor)
				if neighbor == end {
					return reconstructPath(prev, start, end), nil
				}
			}
		}
	}
	return nil, errors.New("Error: No path found between start and end stations")
}

func reconstructPath(prev map[string]string, start, end string) []string {
	var path []string
	for at := end; at != ""; at = prev[at] {
		path = append([]string{at}, path...)
		if at == start {
			break
		}
	}
	return path
}

func MoveTrains(graph *Graph, start, end string, numTrains int) {
	path, err := BFS(graph, start, end)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	trains := make([]string, numTrains)
	for i := 0; i < numTrains; i++ {
		trains[i] = start
	}

	var moves []string
	for _, station := range path[1:] {
		var move []string
		for i := 0; i < numTrains; i++ {
			if trains[i] != end {
				move = append(move, fmt.Sprintf("T%d-%s", i+1, station))
				trains[i] = station
			}
		}
		moves = append(moves, strings.Join(move, " "))
	}

	for _, move := range moves {
		fmt.Println(move)
	}
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

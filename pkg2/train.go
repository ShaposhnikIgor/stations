package train2

import (
	"container/heap"
	"container/list"
	"fmt"
	"math"
	"os"
	"strings"
)

type PriorityQueue []*Node

type Node struct {
	station  string
	priority int
	index    int
}

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	node := x.(*Node)
	node.index = n
	*pq = append(*pq, node)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	node := old[n-1]
	old[n-1] = nil
	node.index = -1
	*pq = old[0 : n-1]
	return node
}

// // Station represents a point on the graph with coordinates (for heuristic calculation).
// type Station struct {
// 	X, Y int
// }

// // Graph represents the train network graph.
// type Graph struct {
// 	AdjList  map[string][]string
// 	Stations map[string]Station
// }

// Heuristic function for A* (using Manhattan distance).
func heuristic(from, to Station) int {
	return int(math.Abs(float64(from.X-to.X)) + math.Abs(float64(from.Y-to.Y)))
}

// AStarSearch algorithm finds the shortest path using A* search algorithm.
func AStarSearch(graph *Graph, start, goal string) ([]string, error) {
	goalStation := graph.Stations[goal]

	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &Node{station: start, priority: 0})

	cameFrom := make(map[string]string)
	costSoFar := make(map[string]int)
	cameFrom[start] = ""
	costSoFar[start] = 0

	for pq.Len() > 0 {
		current := heap.Pop(pq).(*Node).station

		if current == goal {
			path := []string{}
			for current != "" {
				path = append([]string{current}, path...)
				current = cameFrom[current]
			}
			return path, nil
		}

		for _, next := range graph.AdjList[current] {
			newCost := costSoFar[current] + 1
			if cost, ok := costSoFar[next]; !ok || newCost < cost {
				costSoFar[next] = newCost
				priority := newCost + heuristic(graph.Stations[next], goalStation)
				heap.Push(pq, &Node{station: next, priority: priority})
				cameFrom[next] = current
			}
		}
	}

	return nil, fmt.Errorf("no path found from %s to %s", start, goal)
}

// BFS algorithm for comparison or fallback.
func BFS(graph *Graph, start, goal string) ([]string, error) {
	visited := make(map[string]bool)
	queue := list.New()
	queue.PushBack(start)
	visited[start] = true
	cameFrom := make(map[string]string)
	cameFrom[start] = ""

	for queue.Len() > 0 {
		current := queue.Remove(queue.Front()).(string)

		if current == goal {
			path := []string{}
			for current != "" {
				path = append([]string{current}, path...)
				current = cameFrom[current]
			}
			return path, nil
		}

		for _, next := range graph.AdjList[current] {
			if !visited[next] {
				queue.PushBack(next)
				visited[next] = true
				cameFrom[next] = current
			}
		}
	}

	return nil, fmt.Errorf("no path found from %s to %s", start, goal)
}

// HybridSearch combines BFS and A* to find paths based on the number of trains.
func HybridSearch(graph *Graph, start, goal string, numTrains int) ([]string, error) {
	if numTrains > 3 { // Example condition to switch algorithms
		return AStarSearch(graph, start, goal)
	}
	return BFS(graph, start, goal)
}

// WORKS!!!
func MoveTrains(graph *Graph, startStation, endStation string, numTrains int) {
	trains := make(map[string]string)          // Map train ID to its current station
	previousStation := make(map[string]string) // Store previous station of trains
	movedAway := make(map[string]bool)         // Track if a train has moved away from the start

	// Initialize trains at the starting station
	for i := 1; i <= numTrains; i++ {
		trainID := fmt.Sprintf("T%d", i)
		trains[trainID] = startStation
		previousStation[trainID] = "" // Initially, the train has no previous station
		movedAway[trainID] = false    // Track that the train has not moved away yet
	}

	// Station occupation map
	occupiedStations := make(map[string]bool)
	// Connection usage map
	usedConnections := make(map[string]bool)

	turns := 0
	for {
		turns++
		turnMovement := []string{}
		done := true

		// Clear station occupation map, except for the end station
		for k := range occupiedStations {
			if k != endStation {
				delete(occupiedStations, k)
			}
		}

		// Clear used connections map for the new turn
		usedConnections = make(map[string]bool)

		// Move each train
		for i := 1; i <= numTrains; i++ {
			trainID := fmt.Sprintf("T%d", i)
			currentStation := trains[trainID]

			if currentStation == endStation {
				continue // If the train has reached the destination, it no longer moves
			}

			// Find path from current station to the destination
			path, err := HybridSearch(graph, currentStation, endStation, numTrains)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			// Get the next station from the path
			nextStation := path[1]
			connection := fmt.Sprintf("%s-%s", currentStation, nextStation)

			// Check if the connection is already used
			if (nextStation != endStation && occupiedStations[nextStation]) || usedConnections[connection] || nextStation == previousStation[trainID] {
				// Attempt to find an alternative path or wait
				foundAlternative := false
				for _, alternativeStation := range graph.AdjList[currentStation] {
					alternativeConnection := fmt.Sprintf("%s-%s", currentStation, alternativeStation)
					if !occupiedStations[alternativeStation] && !usedConnections[alternativeConnection] && alternativeStation != previousStation[trainID] {
						nextStation = alternativeStation
						connection = alternativeConnection
						foundAlternative = true
						break
					}
				}
				if !foundAlternative {
					// If no alternative path found, the train stays at its current station
					turnMovement = append(turnMovement, fmt.Sprintf("%s-%s", trainID, currentStation))
					continue
				}
			}

			// Check if only one train is left that hasn't reached its destination
			remainingTrains := 0
			var remainingTrainID string
			for id := range trains {
				if trains[id] != endStation {
					remainingTrains++
					// remainingTrainID = id
					// remainingTrainID = "T20" + numTrains
					remainingTrainID = fmt.Sprintf("T%d", numTrains)
				}
			}

			// Move the train to the next station if it's not occupied and not the previous station, except at the end station
			if (!occupiedStations[nextStation] || nextStation == endStation) && nextStation != previousStation[trainID] {
				// If only one train is left and a direct path is possible next time, wait and let it go next time directly
				if remainingTrains == 2 && trainID == remainingTrainID && directPathPossible(graph, currentStation, endStation) {
					turnMovement = append(turnMovement, fmt.Sprintf("%s-%s", trainID, currentStation))
				} else {
					previousStation[trainID] = currentStation
					trains[trainID] = nextStation
					if nextStation != startStation { // Only record if the train has moved away from the starting station
						turnMovement = append(turnMovement, fmt.Sprintf("%s-%s", trainID, nextStation))
						movedAway[trainID] = true // Mark the train as having moved away
					}
					if nextStation != endStation {
						occupiedStations[nextStation] = true // Mark the station as occupied unless it's the end station
					}
					usedConnections[connection] = true // Mark the connection as used
					done = false
				}
			} else {
				// Station is occupied or it's the previous station, the train waits
				turnMovement = append(turnMovement, fmt.Sprintf("%s-%s", trainID, currentStation))
			}
		}

		// Print movements of the current turn, but only if there are any movements
		if len(turnMovement) > 0 {
			// Only include movements for trains that have moved away from the starting station
			filteredMovement := []string{}
			for _, move := range turnMovement {
				trainID := strings.Split(move, "-")[0]
				if movedAway[trainID] {
					filteredMovement = append(filteredMovement, move)
				}
			}
			if len(filteredMovement) > 0 {
				fmt.Println(strings.Join(filteredMovement, " "))
			}
		}

		// End if all trains have reached the destination
		if done {
			break
		}
	}
}

func directPathPossible(graph *Graph, currentStation, endStation string) bool {
	// Get the list of directly connected stations from the current station using the graph's adjacency list
	connectedStations := graph.AdjList[currentStation]

	// Check if the end station is in the list of directly connected stations
	for _, station := range connectedStations {
		if station == endStation {
			return true
		}
	}

	// If not found, return false
	return false
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

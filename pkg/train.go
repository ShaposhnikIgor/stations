package train

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

// Heuristic function for A* (using Manhattan distance)
func heuristic(from, to Station) int {
	return int(math.Abs(float64(from.X-to.X)) + math.Abs(float64(from.Y-to.Y)))
}

// A* search algorithm
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

// BFS algorithm for comparison or fallback
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

// Hybrid algorithm combining BFS and A*
func HybridSearch(graph *Graph, start, goal string, numTrains int) ([]string, error) {
	if numTrains > 5 { // Example condition to switch algorithms
		return AStarSearch(graph, start, goal)
	}
	return BFS(graph, start, goal)
}

func MoveTrains(graph *Graph, startStation, endStation string, numTrains int) {
	trains := make(map[string]string)          // Сопоставление поезда с его текущей станцией
	previousStation := make(map[string]string) // Хранение предыдущих станций поездов

	for i := 1; i <= numTrains; i++ {
		trainID := fmt.Sprintf("T%d", i)
		trains[trainID] = startStation
		previousStation[trainID] = "" // Изначально у поезда нет предыдущей станции
	}

	// Карта занятости станций
	occupiedStations := make(map[string]bool)

	turns := 0
	for {
		turns++
		turnMovement := []string{}
		done := true

		// Очистка карты занятости станций
		for k := range occupiedStations {
			delete(occupiedStations, k)
		}

		// Перемещение каждого поезда
		for trainID, currentStation := range trains {
			if currentStation == endStation {
				continue // Если поезд уже на конечной станции, он больше не двигается
			}

			// Находим путь от текущей станции до конечной
			path, err := HybridSearch(graph, currentStation, endStation, numTrains)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			// Следующая станция на пути
			nextStation := path[1] // Путь начинается с текущей станции, поэтому следующая - это вторая

			// Проверяем занята ли станция и не является ли она предыдущей
			if !occupiedStations[nextStation] && nextStation != previousStation[trainID] {
				previousStation[trainID] = currentStation
				trains[trainID] = nextStation
				turnMovement = append(turnMovement, fmt.Sprintf("%s-%s", trainID, nextStation))
				occupiedStations[nextStation] = true // Станция занята
				done = false
			} else {
				// Станция занята или это предыдущая станция, ждём
				turnMovement = append(turnMovement, fmt.Sprintf("%s-%s", trainID, currentStation))
			}
		}

		// Выводим движения текущего хода
		if len(turnMovement) > 0 {
			fmt.Println(strings.Join(turnMovement, " "))
		}

		// Если все поезда достигли конечной станции, заканчиваем
		if done {
			break
		}
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

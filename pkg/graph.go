package train

type Station struct {
	Name string
	X    int
	Y    int
}

type Connection struct {
	From string
	To   string
}

type Graph struct {
	AdjList  map[string][]string
	Stations map[string]Station
}

func NewGraph(connections []Connection, stations []Station) *Graph {
	adjList := make(map[string][]string)
	stationMap := make(map[string]Station)

	for _, conn := range connections {
		if _, ok := adjList[conn.From]; !ok {
			adjList[conn.From] = []string{}
		}
		if _, ok := adjList[conn.To]; !ok {
			adjList[conn.To] = []string{}
		}
		adjList[conn.From] = append(adjList[conn.From], conn.To)
		adjList[conn.To] = append(adjList[conn.To], conn.From)
	}

	for _, station := range stations {
		stationMap[station.Name] = station
	}

	return &Graph{AdjList: adjList, Stations: stationMap}
}

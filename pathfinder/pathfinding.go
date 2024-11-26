package pathfinder

import (
	"fmt"
	"os"
	"stations-pathfinder/network"
)

// PriorityQueue is a simple implementation of a priority queue to manage stations during pathfinding.
type PriorityQueue struct {
	items []*Item
}

// Item represents a station with a priority, used in the PriorityQueue for pathfinding.
type Item struct {
	Value    *network.Station
	Priority int
}

// Push inserts a new item into the priority queue, keeping it sorted by priority.
func (pq *PriorityQueue) Push(item *Item) {
	index := 0
	// Find the correct position for the new item based on its priority
	for ; index < len(pq.items) && pq.items[index].Priority <= item.Priority; index++ {
	}
	// Insert the item at the found position
	pq.items = append(pq.items[:index], append([]*Item{item}, pq.items[index:]...)...)
}

// Pop removes and returns the item with the highest priority (lowest value).
func (pq *PriorityQueue) Pop() *Item {
	if len(pq.items) == 0 {
		return nil // the queue is empty
	}
	item := pq.items[0]
	pq.items = pq.items[1:]
	return item
}

// Len returns the number of items in the priority queue
func (pq *PriorityQueue) Len() int {
	return len(pq.items)
}

// FindPaths finds paths for multiple trains between the start and end stations
func FindPaths(start, end *network.Station, railMap network.RailLineMap, aStar bool, numTrains int, single []network.Station) ([][]network.Station, int) {
	var paths [][]network.Station
	var uniquePaths int
	var short bool
	if numTrains == 1 {
		path := FindShortestPath(start, end, railMap, aStar, false, single)
		if path != nil {
			paths = append(paths, path)
			uniquePaths = 1
			return paths, uniquePaths // Return immediately after finding the first path
		}
	} else {
		counter := 0
		netTrains := numTrains

		// Check for duplicate station coordinates
		stationCoords := make(map[string]bool)
		for _, station := range railMap.Stations {
			coord := fmt.Sprintf("%d,%d", station.X, station.Y)
			if stationCoords[coord] {
				fmt.Fprintln(os.Stderr, "Error: duplicate station coordinates")
				os.Exit(1)
			}
			stationCoords[coord] = true
		}

		// Check for duplicate station names
		stationNames := make(map[string]bool)
		for _, station := range railMap.Stations {
			if stationNames[station.Name] {
				fmt.Fprintln(os.Stderr, "Error: duplicate station name found:", station.Name)
				os.Exit(1)
			}
			stationNames[station.Name] = true
		}
		// Find the initial path for the first train
		path := FindShortestPath(start, end, railMap, aStar, false, single)
		if path == nil {
			if single != nil {
				return nil, 1
			}
			fmt.Fprintln(os.Stderr, "Error: no path found")
			os.Exit(1)
		}
		numTrains--
		paths = append(paths, path)
		if len(path) == 2 {
			short = true // Marks a path directly from start to end for later pathfinding.
		}

		for {
			path := FindShortestPath(start, end, railMap, aStar, short, single)
			// Dispatch efficiency logic:
			if len(path)-len(paths[0]) < numTrains {
				if len(path) != 0 {
					paths = append(paths, path)
					uniquePaths = len(paths)
				} else if len(paths[counter])-len(paths[0]) < numTrains {
					paths = append(paths, paths[counter]) // Once all the new paths are found all the other trains will be assigned the already existing paths, from most efficient to least.
					counter++
					if counter == uniquePaths {
						counter = 0
					}
				} else {
					paths = append(paths, paths[0])
					counter = 0
				}
				numTrains--
				if numTrains <= 0 {
					break
				}
			} else {
				break
			}
		}

		if uniquePaths == 0 {
			// Check for more optimal multi-pathing options.
			clearStations(railMap)
			multiPaths, uniquePaths := FindPaths(start, end, railMap, aStar, netTrains, paths[0])
			if uniquePaths == 1 {
				return paths, uniquePaths
			}
			if RunSchedule(multiPaths, uniquePaths, true) < (len(paths[0]) + netTrains - 1) {
				return multiPaths, uniquePaths
			}
			uniquePaths = 1
		}
	}
	return paths, uniquePaths // Returns a path for every train and the amount of paths that can be started per turn.
}

// findShortestPath computes the shortest path between two stations using A* or basic pathfinding.
func FindShortestPath(start, end *network.Station, connections network.RailLineMap, aStar bool, short bool, single []network.Station) []network.Station {
	openSet := &PriorityQueue{}

	openSet.Push(&Item{
		Value:    start,
		Priority: 0,
	})

	cameFrom := make(map[*network.Station]*network.Station)
	gScore := make(map[string]int)
	fScore := make(map[*network.Station]int)

	// Initialize scores
	for _, station := range connections.Stations {
		gScore[station.Name] = 1 << 20
		fScore[station] = 1 << 20
	}
	gScore[start.Name] = 0
	fScore[start] = heuristic(start.X, start.Y, end.X, end.Y)

	for openSet.Len() > 0 {
		current := openSet.Pop().Value

		if current == end && !(short && cameFrom[current].Name == start.Name) {
			return reconstructPath(cameFrom, end)
		}

		for _, neighbor := range connections.Connections[current] {
			tempGScore := gScore[current.Name] + 1
			if (tempGScore < gScore[neighbor.Name] && !neighbor.Occupied) || (short && neighbor.Name == end.Name) {
				cameFrom[neighbor] = current
				gScore[neighbor.Name] = tempGScore
				fScore[neighbor] = tempGScore + heuristic(neighbor.X, neighbor.Y, end.X, end.Y)

				// Only add the neighbor if not already in openSet
				found := false
				for _, item := range openSet.items {
					if item.Value == neighbor {
						found = true
						break
					}
				}
				if !found {
					priority := fScore[neighbor]
					if !aStar {
						priority = gScore[neighbor.Name]
					}
					openSet.Push(&Item{Value: neighbor, Priority: priority})
				}
			}
		}
	}

	return nil // No path found
}

// Reconstructs the path from the cameFrom map.
func reconstructPath(cameFrom map[*network.Station]*network.Station, current *network.Station) []network.Station {
	path := make([]network.Station, 0)
	path = append(path, *current)
	current = cameFrom[current]

	for current != nil {
		path = append(path, *current)
		current.Occupied = true
		current = cameFrom[current]
	}
	// Reverse the path before returning
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

// heuristic calculates Manhattan distance for the A* algorithm.
func heuristic(x1, y1, x2, y2 int) int {
	return abs(x1-x2) + abs(y1-y2)
}

// abs returns the absolute value of an integer.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// clearStations resets the occupied status of stations.
func clearStations(railMap network.RailLineMap) {
	for _, station := range railMap.Stations {
		station.Occupied = false
	}
}

package pathfinder

import (
	"fmt"
	"stations-pathfinder/network"
)

func RunSchedule(paths [][]network.Station, uniquePaths int, counting bool) int {
	track := uniquePaths
	active := make([][]string, len(paths))
	var done bool
	turnCount := 0

	for turn := 0; ; turn++ {
		done = true
		if turn != 0 {
			track = track + uniquePaths
			if track > len(paths) {
				track = len(paths)
			}
		}

		anyPrinted := false // Track if any train state was printed

		for i := 0; i < track; i++ {
			if len(active[i]) > 0 {
				if active[i][0] == "*" {
					continue
				}
			}
			active[i] = updateActiveStations(active[i], paths[i], active)
			done = false
		}
		if done {
			break
		}

		for i := 0; i < len(active); i++ {
			if len(active[i]) != 0 {
				if active[i][0] != "*" {
					if !counting {
						fmt.Printf("T%d-%s ", i+1, active[i][len(active[i])-1])
						anyPrinted = true // Set to true if something is printed
					}
				}
			}
		}

		// Only print newline if there was any train activity
		if anyPrinted && !counting {
			fmt.Println("")
		} else if counting {
			turnCount++
		}
	}

	return turnCount
}

// UpdateActiveStations switches the name of a train's station to its neighbor's name while avoiding conflicts.
func updateActiveStations(currentStation []string, path []network.Station, active [][]string) []string {
	if currentStation == nil {
		currentStation = []string{path[0].Name}
	}
	index := findStationIndex(path, currentStation[len(currentStation)-1])

	if index+1 < len(path) {
		currentStation = append(currentStation, path[index+1].Name)

		// Check for identical simultaneous paths
		same := false
		for _, otherPath := range active {
			if len(otherPath) != len(currentStation) {
				continue
			}
			same = true
			for i, stop := range otherPath {
				if currentStation[i] != stop {
					same = false
					break
				}
			}
		}
		if same {
			return nil
		}
		return currentStation
	} else {
		// Current station is the last station
		return []string{"*"}
	}
}

func findStationIndex(path []network.Station, targetName string) int {
	for i, station := range path {
		if station.Name == targetName {
			return i
		}
	}
	return -1 // Not found
}

package pathfinder

import (
	"fmt"
	"stations-pathfinder/network"
)

// RunSchedule simulates train movements along the given paths and outputs their progress each turn.
func RunSchedule(paths [][]network.Station, uniquePaths int, counting bool) int {
	track := uniquePaths 
	active := make([][]string, len(paths)) 
	var done bool
	turnCount := 0

	for turn := 0; ; turn++ {
		done = true        
		anyMoved := false

		// Activate additional paths for trains in later turns if needed
		if turn != 0 {
			track = track + uniquePaths             // Add the unique paths to the active tracks
			if track > len(paths) {                // Limit track count to the number of available paths
				track = len(paths)
			}
		}

		// Process train movements for all active tracks
		for i := 0; i < track; i++ {
			if len(active[i]) > 0 && active[i][0] == "*" {
				// Train already reached its destination, skip further movement
				continue
			}

			// Update the current train's position on its path
			active[i] = updateActiveStations(active[i], paths[i], active)
			if active[i] != nil && active[i][0] != "*" {
				anyMoved = true                     
			}
			done = false                            // Mark that trains are still active
		}

		// If no active trains, terminate the loop
		if done {
			break
		}

		// Print the turn results if we're not in counting mode and at least one train moved
		if !counting && anyMoved {
			fmt.Printf("Turn %d: ", turn+1)
			for i := 0; i < len(active); i++ {
				if len(active[i]) > 0 && active[i][0] != "*" {
					// Print the train's current position in the format "T<number>-station"
					fmt.Printf("T%d-%s ", i+1, active[i][len(active[i])-1])
				}
			}
			fmt.Println()                          
		}

		turnCount++                                 
	}

	return turnCount     // total number of turns taken
}

// updateActiveStations moves a train along its path while avoiding conflicts with other trains.
func updateActiveStations(currentStation []string, path []network.Station, active [][]string) []string {
	if currentStation == nil {
		// Initialize train at the starting station of its path
		currentStation = []string{path[0].Name}
	}

	// Find the train's current station in the path
	index := findStationIndex(path, currentStation[len(currentStation)-1])

	if index+1 < len(path) {
		// Move to the next station in the path
		currentStation = append(currentStation, path[index+1].Name)

		// Check for identical simultaneous paths to avoid conflicts
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
			// Conflict found: return nil to prevent movement this turn
			return nil
		}
		return currentStation
	} else {
		return []string{"*"}
	}
}

// findStationIndex locates the position of a station in the given path by its name.
func findStationIndex(path []network.Station, targetName string) int {
	for i, station := range path {
		if station.Name == targetName {
			return i
		}
	}
	return -1 //if the station is not found
}

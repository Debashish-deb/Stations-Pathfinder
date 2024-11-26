package network

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Station represents a station on the railway map.
type Station struct {
	Name     string
	X        int
	Y        int
	Distance int
	Occupied bool
}

// RailLineMap represents the railway map with stations and connections.
type RailLineMap struct {
	Stations    []*Station
	Connections map[*Station][]*Station
}

// BuildStations reads a map file and builds the railway map.
func BuildStations(filePath string) ([]Station, RailLineMap) {
	var stations []Station
	railMap := RailLineMap{
		Stations:    []*Station{},
		Connections: make(map[*Station][]*Station),
	}

	// Create a map to track connections
	connectionMap := make(map[string]bool)

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	stationSection, connectionsSection := false, false
	for scanner.Scan() {
		line := scanner.Text()
		if line == "stations:" {
			stationSection = true
			continue
		}
		if line == "connections:" {
			stationSection, connectionsSection = false, true
			continue
		}
		line = strings.TrimSpace(strings.Split(line, "#")[0])
		if line == "" {
			continue
		}
		if stationSection {
			station := defineStation(line)
			stations = append(stations, station)
			railMap.Stations = append(railMap.Stations, &stations[len(stations)-1])
		} else if connectionsSection {
			railMap = addConnection(line, stations, railMap, connectionMap)
		}
	}
	return stations, railMap
}

func defineStation(line string) Station {
	parts := strings.Split(line, ",")
	x, _ := strconv.Atoi(parts[1])
	y, _ := strconv.Atoi(parts[2])
	return Station{Name: parts[0], X: x, Y: y, Distance: 1 << 20}
}

func addConnection(line string, stations []Station, railMap RailLineMap, connectionMap map[string]bool) RailLineMap {
	stops := strings.Split(line, "-")
	stop1, stop2 := StationLookup(stops[0], stations), StationLookup(stops[1], stations)

	if stop1 == nil || stop2 == nil {
		fmt.Fprintln(os.Stderr, "Error: invalid connection:", line)
		os.Exit(1)
	}

	// Create connection keys for both directions
	connection1 := fmt.Sprintf("%s-%s", stops[0], stops[1])
	connection2 := fmt.Sprintf("%s-%s", stops[1], stops[0])

	// Check if either direction of the connection already exists
	if connectionMap[connection1] || connectionMap[connection2] {
		fmt.Fprintln(os.Stderr, "Error: duplicate connection found:", line)
		os.Exit(1)
	}

	// Store both directions in the connection map
	connectionMap[connection1] = true
	connectionMap[connection2] = true

	// Add the connections to the railMap
	railMap.Connections[stop1] = append(railMap.Connections[stop1], stop2)
	railMap.Connections[stop2] = append(railMap.Connections[stop2], stop1)
	return railMap
}

func StationLookup(name string, stations []Station) *Station {
	for i := range stations {
		if stations[i].Name == name {
			return &stations[i]
		}
	}
	return nil
}

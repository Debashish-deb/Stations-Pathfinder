package main

import (
	"flag"
	"fmt"
	"os"
	"stations-pathfinder/network"
	"stations-pathfinder/pathfinder"
	"strconv"
)

func main() {
	var aStar bool
	flag.BoolVar(&aStar, "a", false, "use A*")
	flag.Parse()

	args := os.Args
	if !((len(args) == 5 && !aStar) || (len(args) == 6 && aStar)) {
		fmt.Fprintln(os.Stderr, "Error: Missing arguments! Stations-Pathfinder usage:\n"+
			"go run . [path to file containing network map] [start station] [end station] [number of trains]\n"+
			"optional flag -a before other arguments to use distance-based pathfinding")
		os.Exit(1)
	}
	if len(args) > 6 {
		fmt.Fprintln(os.Stderr, "Error: Too many arguments provided")
		os.Exit(1)
	}

	argShift := 0
	if len(args) == 6 {
		argShift = 1
	}

	networkMap, startName, endName, trainsToRun := args[1+argShift], args[2+argShift], args[3+argShift], args[4+argShift]
	numTrains, err := strconv.Atoi(trainsToRun)
	if numTrains <= 0 || err != nil {
		fmt.Fprintln(os.Stderr, "Error: number of trains must be positive int")
		os.Exit(1)
	}
	if startName == endName {
		fmt.Fprintln(os.Stderr, "Error: start and end station are the same")
		os.Exit(1)
	}

	stations, connections := network.BuildStations(networkMap)

	start, end := network.StationLookup(startName, stations), network.StationLookup(endName, stations)
	if start == nil {
		fmt.Fprintln(os.Stderr, "Error: Start station not found.")
		os.Exit(1)
	} else if end == nil {
		fmt.Fprintln(os.Stderr, "Error: End station not found.")
		os.Exit(1)
	}

	paths, uniquePaths := pathfinder.FindPaths(start, end, connections, aStar, numTrains, nil)
	if paths == nil {
		fmt.Fprintln(os.Stderr, "Error: no path exists")
		os.Exit(1)
	}

	pathfinder.RunSchedule(paths, uniquePaths, false)
}

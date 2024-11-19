## Stations-Pathfinder

The Stations-Pathfinder helps calculate the best routes for multiple trains, ensuring no two trains end up at the same station at the same time. It works in turns, moving trains across a railway map you provide. Each turn is shown on a new line, with trains labeled using "T" and their number (like T1, T2, etc.).

                        ##Usage##
The program requires four command-line arguments:

It should look like:

go run . [path to file containing network map] [start station] [end station] [number of trains]. Example:


       go run . network.map waterloo st_pancras 4  

optional flag -a before other arguments to use distance-based pathfinding. Example:
     
     
       go run . network.map waterloo st_pancras 4 | wc -l  

..........................................................................................
..........................................................................................
                        ##How It Works##
The program uses two main algorithms to plan train routes:

Dijkstra's Algorithm: This is used to find the most efficient path by counting the number of stations. It's reliable and ensures trains take the shortest possible routes.
*A Algorithm (optional)**: This focuses on physical distance between stations but isn't always the best choice. While it’s available, it’s not usually recommended because it may overlook important factors like the number of stations between two points.
In most cases, Dijkstra’s algorithm is better suited for the task because it balances the routes more effectively.

Dispatching Trains
To manage the trains, the program assigns a path to each one, from the most efficient route to the least, making sure no two trains are at the same station at once. This involves two steps:

Assigning Paths: Routes are calculated for all trains.
Tracking Movement: During each turn, the program keeps track of which station each train is at and prints the results, so you can see how they move through the map.

                ##Visual representation of the Map:##

    0  1 2 3 4 5 6 7
    1      X <- waterloo
    2     / \
    3    /   \
    4   /     \
    5  X <- euston
    6   \       \
    7    \       X <- victoria
    8     \     /
    9      \   /
    10      \ /
    11       X <- st_pancras
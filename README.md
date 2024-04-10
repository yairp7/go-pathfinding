## Pathfinding Using A* Algorithm

#### This a Golang example of how to use the A* Algorithm for pathfinding in games.

#### Iv'e used the open source ebiten game engine to illustrate the path finding.


#### 1. First we build a graph representation of the map's "walkable" tiles using BFS (Each edge get a weight according to the distance - straight movement = 1, diagonal movement = sqrt(2)).

#### 2. We run our A* on the graph with a heuristic function that returns the "aerial" distance between each tile to a destination tile.

[Live WASM Version](https://pf.pech.work?tileSize=32)
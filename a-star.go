package main

import (
	"math"
)

func estimate_distance(start, target coords) float64 {
	est := math.Sqrt(math.Pow(float64(start.x-target.x), 2) + math.Pow(float64(start.y-target.y), 2))
	return est
}

func build_node(prev *node, x, y, terrain, x_t, y_t int) *node {
	d := estimate_distance(coords{x, y}, coords{x_t, y_t})
	n := node{
		loc:      coords{x, y},
		prev:     prev,
		terrain:  terrain,
		cost:     prev.cost + float64(terrain),
		distance: d,
		estimate: prev.cost + float64(terrain) + d,
	}
	return &n
}

func choose_node(pool map[coords]*node) coords {
	var lowest coords

	for i, v := range pool {
		c, ok := pool[lowest]
		if ok {
			if v.estimate < c.estimate {
				lowest = i
			}
		} else {
			lowest = i
		}

	}
	return lowest
}

// path find
func (e *spawn) pathfind(b *board) {
	if e.target != nil {
		if estimate_distance(e.loc, e.target[0].loc) > e.attack_range {
			open_nodes := make(map[coords]*node)
			closed_nodes := make(map[coords]*node)

			open_nodes[e.loc] = &node{
				loc:      e.loc,
				prev:     nil,
				terrain:  b[terrain_layer][e.loc.y][e.loc.x],
				cost:     float64(b[terrain_layer][e.loc.y][e.loc.x]),
				distance: estimate_distance(e.loc, e.target[0].loc),
				estimate: 0,
			}
			open_nodes[e.loc].estimate = open_nodes[e.loc].cost + open_nodes[e.loc].distance

			// A* begins loop here
			for len(open_nodes) > 0 {
				// get next most likely node to follow
				seed_coords := choose_node(open_nodes)

				// find neighbours
				for x := int(math.Max(0, float64(seed_coords.x-1))); x <= int(math.Min(board_cells_Wide-1, float64(seed_coords.x+1))); x++ {
					for y := int(math.Max(0, float64(seed_coords.y-1))); y <= int(math.Min(board_cells_High-1, float64(seed_coords.y+1))); y++ {
						// validate nodes
						if b[terrain_layer][y][x] < impassable_threshold {
							if previously_visited, ok := closed_nodes[coords{x, y}]; ok {
								if previously_visited.cost > open_nodes[seed_coords].cost+float64(b[terrain_layer][y][x]) {
									// this is a better path to a previously investigated path
									// move the visited node back out to OPEN, updating the predecessor and cost.
									open_nodes[coords{x, y}] = closed_nodes[coords{x, y}]
									open_nodes[coords{x, y}].prev = open_nodes[seed_coords]
									open_nodes[coords{x, y}].cost = open_nodes[coords{x, y}].prev.cost + float64(b[terrain_layer][x][y])
									delete(closed_nodes, coords{x, y})
								} // else, if the previous visit resulted in  a lower cost, there's nothing to do
							} else {
								// the node hasn't been visited before. Perform the same evaluation (is this path better than what's arleady in there.
								// The assumption might be that it's not, or we would have gotten here first. TODO: Validate this assumption.)
								if _, ok := open_nodes[coords{x, y}]; !ok {
									open_nodes[coords{x, y}] = build_node(open_nodes[seed_coords], x, y, b[terrain_layer][y][x], e.target[0].loc.x, e.target[0].loc.y)
								} // else there was already an entry for this node so skip it.
							}
							// TODO: check if this node is the target. SHOULD the loop continue until open_nodes is empty? Or until it reaches target?
							// TODO: What happens if open nodes is emptied but goal hasn't been reached?
						}
					}
				}
				// move seed to closed
				closed_nodes[seed_coords] = open_nodes[seed_coords]
				delete(open_nodes, seed_coords)
			}
			// Target was found, work our way back, populating the path
			next := coords{e.target[0].loc.x, e.target[0].loc.y}
			if _, ok := closed_nodes[next]; ok {
				for closed_nodes[next].prev != nil {
					e.path = append(e.path, next)
					next = closed_nodes[next].prev.loc
				}
			}
		}
	}
}

// path find
func (e *spawner) pathfind(b *board) {
	if e.target != nil {
		//if estimate_distance(e.loc, e.target[0].loc) > e.attack_range {
		open_nodes := make(map[coords]*node)
		closed_nodes := make(map[coords]*node)

		open_nodes[e.loc] = &node{
			loc:      e.loc,
			prev:     nil,
			terrain:  b[terrain_layer][e.loc.y][e.loc.x],
			cost:     float64(b[terrain_layer][e.loc.y][e.loc.x]),
			distance: estimate_distance(e.loc, e.target[0].loc),
			estimate: 0,
		}
		open_nodes[e.loc].estimate = open_nodes[e.loc].cost + open_nodes[e.loc].distance

		// A* begins loop here
		for len(open_nodes) > 0 {
			// get next most likely node to follow
			seed_coords := choose_node(open_nodes)

			// find neighbours
			for x := int(math.Max(0, float64(seed_coords.x-1))); x <= int(math.Min(board_cells_Wide-1, float64(seed_coords.x+1))); x++ {
				for y := int(math.Max(0, float64(seed_coords.y-1))); y <= int(math.Min(board_cells_High-1, float64(seed_coords.y+1))); y++ {
					// validate nodes
					if b[terrain_layer][y][x] < impassable_threshold {
						if previously_visited, ok := closed_nodes[coords{x, y}]; ok {
							if previously_visited.cost > open_nodes[seed_coords].cost+float64(b[terrain_layer][y][x]) {
								// this is a better path to a previously investigated path
								// move the visited node back out to OPEN, updating the predecessor and cost.
								open_nodes[coords{x, y}] = closed_nodes[coords{x, y}]
								open_nodes[coords{x, y}].prev = open_nodes[seed_coords]
								open_nodes[coords{x, y}].cost = open_nodes[coords{x, y}].prev.cost + float64(b[terrain_layer][x][y])
								delete(closed_nodes, coords{x, y})
							} // else, if the previous visit resulted in  a lower cost, there's nothing to do
						} else {
							// the node hasn't been visited before. Perform the same evaluation (is this path better than what's arleady in there.
							// The assumption might be that it's not, or we would have gotten here first. TODO: Validate this assumption.)
							if _, ok := open_nodes[coords{x, y}]; !ok {
								open_nodes[coords{x, y}] = build_node(open_nodes[seed_coords], x, y, b[terrain_layer][y][x], e.target[0].loc.x, e.target[0].loc.y)
							} // else there was already an entry for this node so skip it.
						}
						// TODO: check if this node is the target. SHOULD the loop continue until open_nodes is empty? Or until it reaches target?
						// TODO: What happens if open nodes is emptied but goal hasn't been reached?
					}
				}
			}
			// move seed to closed
			closed_nodes[seed_coords] = open_nodes[seed_coords]
			delete(open_nodes, seed_coords)
		}
		// Target was found, work our way back, populating the path
		next := coords{e.target[0].loc.x, e.target[0].loc.y}
		if _, ok := closed_nodes[next]; ok {
			for closed_nodes[next].prev != nil {
				e.path = append(e.path, next)
				next = closed_nodes[next].prev.loc
			}
		}
	}
	//}
}

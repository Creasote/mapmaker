package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

const (
	wall_threshold = 4 // if < wall_threshold neighbours, floor. Else, wall

)

func load_map() (board, error) {
	b, err := os.ReadFile("map_file.txt")
	var board_file board
	if err != nil {
		fmt.Println("Failed to load file from disk.")
		return board_file, err
	}

	err = json.Unmarshal(b, &board_file)
	if err != nil {
		fmt.Println("Error unmarshalling board file.")
	}

	return board_file, nil
}

func save_map(map_data board) error {
	output_data, err := json.Marshal(map_data)
	if err != nil {
		return err
	}

	err = os.WriteFile("map_file.txt", output_data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func randomise_board(b *board) board {
	rand.Seed(time.Now().Unix())

	for y := 0; y < board_cells_High; y++ {
		for x := 0; x < board_cells_Wide; x++ {
			b[terrain_layer][y][x] = rand.Intn(2)
			if b[terrain_layer][y][x] == 0 {
				b[terrain_layer][y][x] = 2
			} else if b[terrain_layer][y][x] == 1 {
				b[terrain_layer][y][x] = 99
			}
		}
	}

	return *b
}

func count_neighbours(b *board, cell_x, cell_y int) int {
	wall_count, floor_count := 0, 0
	for y := cell_y - 1; y < cell_y+2; y++ {
		for x := cell_x - 1; x < cell_x+2; x++ {
			if y < 0 || y >= board_cells_High || x < 0 || x >= board_cells_Wide {
				wall_count++
			} else if b[terrain_layer][y][x] == grassland {
				floor_count++
			} else {
				wall_count++
			}
		}
	}

	if wall_count > wall_threshold {
		return wall
	}
	return grassland

	// if wall_count < wall_threshold {
	// 	return grassland
	// }
	// return wall

}

func optimise_board(b *board) {
	var temporary_board board
	for y := 0; y < board_cells_High; y++ {
		for x := 0; x < board_cells_Wide; x++ {
			temporary_board[terrain_layer][y][x] = count_neighbours(b, x, y)
		}
	}
	for y := 0; y < board_cells_High; y++ {
		for x := 0; x < board_cells_Wide; x++ {
			b[terrain_layer][y][x] = temporary_board[terrain_layer][y][x]
		}
	}
}

func flood_neighbours(b *board, x, y int) []coords {
	//neighbours := map[coords]int{}
	var open_nodes []coords
	var closed_nodes []coords

	open_nodes = append(open_nodes, coords{x, y})
	counter := 0
	for len(open_nodes) > 0 {
		// Use the first entry as the seed
		loc := open_nodes[0]
		fmt.Println("Starting new iteration, lenght: ", len(open_nodes))
		// Move it to closed nodes
		closed_nodes = append(closed_nodes, open_nodes[0])
		counter++
		if len(open_nodes) > 1 {
			open_nodes = open_nodes[1:]
		} else {
			open_nodes = open_nodes[:0]
		}
		// Check each of its' neighbours.
		// for y_range := maxInt(0, loc.y-1); y_range < minInt(board_cells_High, loc.y+2); y_range++ {
		// 	for x_range := maxInt(0, loc.x-1); x_range < minInt(board_cells_Wide, loc.x+2); x_range++ {
		// 		// If they are the right value AND
		// 		if b[terrain_layer][y_range][x_range] == b[terrain_layer][loc.y][loc.x] {
		// 			// they are not in the closed list,
		// 			if !sliceContains(closed_nodes, coords{x_range, y_range}) && !sliceContains(open_nodes, coords{x_range, y_range}) {
		// 				fmt.Println("Adding new: ", coords{x_range, y_range})
		// 				// add to the open list.
		// 				open_nodes = append(open_nodes, coords{x_range, y_range})
		// 			} else {
		// 				fmt.Println("Skipping repeat: ", coords{x_range, y_range})
		// 			}
		// 		}
		// 	}
		// }
		for y_range := maxInt(0, loc.y-1); y_range < minInt(board_cells_High, loc.y+2); y_range++ {
			if b[terrain_layer][y_range][loc.x] == b[terrain_layer][loc.y][loc.x] {
				// they are not in the closed list,
				if !sliceContains(closed_nodes, coords{loc.x, y_range}) && !sliceContains(open_nodes, coords{loc.x, y_range}) {
					//fmt.Println("Adding new: ", coords{loc.x, y_range})
					// add to the open list.
					open_nodes = append(open_nodes, coords{loc.x, y_range})
				}
			}
		}
		for x_range := maxInt(0, loc.x-1); x_range < minInt(board_cells_Wide, loc.x+2); x_range++ {
			// If they are the right value AND
			if b[terrain_layer][loc.y][x_range] == b[terrain_layer][loc.y][loc.x] {
				// they are not in the closed list,
				if !sliceContains(closed_nodes, coords{x_range, loc.y}) && !sliceContains(open_nodes, coords{x_range, loc.y}) {
					//fmt.Println("Adding new: ", coords{x_range, loc.y})
					// add to the open list.
					open_nodes = append(open_nodes, coords{x_range, loc.y})
				}
			}
		}

	}

	// //neighbours := []coords{}
	// fmt.Println("New seed: ", coords{x, y})

	// for y_range := maxInt(0, y-1); y_range < minInt(cellsHigh, y+1); y_range++ {
	// 	for x_range := maxInt(0, x-1); x_range < minInt(cellsWide, x+1); x_range++ {
	// 		if b[terrain_layer][y_range][x_range] == b[terrain_layer][y][x] {
	// 			// add it to candidate list
	// 			// add it to list to transform
	// 			if !sliceContains(*neighbours, coords{x_range, y_range}) {
	// 				fmt.Println("Adding new: ", coords{x_range, y_range})
	// 				*neighbours = append(*neighbours, coords{x_range, y_range})
	// 				*neighbours = append(*neighbours, flood_neighbours(b, x_range, y_range, neighbours)...)
	// 			} else {
	// 				fmt.Println("Skipping repeated: ", coords{x_range, y_range})
	// 			}
	// 			//neighbours[coords{x_range, y_range}] = 1
	// 		}
	// 	}
	// }
	//for neighbour_cells := range neighbours {
	//for x_range := range neighbours {
	//neighbours = append(neighbours, flood_neighbours(b, x_range, y_range, current_terrain)...)

	//	}
	//}

	//for i:=range
	fmt.Println("I think I'm exporting ", counter)
	return closed_nodes
}

// func flood_neighbours(b *board, x,y){
// 	candidates:=[]coords

// }

func (g *Game) place_terrain(b *board, x, y int) {
	if g.flood_mode {
		// check neighbours and fill
		//current_terrain := b[terrain_layer][y][x]
		//flood_to := flood_neighbours(b, x, y)
		var flood_to []coords
		flood_to = append(flood_to, flood_neighbours(b, x, y)...)
		fmt.Println("Flood will fill ", len(flood_to), " neighbours")
		for _, c := range flood_to {
			b[terrain_layer][c.y][c.x] = g.object_value
		}

	} else {
		b[terrain_layer][y][x] = g.object_value
	}
}

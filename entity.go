package main

func place_entity(x, y int) {

	entity_list = append(entity_list, &entity{
		name:               "Clicker",
		loc:                coords{x, y},
		mob_type:           0,
		sprite_img:         img_player,
		movement_speed:     50,
		health:             0,
		armour:             0,
		damage_per_attack:  0,
		attacks_per_second: 0,
		attack_success_pc:  0,
		attack_range:       2,
		last_attack_time:   0,
		target:             entity_list[0],
		path:               []coords{},
	})
	go entity_list[len(entity_list)-1].move_entity()

}

func set_goal(x, y int) {
	entity_list[0].loc = coords{x, y}
	game_map[entity_layer][y][x] = 1 // Sets goal, able to be retrieved on load.
}

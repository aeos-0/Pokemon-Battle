package main

import (
	"database/sql"
	"sort"

	_ "github.com/mattn/go-sqlite3"
)

// Check if pokemon exists in pokemon list of given player
func CheckPokemon(name string, you Player) bool {
	for _, value := range you.pokemon_list {
		if name == value.species {
			//Pokemon exists so it is ok to return true
			return true
		}
	}
	//Pokemon was not found in this list
	return false
}

// Query for pokemon by name
func GetPokemon(name string, db *sql.DB) (Pokemon, error) {
	var mon Pokemon

	rows, err := db.Query("SELECT * FROM pokemon WHERE name = ?", name)
	if err != nil {
		return Pokemon{}, err
	}
	defer rows.Close()

	if rows.Next() {
		// Use the Scan method to copy values from the row into the struct fields
		err := rows.Scan(&mon.pokedex_number, &mon.species, &mon.abilities, &mon.base_stats[1], &mon.base_stats[2], &mon.base_stats[0], &mon.genderless, &mon.base_stats[3], &mon.base_stats[4], &mon.base_stats[5], &mon.type1, &mon.type2, &mon.weight, &mon.can_dynamax, &mon.has_mega)
		if err != nil {
			return Pokemon{}, err
		}
	}

	return mon, nil
}

// Query for move by name
func GetMove(name string, db *sql.DB) (Moves, error) {
	var move Moves

	rows, err := db.Query("SELECT * FROM moves WHERE name = ?;", name)
	if err != nil {
		return Moves{}, err
	}
	defer rows.Close()

	if rows.Next() {
		// Use the Scan method to copy values from the row into the struct fields
		err := rows.Scan(&move.name, &move.type1, &move.base_power, &move.accuracy, &move.category, &move.priority, &move.pp, &move.crit, &move.description, &move.effect,
			&move.effect_rate, &move.target, &move.contact, &move.sound, &move.punch, &move.bite, &move.can_snatch,
			&move.slice, &move.bullet, &move.wind, &move.powder, &move.can_metro, &move.gravity, &move.defrosts, &move.bounceable, &move.can_block, &move.can_mirror)
		if err != nil {
			return Moves{}, err
		}
	}

	return move, nil
}

// Get all move names currently available
func GetMoveNames(db *sql.DB) []string {
	var moveNames []string
	var move Moves

	rows, err := db.Query("SELECT name FROM moves;")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		// Use the Scan method to copy values from the row into the struct fields
		err := rows.Scan(&move.name)
		moveNames = append(moveNames, move.name)
		if err != nil {
			panic(err)
		}
	}

	return moveNames
}

// Get all pokemon names currently available
func GetPokemonNames(db *sql.DB) []string {
	var monNames []string
	var mon Pokemon

	rows, err := db.Query("SELECT name FROM pokemon;")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		// Use the Scan method to copy values from the row into the struct fields
		err := rows.Scan(&mon.species)
		monNames = append(monNames, mon.species)
		if err != nil {
			panic(err)
		}
	}

	sort.Strings(monNames) //Might be expensive
	return monNames
}

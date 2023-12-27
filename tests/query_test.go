package main

import (
	"testing"
	"database/sql"
	"sort"

	_ "github.com/mattn/go-sqlite3"
)

type Pokemon struct {
	pokedex_number int32
	species        string
	type1          string
	type2          string
	hp             uint
	attack         uint
	defense        uint
	sp_attack      uint
	sp_defense     uint
	speed          uint
	stats          [6]uint //Total hp was removed and added into this array (hp, attack, defense, spa, spd, speed)
	weight         float64 //In kilograms (10)
	accuracy       float64
	evasion        float64
	ability        string //Ability selected (this value can change)
	ability_read   string //This will be the ability value that can change (via skill swap and such)
	abilities      string
	nature         string
	genderless     bool
	gender         string //Randomly selected
	has_mega       bool
	can_dynamax    bool
	can_learn      string //(Not added right now)
	UI             string
	//moveset        [4]Moves
	boosts         uint

	action      string //What will the pokemon do this turn (Swap, Move, Mega) (Maybe "Thunderbolt:TeraNormal" for other actions + megas)
	priority    uint32 //Only used during move priority calculation
	status      string
	temp_status []string //Map[string]?? - Probably not worth since this value changing isn't very common
	item        string
	on_field    bool
}

func GetPokemonNames() []string {
	db, err := sql.Open("sqlite3", "../pokedata.db")
	if err != nil {
		panic(err)
	}

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

func BenchmarkPokemon(b *testing.B) {	
	for i := 0; i < b.N; i++ {
		GetPokemonNames()
	}
}
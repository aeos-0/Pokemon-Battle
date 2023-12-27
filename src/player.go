package main

type Player struct {
	hazards      map[string]uint //Uint is how many layers of it exist
	light_screen uint            //Uint is how many turns are left of it
	reflect      uint            //Uint is how many turns are left of it
	tailwind     uint
	pokemon_list []*Pokemon //List of pokemon owned by player
}

// Returns mon from pokemon name (May not be used)
func (you *Player) GetFromName(name string) *Pokemon {
	for index, mon := range you.pokemon_list {
		if mon.species == name {
			return you.pokemon_list[index] //Need to return reference so its methods can be called
		}
	}

	//If you're here the mon doesn't exist in the list
	defer panic("Are you sure you want to be here?\nIt looks like that name doesn't exist on the team")
	return you.pokemon_list[0]
}

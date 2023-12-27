package main

import (
	"fmt"
	"log"
	"math/rand"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Field struct {
	turn          int
	format        string
	weather       map[string]uint   //Turns left of specified weather condition (only one can be active at once)
	terrain       map[string]uint   //Turns left of specified terrain (only one can be active at once)
	gravity       uint              //Turns left of gravity
	magic_room    uint              //Turns left of magic room
	wonder_room   uint              //Turns left of wonder room
	trick_room    uint              //Turns left of trick room
	mons          []*Pokemon        //Pokemon currently in play
	queue         []*Pokemon        //Action queue (In singles 2 at most, in doubles 4 at most)
	abilities     map[string]string //Pokemon abilities (strings are mapped to event)
	matchup       muMatchup         //Matchup table containing mutex and json map
	yourAliveMons int
	cpuAliveMons  int
}

// Add pokemon to the field
func (field *Field) Add(pokemon *Pokemon) {
	field.mons = append(field.mons, pokemon)
	pokemon.on_field = true
}

// Remove pokemon from the field
func (field *Field) Delete(pokemon *Pokemon) {
	//Get index of passed pokemon by UI and delete that value
	field.mons = slices.DeleteFunc(field.mons, func(p *Pokemon) bool {
		return p.UI == pokemon.UI //could also be substituted by name
	})
	pokemon.on_field = false
}

func (play *Field) GetFromUI(ui string, you *Player, opponent *Player) *Pokemon {
	for index, mon := range you.pokemon_list {
		if mon.UI == ui {
			return you.pokemon_list[index] //Need to return reference so its methods can be called
		}
	}

	for index, mon := range opponent.pokemon_list {
		if mon.UI == ui {
			return opponent.pokemon_list[index] //Need to return reference so its methods can be called
		}
	}

	//If you're here the mon doesn't exist in the list
	panic("Are you sure you want to be here?\nIt looks like that name doesn't exist on the team")
}

// Check events: onswap, ondamagedeal, ondamagerecieve, onfaint, onturnend, onstatlower, onstatus
func (play *Field) CheckEvent(event string, pokemon *Pokemon) bool {
	if play.turn == 0 { //Before first turn so just check onswap abilities
		var monList []*Pokemon
		for _, value := range play.mons {
			if play.abilities[value.ability_read] == event { //Check if ability is "onswap" (Be sure to dereference)
				monList = append(monList, value) //If any mon on the fields ability matches the event then add to the list
			}
		}
		if monList == nil {
			//No pokemon match this event so just return
			return false
		}

		//If you got here some pokemon needs an event check
		sort.Slice(monList, func(i, j int) bool { //Sort priority mons based on speed
			return monList[i].stats[5] > monList[j].stats[5]
		})

		//Before first turn starts you have a (monList) of pokemon with "onswap" abilities ordered by speed
		for _, pokemon := range monList {
			pokemon.ActivateAbility(play) //Pass in playing field, and any pokemon this targets specifically
		}
		return false
	}

	//Battle is turn 1 or later
	if pokemon != nil {
		switch event { //Do checks on each and if they pass then call .ActivateAbility (for example checking if an electric type move was used for lightning rod)
		case "onswap":
			if play.abilities[pokemon.ability_read] == "onswap" {
				pokemon.ActivateAbility(play)
				return true
			}
		case "ondamagedeal":
			if play.abilities[pokemon.ability_read] == "ondamagedeal" {
				pokemon.ActivateAbility(play)
				return true
			}
		case "ondamagerecieve":
			if play.abilities[pokemon.ability_read] == "ondamagerecieve" {
				pokemon.ActivateAbility(play)
				return true
			}
		case "onstatlower":
			if play.abilities[pokemon.ability_read] == "onstatlower" {
				pokemon.ActivateAbility(play)
				return true
			}
		case "onstatus":
			if play.abilities[pokemon.ability_read] == "onstatus" {
				pokemon.ActivateAbility(play)
				return true
			}
		case "onturnend":
			if play.abilities[pokemon.ability_read] == "onturnend" {
				pokemon.ActivateAbility(play)
				return true
			}
			//ADD CHECK FOR POKEMON ON FIELD TO SEE IF MORE NEED TO BE SWAPPED IN
		}
	}
	return false
}

// Get queue of pokemon based on speed stats and priority moves
func (play *Field) GetQueue() {
	if play.turn > 0 { //Queue is empty here but field mons shouldn't be
		var newSlice []*Pokemon
		//Check for swaps
		for _, value := range play.mons { //Should never be null otherwise the battle should end
			if strings.Contains(value.action[0], "Swap") && value.status != "Fainted" {
				newSlice = append(newSlice, value)
			}
		}
		if newSlice != nil {
			sort.Slice(newSlice, func(i, j int) bool { //Sort swaps
				return newSlice[i].stats[5] > newSlice[j].stats[5]
			})
			play.queue = append(play.queue, newSlice...) //Add swaps to actions queue
		}

		/*Megas, teras go here or in another function*/
		//Check for recharging
		newSlice = nil
		for _, value := range play.mons {
			if value.action[0] == "Recharging" && value.status != "Fainted" {
				newSlice = append(newSlice, value)
			}
		}
		if newSlice != nil {
			sort.Slice(newSlice, func(i, j int) bool { //Sort by speed
				return newSlice[i].stats[5] > newSlice[j].stats[5]
			})
			play.queue = append(play.queue, newSlice...) //Add recharging to actions queue
		}

		//Check for moves, after swaps moves are the only thing left that can be done
		newSlice = nil
		for _, mon := range play.mons {
			if !slices.Contains(play.queue, mon) && mon.status != "Fainted" && !mon.has_acted { //Get pokemon that arent swapping and add to newSlice
				newSlice = append(newSlice, mon)
			}
		}

		if newSlice == nil {
			//All pokemon have swapped this turn and/or are recharging
			return
		}

		//Assign each mon its priority value based of its move
		for _, mon := range newSlice {
			index, _ := strconv.ParseInt(mon.action[0], 10, 32)
			mon.priority = mon.moveset[index].priority
		}

		if newSlice != nil {
			//Sort by move priority number instead of speed
			sort.Slice(newSlice, func(i, j int) bool { //Sort priority mons based on speed
				return newSlice[i].priority > newSlice[j].priority
			})

			if len(newSlice) > 1 {
				speed_sort(&newSlice, play) //Sort priority values on speed
			}

			play.queue = append(play.queue, newSlice...) //Add priority moves to actions queue
		}
	} else {
		//Since this is before turn one just sort speed values
		play.queue = append(play.queue, play.mons...) //Add field mons to queue

		sort.Slice(play.queue, func(i, j int) bool { //Sort field mons
			return play.queue[i].stats[5] > play.queue[j].stats[5]
		})
	}
}

// After sorting for priority it sorts those moves based on the speed stat as well
func speed_sort(list *[]*Pokemon, play *Field) {
	//Here the list has sorted pokemon by priority values, now they needed to be sorted by speed in each respective priority bracket
	seen := make(map[int32]bool)
	var tempQueue []*Pokemon
	var newQueue []*Pokemon

	//Get groups of mons that exist in a priority bracket and then sort off that
	for index, pokemon := range *list {
		if !seen[pokemon.priority] { //Has this priority number already been checked
			seen[pokemon.priority] = true
			tempQueue = append(tempQueue, pokemon)
			for j := index + 1; j < len(*list); j++ { //Get all other instances of this value
				if (*list)[j].priority == pokemon.priority {
					//There is a match so add it to tempQueue
					tempQueue = append(tempQueue, (*list)[j])
				}
			}

			if len(tempQueue) > 1 { //If it has more than 1 value in it then sort
				//Sort temp queue
				if play.trick_room > 0 { //Trick room is up so invert sort on each priority bracket
					sort.Slice(tempQueue, func(i, j int) bool { //Sort mons based on speed
						return tempQueue[i].stats[5] < tempQueue[j].stats[5]
					})
				} else { //Trick room is not up so sort normally
					sort.Slice(tempQueue, func(i, j int) bool { //Sort mons based on speed
						return tempQueue[i].stats[5] > tempQueue[j].stats[5]
					})
				}
			}

			//Add temp queue to new queue
			newQueue = append(newQueue, tempQueue...)
		}
	}
	*list = newQueue //List = new sorted list value
}

// The queue should be filled and all the actions in the turn should be executed
func (play *Field) Run(you *Player, opponent *Player) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for _, pokemon := range play.queue {
		//Check if pokemon fainted
		if pokemon.status != "Fainted" && !pokemon.has_acted { //Remove has acted???
			//Do what ever the pokemon intended to do
			action := pokemon.action[0] //Contains move index or swap
			if action == "Swap" {
				name := pokemon.action[1]
				pokemon.Swap(play, you, you.GetFromName(name)) //Gets pokemon to swap from UI
				continue
			}

			//If you got here the pokemon decided to do a move, damage dealing or not
			target := pokemon.action[1]

			if target != "_" { //Print using move text
				if !pokemon.owner {
					fmt.Print("The opposing ")
				}
				index, _ := strconv.ParseInt(action, 10, 32)
				fmt.Println(pokemon.species, "used", pokemon.moveset[index].name)
			} else {
				if !pokemon.owner {
					fmt.Print("The opposing ")
				}
				index, _ := strconv.ParseInt(action, 10, 32)
				fmt.Println(pokemon.species, "used", pokemon.moveset[index].name, "but it failed")
				continue //Move failed so continue to next iteration
			}

			switch target { //Numbers 0 - 5 (also "_")
			case "0":
				pokemon.DamageCalc(pokemon, play, r) //Target is self
			case "1": //Attack opposing pokemon
				for _, mon := range play.mons {
					if mon.owner != pokemon.owner && mon.slot == 1 { //If pokemon is on field and the attacking pokemon isnt its ally
						pokemon.DamageCalc(mon, play, r)
					}
				}
			case "2": //Attack opposing pokemon
				for _, mon := range play.mons {
					if mon.owner != pokemon.owner && mon.slot == 2 { //If pokemon is on field and the attacking pokemon isnt its ally
						pokemon.DamageCalc(mon, play, r)
					}
				}
			case "3": //Attack ally
				for _, mon := range play.mons {
					if mon.owner == pokemon.owner && mon.UI != pokemon.UI {
						pokemon.DamageCalc(mon, play, r)
					}
				}
			case "4": //Attack all opponents
				for _, mon := range play.mons {
					if mon.UI != pokemon.UI && mon.owner != pokemon.owner {
						go pokemon.DamageCalc(mon, play, r)
					}
				}
			case "5": //Attack all pokemon
				for _, mon := range play.mons {
					if mon.UI != pokemon.UI {
						go pokemon.DamageCalc(mon, play, r)
					}
				}
			default:
				log.Fatal("Invalid pokemon target", "\nSpecies:", pokemon.species, "\nOwner: ", pokemon.owner, "\nTarget: ", target)
			}

			//Pokemon moved this turn
			pokemon.has_acted = true
		}
	}
}

// Check for tick damage, reset values
func (play *Field) Reset() {
	for _, mon := range play.queue {
		mon.action[0] = ""
		mon.action[1] = ""
		mon.priority = 0
		mon.has_acted = false

		//Status turns
		if mon.status_turns != 0 {
			mon.status_turns--
		}

		//Tick damage
		switch mon.status {
		case "Burned":
		//Take 1/16th damage
		case "Poisoned":
			//Take 1/8th damage
		case "Badly Poisoned":
			//Take 1/8th damage with multiplier
		}
	}
	play.queue = nil
}

package main

import (
	"fmt"
	"math"
	"math/rand"
	"slices"
	"strconv"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type Pokemon struct {
	pokedex_number int32
	species        string
	type1          string
	type2          string
	base_stats     [6]uint
	stats          [6]uint //Total hp was removed and added into this array (hp, attack, defense, spa, spd, speed) Change to stat changes?
	stat_stages    [7]int8 //Values -6 - +6 for each stat excluding hp (also including accuracy and evasion)
	weight         float64 //In kilograms
	ability        string  //Ability selected (this value can change)
	ability_read   string  //This will be the ability value that can change (via skill swap and such)
	abilities      string
	nature         string
	genderless     bool
	gender         string //Randomly selected
	has_mega       bool
	can_dynamax    bool
	UI             string
	moveset        [4]Moves
	boosts         uint

	action       [2]string //What will the pokemon do this turn (Swap, Move, Mega) (Maybe "Thunderbolt:TeraNormal" for other actions + megas)
	priority     int32     //Only used during move priority calculation
	status       string
	slot         uint8    //Values: 0-5 choose which "slot" the pokemon attacks into: 0 - self, 1 - left target, 2 - right target, 3 - allies, 4 - both enemies, 5 - everyone
	temp_status  []string //Map[string]?? - Probably not worth since this value changing isn't very common
	status_turns uint8    //Sleep turns, poison turns
	item         string
	on_field     bool
	owner        bool //True for you, false for opponent
	has_acted    bool //Has the pokemon moved yet this turn
}

// Calculate a pokemon's stats given nature, evs, and ivs
func (pokemon *Pokemon) CalcStats(ev_hp int32, ev_a int32, ev_d int32, ev_spa int32, ev_spd int32, ev_speed int32) {
	//Randomizer
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	//Constant IVs
	var iv int32 = 31
	//Constant level
	level := 50
	//Evs typically have a overall cap per pokemon but I will be testing how it is without said cap
	var ev int32
	var nature_multiplier float64

	//Calculate hp
	if ev_hp == -1 || ev_hp > 255 {
		ev = r.Int31n(256) //Select random value (not inclusive)
	} else {
		ev = ev_hp //Else select value given
	}

	//Shedninja check (it always has 1 hp)
	if pokemon.species == "Shedninja" {
		pokemon.base_stats[0] = 1
		pokemon.stats[0] = 1
	} else {
		pokemon.base_stats[0] = uint(((2.0*float64(pokemon.base_stats[0]) + float64(iv) + float64((ev / 4))) * float64(level) / 100.0) + float64(level) + 10.0)
		pokemon.stats[0] = pokemon.base_stats[0]
	}

	//Calculate attack
	if ev_a == -1 || ev_a > 255 {
		ev = r.Int31n(256) //Select random value (not inclusive)
	} else {
		ev = ev_a //Else select value given
	}

	switch pokemon.nature {
	case "Lonely", "Adamant", "Naughty", "Brave":
		nature_multiplier = 1.1

	case "Bold", "Modest", "Calm", "Timid":
		nature_multiplier = 0.9

	default:
		nature_multiplier = 1.0
	}

	//fmt.Println("Nature multiplier for attack:", nature_multiplier, "with natures:", pokemon.nature)
	pokemon.base_stats[1] = uint(float64(int32(float64(2*int32(pokemon.base_stats[1])+iv+(ev/4))*float64(level)/100)+5) * nature_multiplier)
	pokemon.stats[1] = pokemon.base_stats[1]

	//Calculate defense
	if ev_d == -1 || ev_d > 255 {
		ev = r.Int31n(256) //Select random value (not inclusive)
	} else {
		ev = ev_d //Else select value given
	}

	switch pokemon.nature {
	case "Bold", "Impish", "Lax", "Relaxed":
		nature_multiplier = 1.1

	case "Lonely", "Mild", "Gentle", "Hasty":
		nature_multiplier = 0.9

	default:
		nature_multiplier = 1.0
	}

	//fmt.Println("Nature multiplier for defense:", nature_multiplier, "with natures:", pokemon.nature)
	pokemon.base_stats[2] = uint(float64(int32(float64(2*int32(pokemon.base_stats[2])+iv+(ev/4))*float64(level)/100)+5) * nature_multiplier)
	pokemon.stats[2] = pokemon.base_stats[2]

	//Calculate special attack
	if ev_spa == -1 || ev_spa > 255 {
		ev = r.Int31n(256) //Select random value (not inclusive)
	} else {
		ev = ev_spa //Else select value given
	}

	switch pokemon.nature {
	case "Modest", "Mild", "Rash", "Quiet":
		nature_multiplier = 1.1

	case "Adamant", "Impish", "Careful", "Jolly":
		nature_multiplier = 0.9

	default:
		nature_multiplier = 1.0
	}

	//fmt.Println("Nature multiplier for special attack:", nature_multiplier, "with natures:", pokemon.nature)
	pokemon.base_stats[3] = uint(float64(int32(float64(2*int32(pokemon.base_stats[3])+iv+(ev/4))*float64(level)/100)+5) * nature_multiplier)
	pokemon.stats[3] = pokemon.base_stats[3]

	//Calculate special defense
	if ev_spd == -1 || ev_spd > 255 {
		ev = r.Int31n(256) //Select random value (not inclusive)
	} else {
		ev = ev_spd //Else select value given
	}

	switch pokemon.nature {
	case "Calm", "Gentle", "Careful", "Sassy":
		nature_multiplier = 1.1

	case "Naughty", "Lax", "Rash", "Naive":
		nature_multiplier = 0.9

	default:
		nature_multiplier = 1.0
	}

	//fmt.Println("Nature multiplier for special defense:", nature_multiplier, "with natures:", pokemon.nature)
	pokemon.base_stats[4] = uint(float64(int32(float64(2*int32(pokemon.base_stats[4])+iv+(ev/4))*float64(level)/100)+5) * nature_multiplier)
	pokemon.stats[4] = pokemon.base_stats[4]

	//Calculate speed
	if ev_speed == -1 || ev_speed > 255 {
		ev = r.Int31n(256) //Select random value (not inclusive)
	} else {
		ev = ev_speed //Else select value given
	}

	switch pokemon.nature {
	case "Timid", "Hasty", "Jolly", "Naive":
		nature_multiplier = 1.1

	case "Brave", "Relaxed", "Quiet", "Sassy":
		nature_multiplier = 0.9

	default:
		nature_multiplier = 1.0
	}

	//fmt.Println("Nature multiplier for speed:", nature_multiplier, "with natures:", pokemon.nature)
	pokemon.base_stats[5] = uint(float64(int32(float64(2*int32(pokemon.base_stats[5])+iv+(ev/4))*float64(level)/100)+5) * nature_multiplier)
	pokemon.stats[5] = pokemon.base_stats[5]
}

// Swaps in a pokemon and swaps another out if needed
func (currentMon *Pokemon) Swap(play *Field, player *Player, pokemon *Pokemon) {
	/*Before this checks should be made
	  to make sure the pokemon can be swapped
	  Things such as already being on the field and shadow tag*/
	//For both (Reset stats, certain status, certain abilities, text)

	if pokemon == nil {
		//No arg: start of game, on death
		if currentMon.owner { //Your pokemon
			play.Add(currentMon)
			fmt.Print("Go ")
			fmt.Print(currentMon.species)
			fmt.Print("!\n") //Make sure this has no spaces
			fmt.Println()
		} else { //Opponents pokemon
			play.Add(currentMon)
			fmt.Print("Opponent sent out ")
			fmt.Print(currentMon.species)
			fmt.Print("!\n") //Make sure this has no spaces
			fmt.Println()
		}

		switch play.format { //Give pokemon its slot
		case "Singles":
			currentMon.slot = 1
		case "Doubles":
			for _, mon := range play.mons {
				if mon.owner == currentMon.owner && mon.species != currentMon.species { //Check for pokemon on your side of the field
					if mon.slot == 1 {
						currentMon.slot = 2
					} else {
						currentMon.slot = 1
					}
				}
			}
		}

		/*Check for hazards damage*/
		if currentMon.item != "Heavy-Duty Boots" {
			if player.hazards["Spikes"] != 0 && !currentMon.HasType("flying") && currentMon.ability_read != "Levitate" {
				//Calculate spikes damage
				var damage uint8

				switch player.hazards["Spikes"] {
				case 1:
					damage = uint8(currentMon.stats[0] / 8.0)
				case 2:
					damage = uint8(currentMon.stats[0] / 6.0)
				case 3:
					damage = uint8(currentMon.stats[0] / 4.0)
				default:
					//The 0 case has already been handled
					panic("Spikes > 2, this should never happen")
				}

				if damage >= uint8(currentMon.stats[0]) { //Check if the pokemon will faint from this
					if !currentMon.owner {
						fmt.Print("The opposing ")
					}

					fmt.Println(currentMon.species, "was afflicted by spikes")
					//Damage animtaion here (maybe use a function for this and pass in "damage" variable)
					fmt.Println(currentMon.species, "fainted")
					currentMon.Faint(play)
					//At the end of the turn the ability to swap will be checked
				} else { //Pokemon wont die so express that it took the damage
					currentMon.stats[0] -= uint(damage)
					//Damage animtaion here (maybe use a function for this and pass in "damage" variable)
					fmt.Println(currentMon.species, "was afflicted by spikes")
				}
			}

			if player.hazards["Toxic Spikes"] != 0 && !currentMon.HasType("steel") && currentMon.ability_read != "Levitate" && pokemon.status == "" {
				//If the pokemon is a poison type, toxic spikes are removed
				if currentMon.HasType("poison") {
					player.hazards["Toxic Spikes"] = 0

					if !currentMon.owner {
						fmt.Print("The opposing ")
					}

					fmt.Println(currentMon.species, "removed toxic spikes from the field")
				}

				//Check for layer of Toxic Spikes
				switch player.hazards["Toxic Spikes"] {
				case 1:
					currentMon.status = "Poisoned"
					fmt.Println(currentMon.species, "has been poisoned")
				case 2:
					currentMon.status = "Badly Poisoned"
					fmt.Println(currentMon.species, "has been badly poisoned")
				default:
					//The 0 case has already been handled
					panic("Toxic Spikes > 2, this should never happen")
				}
			}

			if player.hazards["Stealth Rock"] != 0 {
				effectiveness := play.matchup.matchup_map["rock"][currentMon.type1] * play.matchup.matchup_map["rock"][currentMon.type2]
				var damage uint8

				switch effectiveness {
				//The way the damage is setup the value will always be rounded down to the nearest integer
				case .25:
					/*1/32 of max hp*/
					damage = uint8(currentMon.stats[0] / 32.0)
				case .5:
					/*1/16 of max hp*/
					damage = uint8(currentMon.stats[0] / 16.0)
				case 1.0:
					/*1/8 of max hp*/
					damage = uint8(currentMon.stats[0] / 8.0)
				case 2.0:
					/*1/4 of max hp*/
					damage = uint8(currentMon.stats[0] / 4.0)
				case 4.0:
					/*1/2 of max hp*/
					damage = uint8(currentMon.stats[0] / 2.0)
				}

				if damage >= uint8(currentMon.stats[0]) { //Check if the pokemon will faint from this
					if !currentMon.owner {
						fmt.Print("The opposing ")
					}

					fmt.Println(currentMon.species, "was afflicted by the pointed stones")
					//Damage animtaion here (maybe use a function for this and pass in "damage" variable)
					fmt.Println(currentMon.species, "fainted")
					currentMon.Faint(play)
					//At the end of the turn the ability to swap will be checked
				} else { //Pokemon wont die so express that it took the damage
					currentMon.stats[0] -= uint(damage)
					//Damage animtaion here (maybe use a function for this and pass in "damage" variable)
					fmt.Println(currentMon.species, "was afflicted by the pointed stones")
				}
			}

			if player.hazards["Sticky Web"] != 0 {
				if !currentMon.owner {
					fmt.Print("The opposing ")
				}

				fmt.Println(currentMon.species, "was caught in a sticky web")
				currentMon.StatChange([]string{"speed"}, []int8{-1}, play)
			}
		}

		//If the pokemon is a poison type, toxic spikes are removed (and this pokemon has boots)
		if currentMon.HasType("poison") && player.hazards["Toxic Spikes"] != 0 {
			player.hazards["Toxic Spikes"] = 0

			if !currentMon.owner {
				fmt.Print("The opposing ")
			}

			fmt.Println(currentMon.species, "removed toxic spikes from the field")
		}

		//Check onswap abilities
		if play.turn != 0 {
			play.CheckEvent("onswap", currentMon)
		}
	} else {
		//1 arg, during a turn swap
		currentMon.stats[1] = currentMon.base_stats[1] //Reset and boosts and debuffs
		currentMon.stats[2] = currentMon.base_stats[2]
		currentMon.stats[3] = currentMon.base_stats[3]
		currentMon.stats[4] = currentMon.base_stats[4]
		currentMon.stats[5] = currentMon.base_stats[5]

		currentMon.stat_stages[0] = 0 //Attack
		currentMon.stat_stages[1] = 0 //Defense
		currentMon.stat_stages[2] = 0 //Special Attack
		currentMon.stat_stages[3] = 0 //Special defense
		currentMon.stat_stages[4] = 0 //Speed
		currentMon.stat_stages[5] = 0 //Accuracy
		currentMon.stat_stages[6] = 0 //Defense

		currentMon.boosts = 0                        //May be removed
		currentMon.ability_read = currentMon.ability //Reset ability when swapped
		currentMon.temp_status = nil                 //Remove temp statuses

		pokemon.slot = currentMon.slot //Exchange slots
		currentMon.slot = 0
		play.Delete(currentMon) //Remove current one
		play.Add(pokemon)       //Send out mon

		if pokemon.owner { //Your pokemon
			fmt.Print("Come back ") //Come back text
			fmt.Print(currentMon.species)
			fmt.Print("!\n") //Make sure this has no spaces

			fmt.Print("Go ")
			fmt.Print(pokemon.species)
			fmt.Print("!\n")
			fmt.Println()
		} else { //Opponent text
			fmt.Print("Opponent withdrew ") //Come back text
			fmt.Print(currentMon.species)
			fmt.Print("!\n") //Make sure this has no spaces

			fmt.Print("Opponent sent out ")
			fmt.Print(pokemon.species)
			fmt.Print("!\n")
			fmt.Println()
		}

		/*Check for hazards damage*/
		if pokemon.item != "Heavy-Duty Boots" {
			if player.hazards["Spikes"] != 0 && !pokemon.HasType("flying") && pokemon.ability_read != "Levitate" {
				//Calculate spikes damage
				var damage uint8

				switch player.hazards["Spikes"] {
				case 1:
					damage = uint8(pokemon.stats[0] / 8.0)
				case 2:
					damage = uint8(pokemon.stats[0] / 6.0)
				case 3:
					damage = uint8(pokemon.stats[0] / 4.0)
				default:
					//The 0 case has already been handled
					panic("Spikes > 2, this should never happen")
				}

				if damage >= uint8(pokemon.stats[0]) { //Check if the pokemon will faint from this
					if !pokemon.owner {
						fmt.Print("The opposing ")
					}

					fmt.Println(pokemon.species, "was afflicted by spikes")
					//Damage animtaion here (maybe use a function for this and pass in "damage" variable)
					fmt.Println(pokemon.species, "fainted")
					pokemon.Faint(play)
					//At the end of the turn the ability to swap will be checked
				} else { //Pokemon wont die so express that it took the damage
					pokemon.stats[0] -= uint(damage)
					//Damage animtaion here (maybe use a function for this and pass in "damage" variable)
					fmt.Println(pokemon.species, "was afflicted by spikes")
				}
			}

			if player.hazards["Toxic Spikes"] != 0 && !pokemon.HasType("steel") && pokemon.ability_read != "Levitate" && pokemon.status == "" {
				//If the pokemon is a poison type, toxic spikes are removed
				if pokemon.HasType("poison") {
					player.hazards["Toxic Spikes"] = 0

					if !pokemon.owner {
						fmt.Print("The opposing ")
					}

					fmt.Println(pokemon.species, "removed toxic spikes from the field")
				}

				//Check for layer of Toxic Spikes
				switch player.hazards["Toxic Spikes"] {
				case 1:
					pokemon.status = "Poisoned"
					fmt.Println(pokemon.species, "has been poisoned")
				case 2:
					pokemon.status = "Badly Poisoned"
					fmt.Println(pokemon.species, "has been badly poisoned")
				default:
					//The 0 case has already been handled
					panic("Toxic Spikes > 2, this should never happen")
				}
			}

			if player.hazards["Stealth Rock"] != 0 {
				effectiveness := play.matchup.matchup_map["rock"][pokemon.type1] * play.matchup.matchup_map["rock"][pokemon.type2]
				var damage uint8

				switch effectiveness {
				//The way the damage is setup the value will always be rounded down to the nearest integer
				case .25:
					/*1/32 of max hp*/
					damage = uint8(pokemon.stats[0] / 32.0)
				case .5:
					/*1/16 of max hp*/
					damage = uint8(pokemon.stats[0] / 16.0)
				case 1.0:
					/*1/8 of max hp*/
					damage = uint8(pokemon.stats[0] / 8.0)
				case 2.0:
					/*1/4 of max hp*/
					damage = uint8(pokemon.stats[0] / 4.0)
				case 4.0:
					/*1/2 of max hp*/
					damage = uint8(pokemon.stats[0] / 2.0)
				}

				if damage >= uint8(pokemon.stats[0]) { //Check if the pokemon will faint from this
					if !pokemon.owner {
						fmt.Print("The opposing ")
					}

					fmt.Println(pokemon.species, "was afflicted by the pointed stones")
					//Damage animtaion here (maybe use a function for this and pass in "damage" variable)
					fmt.Println(pokemon.species, "fainted")
					pokemon.Faint(play)
					//At the end of the turn the ability to swap will be checked
				} else { //Pokemon wont die so express that it took the damage
					pokemon.stats[0] -= uint(damage)
					//Damage animtaion here (maybe use a function for this and pass in "damage" variable)
					fmt.Println(pokemon.species, "was afflicted by the pointed stones")
				}
			}

			if player.hazards["Sticky Web"] != 0 {
				if !pokemon.owner {
					fmt.Print("The opposing ")
				}

				fmt.Println(pokemon.species, "was caught in a sticky web")
				pokemon.StatChange([]string{"speed"}, []int8{-1}, play)
			}
		}

		//If the pokemon is a poison type, toxic spikes are removed (and this pokemon has boots)
		if pokemon.HasType("poison") && player.hazards["Toxic Spikes"] != 0 {
			player.hazards["Toxic Spikes"] = 0

			if !pokemon.owner {
				fmt.Print("The opposing ")
			}

			fmt.Println(pokemon.species, "removed toxic spikes from the field")
		}

		//Check onswap abilities
		if play.turn != 0 {
			play.CheckEvent("onswap", currentMon)
		}
	}
}

// Returns a pokemon's base stats
func (pokemon Pokemon) GetStats() {
	fmt.Println("Stats for", pokemon.species)
	fmt.Println("Hp:", pokemon.base_stats[0])
	fmt.Println("Attack:", pokemon.base_stats[1])
	fmt.Println("Defense:", pokemon.base_stats[2])
	fmt.Println("Special attack:", pokemon.base_stats[3])
	fmt.Println("Special defense:", pokemon.base_stats[4])
	fmt.Println("Speed:", pokemon.base_stats[5])
	fmt.Println()
}

// Activates a pokemon's ability
func (pokemon *Pokemon) ActivateAbility(play *Field, opponent ...*Pokemon) { //Opponent may be ommitted as all abilities dont have to do with iterations with others
	var yourStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("87"))
	var opponentStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("197"))

	switch pokemon.ability_read { //Long ass switch statement for every ability (do I program all abilities here?)
	case "Intimidate":
		var effect []*Pokemon
		//Get pokemon on the field currently on the opponents side
		for _, mon := range play.mons {
			if mon.owner != pokemon.owner { //Gets opposing pokemon no matter what side the pokemon is on
				effect = append(effect, mon)
			}
		}

		pokemon.PokemonSStyle() //Ability activation text
		if pokemon.owner {
			fmt.Println(yourStyle.Render("Intimidate"))
		} else {
			fmt.Println(opponentStyle.Render("Intimidate"))
		}
		//Maybe do something fancy with ability activations with lip gloss
		for _, mon := range effect {
			mon.StatChange([]string{"attack"}, []int8{-1}, play)
		}
		fmt.Println()
	}
}

// What happens when the attacker hits the defender
func (attacker *Pokemon) DamageCalc(defender *Pokemon, play *Field, r *rand.Rand) {
	//Get move that its using
	index, _ := strconv.ParseInt(attacker.action[0], 10, 32)
	move := attacker.moveset[index]

	//Check abilties for this (come back to the effects of this as you add abilities)
	play.CheckEvent("ondamagerecieve", defender)

	//Check if the move missed
	if move.accuracy < 101 { //The move has a chance to miss and must be calculated
		//Get adjusted stages = your acc stages - their evasion stages
		adjustedStages := attacker.stat_stages[5] - defender.stat_stages[6]

		//Adjustment if needed
		if adjustedStages > 6 {
			adjustedStages = 6
		} else if adjustedStages < -6 {
			adjustedStages = -6
		}

		accuracyMap := map[int8][]float64{
			6:  {9.0, 3.0},
			5:  {8.0, 3.0},
			4:  {7.0, 3.0},
			3:  {6.0, 3.0},
			2:  {5.0, 3.0},
			1:  {4.0, 3.0},
			0:  {3.0, 3.0},
			-1: {3.0, 4.0},
			-2: {3.0, 5.0},
			-3: {3.0, 6.0},
			-4: {3.0, 7.0},
			-5: {3.0, 8.0},
			-6: {3.0, 9.0},
		}

		//Get modifiers
		var modifier = 4096.0

		//Check for gravity
		if play.gravity > 0 {
			modifier *= 6840.0 / 4096.0
		}

		//Check for defender abilities
		if defender.ability_read == "Tangled Feet" && slices.Contains(defender.temp_status, "Confused") { //Mold breaker matters here
			modifier *= 0.5
		} else if defender.ability_read == "Sand Veil" && play.weather["sandstorm"] > 0 { //Mold breaker matters here + air lock/cloud nine
			modifier *= 3277.0 / 4096.0
		} else if defender.ability_read == "Snow Cloak" && play.weather["snow"] > 0 { //Mold breaker matters here + air lock/cloud nine
			modifier *= 3277.0 / 4096.0
		}

		//Check for attacker abilities
		if attacker.ability_read == "Hustle" && move.category == "physical" { //Category types are lower case
			modifier *= 3277.0 / 4096.0
		} else if attacker.ability_read == "Victory Star" {
			modifier *= 6840.0 / 4096.0

			if play.format == "Doubles" { //Check for other mons that may also have this ability
				for _, mon := range play.mons {
					if mon.owner == attacker.owner && mon.UI != attacker.UI {
						modifier *= 6840.0 / 4096.0
					}
				}
			}

		} else if attacker.ability_read == "Compound Eyes" {
			modifier *= 5325.0 / 4096.0
		}

		//Check for bright powder or lax incense item
		if defender.item == "Bright Powder" || defender.item == "Lax Incense" {
			modifier *= 3686.0 / 4096.0
		}
		//Check for wide lens or zoom lens
		if attacker.item == "Wide Lens" {
			modifier *= 4505.0 / 4096.0
		}
		//Check for zoom lens and has the defender moved yet this turn
		if attacker.item == "Zoom Lens" && defender.has_acted {
			modifier *= 4915.0 / 4096.0
		}

		modifier = math.Round(modifier / 4096)
		//Modified accuracy = moveaccuracy * modifier * adjustedstages (* micle berry)
		accuracy := int32(float64(move.accuracy) * modifier * (accuracyMap[adjustedStages][0] * accuracyMap[adjustedStages][1])) //* micle berry if added

		//Check if the move hits
		if random := r.Int31n(101); random > accuracy {
			//Target avoided the attack
			if !defender.owner {
				fmt.Print("The opposing")
			}

			fmt.Println(defender.species, "avoided the attack")
			return //The move missed so you don't need to deal with anything else
		}
	}
	//The move hit so figure out what happened (All self used moves should be here as they do not miss)

	//Check move effectiveness
	/*THIS CONTAINS A BUNCH OF EDGE CASES*/
	var effectiveness float64
	if move.target != "Self" {
		//Lock map mutex
		play.matchup.locker.RLock()
		effectiveness = play.matchup.matchup_map[move.type1][defender.type1] * play.matchup.matchup_map[move.type1][defender.type2]
		//Unlock map mutex
		play.matchup.locker.RUnlock()
		if effectiveness == 0.0 {
			fmt.Print("But it had no effect on ")

			if !defender.owner {
				fmt.Print(" the opposing ")
			}

			fmt.Println(defender.species)
			return //The move did nothing so return
		}
	}

	//If status move damage calc can be ignored
	if move.category != "status" {
		//Spread damage
		var targets float64
		if move.target == "All Adjacent Foes" || move.target == "All Adjacent" {
			targets = 0.75
		} else {
			targets = 1.0
		}

		//Weather multiplier
		var weather float64
		if play.weather["rain"] > 0 { //Kyogre not covered
			if move.type1 == "water" {
				weather = 1.5
			} else if move.type1 == "fire" {
				weather = 0.5
			}
		} else if play.weather["harsh sunlight"] > 0 { //Groudon not covered
			if move.type1 == "water" {
				weather = 0.5
			} else if move.type1 == "fire" {
				weather = 1.5
			}
		} else {
			weather = 1.0
		}

		//Glaive rush multiplier
		var glaiveRush float64
		if slices.Contains(defender.temp_status, "Glaive Rush") {
			glaiveRush = 2.0
		} else {
			glaiveRush = 1.0
		}

		//Crit multiplier
		var critical float64
		if defender.ability_read == "Battle Armor" || defender.ability_read == "Shell Armor" || slices.Contains(defender.temp_status, "Lucky Chant") {
			critical = 1.0
		} else { //Check for critical hit chance
			if random := r.Float64() * 100; random < move.crit {
				critical = 1.5 //A critical hit
			} else {
				critical = 1.0
			}
		}

		//Damage roll
		roll := (float64(r.Int31n(16)) + 85.0) / 100.0 //Random is now a random number between 85 and 100 inclusive

		//Stab
		var stab float64
		if attacker.HasType(move.type1) {
			if attacker.ability_read == "Adaptability" {
				stab = 2.0
			} else {
				stab = 1.5
			}
		} else {
			stab = 1.0
		}

		play.CheckEvent("ondamagedeal", attacker) //Should this be here???

		var damage uint
		switch move.category {
		case "physical":
			var burn float64
			if attacker.status == "Burned" {
				burn = 0.5
			} else {
				burn = 1.0
			}

			damage = uint(math.Floor(((2.0*50.0/5.0+2.0)*float64(move.base_power)*(float64(attacker.stats[1])/float64(defender.stats[2]))/50.0 + 2.0) * targets * weather * glaiveRush * critical * roll * stab * effectiveness * burn)) //* other (this is an important edge case)
		case "special":

			damage = uint(math.Floor((2.0*50.0/5.0+2.0)*float64(move.base_power)*(float64(attacker.stats[3])/float64(defender.stats[4]))/50.0+2.0) * targets * weather * glaiveRush * critical * roll * stab * effectiveness) //* other (this is an important edge case)
		}

		//We have the damage the move did so now apply it

		if effectiveness >= 2.0 { //Check for effectiveness
			fmt.Println("It's super effective!")
		} else if effectiveness == 0.5 {
			fmt.Println("It's not very effective!")
		}

		//Check for a critical hit
		if critical == 1.5 {
			fmt.Println("A critical hit!")
		}

		//Apply damage to pokemon
		//Damage effect animation or something here
		if damage >= defender.stats[0] {
			//Pokemon will faint
			if !defender.owner {
				fmt.Print("The opposing ")
			}

			fmt.Print(defender.species, " fainted\n")
			defender.Faint(play)
			return
		} else {
			if !defender.owner {
				fmt.Print("The opposing ")
			}
			fmt.Print(defender.species, " lost ", math.Floor((float64(damage)/float64(defender.base_stats[0]))*100)) //How much % hp did it lose
			fmt.Println("%", "of its hp")

			defender.stats[0] -= damage //The pokemon will take damage
			fmt.Print(defender.species, " has ", math.Floor((float64(defender.stats[0])/float64(defender.base_stats[0]))*100), " % hp remaining\n")
			fmt.Println()
		}
		//Calculate secondary effects
	} else {
		//If you got here your move is a status move and doesnt require damage calculation
		fmt.Println("This was a status move and nothing happened (yet)")
	}
}

// Boosts or lowers a stat given the stat in question, and the amount to be adjusted by
func (pokemon *Pokemon) StatChange(statsList []string, value []int8, play *Field) {
	statMap := map[string]int{
		"attack":          0,
		"defense":         1,
		"special attack":  2,
		"special defense": 3,
		"speed":           4,
		"accuracy":        5,
		"evasion":         6,
	}

	for index, stat := range statsList { //Goes over each stat to be changed
		//Check if ability or item is triggered
		if value[index] < 0 {
			if play.CheckEvent("onstatlower", pokemon) {
				continue
			}
		}
		//First check if the stat can go any higher or lower
		if pokemon.stat_stages[statMap[stat]] == 6 && value[index] > 0 { //Stat at its highest
			pokemon.PokemonS() //Handle grammar

			fmt.Print(value, "won't go any higher!")
			continue //Check next stat
		} else if pokemon.stat_stages[statMap[stat]] == -6 && value[index] < 0 { //Stat at its lowest
			pokemon.PokemonS() //Handle grammar

			fmt.Print(value, "won't go any lower!")
			continue //Check next stat
		} else { //Continue with calculating stat changes
			//See if stat change makes it reach its cap (remember this is still in a for loop so this will be done for each stat that needs to change)
			var change int8                                          //How many stages the stat will be adjusted by
			if pokemon.stat_stages[statMap[stat]]+value[index] > 6 { //Stat boost will raise the stat past + 6
				change = int8(math.Abs(float64(pokemon.stat_stages[statMap[stat]] - 6)))
			} else if pokemon.stat_stages[statMap[stat]]+value[index] < -6 { //Stat lower will lower the stat past - 6
				change = -1 * (pokemon.stat_stages[statMap[stat]] + 6)
			} else {
				//Stat change will be normal and not go out of bounds
				change = value[index]
			}

			//Calculate stat change
			switch stat { //Find out accuracy and evasion formula
			case "accuracy": /*
					multiplierMap := map[int][]float64{
						6:  {9.0, 3.0},
						5:  {8.0, 3.0},
						4:  {7.0, 3.0},
						3:  {6.0, 3.0},
						2:  {5.0, 3.0},
						1:  {4.0, 3.0},
						0:  {3.0, 3.0},
						-1: {3.0, 4.0},
						-2: {3.0, 5.0},
						-3: {3.0, 6.0},
						-4: {3.0, 7.0},
						-5: {3.0, 8.0},
						-6: {3.0, 9.0},
					}*/
				pokemon.stat_stages[5] += change
			case "evasion": /*
					multiplierMap := map[int][]float64{
						6:  {3.0, 9.0},
						5:  {3.0, 8.0},
						4:  {3.0, 7.0},
						3:  {3.0, 6.0},
						2:  {3.0, 5.0},
						1:  {3.0, 4.0},
						0:  {3.0, 3.0},
						-1: {4.0, 3.0},
						-2: {5.0, 3.0},
						-3: {6.0, 3.0},
						-4: {7.0, 3.0},
						-5: {8.0, 3.0},
						-6: {9.0, 3.0},
					}*/
				pokemon.stat_stages[6] += change
			default:
				//Every other stat
				multiplierMap := map[int][]float64{
					6:  {8.0, 2.0},
					5:  {7.0, 2.0},
					4:  {6.0, 2.0},
					3:  {5.0, 2.0},
					2:  {4.0, 2.0},
					1:  {3.0, 2.0},
					0:  {2.0, 2.0},
					-1: {2.0, 3.0},
					-2: {2.0, 4.0},
					-3: {2.0, 5.0},
					-4: {2.0, 6.0},
					-5: {2.0, 7.0},
					-6: {2.0, 8.0},
				}

				//Calculation done here
				pokemon.stats[1+statMap[stat]] = uint(float64(pokemon.base_stats[1+statMap[stat]]) * multiplierMap[int(change)][0] / multiplierMap[int(change)][1])
			}

			//Print message of stat change
			switch change {
			case 6:
				//Should only be for belly drum
				/*NEW SOUND TEST*/
				PlaySound("../sounds/statrise.mp3")
				pokemon.PokemonS()
				fmt.Println(stat, "has been maxed out!")
			case 3, 4, 5:
				/*NEW SOUND TEST*/
				PlaySound("../sounds/statrise.mp3")
				pokemon.PokemonS()
				fmt.Println(stat, "rose drastically!")
			case 2:
				/*NEW SOUND TEST*/
				PlaySound("../sounds/statrise.mp3")
				pokemon.PokemonS()
				fmt.Println(stat, "rose sharply!")
			case 1:
				/*NEW SOUND TEST*/
				PlaySound("../sounds/statrise.mp3")
				pokemon.PokemonS()
				fmt.Println(stat, "rose!")
			case -1:
				/*NEW SOUND TEST*/
				PlaySound("../sounds/statfall.mp3")
				pokemon.PokemonS()
				fmt.Println(stat, "fell!")
			case -2:
				/*NEW SOUND TEST*/
				PlaySound("../sounds/statfall.mp3")
				pokemon.PokemonS()
				fmt.Println(stat, "harshly!")
			case -3, -4, -5, -6:
				/*NEW SOUND TEST*/
				PlaySound("../sounds/statfall.mp3")
				pokemon.PokemonS()
				fmt.Println(stat, "severely fell!")
			default:
				panic("Stat change value is invalid: out of bounds -6 - +6")
			}

		}
	}
}

// Checks if a pokemon has a certain type
func (pokemon *Pokemon) HasType(monType string) bool {
	if pokemon.type1 == monType || pokemon.type2 == monType {
		return true
	}
	return false
}

// Check if pokemon is grounded via iron ball, gravity, magnet rise, telekenesis
func (pokemon *Pokemon) IsGrounded() bool {
	for _, value := range pokemon.temp_status {
		if value == "Grounded" {
			return true
		}
	}

	if pokemon.item == "Iron Ball" {
		return true
	}

	if pokemon.item == "Air Balloon" { //May be removed
		return false
	}

	return false
}

// Handles pokemon names that end with 's' ex: Solosis' vs Incineroar's (with lip gloss)
func (pokemon *Pokemon) PokemonSStyle() {
	var yourStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("87"))
	var opponentStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("197"))

	if !pokemon.owner { //Lip gloss pokemon that isnt yours
		fmt.Print(opponentStyle.Render("The opposing "))

		fmt.Print(opponentStyle.Render(pokemon.species))

		if string(pokemon.species[len(pokemon.species)-1]) == "s" {
			fmt.Print(opponentStyle.Render("'")) //Name ends with an s so adjust accordingly
		} else {
			fmt.Print(opponentStyle.Render("'s"))
		}
	} else { //Lip gloss pokemon that is yours
		fmt.Print(yourStyle.Render(pokemon.species))

		if string(pokemon.species[len(pokemon.species)-1]) == "s" {
			fmt.Print(yourStyle.Render("'")) //Name ends with an s so adjust accordingly
		} else {
			fmt.Print(yourStyle.Render("'s"))
		}
	}
	fmt.Print(" ") //Space added here for both
}

// Handles pokemon names that end with 's' ex: Solosis' vs Incineroar's
func (pokemon *Pokemon) PokemonS() {
	if !pokemon.owner {
		fmt.Print("The opposing ")
	}
	fmt.Print(pokemon.species)

	if string(pokemon.species[len(pokemon.species)-1]) == "s" {
		fmt.Print("'") //Name ends with an s so adjust accordingly
	} else {
		fmt.Print("'s")
	}
	fmt.Print(" ") //Space added here for both
}

// Make sure pokemon faints
func (pokemon *Pokemon) Faint(play *Field) {
	pokemon.stat_stages[0] = 0
	pokemon.stat_stages[1] = 0
	pokemon.stat_stages[2] = 0
	pokemon.stat_stages[3] = 0
	pokemon.stat_stages[4] = 0
	pokemon.stat_stages[5] = 0
	pokemon.stat_stages[6] = 0
	pokemon.boosts = 0

	pokemon.stats[0] = 0                   //HP
	pokemon.ability_read = pokemon.ability //Reset ability
	pokemon.priority = 0                   //Reset priority
	pokemon.temp_status = nil              //May cause errors somewhere
	pokemon.status_turns = 0
	pokemon.on_field = false
	pokemon.has_acted = false
	pokemon.status = "Fainted"
	play.Delete(pokemon)

	select {
	case fainted <- pokemon.UI:
		//fmt.Println("Pokemon channel", pokemon.UI, "sent to fainted channel")
	default:
		//log.Fatal("Unable to send to fatal channel")
	}

	//Remove 1 from counter
	if pokemon.owner {
		play.yourAliveMons--
		if play.yourAliveMons == 0 { //Check for winner
			select {
			case winner <- "Opponent":
				 //If you have no alive mons the opponent wins
			default:
			}
			
			close(winner)
		}
	} else {
		play.cpuAliveMons--
		if play.cpuAliveMons == 0 { //Check for winner
			select {
			case winner <- "You":
				//If the opponent has no alive mons you win
			default:
			}

			close(winner)
		}
	}
}

package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"log"
	"math/rand"
	"runtime/debug"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/huh"
	//tea "github.com/charmbracelet/bubbletea"

	_ "github.com/mattn/go-sqlite3"
)

// Concurrent reads or writes of maps in Go may have unintended behavior
type muMatchup struct {
	locker      sync.RWMutex                  //Mutex for json map
	matchup_map map[string]map[string]float64 //Variable to access matchup.json
}

var winner chan string  //Contains "You" or "Opponent"
var fainted chan string //Contains pokemon UI

func main() {
	debug.SetGCPercent(120) //More memory, less overall cpu cost

	winner = make(chan string, 1)
	fainted = make(chan string, 4) //Mo more than 4 pokemon fainted in a single turn

	you := Player{}
	opponent := Player{}
	//Setup entry hazards maps
	you.hazards = map[string]uint{
		"Spikes":       0,
		"Toxic Spikes": 0,
		"Stealth Rock": 0,
		"Sticky Web":   0,
		//"Steelsurge": 0,
	}
	opponent.hazards = map[string]uint{
		"Spikes":       0,
		"Toxic Spikes": 0,
		"Stealth Rock": 0,
		"Sticky Web":   0,
		//"Steelsurge": 0,
	}

	play := Field{}
	play.matchup = muMatchup{}           //Matchup table with read/write mutex
	play.matchup.matchup_map = GetJson() //Create json matchup table
	play.abilities = GetAbilities()
	var wg sync.WaitGroup

	db, err := sql.Open("sqlite3", "../pokedata.db")
	if err != nil {
		panic(err)
	}
	defer db.Close() //Queries may be needed later but might be removed

	wg.Add(1)
	go generate_mons(&opponent, db, &wg) //Generate cpu team

	fmt.Println("Cpu team: ")
	wg.Wait()
	for _, mon := range opponent.pokemon_list {
		fmt.Print("    ")
		fmt.Println(mon.species)
		for _, move := range mon.moveset {
			fmt.Print("        ")
			fmt.Println(move.name)
		}
	}

	iteration := map[int]string{
		0: "first",
		1: "second",
		2: "third",
		3: "fourth",
		4: "fifth",
		5: "last",
	}

	play.format = getFormat()

	for i := 0; i < 1; i++ { //Adjust second number to pick how many pokemon you want to add
		create_team(&you, db, iteration[i]) //Create team and assign format
	}
	play.yourAliveMons = len(you.pokemon_list) //Add fainted count
	play.cpuAliveMons = len(opponent.pokemon_list)

	battle_start(play.format, &you, &play) //Select pokemon to start the battle

	switch play.format {
	//Singles
	case "Singles":
		opponent.pokemon_list[0].Swap(&play, &opponent, nil)
	//Doubles
	case "Doubles":
		opponent.pokemon_list[0].Swap(&play, &opponent, nil)
		opponent.pokemon_list[1].Swap(&play, &opponent, nil)
	}

	//Battle logic actually starts
	//Check for on swap abilities with speed order (priority moves dont apply yet)
	play.turn = 0
	play.GetQueue() //Get queue to check for abilities
	play.CheckEvent("onswap", nil)

	//Battle loop
	var turnStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("141"))
OuterLoop:
	for { //Checks if one side has all fainted pokemon (check this on every pokemon death???)
		play.turn++
		fmt.Println()
		fmt.Println(turnStyle.Render("Turn", strconv.Itoa(play.turn)))//Lip gloss new turn text
		wg.Add(1)
		go PickCpuMoves(&opponent, &play, &wg)

		PickYourMoves(&you, &play)
		wg.Wait()

		//Get action queue
		play.GetQueue()
		//Run the turn
		play.Run(&you, &opponent)

		play.CheckEvent("onturnend", nil) //Check for abilities

		play.Reset() //Reset values that may have been changed mid-turn

		select { //Check for winner
		case play := <-winner:
			switch play { //Check values of play
			case "You":
				fmt.Println(turnStyle.Render(play, "win!"))
			default:
				fmt.Println(turnStyle.Render("The opponent wins!"))
			}
			break OuterLoop
		default:

		}

		CheckFainted(&play, &you, &opponent)

		play.queue = nil
	}
	close(fainted)

	//Wait so user can view results
	fmt.Println(turnStyle.Render("Program will close in 10 seconds..."))
	time.Sleep(10 * time.Second)
}

// Generate pokemon for the cpu opponent based on random values
func generate_mons(opponent *Player, db *sql.DB, wg *sync.WaitGroup) {
	removedPokemon := []int{19, 20, 26, 27, 28, 37, 38, 50, 51, 52, 53, 74, 75, 76, 88, 89, 105}

	natures := [25]string{"Adamant", "Bashful", "Bold", "Brave", "Calm",
		"Careful", "Docile", "Gentle", "Hardy", "Hasty",
		"Impish", "Jolly", "Lax", "Lonely", "Mild",
		"Modest", "Naive", "Naughty", "Quiet", "Quirky",
		"Rash", "Relaxed", "Sassy", "Serious", "Timid"}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	//Get all moves to be randomly selected
	var move Moves
	var moveList []Moves
	moveRows, err := db.Query("SELECT * FROM moves") //'select count(*) from moves' to check how many moves (value should be the amount of moves + 1)
	if err != nil {
		panic(err)
	}
	defer moveRows.Close()

	for moveRows.Next() {
		// Use the Scan method to copy values from the row into the struct fields
		err := moveRows.Scan(&move.name, &move.type1, &move.base_power, &move.accuracy, &move.category, &move.priority, &move.pp, &move.crit, &move.description, &move.effect,
			&move.effect_rate, &move.target, &move.contact, &move.sound, &move.punch, &move.bite, &move.can_snatch,
			&move.slice, &move.bullet, &move.wind, &move.powder, &move.can_metro, &move.gravity, &move.defrosts, &move.bounceable, &move.can_block, &move.can_mirror)
		if err != nil {
			panic(err)
		}

		moveList = append(moveList, move)
	}

	//Loop 6 times for each pokemon in the party
	for x := 0; x < 6; x++ {
	Party:
		//Query for mon by number index
		var pokemon Pokemon

		random := r.Intn(800) + 1
		if slices.Contains(removedPokemon, random) {
			//That pokemon has been removed so redo query
			goto Party
		}

		pokemonRows, err := db.Query("SELECT * FROM pokemon WHERE pokedex_number = ?", random) //801 pokemon (+1 because you don't want 0 to be a value)
		if err != nil {
			panic(err)
		}

		defer pokemonRows.Close()

		if pokemonRows.Next() {
			// Use the Scan method to copy values from the row into the struct fields
			err := pokemonRows.Scan(&pokemon.pokedex_number, &pokemon.species, &pokemon.abilities, &pokemon.base_stats[1], &pokemon.base_stats[2], &pokemon.base_stats[0], &pokemon.genderless, &pokemon.base_stats[3], &pokemon.base_stats[4], &pokemon.base_stats[5], &pokemon.type1, &pokemon.type2, &pokemon.weight, &pokemon.can_dynamax, &pokemon.has_mega)
			if err != nil {
				panic(err)
			}
		}

		for _, mon := range opponent.pokemon_list {
			if pokemon.species == mon.species {
				//This pokemon already exists in this list so try again
				goto Party
			}
		}

		for y := 0; y < 4; y++ {
			//Randomly pick move and assign it to slot one (Currently 12 moves) but first check if it exists in the list
		Loop:
			for {
				newMove := moveList[r.Int31n(11)] //Randomly pick move, 'select count(*) from moves' to check how many moves (value should be the amount of moves + 1)
				for _, value := range pokemon.moveset {
					if newMove.name == value.name {
						goto Loop //Move is already in list
					}
				}
				//Move isnt in list so add it to be checked later
				pokemon.moveset[y] = newMove
				break
			}
		}

		//Pick first ability
		abilityList := strings.Split(pokemon.abilities, ", ")
		pokemon.ability = abilityList[0]
		pokemon.ability_read = abilityList[0]
		//Pick random nature
		pokemon.nature = natures[r.Intn(25)]

		//Get stats ready
		pokemon.CalcStats(-1, -1, -1, -1, -1, -1)

		//Add ui (unique identifier)
		pokemon.UI = "o" + strconv.Itoa(x)

		//Random gender
		if !pokemon.genderless {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			if random := r.Int31n(2); random == 1 {
				pokemon.gender = "Male"
			} else {
				pokemon.gender = "Female"
			}
		} else { //Pokemon has no gender
			pokemon.gender = "None"
		}

		//Set owner
		pokemon.owner = false

		//Add pokemon to team
		opponent.pokemon_list = append(opponent.pokemon_list, &pokemon)
	}
	wg.Done() //Last line of function
}

// Add pokemon to be used on your team
func create_team(you *Player, db *sql.DB, iteration string) {
	var (
		pokeName        string
		pokeNature      string
		pokeAbility     string
		pokeAbilityList []string
		pokemon         Pokemon
		ev_hp           string
		ev_a            string
		ev_d            string
		ev_spa          string
		ev_spd          string
		ev_speed        string
	)
	allMoveNames := GetMoveNames(db)
	natures := []string{"Adamant", "Bashful", "Bold", "Brave", "Calm",
		"Careful", "Docile", "Gentle", "Hardy", "Hasty",
		"Impish", "Jolly", "Lax", "Lonely", "Mild",
		"Modest", "Naive", "Naughty", "Quiet", "Quirky",
		"Rash", "Relaxed", "Sassy", "Serious", "Timid"}
	var moveset [4]string

	//Pokemon1
	pokeForm := huh.NewForm(
		//Name
		huh.NewGroup(
			huh.NewInput().
				Title("Enter the name of the "+iteration+" pokemon you'd like to add to your team").
				Value(&pokeName).
				// Validating fields is easy. The form will mark erroneous fields
				// and display error messages accordingly.
				Validate(func(str string) error {
					if str == " " {
						return errors.New("Sorry, that value is invalid")
					}
					//Check for duplicates (may be removed)
					if iteration != "first" {
						for _, value := range you.pokemon_list {
							if GetBubble(str) == value.species {
								//This pokemon already exists
								return errors.New("You are only allowed 1 of each species of pokemon per team")
							}
						}
					}
					pokemon, err := GetPokemon(GetBubble(str), db) //New variable here for scope reasons
					if err != nil {
						return err
					}
					if pokemon.species == "" {
						//Struct is empty so that value doesn't exist
						return errors.New("That pokemon species doesn't exist in our database")
					}
					return nil
				}),
		),
		//Moves
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select first move").
				Options(
					huh.NewOptions(allMoveNames...)...).
				Value(&moveset[0]),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select second move").
				Options(
					huh.NewOptions(allMoveNames...)...).
				Value(&moveset[1]).
				Validate(func(str string) error {
					//Prevent move duplicates
					if str == moveset[0] {
						return errors.New("A pokemon can't learn the same move twice!")
					}
					return nil
				}),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select third move").
				Options(
					huh.NewOptions(allMoveNames...)...).
				Value(&moveset[2]).
				Validate(func(str string) error {
					//Prevent move duplicates
					if str == moveset[1] || str == moveset[0] {
						return errors.New("A pokemon can't learn the same move twice!")
					}
					return nil
				}),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select last move").
				Options(
					huh.NewOptions(allMoveNames...)...).
				Value(&moveset[3]).
				Validate(func(str string) error {
					//Prevent move duplicates
					if str == moveset[2] || str == moveset[1] || str == moveset[0] {
						return errors.New("A pokemon can't learn the same move twice!")
					}
					return nil
				}),
		),
	)
	err := pokeForm.Run()
	if err != nil {
		panic(err)
	}
	pokemon, err = GetPokemon(GetBubble(pokeName), db) //ALL VALUES HERE ARE SET CORRECTLY
	if err != nil {
		log.Fatal(err) //Should not fail here
	}
	pokeAbilityList = strings.Split(pokemon.abilities, ", ")

	pokeForm2 := huh.NewForm(
		//Nature
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select nature").
				Options(
					huh.NewOptions(natures...)...).
				Value(&pokeNature),
		),

		//Abilities
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select pokemon ability").
				Options(
					huh.NewOptions(pokeAbilityList...)...).
				Value(&pokeAbility),
		),

		//EVs
		huh.NewGroup(
			huh.NewInput().
				Title("Values out of this range will be randomized\nEnter HP Evs: 0-255").
				Value(&ev_hp).
				// Validating fields is easy. The form will mark erroneous fields
				// and display error messages accordingly.
				Validate(func(str string) error {
					if _, err := strconv.ParseInt(str, 10, 32); err != nil {
						return errors.New("That value is invalid, please be sure to enter a number")
					}
					return nil
				}),

			huh.NewInput().
				Title("Enter Attack Evs: 0-255").
				Value(&ev_a).
				// Validating fields is easy. The form will mark erroneous fields
				// and display error messages accordingly.
				Validate(func(str string) error {
					if _, err := strconv.ParseInt(str, 10, 32); err != nil {
						return errors.New("That value is invalid, please be sure to enter a number")
					}
					return nil
				}),

			huh.NewInput().
				Title("Enter Defense Evs: 0-255").
				Value(&ev_d).
				// Validating fields is easy. The form will mark erroneous fields
				// and display error messages accordingly.
				Validate(func(str string) error {
					if _, err := strconv.ParseInt(str, 10, 32); err != nil {
						return errors.New("That value is invalid, please be sure to enter a number")
					}
					return nil
				}),

			huh.NewInput().
				Title("Enter Special Attack Evs: 0-255").
				Value(&ev_spa).
				// Validating fields is easy. The form will mark erroneous fields
				// and display error messages accordingly.
				Validate(func(str string) error {
					if _, err := strconv.ParseInt(str, 10, 32); err != nil {
						return errors.New("That value is invalid, please be sure to enter a number")
					}
					return nil
				}),

			huh.NewInput().
				Title("Enter Special Defense Evs: 0-255").
				Value(&ev_spd).
				// Validating fields is easy. The form will mark erroneous fields
				// and display error messages accordingly.
				Validate(func(str string) error {
					if _, err := strconv.ParseInt(str, 10, 32); err != nil {
						return errors.New("That value is invalid, please be sure to enter a number")
					}
					return nil
				}),

			huh.NewInput().
				Title("Enter Speed Evs: 0-255").
				Value(&ev_speed).
				// Validating fields is easy. The form will mark erroneous fields
				// and display error messages accordingly.
				Validate(func(str string) error {
					if _, err := strconv.ParseInt(str, 10, 32); err != nil {
						return errors.New("That value is invalid, please be sure to enter a number")
					}
					return nil
				}),
		),
	)
	err = pokeForm2.Run()
	if err != nil {
		panic(err)
	}
	//Add EVs to pokemon
	var evInts []int64
	ev, _ := strconv.ParseInt(ev_hp, 10, 32)
	evInts = append(evInts, ev)
	ev, _ = strconv.ParseInt(ev_a, 10, 32)
	evInts = append(evInts, ev)
	ev, _ = strconv.ParseInt(ev_d, 10, 32)
	evInts = append(evInts, ev)
	ev, _ = strconv.ParseInt(ev_spa, 10, 32)
	evInts = append(evInts, ev)
	ev, _ = strconv.ParseInt(ev_spd, 10, 32)
	evInts = append(evInts, ev)
	ev, _ = strconv.ParseInt(ev_speed, 10, 32)
	evInts = append(evInts, ev)
	pokemon.CalcStats(int32(evInts[0]), int32(evInts[1]), int32(evInts[2]), int32(evInts[3]), int32(evInts[4]), int32(evInts[5]))

	//Add nature
	pokemon.nature = pokeNature

	//Add ability
	pokemon.ability_read = pokeAbility
	pokemon.ability = pokeAbility

	//Add ui
	pokemon.UI = "p" + iteration

	//Random gender
	if !pokemon.genderless {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		if random := r.Int31n(2); random == 1 {
			pokemon.gender = "Male"
		} else {
			pokemon.gender = "Female"
		}
	} else { //Pokemon has no gender
		pokemon.gender = "None"
	}

	//Give moves
	for index, value := range moveset {
		move, err := GetMove(value, db)
		if err != nil {
			log.Fatal(err)
		}
		pokemon.moveset[index] = move
	}

	//Set owner
	pokemon.owner = true

	//Add pokemon to team
	you.pokemon_list = append(you.pokemon_list, &pokemon)
}

// Gets format from a "huh" form
func getFormat() string {
	var format string
	formatForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select your format!").
				Options(
					huh.NewOption("Singles", "Singles"),
					huh.NewOption("Doubles", "Doubles"),
				).
				Value(&format),
		),
	)
	err := formatForm.Run()
	if err != nil {
		panic(err)
	}
	return format
}

// The first pokemon of the battle get sent out
func battle_start(format string, you *Player, play *Field) {
	var name []string
	var teamList []string

	for _, value := range you.pokemon_list {
		teamList = append(teamList, value.species)
	}

	switch format {
	case "Singles":
		start := huh.NewForm(
			huh.NewGroup(
				huh.NewMultiSelect[string]().
					Options(
						huh.NewOptions(teamList...)...).
					Title("Who do you want to send out?").
					Limit(1).
					Validate(func(str []string) error {
						if len(str) < 1 { //Make sure 2 pokemon were selected
							return errors.New("Please be sure to select a pokemon")
						}
						return nil
					}).
					Value(&name),
			),
		)
		err := start.Run()
		if err != nil {
			panic(err)
		}
	case "Doubles":
		start := huh.NewForm(
			huh.NewGroup(
				huh.NewMultiSelect[string]().
					Options(
						huh.NewOptions(teamList...)...).
					Title("Who do you want to send out? Choose 2!").
					Limit(2).
					Validate(func(str []string) error {
						if len(str) < 2 { //Make sure 2 pokemon were selected
							return errors.New("You must select at least 2 pokemon")
						}
						return nil
					}).
					Value(&name),
			),
		)
		err := start.Run()
		if err != nil {
			panic(err)
		}
	}

	//Name now contains 1 or 2 values of the names of mons to send out
	if len(name) == 2 {
		//Doubles
		you.GetFromName(name[0]).Swap(play, you, nil)
		you.GetFromName(name[1]).Swap(play, you, nil)
	} else {
		//Singles
		you.GetFromName(name[0]).Swap(play, you, nil)
	}
}

// Randomly select what the CPUs pokemon will do
func PickCpuMoves(opponent *Player, play *Field, wg *sync.WaitGroup) {
	defer wg.Done()
	var playerMons []*Pokemon           //Pokemon on the field that belong to the player
	var opponentMons []*Pokemon         //Pokemon on the field that belong to the computer
	for _, pokemon := range play.mons { //Get pokemon on field that belong to you and the user
		if pokemon.owner {
			playerMons = append(playerMons, pokemon)
		} else {
			opponentMons = append(opponentMons, pokemon)
		}
	}

	//For each pokemon (minimum 1, max 2) tell it what to do
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var moveIndex int
	for _, pokemon := range opponentMons { //For each pokemon the cpu has on the field: select a move
		moveIndex = r.Intn(4)
		selectedMove := pokemon.moveset[moveIndex]
		if play.format == "Singles" && selectedMove.target == "Adjacent Ally" { //There is no adjacent ally in singles so pick another move
			for _, move := range pokemon.moveset {
				if move.target != "Adjacent Ally" {
					selectedMove = move
					break
				}
			} //If there is no move that isnt adjacent ally then the move will be selected anyway
		}

		switch selectedMove.target {
		//Move index + target ("_" as target means there is none and the move will fail)
		/*
			EXAMPLES:
				See pokemon.slot in pokemon.go file
		*/
		case "Selected Target":
			if play.format == "Doubles" { //In doubles pick a random target
				pokemon.action[0] = strconv.Itoa(moveIndex) //Slot is selected here so allies are not hit with moves
				pokemon.action[1] = strconv.Itoa(int(playerMons[r.Intn(len(playerMons))].slot))
			} else { //Singles
				pokemon.action[0] = strconv.Itoa(moveIndex)
				pokemon.action[1] = "1"
			}
		case "Self":
			pokemon.action[0] = strconv.Itoa(moveIndex)
			pokemon.action[1] = "0"
		case "Adjacent Ally":
			if play.format == "Singles" { //In singles there is no adjacent ally
				pokemon.action[0] = strconv.Itoa(moveIndex)
				pokemon.action[1] = "_"
			} else { //In doubles there is ONLY ONE adjacent ally
				for _, mon := range opponentMons {
					if mon.UI != pokemon.UI {
						pokemon.action[0] = strconv.Itoa(moveIndex)
						pokemon.action[1] = string(mon.slot)
						break
					}
				}
				if pokemon.action[0] == "" { //In doubles but all allies have fainted
					pokemon.action[0] = strconv.Itoa(moveIndex)
					pokemon.action[1] = "_"
				}
			}
		case "All Adjacent Foes":
			pokemon.action[0] = strconv.Itoa(moveIndex)
			pokemon.action[1] = "4"
		case "All Adjacent":
			pokemon.action[0] = strconv.Itoa(moveIndex)
			pokemon.action[1] = "5"
		}
	} //All pokemon on the opponents side should have "action" filled
}

// Select your moves for each pokemon on the field
func PickYourMoves(you *Player, play *Field) {
	//Run "huh" form to pick moves, make swaps, and MAYBE check field status like weather terrain and such with bubbles tables
	var yourMons []*Pokemon
	var theirMons []*Pokemon
	var swappable []*Pokemon
	for _, mon := range play.mons { //Get your pokemon and their pokemon that are on the field
		if mon.owner {
			yourMons = append(yourMons, mon)
		} else {
			theirMons = append(theirMons, mon)
		}
	}

	for _, mon := range you.pokemon_list { //Get pokemon that can be swapped in (put this in as an error in the form)
		if !mon.on_field && mon.stats[0] > 0 { //If pokemon is in the back and not fainted
			swappable = append(swappable, mon)
		}
	}

	for _, pokemon := range yourMons { //What will each pokemon do
		prompt := "What will " + pokemon.species + " do?"
		var confirm bool
		//Run the form to pick moves and the select the target
		var action string
		var actionList []string
		for _, move := range pokemon.moveset {
			actionList = append(actionList, move.name)
		}
		actionList = append(actionList, "Switch Out")

	ActionSelect: //If confirmation is false you should get here
		//Form 1 move or swap select
		actionForm := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title(prompt).
					Options(
						huh.NewOptions(actionList...)...).
					Value(&action).
					Validate(func(str string) error {
						for _, mon := range theirMons { //Check for trapping abilities
							if str == "Switch Out" && (mon.ability_read == "Arena Trap" || mon.ability_read == "Shadow Tag") {
								return errors.New("You can't escape!")
							}
						}

						if str == "Switch Out" && swappable == nil {
							return errors.New("There is no one who you can swap to")
						}

						return nil
					}),
				huh.NewConfirm().
					Title("Are you sure?").
					Affirmative("Yes!").
					Negative("No.").
					Value(&confirm),
			),
		)

		err := actionForm.Run() //Choose the pokemon's action
		if err != nil {
			panic(err)
		}

		if !confirm { //If there was a mistake then retry
			goto ActionSelect
		}

		//Choose move target or choose swap target
		if action != "Switch Out" {
			//Get move that they picked by name
			var selectedMove Moves

			//Setup .action based on target
			if selectedMove.target == "Self" { //"move index:slot" if the move is a self move then it will be slot 0
				for index, move := range pokemon.moveset {
					if move.name == action {
						selectedMove = move
						pokemon.action[0] = strconv.Itoa(index)
						pokemon.action[1] = "0"
						return
					}
				}
			} else if selectedMove.target == "Adjacent Ally" {
				for index, move := range pokemon.moveset {
					if move.name == action {
						selectedMove = move
						pokemon.action[0] = strconv.Itoa(index)
						pokemon.action[1] = "3"
						return
					}
				}
			} else if selectedMove.target == "All Adjacent Foes" {
				for index, move := range pokemon.moveset {
					if move.name == action {
						selectedMove = move
						pokemon.action[0] = strconv.Itoa(index)
						pokemon.action[1] = "4"
						return
					}
				}
			} else if selectedMove.target == "All Adjacent" {
				for index, move := range pokemon.moveset {
					if move.name == action {
						selectedMove = move
						pokemon.action[0] = strconv.Itoa(index)
						pokemon.action[1] = "5"
						return
					}
				}
			} else {
				for index, move := range pokemon.moveset { //Single target move and not targeted self
					if move.name == action {
						pokemon.action[0] = strconv.Itoa(index)
						selectedMove = move
						if play.format == "Singles" {
							pokemon.action[1] = "1" //Singles format so theres only one target
							return
						}
					}
				}
				//If you got here its doubles and you need another target
			}

			//Single target move so prompt the user for targets (if you get here this is the doubles format)
			var targets []string
			for _, mon := range yourMons {
				if mon.UI != pokemon.UI {
					targets = append(targets, mon.species)
				}
			}
			for _, mon := range theirMons {
				text := "Opponent's " + mon.species
				targets = append(targets, text)
			}

			//Create form to pick target
			targetForm := huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[string]().
						Title(prompt).
						Options(
							huh.NewOptions(targets...)...).
						Value(&action),

					huh.NewConfirm().
						Title("Are you sure?").
						Affirmative("Yes!").
						Negative("No.").
						Value(&confirm),
				),
			)
			err := targetForm.Run() //Choose the pokemon's target
			if err != nil {
				panic(err)
			}

			//Process selected targets (target will be some pokemon on opponents side)
			if string(action[0]) == "O" {
				//Get opponents pokemon.slot
				name := strings.Split(action, "Opponent's ")[1]
				for _, mon := range theirMons {
					if mon.species == name {
						pokemon.action[1] = strconv.Itoa(int(mon.slot))
					}
				}
			} else {
				//Get allies slot
				pokemon.action[1] = "3"
			}
			return
		}

		//Handle switch outs
		var targets []string
		for _, mon := range yourMons {
			if !mon.on_field {
				targets = append(targets, mon.species)
			}
		}
		//Create form to pick target
		targetForm := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Choose who to swap into...").
					Options(
						huh.NewOptions(targets...)...).
					Value(&action),

				huh.NewConfirm().
					Title("Are you sure?").
					Affirmative("Yes!").
					Negative("No.").
					Value(&confirm),
			),
		)
		err = targetForm.Run() //Choose the pokemon's swap target
		if err != nil {
			panic(err)
		}

		pokemon.action[0] = "Swap"
		pokemon.action[1] = action //Swap + pokemon.name
	}
}

// Huh form to swap out your pokemon at the end of a turn or during (U-Turn)
func SwapForm(play *Field, you *Player, pokemon *Pokemon) { //Pokemon can be nil
	var confirm bool
	var swapList []string
	var response string
	prompt := "Who would you like to swap " + pokemon.species + " into?"

	for _, mon := range play.mons {
		if mon.owner && mon.status != "Fainted" && !mon.on_field { //You are the owner, not fainted, and not on the field
			swapList = append(swapList, mon.species)
		}
	}

	if swapList == nil {
		log.Fatal("Swap list is empty, this should be checked prior to function being called")
	}

	swapForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(prompt).
				Options(
					huh.NewOptions(swapList...)...).
				Value(&response).
				Validate(func(str string) error {
					return nil
				}),
			huh.NewConfirm().
				Title("Are you sure?").
				Affirmative("Yes!").
				Negative("No.").
				Value(&confirm),
		),
	)

	err := swapForm.Run() //Choose the pokemon's action
	if err != nil {
		panic(err)
	}

	//Call pokemon.swap with new pokemon
	pokemon = you.GetFromName(response)
	pokemon.Swap(play, you, nil)
}

// Checks the fainted channel and reads its values
func CheckFainted(play *Field, you *Player, opponent *Player) {
	for {
		select { //Check for swaps because pokemon fainted
		case monUI, ok := <-fainted: //Recieve pokemon.UI
			if !ok {
				fmt.Println("Fainted channel is empty")
				return
			}
			pokemon := play.GetFromUI(monUI, you, opponent)
			if pokemon.owner {
				//Check if the pokemon is yours, if not the run swap on random pokemon
				SwapForm(play, you, pokemon)
			} else {
				fmt.Println("Opponent is choosing who to send out next...")
				for _, mon := range opponent.pokemon_list {
					if !mon.on_field && mon.status != "Fainted" {
						time.Sleep(1 * time.Second)
						mon.Swap(play, opponent, nil)
						break //Breaks out of for loop and continues on
					}
				}
			}
		default:
			return //No value in here so just return
		}
	}
}

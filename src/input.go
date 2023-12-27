package main

import (
	"io"
	"os"
	"strings"

	json "github.com/json-iterator/go"
)

// Returns uppercased first char from string input
func UpperFirst(input string) string {
	//Uppercase, combine, and return
	s := strings.ToUpper(string(input[0]))
	return s + input[1:]
}

// Return type matchup json
func GetJson() map[string]map[string]float64 {
	var json = json.ConfigCompatibleWithStandardLibrary
	var newMap map[string]map[string]float64

	// Open the file
	//When main is in src folder: file, err := os.Open("../matchups.json")
	file, err := os.Open("../matchups.json")
	if err != nil {
		panic(err)
	}
	defer file.Close() // Ensure the file is closed when we're done with it

	// Read all bytes from the file
	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	//Assign data to map from json
	if err := json.Unmarshal(data, &newMap); err != nil {
		panic(err)
	}

	return newMap
}

// Get input function but should work with new huh framework
func GetBubble(input string) string {
	if len(strings.TrimSpace(input)) == 0 {
		input = " "
	}

	//Make sure everything is lowered
	input = strings.ToLower(input)
	//Fix input problem removing input spaces
	input = strings.TrimSpace(input)

	if strings.Contains(input, " ") {
		//Handle 2 word moves (Check for 3 word moves to see if they exist in db)
		words := strings.Split(input, " ")
		word1 := UpperFirst(words[0])
		word2 := UpperFirst(words[1])
		input = ""
		input = word1 + " " + word2
		return input
	}

	//Handle 1 word input
	//Uppercase first letter
	return UpperFirst(input)
}

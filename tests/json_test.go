package main 

import (
	"testing"
	"os"
	"io"
	"encoding/json"
	json2 "github.com/json-iterator/go"
)

func BenchmarkJson(b *testing.B) {
	for i := 0; i < b.N; i++ {
		std_json()
	}
}

//The other json library really is faster
func BenchmarkJsonGithub(b *testing.B) {
	for i := 0; i < b.N; i++ {
		json_git()
	}
}

func json_git() map[string]map[string]float64 {
	var json = json2.ConfigCompatibleWithStandardLibrary
	var newMap map[string]map[string]float64
	// Open the file
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

func std_json() map[string]map[string]float64 {
	var newMap map[string]map[string]float64
	// Open the file
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
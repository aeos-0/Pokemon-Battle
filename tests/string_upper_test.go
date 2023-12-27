package main

import (
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"
)

// This will keep the same rune method as before
func BenchmarkUpperFirstRune(b *testing.B) {
	for i := 0; i < b.N; i++ {
		UpperFirstRune("thebatman")
	}
}

// This will try to use strings
func BenchmarkUpperFirstString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		UpperFirstString("thebatman")
	}
}

func BenchmarkUpperFirstString2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		UpperFirstString2("thebatman")
	}
}

func UpperFirstString2(input string) string { //Best performance
	//Uppercase, combine, and return
	s := strings.ToUpper(string(input[0]))
	return s + input[1:]
}

func UpperFirstString(input string) string {
	//Uppercase, combine, and return (without making new var)
	return strings.ToUpper(string(input[0])) + input[1:]
}

func UpperFirstRune(input string) string {
	//For some reason I thought you couldn't access string characters with input = input[0] but this doesn't work because they are immutable
	/*Everything down here is new (contains lots of allocations)*/
	firstRune, _ := utf8.DecodeRuneInString(input)      //Get first rune
	firstRune = unicode.ToUpper(firstRune)              //Uppercase it
	runeSlice := []rune(input)                          //Create slice of runes from string
	runeSlice = runeSlice[1:]                           //Remove first character
	runeSlice = append([]rune{firstRune}, runeSlice...) //Add new rune to the front

	// Change array of runes to string and return
	input = ""
	for _, value := range runeSlice {
		input += string(value)
	}
	return input
}

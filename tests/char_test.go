package main

import (
	"strings"
	"testing"
)

// This will keep the same rune method as before
func BenchmarkCharRune(b *testing.B) {
	for i := 0; i < b.N; i++ {
		get_char_rune("thebatman")
	}
}

// This will try to use strings
func BenchmarkCharString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		get_char_string("thebatman")
	}
}

func BenchmarkCharString2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		get_char_string2("thebatman")
	}
}

func get_char_rune(input string) string {
	input = strings.TrimSpace(input)
	//Weird input fix
	runes := []rune(input)
	input = string(runes[0])

	//Make sure strings are lowercase for consistency
	input = strings.ToLower(input)
	return input
}

func get_char_string(input string) string {
	input = strings.TrimSpace(input)
	input = string(input[0])

	//Make sure strings are lowercase for consistency
	input = strings.ToLower(input)
	return input
}

func get_char_string2(input string) string {
	input = strings.TrimSpace(input)
	return strings.ToLower(string(input[0]))
}

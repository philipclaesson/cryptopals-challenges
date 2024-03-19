package lib

import (
	"math"
	"strings"
)

type CharFreq struct {
	Char rune
	Freq float64
}

func getFrequencyScore(s string) float64 {
	// Letter frequencies according to https://www3.nd.edu/~busiforc/handouts/cryptography/letterfrequencies.html
	var charFreqs = []CharFreq{
		{'E', 0.0111607},
		{'M', 0.030129},
		{'A', 0.084966},
		{'H', 0.030034},
		{'R', 0.075809},
		{'G', 0.024705},
		{'I', 0.075448},
		{'B', 0.020720},
		{'O', 0.071635},
		{'F', 0.018121},
		{'T', 0.069509},
		{'Y', 0.017779},
		{'N', 0.066544},
		{'W', 0.012899},
		{'S', 0.057351},
		{'K', 0.011016},
		{'L', 0.054893},
		{'V', 0.010074},
		{'C', 0.045388},
		{'X', 0.002902},
		{'U', 0.036308},
		{'Z', 0.002722},
		{'D', 0.033844},
		{'J', 0.001965},
		{'P', 0.031671},
		{'Q', 0.001962},
		{'@', 0.0}, // me punishing strings for containing chars that are unlikely in actual text
		{'/', 0.0},
		{'\\', 0.0},
		{':', 0.0},
		{'}', 0.0},
		{'{', 0.0},
		{'(', 0.0},
		{')', 0.0},
		{'*', 0.0},
		{'|', 0.0},
		{']', 0.0},
		{'[', 0.0},
		{';', 0.0},
		{':', 0.0},
		{'-', 0.0},
		{'\'', 0.0},
		{'`', 0.0},
		{'^', 0.0},
		{'$', 0.0},
		{'1', 0.0},
		{'2', 0.0},
		{'3', 0.0},
		{'4', 0.0},
		{'5', 0.0},
		{'6', 0.0},
		{'7', 0.0},
		{'8', 0.0},
		{'9', 0.0},
		{'0', 0.0},
		{'~', 0.0},
		{'>', 0.0},
		{'<', 0.0},
		{' ', 0.2}, // a guess
	}
	errSum := 0.0
	for _, charFreq := range charFreqs {
		errSum += (math.Abs(getCharFrequency(s, charFreq.Char) - charFreq.Freq))
	}
	return errSum
}

func getCharFrequency(s string, c rune) float64 {
	count := 0.0
	for _, cAtPos := range strings.ToUpper(s) {
		if c == cAtPos {
			count++
		}
	}
	return count / float64(len(s))
}

package main

import (
	"fmt"
	"math"
	"os"
)

const hexChars = "0123456789abcdef"

// Single-byte XOR cipher
// https://cryptopals.com/sets/1/challenges/3

func main() {
	// This string has been xored against some single character.
	// Task: find key, decrypt message.
	in := "1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736"
	testHexStringXorChar()

	// We could try xoring with every possible byte and try to decode it
	// It gives us 32 possible candidates which all decode properly into ascii.
	candidates := []byte{}
	candidatesDecoded := []string{}
	for i := 0; i <= 255; i++ {
		str := hexStringXorChar(in, byte(i))
		decodedStr, err := hexDecodeString(str)
		if err == nil {
			candidates = append(candidates, byte(i))
			candidatesDecoded = append(candidatesDecoded, decodedStr)
		}
	}

	if len(candidates) > 1 {
		fmt.Printf("%d candidates found\n", len(candidates))
	} else if len(candidates) == 0 {
		fmt.Printf("No candidates found")
		os.Exit(1)
	} else {
		fmt.Printf("Solution found? %s", candidatesDecoded[0])
		os.Exit(0)
	}

	// We have 32 candidates, let's try do use letter frequencies in order to
	// find the actual text among the noise.
	bestCandidateDecoded := ""
	bestCandidateFreqDiff := 1.0

	for _, candidate := range candidatesDecoded {
		frequencyDiff := getFrequencyDiff(candidate)
		if frequencyDiff < bestCandidateFreqDiff {
			bestCandidateDecoded = candidate
			bestCandidateFreqDiff = frequencyDiff
		}
	}
	// Works! Prints: Cooking MC's like a pound of bacon
	fmt.Printf("best candidate: %s (frequency err: %f)", bestCandidateDecoded, bestCandidateFreqDiff)
}

type CharFreq struct {
	char rune
	freq float64
}

func getFrequencyDiff(s string) float64 {
	// Letter frequencies according to https://www3.nd.edu/~busiforc/handouts/cryptography/letterfrequencies.html
	var charFreqs = []CharFreq{
		{'e', 0.111}, {'a', 0.085}, {'r', 0.075}, {'i', 0.075},
		{'o', 0.075}, {'t', 0.069}, {'n', 0.065}, {'s', 0.057},
		{'k', 0.55}, {'c', 0.045}, {'u', 0.036}, {'d', 0.038},
		{'p', 0.031}, {'m', 0.03},
	}
	errSum := 0.0
	for _, charFreq := range charFreqs {
		errSum += (math.Abs(getCharFrequency(s, charFreq.char) - charFreq.freq))
	}
	return errSum
}

func getCharFrequency(s string, c rune) float64 {
	count := 0.0
	for _, cAtPos := range s {
		if c == cAtPos {
			count++
		}
	}
	return count / float64(len(s))
}

/* xor */
// this function takes a hex encoded string and a char, performs bytewise xor and returns the hex encoded result
func hexStringXorChar(s string, char byte) string {
	out := ""
	for i := 1; i < len(s); i += 2 {
		encodedChars := s[i-1 : i+1]
		b := stringToByte(encodedChars)
		xoredByte := b ^ char
		out += byteToHexString(xoredByte)
	}
	if len(out) != len(s) {
		panic("length mismatch in hexStringXorChar")
	}
	return out
}

// the inverse of xor is... xor: https://stackoverflow.com/questions/14279866/what-is-inverse-function-to-xor
// so we can self test the implementation
// make sure c = a^b, a=b^c holds
func testHexStringXorChar() {
	a := "aaaaddddeeff112341"
	b := byte(92)
	c := hexStringXorChar(a, b)
	a2 := hexStringXorChar(c, b)

	if a != a2 {
		fmt.Println(a)
		fmt.Println(c)
		fmt.Println(a2)
		panic("not working as we thought")
	}
}

/* String decoding stuff */
func hexDecodeString(s string) (string, error) {
	out := ""
	for i := 1; i < len(s); i += 2 {
		encodedChars := s[i-1 : i+1]

		// convert to a byte
		b := stringToByte(encodedChars)

		// convert the 1-byte integer to a ascii char
		decodedChar, err := byteToAscii(b)
		if err != nil {
			// we did not manage to convert it to ascii, which tells us that we've deciphered it wrong. send upwards.
			return "", err
		}

		// append to out
		out += decodedChar
	}
	return out, nil
}

func stringToByte(s string) byte {
	if len(s) != 2 {
		panic("expecting two chars")
	}
	i1 := hexCharToByte(rune(s[0]))
	i2 := hexCharToByte(rune(s[1]))

	return i1<<4 | i2
}

func hexCharToByte(s rune) byte {
	for i, x := range hexChars {
		if s == x {
			return byte(i)
		}
	}
	panic("bad rune")
}

// Takes a byte, returns hex encoded (two character string)
func byteToHexString(b byte) string {
	c1 := hexChars[b>>4]
	c2 := hexChars[b&0xf]
	return string([]byte{c1, c2})
}

func byteToAscii(b byte) (string, error) {
	// lets assume that the encrypted string only contains printable ascii, if we get anything else we can exit early
	printableAscii := []string{" ", "!", "\"", "#", "$", "%", "&", "'", "(", ")", "*", "+", ",", "-", ".", "/", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", ":", ";", "<", "=", ">", "?", "@", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "[", "\\", "]", "^", "_", "`", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "{", "|", "}", "~", "DEL"}
	const firstPrintableAscii byte = 32

	if b < firstPrintableAscii || b >= (firstPrintableAscii+byte(len(printableAscii))) {
		return "", fmt.Errorf("byte %d is not a printable ascii", b)
	}
	return printableAscii[b-firstPrintableAscii], nil
}

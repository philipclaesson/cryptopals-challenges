package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

const hexChars = "0123456789abcdef"

type Candidate struct {
	decoded string
	freq    float64
}

type CharFreq struct {
	char rune
	freq float64
}

// Detect single-character XOR
// https://cryptopals.com/sets/1/challenges/4
// this is a continuation of 3, code copied over.
//
//	One of the 60-character strings in data.txt has been encrypted by single-character XOR - find it.
func main() {
	in := []string{}
	file, err := os.Open("data.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		in = append(in, line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	candidates := []Candidate{}
	for _, line := range in {
		candidatesDecoded := []string{}
		for i := 0; i <= 255; i++ {
			str := hexStringXorChar(line, byte(i))
			decodedStr, err := hexDecodeString(str)
			if err == nil {
				candidatesDecoded = append(candidatesDecoded, decodedStr)
			}
		}

		if len(candidatesDecoded) == 0 {
			// no candidates found for this line - go on with the next line
			continue
		}

		// Use letter frequencies in order to
		// find the actual text among the noise.
		for _, candidate := range candidatesDecoded {
			frequencyDiff := getFrequencyDiff(candidate)
			candidates = append(candidates, Candidate{candidate, frequencyDiff})
		}
		fmt.Printf("line %s had %d candidates\n", line, len(candidatesDecoded))
	}
	fmt.Printf("found %d candidates in %d lines\n", len(candidates), len(in))
	var bestCandidate Candidate = Candidate{
		decoded: "",
		freq:    1000.0,
	}
	for _, c := range candidates {
		if c.freq < bestCandidate.freq {
			bestCandidate = c
		}
		fmt.Println(c.decoded, c.freq)
	}
	fmt.Printf("best candidate: %s (frequency err: %f)", bestCandidate.decoded, bestCandidate.freq)
}

func getFrequencyDiff(s string) float64 {
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
		errSum += (math.Abs(getCharFrequency(s, charFreq.char) - charFreq.freq))
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

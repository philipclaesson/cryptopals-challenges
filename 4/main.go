package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const hexChars = "0123456789abcdef"

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

	candidates := Candidates{}
	for _, line := range in {
		candidates = append(candidates, breakSingleByteXOR(line))
	}
	fmt.Printf("found %d candidates in %d lines\n", len(candidates), len(in))
	sort.Sort(candidates)
	best := candidates[0]
	fmt.Printf("best candidate: %s (frequency err: %f)\n", best.decoded, best.freq)
	fmt.Printf("from the line: %s\n", best.encoded)
	fmt.Printf("encrypted using key: 0x%02x\n", best.key)
	// Prints: Now that the party is jumping, key = 0x35
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
	var retErr error
	for i := 1; i < len(s); i += 2 {
		encodedChars := s[i-1 : i+1]

		// convert to a byte
		b := stringToByte(encodedChars)

		// convert the 1-byte integer to a ascii char
		decodedChar, err := byteToAscii(b)
		if err != nil {
			// we did not manage to convert the byte to ascii, ignore it but pass an error.
			retErr = err
		}

		// append to out
		out += decodedChar
	}
	return out, retErr
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

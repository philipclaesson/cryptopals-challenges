package main

import "fmt"

// https://cryptopals.com/sets/1/challenges/1
func main() {
	// Hex encoded string
	const in = "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"
	// Base64 encoded string
	const out = "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t"
	var res = ""

	// Hex encoded means that every two characters in the input string represent a byte - 4 bytes per character.
	// Base64 encoded store 6 bytes per character.

	// Let's step through in three chars at the time and turn that into two base64 chars
	for i := range in {
		if i%3 != 2 {
			continue
		}
		threeChars := in[i-2 : i+1]

		// lets convert each char to a byte (we only care about the lower 4 bits of each byte but yeah)
		var b1 int = hexToInt(string(threeChars[0]))
		var b2 int = hexToInt(string(threeChars[1]))
		var b3 int = hexToInt(string(threeChars[2]))

		// let's now form a 12 bit integer using 4 bits from each byte
		var b12 int = (b1 << 8) | (b2 << 4) | b3

		// now we can convert our 12 bits into a base64 char
		chars := intToBase64Chars(b12)
		res += chars
	}
	fmt.Printf("Result: %s\n", res)
	fmt.Printf("Expect: %s\n", out)
	if res == out {
		fmt.Println("It's a match 🫶")
	} else {
		fmt.Println("No Michael! That was so not right")
	}
}

func hexToInt(s string) int {
	switch s {
	case "0":
		return 0x00
	case "1":
		return 0x01
	case "2":
		return 0x02
	case "3":
		return 0x03
	case "4":
		return 0x04
	case "5":
		return 0x05
	case "6":
		return 0x06
	case "7":
		return 0x07
	case "8":
		return 0x08
	case "9":
		return 0x09
	case "a":
		return 0x0a
	case "b":
		return 0x0b
	case "c":
		return 0x0c
	case "d":
		return 0x0d
	case "e":
		return 0x0e
	case "f":
		return 0x0f
	default:
		panic("invalid hex char")
	}
}

func intToBase64Chars(b int) string {
	base64Chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	if b > 1<<12 {
		panic("noooooo to big")
	}
	b1 := b >> 6
	b2 := b & 0b111111

	// now convert these to chars
	c1 := base64Chars[b1]
	c2 := base64Chars[b2]
	return string([]byte{c1, c2})
}

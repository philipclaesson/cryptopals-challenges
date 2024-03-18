package main

import "fmt"

const hexChars = "0123456789abcdef"

// https://cryptopals.com/sets/1/challenges/2
func main() {
	in1 := "1c0111001f010100061a024b53535009181c"
	in2 := "686974207468652062756c6c277320657965"
	out := "746865206b696420646f6e277420706c6179"
	res := ""

	// let's loop through the string byte by byte, two chars at the time
	for i := 1; i < len(in1); i += 2 {
		s1 := in1[i-1 : i+1]
		s2 := in2[i-1 : i+1]

		// convert to bytes
		b1 := stringToByte(s1)
		b2 := stringToByte(s2)

		// bitwise xor
		bres := b1 ^ b2

		res += byteToString(bres)
	}

	fmt.Printf("Result: %s\n", res)
	fmt.Printf("Expect: %s\n", out)
	if res == out {
		fmt.Println("It's a match 🫶")
	} else {
		fmt.Println("No Michael! That was so not right")
	}
}

func stringToByte(s string) byte {
	if len(s) != 2 {
		panic("expecting two chars")
	}
	i1 := hexRuneToInt(rune(s[0]))
	i2 := hexRuneToInt(rune(s[1]))

	return byte(i1<<4 | i2)
}

func hexRuneToInt(s rune) int {
	for i, x := range hexChars {
		if s == x {
			return i
		}
	}
	panic("bad rune")
}

func byteToString(b byte) string {
	c1 := hexChars[b>>4]
	c2 := hexChars[b&0xf]
	return string([]byte{c1, c2})
}

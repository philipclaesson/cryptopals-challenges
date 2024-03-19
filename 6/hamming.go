package main

import "fmt"

/* HAMMING DISTANCE */

// Returns the hamming distance between two strings of equal length
func hammingDistance(s1 string, s2 string) int {
	if len(s1) != len(s2) {
		panic("hammingDistance assumes same length of s1 and s2")
	}
	b1 := []byte(s1)
	b2 := []byte(s2)
	distance := 0
	for i := range b1 {
		// the number of high bits in the xor should be the hamming distance
		distance += numberOfHighBits(b1[i] ^ b2[i])
	}
	return distance
}

func numberOfHighBits(b byte) int {
	count := 0
	for i := 0; i < 8; i++ {
		count += (int(b) >> i) & 1
	}
	return count
}

func testHammingDistance() {
	distance := hammingDistance("this is a test", "wokka wokka!!!")
	if distance != 37 {
		panic(fmt.Sprintf("Hamming distance does not work %d != %d", distance, 37))
	}
}

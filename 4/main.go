package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"

	"github.com/philipclaesson/cryptopals-challenges/lib"
)

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

	candidates := lib.Candidates{}
	for _, line := range in {
		candidates = append(candidates, lib.BreakSingleByteXOR(line))
	}
	fmt.Printf("found %d candidates in %d lines\n", len(candidates), len(in))
	sort.Sort(candidates)
	best := candidates[0]
	fmt.Printf("best candidate: %s (frequency err: %f)\n", best.Decoded, best.FreqScore)
	fmt.Printf("from the line: %s\n", best.Encoded)
	fmt.Printf("encrypted using key: 0x%02x\n", best.Key)
	// Prints: Now that the party is jumping, key = 0x35
}

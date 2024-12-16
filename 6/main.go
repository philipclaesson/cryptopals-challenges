package main

import (
	"fmt"

	"github.com/philipclaesson/cryptopals-challenges/lib"
)

// Break repeating-key XOR
// https://cryptopals.com/sets/1/challenges/6
// data.txt has been base64'd after being encrypted with repeating-key XOR.
// Decrypt it.

func main() {
	testHammingDistance()

	// Read data from file
	payload, err := lib.ReadDataFile("data.txt")
	if err != nil {
		fmt.Println("Error reading data file:", err)
		return
	}
	testCreateBlocks(payload)

	lib.Testb64Encode()

	rawPayload := lib.B64Decode([]byte(payload))

	// This was the trickiest part - depending on how many blocks to sample upon
	// I got different key sizes. more chunks seemed to converge on 29 which was the answer.
	keysize := findKeySize(string(rawPayload), 2, 40, 40)
	fmt.Println("got keysize", keysize)
	blocks := createBlocks(string(rawPayload), keysize)
	key := []byte{}
	for _, block := range blocks {
		// fmt.Printf("%02x %02x %02x ... %02x %02x\n", block[0], block[1], block[2], block[len(block)-2], block[len(block)-1])
		cand := lib.BreakSingleByteXOR(block)
		key = append(key, cand.Key)
	}
	fmt.Printf("key is %s (%v)\n", string(key), key)

	// let's try to decypt using the found key
	decrypted := lib.XORRepeatingKey(rawPayload, key)
	fmt.Println(string(decrypted)) // Prints Vanilla Ice lyrics as expected
}

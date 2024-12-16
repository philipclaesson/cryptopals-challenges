package main

import (
	"fmt"

	"github.com/philipclaesson/cryptopals-challenges/lib"
)

// AES in ECB mode
// https://cryptopals.com/sets/1/challenges/7
// The Base64-encoded content in this file has been encrypted via AES-128 in ECB mode under the key "YELLOW SUBMARINE" (128 bits).
// Decrypt it.
func main() {

	payload, err := lib.ReadDataFile("data.txt")
	if err != nil {
		fmt.Println("Error reading data file:", err)
		return
	}
	rawPayload := lib.B64Decode([]byte(payload))
	key := []byte{'Y', 'E', 'L', 'L', 'O', 'W', ' ', 'S', 'U', 'B', 'M', 'A', 'R', 'I', 'N', 'E'}
	lib.AES128ECBDecrypt(rawPayload, key)
}

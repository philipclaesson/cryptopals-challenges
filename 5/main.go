package main

import (
	"fmt"
	"os"

	"github.com/philipclaesson/cryptopals-challenges/lib"
)

// https://cryptopals.com/sets/1/challenges/5
// encrypt the payload using repeating-key XOR with the key "ICE"
// you can use this to encrypt and encrypt whatever string payload
// go run . encode/decode payload key
func main() {
	lib.TestHexCodec("asdfg")
	lib.TestXORRepeatingKey()
	lib.TestXORRepeatingKeyEncryptDecrypt([]byte("asdfg"), []byte("ice"))
	if len(os.Args) != 4 {
		fmt.Println("Usage: mode payload key")
		fmt.Println("Ex: go run . encrypt hello ice")
		os.Exit(1)
	}
	mode := os.Args[1]
	payload := os.Args[2]
	key := os.Args[3]
	if mode == "encrypt" {
		fmt.Println(lib.HexEncode(string(lib.XORRepeatingKey([]byte((payload)), []byte(key)))))
	} else if mode == "decrypt" {
		fmt.Println(string(lib.XORRepeatingKey([]byte(lib.HexDecode(payload)), []byte(key))))
	} else {
		panic("what mode")
	}
}

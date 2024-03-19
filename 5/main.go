package main

import (
	"fmt"
	"os"
)

// https://cryptopals.com/sets/1/challenges/5
// encrypt the payload using repeating-key XOR with the key "ICE"
// you can use this to encrypt and encrypt whatever string payload
// go run . encode/decode payload key
func main() {
	testHexCodec("asdfg")
	testEncrypt()
	testEncryptDecrypt([]byte("asdfg"), []byte("ice"))
	if len(os.Args) != 4 {
		fmt.Println("Usage: mode payload key")
		fmt.Println("Ex: go run . encrypt hello ice")
		os.Exit(1)
	}
	mode := os.Args[1]
	payload := os.Args[2]
	key := os.Args[3]
	if mode == "encrypt" {
		fmt.Println(hexEncode(string(encrypt([]byte((payload)), []byte(key)))))
	} else if mode == "decrypt" {
		fmt.Println(string(encrypt([]byte(hexDecode(payload)), []byte(key))))
	} else {
		panic("what mode")
	}
}

// XOR Encrypt an arbitrary payload with an arbitrary key
func encrypt(payload []byte, key []byte) []byte {
	out := payload
	for i := 0; i < len(payload); i++ {
		pb := payload[i]
		kb := key[i%len(key)]
		out[i] = pb ^ kb
	}
	return out
}

func testEncrypt() {
	payload := "Burning 'em, if you ain't quick and nimble\nI go crazy when I hear a cymbal"
	key := "ICE"
	expected := "0b3637272a2b2e63622c2e69692a23693a2a3c6324202d623d63343c2a26226324272765272a282b2f20430a652e2c652a3124333a653e2b2027630c692b20283165286326302e27282f"
	encrypted := encrypt([]byte(payload), []byte(key))
	encryptedAsHex := hexEncode(string(encrypted))
	if expected == encryptedAsHex {
		// fmt.Printf("Encrypted correctly to %s\n", encryptedAsHex)
	} else {
		fmt.Printf("Expected: 	%s\n", expected)
		fmt.Printf("Actual:		%s\n", encryptedAsHex)
		panic("Encryption incorrect\n")
	}
}

func testEncryptDecrypt(payload []byte, key []byte) {
	encrypted := encrypt(payload, key)
	decrypted := encrypt(encrypted, key)
	if string(decrypted) != string(payload) {
		panic(fmt.Sprintf("encrypt decrypt not working as it should %s != %s", decrypted, payload))
	}
}

package lib

import "fmt"

const hexChars = "0123456789abcdef"

// XOR Encrypt/Decrypt an arbitrary payload with an arbitrary rotating key
func XORRepeatingKey(payload []byte, key []byte) []byte {
	out := payload
	for i := 0; i < len(payload); i++ {
		pb := payload[i]
		kb := key[i%len(key)]
		out[i] = pb ^ kb
	}
	return out
}

func TestXORRepeatingKey() {
	payload := "Burning 'em, if you ain't quick and nimble\nI go crazy when I hear a cymbal"
	key := "ICE"
	expected := "0b3637272a2b2e63622c2e69692a23693a2a3c6324202d623d63343c2a26226324272765272a282b2f20430a652e2c652a3124333a653e2b2027630c692b20283165286326302e27282f"
	encrypted := XORRepeatingKey([]byte(payload), []byte(key))
	encryptedAsHex := HexEncode(string(encrypted))
	if expected == encryptedAsHex {
		// fmt.Printf("Encrypted correctly to %s\n", encryptedAsHex)
	} else {
		fmt.Printf("Expected: 	%s\n", expected)
		fmt.Printf("Actual:		%s\n", encryptedAsHex)
		panic("Encryption incorrect\n")
	}
}

func TestXORRepeatingKeyEncryptDecrypt(payload []byte, key []byte) {
	encrypted := XORRepeatingKey(payload, key)
	decrypted := XORRepeatingKey(encrypted, key)
	if string(decrypted) != string(payload) {
		panic(fmt.Sprintf("encrypt decrypt not working as it should %s != %s", decrypted, payload))
	}
}

// this function takes a hex encoded string and a char, performs bytewise xor and returns the hex encoded result
func hexStringXorChar(s string, char byte) string {
	for _, c := range s {
		indexOf(hexChars, c)
	}
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
	panic(fmt.Sprintf("bad rune %v", s))
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

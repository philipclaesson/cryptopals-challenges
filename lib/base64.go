package lib

import "fmt"

const base64Chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

// A probably overcomplicated implementation of base64 decode
func B64Decode(payload []byte) []byte {
	out := []byte{}
	var current24bits uint32 = 0
	for i := 0; i < len(payload); i++ {
		// get a byte from payload
		p := payload[i]

		// = means padding means end of string
		// it is time to end
		if rune(p) == '=' {
			// we need to flush whatever is in current24bits
			// how much is in there is know how far into the quartet we are
			var o byte = 0
			for j := uint32(0); j < uint32(i)%uint32(4); j++ {
				o = byte((current24bits >> (16 - (8 * j))) & 0xff)
				if o != 0 {
					out = append(out, o)
				}
			}
			return out
		}

		// convert from base64 to int.
		// integers will be 6 bit but we want to store them
		// in 32 bit so we can shift them upwards later
		b := uint32(base64CharToByte(rune(p)))

		// store 6 bits per byte into current24bits
		// this is similar to having all 4 bytes and storing them like this
		// t24 := b1<<18 | b2<<12 | b3<<6 | b4
		current24bits |= b << uint32(18-(6*(i%4)))

		// for every 4th byte, flush
		if i%4 == 3 {
			// split up into 3 bytes
			o1 := byte(current24bits >> 16 & 0xff)
			o2 := byte(current24bits >> 8 & 0xff)
			o3 := byte(current24bits & 0xff)

			// add to output
			out = append(out, o1)
			out = append(out, o2)
			out = append(out, o3)
			current24bits = 0
		}
	}
	return out
}

// takes a base64 encoded char
// returns its decoded value in a byte
func base64CharToByte(c rune) byte {
	// = is padding in the end, return 0
	if c == '=' {
		return 0
	}
	return indexOf(base64Chars, rune(c))
}

func Testb64Encode() {
	decoded := string(B64Decode([]byte("YWFh")))
	if decoded != "aaa" {
		panic(fmt.Sprintf("b64decode is incorrect, got %s expected aaa", decoded))
	}

	decoded2 := B64Decode([]byte("aGVsbG8="))
	decoded2s := string(decoded2)
	if len(decoded2) != len("hello") {
		panic(fmt.Sprintf("expected len %d got %d (%s)", len("hello"), len(decoded2), decoded2s))
	}
	if decoded2s != "hello" {
		panic(fmt.Sprintf("b64decode is incorrect, got %s expected hello", decoded2s))
	}

	decoded3 := string(B64Decode([]byte("SW4gY29tcHV0ZXIgcHJvZ3JhbW1pbmcsIEJhc2U2NCBpcyBhIGdyb3VwIG9mIGJpbmFyeS10by10ZXh0IGVuY29kaW5nIHNjaGVtZXMgdGhhdCB0cmFuc2Zvcm1zIGJpbmFyeSBkYXRhIGludG8gYSBzZXF1ZW5jZSBvZiBwcmludGFibGUgY2hhcmFjdGVycy4=")))
	expected3 := "In computer programming, Base64 is a group of binary-to-text encoding schemes that transforms binary data into a sequence of printable characters."
	if decoded3 != expected3 {
		panic(fmt.Sprintf("got %s\nexp %s", decoded3, expected3))
	}
}

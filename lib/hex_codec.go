package lib

func TestHexCodec(payload string) {
	encoded := HexEncode(payload)
	decoded := HexDecode(encoded)
	if payload != decoded {
		panic("hexcodec does not work")
	}
}

func HexEncode(payload string) string {
	var out string
	for i := 0; i < len(payload); i++ {
		h1 := hexChars[payload[i]>>4]
		h2 := hexChars[payload[i]&0xf]
		out += string(h1)
		out += string(h2)
	}
	return out
}

func HexDecode(s string) string {
	var out string
	for i := 1; i < len(s); i += 2 {
		upper := indexOf(hexChars, rune(s[i-1]))
		lower := indexOf(hexChars, rune(s[i]))
		out += string(upper<<4 | lower)
	}
	return out
}

func indexOf(s string, c rune) byte {
	for i, r := range s {
		if c == r {
			return byte(i)
		}
	}
	panic("char not in string")
}

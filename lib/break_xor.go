package lib

import (
	"sort"
)

type Candidate struct {
	Encoded   string
	Decoded   string
	FreqScore float64
	Key       byte
}

type Candidates []Candidate

// Implement sort.Interface for Candidate
func (c Candidates) Len() int           { return len(c) }
func (c Candidates) Less(i, j int) bool { return c[i].FreqScore < c[j].FreqScore }
func (c Candidates) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

// Takes a hex encoded payload, tries XORing with every single byte and returns the most likely key and result
func BreakSingleByteXORHexEncoded(hexPayload string) Candidate {
	candidates := Candidates{}
	for i := 0; i <= 255; i++ {
		str := hexStringXorChar(hexPayload, byte(i))
		decodedStr, _ := hexDecodeString(str)
		if decodedStr != "" {
			candidates = append(candidates, Candidate{
				str,
				decodedStr,
				getFrequencyScore(decodedStr),
				byte(i),
			})
		}
	}

	sort.Sort(candidates)
	return candidates[0]
}

// Takes a byte-array payload, tries XORing with every single byte and decode to string
// returns the most likely key and result
func BreakSingleByteXOR(payload []byte) Candidate {
	candidates := Candidates{}
	for i := 0; i <= 255; i++ {
		// str := hexStringXorChar(hexPayload, byte(i))
		decrypted := []byte{}
		for _, c := range payload {
			decrypted = append(decrypted, c^byte(i))
		}
		str := string(decrypted)
		candidates = append(candidates, Candidate{
			"",
			str,
			getFrequencyScore(str),
			byte(i),
		})
		// fmt.Println(i, getFrequencyScore(str), str[0], str[1], str[2], payload[0], payload[1], payload[2], payload[0]^byte(i), payload[1]^byte(i), payload[2]^byte(i))
	}

	sort.Sort(candidates)
	return candidates[0]
}

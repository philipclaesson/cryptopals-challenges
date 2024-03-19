package lib

import "sort"

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
		str := ""
		for i, c := range payload {
			str += string(c ^ byte(i))
		}
		if str != "" {
			candidates = append(candidates, Candidate{
				str,
				str,
				getFrequencyScore(str),
				byte(i),
			})
		}
	}

	sort.Sort(candidates)
	return candidates[0]
}

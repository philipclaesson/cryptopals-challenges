package main

import "sort"

type Candidate struct {
	encoded string
	decoded string
	freq    float64
	key     byte
}

type Candidates []Candidate

// Implement sort.Interface for Candidate
func (c Candidates) Len() int           { return len(c) }
func (c Candidates) Less(i, j int) bool { return c[i].freq < c[j].freq }
func (c Candidates) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

// Takes an arbitrary payload, tries XORing with every single byte and returns the most likely byte and result
func breakSingleByteXOR(payload string) Candidate {
	candidates := Candidates{}
	for i := 0; i <= 255; i++ {
		str := hexStringXorChar(payload, byte(i))
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

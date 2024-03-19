package main

type KeySize struct {
	size     int
	distance float64
}

func findKeySize(payload string, min int, max int, nChunks int) int {
	// For each KEYSIZE, take the first KEYSIZE worth of bytes, and the second KEYSIZE worth of bytes, and find the edit distance between them.
	// Normalize this result by dividing by KEYSIZE.
	keySizes := []KeySize{}
	for i := min; i <= max; i++ {
		keySizes = append(keySizes, KeySize{
			i,
			getKeySizeDistance(payload, i, nChunks),
		})
		// fmt.Printf("Size: %d, Distance: %f\n", i, getKeySizeDistance(payload, i))
	}

	// Sort keySizes by distance
	for i := 0; i < len(keySizes)-1; i++ {
		for j := i + 1; j < len(keySizes); j++ {
			if keySizes[i].distance > keySizes[j].distance {
				keySizes[i], keySizes[j] = keySizes[j], keySizes[i]
			}
		}
	}
	return keySizes[0].size
}

func getKeySizeDistance(payload string, keysize int, nChunks int) float64 {
	totalDistance := 0.0
	for i := 0; i < nChunks*keysize; i += keysize * 2 {
		chunk1 := payload[i : i+keysize]
		chunk2 := payload[i+keysize : i+keysize*2]
		totalDistance += (float64(hammingDistance(chunk1, chunk2)))
	}
	return totalDistance / float64(keysize)
}

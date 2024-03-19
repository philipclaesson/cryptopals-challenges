package main

// Creates {keysize} blocks.
// This function does the following but in one step:
// Now that you probably know the KEYSIZE: break the ciphertext into blocks of KEYSIZE length.
// Now transpose the blocks: make a block that is the first byte of every block, and a block that is the second byte of every block, and so on.
func createBlocks(payload string, keysize int) [][]byte {
	// initialize block structure
	blocks := [][]byte{}
	for i := 0; i < keysize; i++ {
		blocks = append(blocks, []byte{})
	}

	// The rationale here is to put every byte that has been XOR encrypted with the same
	// letter into the same block
	// For instance, a 3-letter repeating key means every 3rd byte has been encrypted with the same letter.
	for i, p := range payload {
		blockIdx := i % keysize
		blocks[blockIdx] = append(blocks[blockIdx], byte(p))
	}
	return blocks
}

func testCreateBlocks(payload string) {
	blocks := createBlocks(payload, 3)
	if blocks[0][0] != payload[0] {
		panic("Expecting first char in first block to equal first char in file")
	}
	if blocks[2][0] != payload[2] {
		panic("first char in 3rd block should equal 3rd char")
	}
}

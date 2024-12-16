package lib

import "fmt"

func AES128ECBDecrypt(payload []byte, key []byte) []byte {
	/*
		AES:
		From https://en.wikipedia.org/wiki/Advanced_Encryption_Standard
		High-level description of the algorithm
			1. KeyExpansion – round keys are derived from the cipher key using the AES key schedule. AES requires a separate 128-bit round key block for each round plus one more.
			2. Initial round key addition:
				AddRoundKey – each byte of the state is combined with a byte of the round key using bitwise xor.
			3. 9, 11 or 13 rounds:
				1. SubBytes – a non-linear substitution step where each byte is replaced with another according to a lookup table.
				2. ShiftRows – a transposition step where the last three rows of the state are shifted cyclically a certain number of steps.
				3. MixColumns – a linear mixing operation which operates on the columns of the state, combining the four bytes in each column.
				4. AddRoundKey
			4. Final round (making 10, 12 or 14 rounds in total):
				1. SubBytes
				2. ShiftRows
				3. AddRoundKey
		ECB:
		From https://en.wikipedia.org/wiki/Block_cipher_mode_of_operation#Electronic_codebook_(ECB)
		The simplest of the encryption modes is the electronic codebook (ECB) mode (named after conventional physical codebooks). The message is divided into blocks, and each block is encrypted separately.
	*/
	testSplitToECBBlocks()
	if len(key) != 16 {
		panic(fmt.Sprintf("expected key to be 16 bytes, was %d", len(key)))
	}
	const nRounds int = 10 // or 12, or 14

	// 0. Split up into blocks
	splitToECBBlocks(payload)
	/*
		// For each block
		for i, block := range blocks {
			// 	1. KeyExpansion
			roundKeys := KeyExpansion() // TODO: should this be done once per block or not?
			// 	2. AddRoundKey
			block.AddRoundKey()
			// 	3. for i in 10, 12 or 14:
			for round := 0; round < nRounds; round++ {
				// 1. SubBytes
				block.SubBytes()
				// 2. ShiftRows
				block.ShiftRows()
				// 3. MixColumns (except for last round)
				block.MixColumns()
				// 4. AddRoundKey
				block.AddRoundKey()
			}
		}
	*/
	return nil
}

type Block [][]byte

// Creates an n by n block
func NewBlock(n int) Block {
	block := make(Block, n)
	for i := range block {
		block[i] = make([]byte, n)
	}
	return block
}

func splitToECBBlocks(payload []byte) []Block {
	out := []Block{}
	if len(payload)%16 != 0 {
		panic("assumed payload length to be multiple of 16")
	}
	for i := 0; i < len(payload); i += 16 {
		block := NewBlock(4)
		for col := 0; col < 4; col++ {
			for row := 0; row < 4; row++ {
				block[row][col] = payload[i+col*4+row]
			}
		}
		out = append(out, block)
	}
	return out
}

func testSplitToECBBlocks() {
	payload := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	blocks := splitToECBBlocks(payload)
	if len(blocks) != 1 {
		panic("returned too many of few blocks")
	}
	block := blocks[0]

	if block[0][0] != 0 {
		panic(fmt.Sprintf("0, 0 should be 0 is %d", block[0][0]))
	}
	if block[2][3] != 14 {
		panic(fmt.Sprintf("2, 3 should be 14 is %d", block[2][3]))
	}
	if block[2][0] != 2 {
		panic(fmt.Sprintf("2, 0 should be 2 is %d", block[2][0]))
	}
	if block[3][3] != 15 {
		panic(fmt.Sprintf("2, 0 should be 15 is %d", block[3][3]))
	}
	if block[1][2] != 9 {
		panic(fmt.Sprintf("2, 0 should be 9 is %d", block[1][2]))
	}
	fmt.Println("splitToECBBlocks works as expected")
}

// https://en.wikipedia.org/wiki/AES_key_schedule
func KeyExpansion(key []byte) [][]byte {
	if len(key) != 16 {
		panic("bad key")
	}
	sbox := CreateSbox()

	// Key Schedule

	// define N as the length of the key in 32-bit words: 4 words for AES-128, 6 words for AES-192, and 8 words for AES-256
	n := 4
	// define K0, K1, ... KN-1 as the 32-bit words of the original key
	k := []uint32{}
	for i := 0; i < len(key); i += n {
		k = append(k, uint32(key[i])<<24|uint32(key[i+1])<<16|uint32(key[i+2])<<8|uint32(key[i+3]))
	}

	// R as the number of round keys needed: 11 round keys for AES-128, 13 keys for AES-192, and 15 keys for AES-256[note 4]
	r := 11

	// W0, W1, ... W4R-1 as the 32-bit words of the expanded key
	w := []uint32{}
	for i := 0; i < 4*r; i++ {
		var wi uint32
		if i < n {
			wi = k[i]
		} else if i >= n && i%n == 0 {
			wi = w[i-n] ^ SubWord(RotWord(w[i-1]), sbox)
		}
		w = append(w, wi)
	}
	return nil
}

// a one-byte left circular shift
func RotWord(w uint32) uint32 {
	t := w >> 24
	return w<<8 | t
}

func TestRotWord() {
	b1 := uint32(0b11110101)
	b2 := uint32(0b11101101)
	b3 := uint32(0b00000000)
	b4 := uint32(0b11011010)
	out := RotWord(b1<<24 | b2<<16 | b3<<8 | b4)
	o1 := out >> 24 & 0xff
	o2 := out >> 16 & 0xff
	o3 := out >> 8 & 0xff
	o4 := out & 0xff
	if (b2 != o1) || (b3 != o2) || (b4 != o3) || (b1 != o4) {
		panic("RotWord incorrectly implemented")
	}
}

// an application of the AES S-box to each of the four bytes of the word
// https://en.wikipedia.org/wiki/Rijndael_S-box
func SubWord(w uint32, sbox [256]byte) uint32 {
	b1 := uint32(sbox[w>>24]) << 24
	b2 := uint32(sbox[w>>16&0xff]) << 16
	b3 := uint32(sbox[w>>8&0xff]) << 8
	b4 := uint32(sbox[w&0xff])
	return (b1 | b2 | b3 | b4)
}

func TestSubWord() {
	sbox := CreateSbox()
}

func ROTL8(x uint8, shift uint8) uint8 {
	return uint8(x<<(shift)) | (x >> (8 - (shift)))
}

// this implementation is my golang translation of the c reference implementation in
// https://en.wikipedia.org/wiki/Rijndael_S-box
func CreateSbox() [256]byte {
	p, q := uint8(1), uint8(1)

	sbox := [256]byte{}
	for {
		// multiply p by 3
		if p&0x80 != 0 {
			p = p ^ (p << 1) ^ 0x1B
		} else {
			p = p ^ (p << 1)
		}

		// divide q by 3 (equals multiplication by 0xf6)
		q ^= q << 1
		q ^= q << 2
		q ^= q << 4
		if q&0x80 != 0 {
			q ^= 0x09
		}

		// compute the affine transformation
		xformed := q ^ ROTL8(q, 1) ^ ROTL8(q, 2) ^ ROTL8(q, 3) ^ ROTL8(q, 4)
		sbox[p] = xformed ^ 0x63

		if p == 1 {
			break
		}
	}
	sbox[0] = 0x63
	return sbox
}

func (b *Block) AddRoundKey() {
}

func (b *Block) SubBytes() {
}

func (b *Block) ShiftRows() {
}

func (b *Block) MixColumns() {
}

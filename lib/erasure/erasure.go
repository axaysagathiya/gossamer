// Copyright 2021 ChainSafe Systems (ON)
// SPDX-License-Identifier: LGPL-3.0-only

package erasure

import (
	"bytes"
	"errors"
	"fmt"
	"math"

	"github.com/ChainSafe/gossamer/lib/common"
	"github.com/ChainSafe/gossamer/lib/parachain"
	"github.com/ChainSafe/gossamer/lib/trie"
	"github.com/ChainSafe/gossamer/pkg/scale"
	"github.com/klauspost/reedsolomon"
)

type Additive uint16

// ErrNotEnoughValidators cannot encode something for zero or one validator
var ErrNotEnoughValidators = errors.New("expected at least 2 validators")

var MAX_VALIDATORS = uint(65536)

// Reconstruct the missing data from a set of chunks
func Reconstruct(validatorsQty, originalDataLen int, chunks [][]byte) ([]byte, error) {
	recoveryThres, err := recoveryThreshold(validatorsQty)
	if err != nil {
		return nil, err
	}

	enc, err := reedsolomon.New(validatorsQty, recoveryThres)
	if err != nil {
		return nil, err
	}
	err = enc.Reconstruct(chunks)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	err = enc.Join(buf, chunks, originalDataLen)
	return buf.Bytes(), err
}

func recoveryThreshold(validators uint) (uint, error) {
	if validators <= 1 {
		return 0, ErrNotEnoughValidators
	}

	needed := (validators - 1) / 3

	return needed + 1, nil
}

func deriveParameters(totalShards uint, shardsNeededToRecover uint) (uint, uint, error) {
	if totalShards < 2 {
		return totalShards, shardsNeededToRecover, fmt.Errorf("Number of wanted shards must be at least 2, but is %d", totalShards)
	}
	if shardsNeededToRecover < 1 {
		return totalShards, shardsNeededToRecover, fmt.Errorf("Number of wanted payload shards must be at least 1, but is %d", shardsNeededToRecover)
	}

	dataShards := NextHigherPowerOfTwo(totalShards)
	parityShards := NextLowerPowerOfTwo(shardsNeededToRecover)

	if totalShards*parityShards > dataShards*shardsNeededToRecover {
		panic(totalShards*parityShards > dataShards*shardsNeededToRecover)
	}

	if dataShards > MAX_VALIDATORS {
		return dataShards, parityShards, fmt.Errorf("Number of wanted shards %d exceeds max of 2^16", dataShards)
	}

	return dataShards, parityShards, nil
}

// isPowerOfTwo checks if a number is a power of 2
func isPowerOfTwo(n uint) bool {
	return n > 0 && (n&(n-1)) == 0
}

// NextHigherPowerOfTwo ticks down to the next higher power of 2
func NextHigherPowerOfTwo(n uint) uint {
	if isPowerOfTwo(n) {
		return n
	}

	power := uint(1)
	for power < n {
		power <<= 1
	}
	return power
}

// NextLowerPowerOfTwo ticks down to the next lower power of 2
func NextLowerPowerOfTwo(n uint) uint {
	if isPowerOfTwo(n) {
		return n
	}

	power := uint(1)
	for power < n {
		power <<= 1
	}
	return power >> 1
}

func ChunksToTrie(chunks [][]byte) *trie.Trie {
	chunkTrie := trie.NewEmptyTrie()
	for i, chunk := range chunks {
		encodedI := scale.MustMarshal(uint(i))
		chunkHash := common.MustBlake2bHash(chunk)

		if err := chunkTrie.Put(encodedI, chunkHash[:]); err != nil {
			panic("a fresh trie stored in memory cannot have errors loading nodes")
		}
	}
	return chunkTrie
}

// ObtainChunks obtains erasure-coded chunks, divides data into number of validatorsQty chunks and
// creates parity chunks for reconstruction
func ObtainChunks(validatorsQty uint, data parachain.AvailableData) ([][]byte, error) {
	if validatorsQty > MAX_VALIDATORS {
		return nil, fmt.Errorf("TooManyValidators")
	}

	encodedData := scale.MustMarshal(data)
	if len(encodedData) == 0 {
		return nil, fmt.Errorf("BadPayload")
	}

	recoveryThres, err := recoveryThreshold(validatorsQty)
	if err != nil {
		return nil, err
	}

	dataShards, parityShards, err := deriveParameters(validatorsQty, recoveryThres)

	shardLen := ShardLen(parityShards, uint(len(encodedData)))
	if shardLen <= 0 {
		return nil, fmt.Errorf("invalid shardLen")
	}

	parityShards2 := parityShards * 2

	shards := make([][]byte, validatorsQty)
	for i := uint(0); i < validatorsQty; i++ {
		v := make([]byte, shardLen)
		shards[i] = v
	}

	for chunkIdx, i := uint(0), uint(0); i < uint(len(encodedData)); i += parityShards2 {
		end := uint(math.Min(float64(i+parityShards2), float64(len(encodedData))))
		if i != end {
			dataPiece := encodedData[i:end]
			if len(dataPiece) == 0 {
				// Handle the assertion logic here
			}
			if uint(len(dataPiece)) > parityShards2 {
				// Handle the assertion logic here
			}

			encodingRun, err := encodeSub(dataPiece, dataShards, parityShards2)
			if err != nil {
				// Handle the encoding error
			}
			// for valIdx := uint(0); valIdx < validatorsQty; valIdx++ {
			// 	shards[valIdx] = append(shards[valIdx], [2]byte{
			// 		byte((encodingRun[valIdx].Value >> 8) & 0xFF),
			// 		byte(encodingRun[valIdx].Value & 0xFF),
			// 	})
			// }
			chunkIdx++
		}

	}

	return nil, nil
}

func ShardLen(parityShards, payloadSize uint) uint {
	payloadSymbols := (payloadSize + 1) / 2
	shardSymbolsCeil := (payloadSymbols + parityShards - 1) / parityShards
	shardBytes := shardSymbolsCeil * 2
	return shardBytes
}

func encodeSub(bytes []byte, dataShards, parityShards uint) ([]Additive, error) {
	// Check that the input parameters are valid.
	if !isPowerOfTwo(dataShards) {
		return nil, errors.New("Algorithm only works for 2^i sizes for N")
	}
	if !isPowerOfTwo(parityShards) {
		return nil, errors.New("Algorithm only works for 2^i sizes for K")
	}
	if uint(len(bytes)) > (parityShards << 1) {
		return nil, errors.New("Number of bytes must be less than or equal to k*2")
	}
	if parityShards > (dataShards / 2) {
		return nil, errors.New("k must be less than or equal to n/2")
	}

	// Calculate the length of the data buffer, which must be a power of 2.
	bytesLen := uint(len(bytes))
	var bufferLen uint
	if isPowerOfTwo(bytesLen) {
		bufferLen = bytesLen
	} else {
		lLocal := int(math.Log2(float64(bytesLen)))
		lLocal = 1 << lLocal
		if uint(lLocal) >= bytesLen {
			bufferLen = uint(lLocal)
		} else {
			bufferLen = uint(lLocal << 1)
		}
	}

	if !isPowerOfTwo(bufferLen) {
		return nil, fmt.Errorf("bufferLen is not power of 2")
	}
	if bufferLen < bytesLen {
		return nil, fmt.Errorf("bufferLen can not be less than bytesLen")
	}

	// Pad the incoming bytes with trailing 0s.
	zeroBytesToAdd := dataShards*2 - bytesLen

	bytesWithPadding := bytes[:]
	for i := uint(0); i < zeroBytesToAdd; i++ {
		bytesWithPadding = append(bytesWithPadding, 0)
	}

	data := []Additive{}
	for i := 0; i < len(bytesWithPadding)-1; i = i + 2 {
		data = append(data, Additive(bytes[i])<<8|Additive(bytes[i+1]))
	}

	// let l = data.len();
	// assert_eq!(l, n);
	bufferLen = uint(len(data))
	if bufferLen != dataShards {
		return nil, fmt.Errorf("buffer len and number of data shards must be same.")
	}

	firstHalfOfData := data[:parityShards]
	secondHalfOfData := data[parityShards:]

	// Encode the data.
	codeword := make([]Additive, bufferLen)
	encodeLow(data, parityShards, codeword, dataShards)

	return codeword, nil
}

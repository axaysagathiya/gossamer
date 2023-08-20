// Copyright 2021 ChainSafe Systems (ON)
// SPDX-License-Identifier: LGPL-3.0-only

package erasure

import (
	"fmt"
	"testing"

	"github.com/ChainSafe/gossamer/lib/common"
	"github.com/ChainSafe/gossamer/lib/parachain"
	"github.com/ChainSafe/gossamer/pkg/scale"
	"github.com/klauspost/reedsolomon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testData = []byte("this is a test of the erasure coding")
var expectedChunks = [][]byte{{116, 104, 105, 115}, {32, 105, 115, 32}, {97, 32, 116, 101}, {115, 116, 32, 111},
	{102, 32, 116, 104}, {101, 32, 101, 114}, {97, 115, 117, 114}, {101, 32, 99, 111}, {100, 105, 110, 103},
	{0, 0, 0, 0}, {133, 189, 154, 178}, {88, 245, 245, 220}, {59, 208, 165, 70}, {127, 213, 208, 179}}

// erasure data missing chunks
var missing2Chunks = [][]byte{{116, 104, 105, 115}, {32, 105, 115, 32}, {}, {115, 116, 32, 111},
	{102, 32, 116, 104}, {101, 32, 101, 114}, {}, {101, 32, 99, 111}, {100, 105, 110, 103},
	{0, 0, 0, 0}, {133, 189, 154, 178}, {88, 245, 245, 220}, {59, 208, 165, 70}, {127, 213, 208, 179}}
var missing3Chunks = [][]byte{{116, 104, 105, 115}, {32, 105, 115, 32}, {}, {115, 116, 32, 111},
	{}, {101, 32, 101, 114}, {}, {101, 32, 99, 111}, {100, 105, 110, 103}, {0, 0, 0, 0}, {133, 189, 154, 178},
	{88, 245, 245, 220}, {59, 208, 165, 70}, {127, 213, 208, 179}}
var missing5Chunks = [][]byte{{}, {}, {}, {115, 116, 32, 111},
	{}, {101, 32, 101, 114}, {}, {101, 32, 99, 111}, {100, 105, 110, 103}, {0, 0, 0, 0}, {133, 189, 154, 178},
	{88, 245, 245, 220}, {59, 208, 165, 70}, {127, 213, 208, 179}}

func TestObtainChunks(t *testing.T) {
	t.Parallel()
	type args struct {
		validatorsQty int
		data          []byte
	}
	tests := map[string]struct {
		args          args
		expectedValue [][]byte
		expectedError error
	}{
		"happy_path": {
			args: args{
				validatorsQty: 10,
				data:          testData,
			},
			expectedValue: expectedChunks,
		},
		"nil_data": {
			args: args{
				validatorsQty: 10,
				data:          nil,
			},
			expectedError: reedsolomon.ErrShortData,
		},
		"not_enough_validators": {
			args: args{
				validatorsQty: 1,
				data:          testData,
			},
			expectedError: ErrNotEnoughValidators,
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, err := ObtainChunks(tt.args.validatorsQty, tt.args.data)
			expectedThreshold, _ := recoveryThreshold(tt.args.validatorsQty)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.args.validatorsQty+expectedThreshold, len(got))
			}
			assert.Equal(t, tt.expectedValue, got)
		})
	}
}

func TestReconstruct(t *testing.T) {
	t.Parallel()
	type args struct {
		validatorsQty int
		chunks        [][]byte
	}
	tests := map[string]struct {
		args
		expectedData   []byte
		expectedChunks [][]byte
		expectedError  error
	}{
		"missing_2_chunks": {
			args: args{
				validatorsQty: 10,
				chunks:        missing2Chunks,
			},
			expectedData:   testData,
			expectedChunks: expectedChunks,
		},
		"missing_2_chunks,_validator_qty_3": {
			args: args{
				validatorsQty: 3,
				chunks:        missing2Chunks,
			},
			expectedError:  reedsolomon.ErrTooFewShards,
			expectedChunks: expectedChunks,
		},
		"missing_3_chunks": {
			args: args{
				validatorsQty: 10,
				chunks:        missing3Chunks,
			},
			expectedData:   testData,
			expectedChunks: expectedChunks,
		},
		"missing_5_chunks": {
			args: args{
				validatorsQty: 10,
				chunks:        missing5Chunks,
			},
			expectedChunks: missing5Chunks,
			expectedError:  reedsolomon.ErrTooFewShards,
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			data, err := Reconstruct(tt.args.validatorsQty, len(testData), tt.args.chunks)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedChunks, tt.args.chunks)
			}
			assert.Equal(t, tt.expectedData, data)
		})
	}
}

func TestChunksToTrie(t *testing.T) {

	qty := 2
	availableData := parachain.AvailableData{
		PoV: parachain.PoV{
			BlockData: parachain.BlockData{2},
		},
		ValidationData: parachain.PersistedValidationData{
			ParentHead:             []byte{},
			RelayParentNumber:      0,
			RelayParentStorageRoot: common.Hash{},
			MaxPovSize:             0,
		},
	}
	dataBytes := scale.MustMarshal(availableData)

	chunks, err := ObtainChunks(qty, dataBytes)
	require.NoError(t, err)

	fmt.Printf("\n\nchunks =>\n%+v\n\n", chunks)

}

func newObtainChunks(validatorsQty int, originalData []byte) ([][]byte, error) {
	dataShards := 2
	parityShards := 1
	enc, err := reedsolomon.New(dataShards, parityShards)
	if err != nil {
		// Handle error
		return nil, err
	}

	// Calculate the size of each data chunk
	chunkSize := (len(originalData) + dataShards - 1) / dataShards // Round up division

	// Split the data into individual chunks
	var dataChunks [][]byte
	for i := 0; i < dataShards; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > len(originalData) {
			end = len(originalData)
		}
		dataChunks = append(dataChunks, originalData[start:end])
	}

	// Encode the data to generate the parity shard
	err = enc.Encode(dataChunks)
	if err != nil {
		// Handle error
		return nil, err
	}

	// Now, the encoded variable will contain the data chunks and the parity shard
	fmt.Printf("**************************************************\n")
	for i, chunk := range dataChunks {
		fmt.Printf("Data Chunk %d: %v\n", i+1, chunk)
	}

	parityChunk := dataChunks[dataShards]
	fmt.Printf("Parity Chunk: %v\n", parityChunk)
	fmt.Printf("**************************************************\n")
	return dataChunks, nil
}

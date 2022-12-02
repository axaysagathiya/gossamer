// Copyright 2021 ChainSafe Systems (ON)
// SPDX-License-Identifier: LGPL-3.0-only

package node

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Node_String(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		node *Node
		s    string
	}{
		"leaf with storage value smaller than 1024": {
			node: &Node{
				PartialKey:   []byte{1, 2},
				StorageValue: []byte{3, 4},
				Dirty:        true,
			},
			s: `Leaf
├── Generation: 0
├── Dirty: true
├── Key: 0x0102
├── Storage value: 0x0304
└── Merkle value: nil`,
		},
		"leaf with storage value higher than 1024": {
			node: &Node{
				PartialKey:   []byte{1, 2},
				StorageValue: make([]byte, 1025),
				Dirty:        true,
			},
			s: `Leaf
├── Generation: 0
├── Dirty: true
├── Key: 0x0102
├── Storage value: 0x0000000000000000...0000000000000000
└── Merkle value: nil`,
		},
		"branch with storage value smaller than 1024": {
			node: &Node{
				PartialKey:   []byte{1, 2},
				StorageValue: []byte{3, 4},
				Dirty:        true,
				Descendants:  3,
				Children: []*Node{
					nil, nil, nil,
					{},
					nil, nil, nil,
					{
						Descendants: 1,
						Children:    padRightChildren([]*Node{{}}),
					},
					nil, nil, nil,
					{},
					nil, nil, nil, nil,
				},
			},
			s: `Branch
├── Generation: 0
├── Dirty: true
├── Key: 0x0102
├── Storage value: 0x0304
├── Descendants: 3
├── Merkle value: nil
├── Child 3
|   └── Leaf
|       ├── Generation: 0
|       ├── Dirty: false
|       ├── Key: nil
|       ├── Storage value: nil
|       └── Merkle value: nil
├── Child 7
|   └── Branch
|       ├── Generation: 0
|       ├── Dirty: false
|       ├── Key: nil
|       ├── Storage value: nil
|       ├── Descendants: 1
|       ├── Merkle value: nil
|       └── Child 0
|           └── Leaf
|               ├── Generation: 0
|               ├── Dirty: false
|               ├── Key: nil
|               ├── Storage value: nil
|               └── Merkle value: nil
└── Child 11
    └── Leaf
        ├── Generation: 0
        ├── Dirty: false
        ├── Key: nil
        ├── Storage value: nil
        └── Merkle value: nil`,
		},
		"branch with storage value higher than 1024": {
			node: &Node{
				PartialKey:   []byte{1, 2},
				StorageValue: make([]byte, 1025),
				Dirty:        true,
				Descendants:  3,
				Children: []*Node{
					nil, nil, nil,
					{},
					nil, nil, nil,
					{
						Descendants: 1,
						Children:    padRightChildren([]*Node{{}}),
					},
					nil, nil, nil,
					{},
					nil, nil, nil, nil,
				},
			},
			s: `Branch
├── Generation: 0
├── Dirty: true
├── Key: 0x0102
├── Storage value: 0x0000000000000000...0000000000000000
├── Descendants: 3
├── Merkle value: nil
├── Child 3
|   └── Leaf
|       ├── Generation: 0
|       ├── Dirty: false
|       ├── Key: nil
|       ├── Storage value: nil
|       └── Merkle value: nil
├── Child 7
|   └── Branch
|       ├── Generation: 0
|       ├── Dirty: false
|       ├── Key: nil
|       ├── Storage value: nil
|       ├── Descendants: 1
|       ├── Merkle value: nil
|       └── Child 0
|           └── Leaf
|               ├── Generation: 0
|               ├── Dirty: false
|               ├── Key: nil
|               ├── Storage value: nil
|               └── Merkle value: nil
└── Child 11
    └── Leaf
        ├── Generation: 0
        ├── Dirty: false
        ├── Key: nil
        ├── Storage value: nil
        └── Merkle value: nil`,
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := testCase.node.String()

			assert.Equal(t, testCase.s, s)
		})
	}
}

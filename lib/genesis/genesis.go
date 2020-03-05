// Copyright 2019 ChainSafe Systems (ON) Corp.
// This file is part of gossamer.
//
// The gossamer library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The gossamer library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the gossamer library. If not, see <http://www.gnu.org/licenses/>.

package genesis

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/ChainSafe/gossamer/lib/common"
)

// Genesis stores the data parsed from the genesis configuration file
type Genesis struct {
	Name       string
	ID         string
	Bootnodes  []string
	ProtocolID string
	Genesis    Fields
}

// Data defines the genesis file data formatted for trie storage
type Data struct {
	Name       string
	ID         string
	Bootnodes  [][]byte
	ProtocolID string
}

// Fields stores genesis raw data
type Fields struct {
	Raw [2]map[string]string
}

// LoadGenesisFromJSON parses a JSON formatted genesis file
func LoadGenesisFromJSON(file string) (*Genesis, error) {
	fp, err := filepath.Abs(file)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(filepath.Clean(fp))
	if err != nil {
		return nil, err
	}

	g := new(Genesis)
	err = json.Unmarshal(data, g)
	return g, err
}

// GenesisData formats genesis for trie storage
func (g *Genesis) GenesisData() *Data {
	return &Data{
		Name:       g.Name,
		ID:         g.ID,
		Bootnodes:  common.StringArrayToBytes(g.Bootnodes),
		ProtocolID: g.ProtocolID,
	}
}

// GenesisFields returns the genesis fields including genesis raw data
func (g *Genesis) GenesisFields() Fields {
	return g.Genesis
}
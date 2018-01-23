// Copyright 2018 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package observer

import (
	"testing"

	"github.com/ethereum/go-ethereum/ethdb"
)

func TestNewChainHasNoFirstBlock(t *testing.T) {
	testdb, _ := ethdb.NewMemDatabase()
	c, err := NewChain(testdb)
	if err != nil {
		t.Errorf("NewChain() error = %v", err)
		return
	}
	b, err := c.Block(0)
	if err != nil {
		t.Errorf("Block retrieval returned error %s", err)
	}
	if b == nil {
		t.Errorf("Empty chain has a block error")
	}
}

// func TestNewChain(t *testing.T) {

// 	var testdb, _ = ethdb.NewMemDatabase()

// 	type args struct {
// 		db ethdb.Database
// 	}

// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    int64
// 		wantErr bool
// 	}{
// 		{"No initial blocks", args{db: testdb}, 1, false},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			c, err := NewChain(tt.args.db)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("NewChain() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			b, _ := c.Block(0)
// 			if b != nil {
// 				t.Errorf("NewChain() has block :(")
// 			}
// 			//if !reflect.DeepEqual(got, tt.want) {
// 			//	t.Errorf("NewChain() = %v, want %v", got, tt.want)
// 			//}
// 		})
// 	}
// }

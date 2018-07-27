package types

import (
	"datx_chain/utils/common"
	"datx_chain/utils/crypto"
)

type MerkleTree struct {
	nodeCount  int64
	activeNode []common.Hash
}

//Merkle is to calculate the hash
func Merkle(ids []common.Hash) common.Hash {
	if ids == nil {
		return common.HexToHash("")
	}
	for len(ids) > 1 {
		if len(ids)%2 == 0 {
			ids = append(ids, ids[len(ids)-1])
		}
		for i := 0; i < len(ids); i++ {
			ids[i] = crypto.Keccak256Hash(ids[2*i].Bytes(), ids[(2*i)+1].Bytes())
		}
	}
	return ids[0]
}

//given an unsigned integral number return the smallest power-of-2 which is greater than or equal to the given number
func NextPowerOf2(value int64) int64 {
	value--
	value |= value >> 1
	value |= value >> 2
	value |= value >> 4
	value |= value >> 8
	value |= value >> 16
	value |= value >> 32
	value++
	return value
}

//GetMroot is to calculate the root of the merkle tree
func (mt *MerkleTree) GetMroot() common.Hash {
	if mt.nodeCount > 0 {
		return mt.activeNode[len(mt.activeNode)-1]
	}
	return common.HexToHash("")
}

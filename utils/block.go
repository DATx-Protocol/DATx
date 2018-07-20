package utils

import (
	"DATx/utils/common"
	"DATx/utils/crypto/sha3"
	"DATx/utils/rlp"
)

type BlockHeader struct {
	//current block num
	BlockNum uint32

	//current block header hash
	ID common.Hash

	//Previous block id
	Previous uint32

	//producer account
	Producer string

	//producer signature
	Signature string

	//the root of action merkle tree
	ActionMroot common.Hash

	//the root of transaction merkle tree
	TransactionMroot common.Hash
}

type Block struct {
	//block header
	BlockHeader

	//transcation pool
	Transactions []*Transaction
}

func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

func MakeBlockHeader(num, pre uint32, prod string) *BlockHeader {
	return &BlockHeader{
		BlockNum: num,
		Previous: pre,
		Producer: prod,
	}
}

func (self *BlockHeader) Hash() common.Hash {
	self.ID = rlpHash(&self)
	return self.ID
}

func MakeBlock(num, pre uint32, prod string) *Block {
	return &Block{
		BlockHeader: BlockHeader{
			BlockNum: num,
			Previous: pre,
			Producer: prod,
		},
	}
}

func (self *Block) GetID() *common.Hash {
	return &self.ID
}

func (self *Block) GetNum() uint32 {
	return self.BlockNum
}

func (self *Block) GetPrevious() uint32 {
	return self.Previous
}

func (self *Block) GetHead() *BlockHeader {
	return &self.BlockHeader
}

// func (self *Block) TransactionMerkle() common.Hash {
// 	var ids = self.Transactions
// 	if ids == nil {
// 		return common.HexToHash("")
// 	}
// 	for len(ids) > 1 {
// 		if len(ids)%2 == 0 {
// 			ids = append(ids, ids[len(ids)-1])
// 		}
// 		for i := 0; i < len(ids)/2; i++ {
// 			ids[i].TransactionHash = crypto.Keccak256Hash(ids[2*i].TransactionHash.Bytes(), ids[(2*i)+1].TransactionHash.Bytes())
// 		}
// 	}
// 	return ids[0].TransactionHash
// }

// func (self *Block) GetTransactionMroot() common.Hash {
// 	self.TransactionMroot = self.TransactionMerkle()
// 	return self.TransactionMroot
// }

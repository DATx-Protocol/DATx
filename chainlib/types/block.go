package types

import (
	"datx_chain/chainlib/chainobject"
	"datx_chain/utils/common"
	"encoding/binary"
	"math/big"
)

//transaction status
const (
	Executed uint8 = iota //succeed,no error
	SoftFail              //objectively failed (not executed), error handler executed
	HardFail              //objectively failed and error handler objectively failed thus no state change
	Delayed               //transaction delayed/deferred/scheduled for future execution
	Expired               //transaction expired and storage space refuned to user
)

//TransactionReceiptHeader receipt header
type TransactionReceiptHeader struct {
	Status        uint8
	CPUUsageUS    uint32 //total billed CPU usage
	NetUsageWords uint   //total billed NET usage, so we can reconstruct resource state when skipping context free data... hard failures...
}

//TransactionReceipt struct
type TransactionReceipt struct {
	TransactionReceiptHeader

	TrxID common.Hash

	PackedTrx *PackedTransaction
}

//NewTrxReceiptID new trx receipt by trx id
func NewTrxReceiptID(id common.Hash) *TransactionReceipt {
	return &TransactionReceipt{
		TransactionReceiptHeader: TransactionReceiptHeader{Status: Executed},
		TrxID:     id,
		PackedTrx: nil,
	}
}

//NewTrxReceiptPacked new trx receipt by packed trx
func NewTrxReceiptPacked(pt *PackedTransaction) *TransactionReceipt {
	return &TransactionReceipt{
		TransactionReceiptHeader: TransactionReceiptHeader{Status: Executed},
		TrxID:     common.Hash{},
		PackedTrx: pt,
	}
}

type BlockIdType common.Hash

//Block struct
type Block struct {
	//block header
	BlockHeader

	//transcation pool
	Transactions []*TransactionReceipt
}

//NewBlock new
func NewBlock(h BlockHeader) *Block {
	return &Block{
		BlockHeader:  h,
		Transactions: make([]*TransactionReceipt, 0),
	}
}

func NewBlockEmpty() *Block {
	return &Block{
		BlockHeader: BlockHeader{
			Signature:    make([]byte, 0),
			TimeStamp:    &BlockTime{Time: big.NewInt(0)},
			NewProducers: chainobject.ProducerSchedule{Producers: make([]chainobject.ProducerKey, 0)},
		},
		Transactions: make([]*TransactionReceipt, 0),
	}
}

// //TransactionMerkle merkle
// func (block *Block) TransactionMerkle() common.Hash {
// 	var ids = block.Transactions
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

// //GetTransactionMroot get
// func (block *Block) GetTransactionMroot() common.Hash {
// 	block.TransactionMroot = block.TransactionMerkle()
// 	return block.TransactionMroot
// }

//GetNum get block num
func (block *Block) GetNum() []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, block.BlockHeader.BlockNum)

	return b
}

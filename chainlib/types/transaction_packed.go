package types

import (
	"bytes"
	"compress/zlib"
	"datx_chain/utils/common"
	"datx_chain/utils/helper"
	"datx_chain/utils/rlp"
	"io"
)

//enum compression type
const (
	None uint8 = iota
	Zlib
)

//PackedTransaction struct
type PackedTransaction struct {
	signatures []bdata

	Compression uint8 //compression type

	PackedContextFreeData []byte

	PackedTrx []byte

	unPackedTrx *Transaction
}

//NewPackedTransaction new
func NewPackedTransaction(t *SignedTransaction, compression uint8) *PackedTransaction {
	var res PackedTransaction
	res.signatures = t.Signatures

	res.SetTransaction(&t.Transaction, t.ContextFreeData, compression)
	return &res
}

//Expiration time
func (pt PackedTransaction) Expiration() uint64 {
	return pt.unPackedTrx.Expiration
}

//ID hash
func (pt *PackedTransaction) ID() common.Hash {
	pt.localUnpack()
	trx := pt.getTransaction()
	return trx.ID()
}

//SetTransaction set transaction
func (pt *PackedTransaction) SetTransaction(t *Transaction, tbytes []byte, compression uint8) {
	errHandle := func() {
		panic("[PackedTransaction] SetTransaction failed")
	}

	helper.CatchException(nil, errHandle)

	switch compression {
	case None:
		pt.PackedTrx, _ = rlp.EncodeToBytes(t)
		pt.PackedContextFreeData = tbytes
	case Zlib:
		pt.PackedTrx = zlibCompressionTransaction(t)
		pt.PackedContextFreeData = zlibCompressionContextFreeData(tbytes)
	default:
		panic("Unknown transaction compression algorithm")
	}

	pt.Compression = compression
}

func (pt *PackedTransaction) localUnpack() {
	helper.CatchException(nil, func() {
		panic("[PackedTransaction] localUnpack failed")
	})

	if pt.unPackedTrx == nil {
		switch pt.Compression {
		case None:
			pt.unPackedTrx = unpackTransaction(pt.PackedTrx)
		case Zlib:
			pt.unPackedTrx = zlibUnCompressionTransaction(pt.PackedTrx)
		default:
			panic("Unknown transaction compression algorithm")
		}
	}
}

func (pt *PackedTransaction) getTransaction() *Transaction {
	pt.localUnpack()
	return pt.unPackedTrx
}

//GetSignTransaction get signed trx
func (pt *PackedTransaction) GetSignTransaction() SignedTransaction {
	errHandle := func() {
		panic("[PackedTransaction] GetSignTransaction failed")
	}

	helper.CatchException(nil, errHandle)

	var result SignedTransaction
	result.Signatures = pt.signatures
	if pt.unPackedTrx != nil {
		result.Transaction = *pt.unPackedTrx
	}

	switch pt.Compression {
	case None:
		result.ContextFreeData = pt.PackedContextFreeData
	case Zlib:
		result.ContextFreeData = zlibUnCompressionContextFreeData(pt.PackedContextFreeData)
	default:
		panic("Unknown transaction compression algorithm")
	}
	return result
}

func unpackTransaction(src []byte) *Transaction {
	var t Transaction
	if err := rlp.DecodeBytes(src, &t); err != nil {
		return nil
	}

	return &t
}

func zlibCompressionTransaction(t *Transaction) []byte {
	data, err := rlp.EncodeToBytes(t)
	if err != nil {
		return nil
	}

	//zlib compress
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write(data)
	w.Close()
	return in.Bytes()
}

func zlibCompressionContextFreeData(src []byte) []byte {
	//zlib compress
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

func zlibUnCompressionContextFreeData(src []byte) []byte {
	//zlib decompress
	var out bytes.Buffer
	w := bytes.NewReader(src)
	r, _ := zlib.NewReader(w)
	io.Copy(&out, r)
	return out.Bytes()
}

func zlibUnCompressionTransaction(src []byte) *Transaction {
	//zlib decompress
	var out bytes.Buffer
	w := bytes.NewReader(src)
	r, _ := zlib.NewReader(w)
	io.Copy(&out, r)
	return unpackTransaction(out.Bytes())
}

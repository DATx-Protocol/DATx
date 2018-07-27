package controller

import (
	"datx_chain/chainlib/types"
	"datx_chain/utils/db"
	"datx_chain/utils/helper"
	"datx_chain/utils/rlp"
	"encoding/binary"
	"log"
	"os"
)

const (
	//SupportedVersion version
	SupportedVersion uint32 = 1
)

//Blog is to save block data in chain
type Blog struct {
	file string

	db *datxdb.LDBDatabase //levelDB instance

	batch datxdb.Batch //levelDB batch

	head *types.Block
}

//NewBlog create Blog
func NewBlog(abspath string, cache, handles int) (b *Blog, err error) {
	//catch exception,return nil and error
	helper.CatchException(err, func() {
		b = nil
	})

	var res Blog

	res.file = abspath + string(os.PathSeparator) + "block_log"

	//check for the directory's existence and create it if it doesn't exist
	if err := os.MkdirAll(res.file, os.ModePerm); err != nil {
		return nil, err
	}

	//open levelDB. the db will recover when db crashed or file is existence
	if res.db, err = datxdb.NewLDBDatabase(res.file, cache, handles); err != nil {
		log.Printf("NewBlog err #{%s}", err)
		return nil, err
	}

	//create db batch for write batched
	res.batch = res.db.NewBatch()

	//unmarshal index map
	iter := res.db.NewIterator()
	defer iter.Release()

	//unmarshal last block to head
	if iter.Last() == true && iter.Valid() {
		//var block types.Block
		block := types.NewBlockEmpty()
		err := rlp.DecodeBytes(iter.Value(), block)
		if err != nil {
			log.Printf("Blog::NewBlog err={%v} ", err)
			return nil, err
		}
		res.head = block
	}

	return &res, nil
}

func (blog *Blog) Write(b *types.Block) {
	//encode data
	data, err := rlp.EncodeToBytes(b)
	if err != nil {
		log.Printf("Blog::Write encode block num={%d} err={%v}\n", b.BlockNum, err)
		return
	}

	//batch write
	err = blog.batch.Put(b.GetNum(), data)
	if err != nil {
		log.Printf("Blog::Write batch.Put block num={%d} err={%v}", b.BlockNum, err)
		return
	}

	err = blog.batch.Write()
	if err != nil {
		log.Printf("Blog::Write batch.Write block num={%d} err={%v}", b.BlockNum, err)
		return
	}

	//update head
	blog.head = b

}

//Clear clear data in levelDB
func (blog *Blog) Clear() {
	iter := blog.db.NewIterator()
	defer iter.Release()

	//clear all data in levelDB
	for iter.First(); iter.Valid(); iter.Next() {
		key := iter.Key()
		blog.db.Delete(key)
	}
}

//Close close the levelDB
func (blog *Blog) Close() {
	blog.db.Close()
}

//ReadByBlockNum read block by block num from levelDB
func (blog *Blog) ReadByBlockNum(num uint32) *types.Block {
	bnum := make([]byte, 4)
	binary.BigEndian.PutUint32(bnum, num)

	//if not exist,return nil
	data, err := blog.db.Get(bnum)
	if err != nil {
		log.Printf("Blog::ReadByBlockNum block num={%d} err={%v}", num, err)
		return nil
	}

	var block types.Block
	if err := rlp.DecodeBytes(data, &block); err != nil {
		log.Printf("Blog::ReadByBlockNum decode block num={%v} err={%v}", num, err)
		return nil
	}

	return &block
}

//ReadHead read first block in levelDB
func (blog *Blog) ReadHead() *types.Block {
	iter := blog.db.NewIterator()
	defer iter.Release()

	var block types.Block

	//unmarshal last block to head
	if iter.Last() == true && iter.Valid() {
		err := rlp.DecodeBytes(iter.Value(), &block)
		if err != nil {
			log.Printf("Blog::NewBlog err={%v} ", err)
			return nil
		}

	}

	return &block
}

//Head return head
func (blog *Blog) Head() *types.Block {
	return blog.head
}

//ResetGenesis set genesis block
func (blog *Blog) ResetGenesis(b *types.Block) {
	//clear
	blog.Clear()

	//add data
	blog.Write(b)
}

// func (blog *Blog) ExtractGenesisState(path string) GenesisState {

// }

package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

// 迭代器

type BlockChainIterator struct {
	DB        *bolt.DB
	Curr_Hash []byte
}

// 返回迭代器对象
func (bc *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{bc.DB, bc.Top}
}

// 遍历
func (it *BlockChainIterator) Next() *Block {
	var block *Block
	err := it.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BLOCKTABLENAME))
		if b != nil {
			// 获取指定hash区块
			cur_BlockBytes := b.Get(it.Curr_Hash)
			block = DeserializeBlock(cur_BlockBytes)
			//log.Println("current height: ", block.Height)
			// update hash info
			it.Curr_Hash = block.Pre_Hash
		} else {
			log.Panicf("access the bucket error!")
		}
		return nil
	})
	if err != nil {
		log.Panicf("iterator the db failed. %v\n", err)
	}
	return block
}

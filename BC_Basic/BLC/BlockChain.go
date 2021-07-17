package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

// 区块链
type BlockChain struct {
	DB  *bolt.DB // 数据库
	Top []byte   //最新区块的hash值
}

// 由创世区块初始化区块链

func CreateBlockChainWithGenesisBlock() *BlockChain {
	// 打开数据库
	db, err := bolt.Open(DBNAME, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()

	var blockHash []byte

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BLOCKTABLENAME))
		// 如果为空则需要创建Bucket
		if b == nil {
			b, err = tx.CreateBucket([]byte(BLOCKTABLENAME))
			if err != nil {
				log.Panicf("Create the bucket [%s] failed. %v\n", BLOCKTABLENAME, err)
			}
		}
		if b != nil {
			// 不为空则创建创世区块
			genesisBlock := CreateGenesisBlock("the genesis block")
			// 用hash值作为唯一key
			err := b.Put(genesisBlock.Self_Hash, genesisBlock.Serialize())
			if err != nil {
				log.Panicf("put the data failed %v\n", err)
			}

			// 存储最新区块的hash值
			err = b.Put([]byte(LATEST_HASH), genesisBlock.Self_Hash)
			if err != nil {
				log.Panicf("put the hash failed %v\n", err)
			}
			blockHash = genesisBlock.Self_Hash
		}
		return nil
	})

	if err != nil {
		log.Panicf("update the db failed. %v\n", err)
	}

	return &BlockChain{db, blockHash}
}

// 添加新的区块
func (bc *BlockChain) AddBlock(data []byte) {
	//newBlock := NewBlock(height, pre_hash, data)
	//bc.Blocks = append(bc.Blocks, newBlock)

	err := bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BLOCKTABLENAME))
		if b != nil {
			// 获取最新块的hash
			latest_hash := b.Get([]byte(LATEST_HASH))
			blockBytes := b.Get(latest_hash)
			latest_Block := DeserializeBlock(blockBytes)
			// 创建新区块
			newblock := NewBlock(latest_Block.Height+1, latest_Block.Pre_Hash, data)
			// 存入数据库
			err := b.Put(newblock.Self_Hash, newblock.Serialize())
			if err != nil {
				log.Panicf("Put latest block to db failed %v\n", err)
			}
			// 更新最新区块hash
			err = b.Put([]byte(LATEST_HASH), newblock.Self_Hash)
			if err != nil {
				log.Panicf("Put latest hash to db failed %v\n", err)
			}
		}

		return nil
	})

	if err != nil {
		log.Panicf("update the db failed. %v\n", err)
	}
}

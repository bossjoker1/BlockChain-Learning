package BLC

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
	"strconv"
)

// 区块链
type BlockChain struct {
	DB  *bolt.DB // 数据库
	Top []byte   //最新区块的hash值
}

// 判断数据库文件是否存在
func dbExists() bool {
	if _, err := os.Stat(DBNAME); os.IsNotExist(err) {
		return false
	}
	return true
}

// 由创世区块初始化区块链
func CreateBlockChainWithGenesisBlock(addr string) *BlockChain {

	if dbExists() == true {
		fmt.Println("genesis block has existed.")
		//
		//
		//var blockChain *BlockChain
		//// 取出存在的区块
		//
		//err = db.View(func(tx *bolt.Tx) error {
		//	b := tx.Bucket([]byte(BLOCKTABLENAME))
		//	hash := b.Get([]byte(LATEST_HASH))
		//	blockChain = &BlockChain{db, hash}
		//	return nil
		//})
		//if err != nil{
		//	log.Panicf("get the block from db failed %v\n", err)
		//}
		//return blockChain
		os.Exit(1)
	}
	//fmt.Println(dbExists())

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

			// 生成交易
			txCoinbase := NewCoinBaseTX(addr)

			// 不为空则创建创世区块
			genesisBlock := CreateGenesisBlock([]*TX{txCoinbase})
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
func (bc *BlockChain) AddBlock(txs []*TX) {
	//newBlock := NewBlock(height, pre_hash, data)
	//bc.Blocks = append(bc.Blocks, newBlock)

	err := bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BLOCKTABLENAME))
		if b != nil {
			// 获取最新块的hash
			// latest_hash := b.Get([]byte(LATEST_HASH))
			blockBytes := b.Get(bc.Top)
			latest_Block := DeserializeBlock(blockBytes)
			// 创建新区块
			newblock := NewBlock(latest_Block.Height+1, latest_Block.Self_Hash, txs)
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
			bc.Top = newblock.Self_Hash
		}

		return nil
	})

	if err != nil {
		log.Panicf("update the db failed. %v\n", err)
	}
}

// 从最新区块开始倒序打印

func (bc *BlockChain) PrintChainInfo() {
	var curBlock *Block

	// 创建对象
	it := bc.Iterator()

	for {
		fmt.Println("-------------------------------------------------------------------------------")

		curBlock = it.Next()
		fmt.Printf("\tHeight :		%d\n", curBlock.Height)
		fmt.Printf("\tTimeStamp:		%d\n", curBlock.TimeStamp)
		fmt.Printf("\tPre_Hash:		%x\n", curBlock.Pre_Hash)
		fmt.Printf("\tSelf_Hash:		%x\n", curBlock.Self_Hash)
		fmt.Printf("\tnonce:			%d\n", curBlock.Nonce)
		fmt.Printf("\tTransaction:	\n")
		for i, tx := range curBlock.Txs {
			fmt.Printf("\t\tTX_hash_%d:	%x\n", i+1, tx.Tx_hash)
			fmt.Println("\t\tinput......")
			for _, tins := range tx.Tins {
				fmt.Printf("\t\t\tTX_in_hash:	%x\n", tins.Tx_hash)
				fmt.Printf("\t\t\tTX_index_out:	%d\n", tins.Index_out)
				fmt.Printf("\t\t\tTX_scriptsig:	%s\n", tins.ScriptSig)
			}
			fmt.Println("\t\toutput......")
			for _, touts := range tx.Touts {
				fmt.Printf("\t\t\tTX_out_values:	%d\n", touts.Value)
				fmt.Printf("\t\t\tTX_ScriptPubkey: %s\n", touts.ScriptPubkey)
			}
		}

		// 遍历到创世区块退出，Pre_hash为空
		//var hashInt big.Int
		//hashInt.SetBytes(curBlock.Pre_Hash)
		//if big.NewInt(0).Cmp(&hashInt) == 0{
		//	break
		//}
		if curBlock.Pre_Hash == nil {
			break
		}
	}
}

// 获得区块链对象
func GetBCObject() *BlockChain {
	// 读取数据库
	db, err := bolt.Open(DBNAME, 0600, nil)

	if err != nil {
		log.Panicf("get the object failed. %v\n", err)
	}

	var top_hash []byte

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BLOCKTABLENAME))

		if b != nil {
			top_hash = b.Get([]byte(LATEST_HASH))
		}
		return nil
	})
	return &BlockChain{db, top_hash}
}

// 挖矿函数, 生成新的区块，不同于AddBlock
// 通过接受交易，进行打包确认(PoW)，然后生成新的区块
func (bc *BlockChain) MineNewBlock(from, to, amount []string) {
	// 接受交易
	var txs []*TX
	value, _ := strconv.Atoi(amount[0])
	tx := NewSimpleTX(from[0], to[0], int64(value))
	// 多笔交易只是个 for 循环的事情
	txs = append(txs, tx)
	// 打包交易

	var block *Block
	// 从db中获取最新数据库
	bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BLOCKTABLENAME))
		if b != nil {
			hash := b.Get([]byte(LATEST_HASH))
			blockBytes := b.Get(hash) //  ? -> 为了获取区块高度
			block = DeserializeBlock(blockBytes)
		}
		return nil
	})

	// 生成新的区块
	block = NewBlock(block.Height+1, block.Self_Hash, txs)

	// 持久化新区块
	bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BLOCKTABLENAME))
		if b != nil {
			err := b.Put(block.Self_Hash, block.Serialize())
			if err != nil {
				log.Panicf("update the new block failed. %v\n", err)
			}
			b.Put([]byte(LATEST_HASH), block.Self_Hash)
			bc.Top = block.Self_Hash
		}
		return nil
	})
}

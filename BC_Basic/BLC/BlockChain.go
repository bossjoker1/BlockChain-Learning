package BLC

import (
	"BlockChain-Learning/BC_Basic/Utils"
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"os"
	"strconv"
)

// 区块链
type BlockChain struct {
	DB  *bolt.DB // 数据库
	Top []byte   //最新区块的hash值
}

// 判断数据库文件是否存在
func DbExists(nodeId string) bool {
	dbName := fmt.Sprintf(Utils.DBNAME, nodeId)
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		return false
	}
	return true
}

// 由创世区块初始化区块链
func CreateBlockChainWithGenesisBlock(addr string, nodeId string) *BlockChain {

	if DbExists(nodeId) == true {
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
	dbName := fmt.Sprintf(Utils.DBNAME, nodeId)
	// 打开数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()

	var blockHash []byte

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(Utils.BLOCKTABLENAME))
		// 如果为空则需要创建Bucket
		if b == nil {
			b, err = tx.CreateBucket([]byte(Utils.BLOCKTABLENAME))
			if err != nil {
				log.Panicf("Create the bucket [%s] failed. %v\n", Utils.BLOCKTABLENAME, err)
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
			err = b.Put([]byte(Utils.LATEST_HASH), genesisBlock.Self_Hash)
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

//// 添加新的区块
//func (bc *BlockChain) AddBlock(txs []*TX) {
//	//newBlock := NewBlock(height, pre_hash, data)
//	//bc.Blocks = append(bc.Blocks, newBlock)
//
//	err := bc.DB.Update(func(tx *bolt.Tx) error {
//		b := tx.Bucket([]byte(BLOCKTABLENAME))
//		if b != nil {
//			// 获取最新块的hash
//			// latest_hash := b.Get([]byte(LATEST_HASH))
//			blockBytes := b.Get(bc.Top)
//			latest_Block := DeserializeBlock(blockBytes)
//			// 创建新区块
//			newblock := NewBlock(latest_Block.Height+1, latest_Block.Self_Hash, txs)
//			// 存入数据库
//			err := b.Put(newblock.Self_Hash, newblock.Serialize())
//			if err != nil {
//				log.Panicf("Put latest block to db failed %v\n", err)
//			}
//			// 更新最新区块hash
//			err = b.Put([]byte(LATEST_HASH), newblock.Self_Hash)
//			if err != nil {
//				log.Panicf("Put latest hash to db failed %v\n", err)
//			}
//			bc.Top = newblock.Self_Hash
//		}
//
//		return nil
//	})
//
//	if err != nil {
//		log.Panicf("update the db failed. %v\n", err)
//	}
//}

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
				// fmt.Printf("\t\t\tTX_scriptsig:	%s\n", tins.ScriptSig)
			}
			fmt.Println("\t\toutput......")
			for _, touts := range tx.Touts {
				fmt.Printf("\t\t\tTX_out_values:	%d\n", touts.Value)
				fmt.Printf("\t\t\tTX_ScriptPubkey: %x\n", touts.PubkeyHash)
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
// 也因节点而异
func GetBCObject(nodeId string) *BlockChain {
	dbName := fmt.Sprintf(Utils.DBNAME, nodeId)
	// 读取数据库
	db, err := bolt.Open(dbName, 0600, nil)

	if err != nil {
		log.Panicf("get the object failed. %v\n", err)
	}

	var top_hash []byte

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(Utils.BLOCKTABLENAME))

		if b != nil {
			top_hash = b.Get([]byte(Utils.LATEST_HASH))
		}
		return nil
	})
	return &BlockChain{db, top_hash}
}

// 生成Wallet_3000.dat
// 挖矿函数, 生成新的区块，不同于AddBlock
// 通过接受交易，进行打包确认(PoW)，然后生成新的区块
func (bc *BlockChain) MineNewBlock(from []string, to []string, amount []string, node_id string) {
	// 接收交易
	var txs []*TX
	for index, addr := range from {
		value, _ := strconv.Atoi(amount[index])
		// 生成交易
		tx := NewSimpleTX(addr, to[index], int64(value), bc, txs, &UTXOSet{bc}, node_id)
		// 多笔交易只是个 for 循环的事情
		txs = append(txs, tx)
		// 打包交易
	}

	// 给矿工一定的奖励
	// 默认设置地址列表中的第一个地址为抢夺到几张权的矿工地址
	// 所以这里取from[0]
	tx := NewCoinBaseTX(from[0])
	txs = append(txs, tx)

	var block *Block
	// 从db中获取最新数据库
	_ = bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(Utils.BLOCKTABLENAME))
		if b != nil {
			hash := b.Get([]byte(Utils.LATEST_HASH))
			blockBytes := b.Get(hash) //  ? -> 为了获取区块高度
			block = DeserializeBlock(blockBytes)
		}
		return nil
	})

	var txsc []*TX // 关联交易
	// 在生成新区块之前需要进行验证签名
	for _, tx := range txs {
		// 每一笔交易都得验证
		// 如果后续交易存在前面的关联输入
		// 则对应的前面交易都可能需要前面，因此放入是此时交易的缓存
		fmt.Printf("txHash : %x\n", tx.Tx_hash)
		if !bc.VerifyTX(tx, txs) {
			log.Panicf("error tx [%x] verify failed! \n", tx.Tx_hash)
		}
		txsc = append(txsc, tx)
	}

	// 生成新的区块
	block = NewBlock(block.Height+1, block.Self_Hash, txs)

	// 持久化新区块
	bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(Utils.BLOCKTABLENAME))
		if b != nil {
			err := b.Put(block.Self_Hash, block.Serialize())
			if err != nil {
				log.Panicf("update the new block failed. %v\n", err)
			}
			_ = b.Put([]byte(Utils.LATEST_HASH), block.Self_Hash)
			bc.Top = block.Self_Hash
		}
		return nil
	})
}

// 返回指定地址的余额

func (bc *BlockChain) GetUTXOs(addr string, txs []*TX) []*UTXO {
	fmt.Printf("the address is %s\n", addr)

	// 存储未花费的输出
	var utxos []*UTXO
	// 存储所有已花费的输出
	// key: 每个input所”引用“的hash
	// value: 存 index 用 []int
	spentTXOutput := make(map[string][]int)
	// 查找缓存中(未打包)的交易中是否有该地址的UTXO
	for _, tx := range txs { // 遍历区块中的交易
		// 先查找输入
		// 再通过索引找到相关的output
		if !tx.IsCoinbase() {
			for _, in := range tx.Tins {
				// 验证公钥hash
				publicKeyHash := Base58Decode([]byte(addr))
				ripemd_sha := publicKeyHash[1 : len(publicKeyHash)-Utils.CHECKSUMLEN]
				if in.UnLockWithRipemd_SHA(ripemd_sha) {
					// 添加到已花费输出map中
					key := hex.EncodeToString(in.Tx_hash)
					spentTXOutput[key] = append(spentTXOutput[key], in.Index_out)
				}
			}
		}
		// 查找缓存输出
	workCache:
		for index, tout := range tx.Touts {
			if tout.UnLockPubKeyWithAddr(addr) {
				// 判断是否是未花费的输出
				// 如果已花费输出不为空
				if len(spentTXOutput) != 0 {
					var isSpentTXOutput bool
					for txHash, indexArray := range spentTXOutput {
						for _, i := range indexArray {
							if txHash == hex.EncodeToString(tx.Tx_hash) && i == index {
								// 已被花费
								isSpentTXOutput = true
								// 只要满足即跳到下一个output而不是遍历spentTXoutput中是否有还有
								continue workCache
							}
						}
					}
					// 遍历完所有都没有被花费，才能加入到utxos中
					if isSpentTXOutput == false {
						utxo := &UTXO{tx.Tx_hash, index, tout}
						utxos = append(utxos, utxo)
					}
				} else {
					utxo := &UTXO{tx.Tx_hash, index, tout}
					utxos = append(utxos, utxo)
				}
			}
		}

	}
	blockIterator := bc.Iterator()

	for {
		block := blockIterator.Next()
		for _, tx := range block.Txs { // 遍历区块中的交易
			// 先查找输入
			// 再通过索引找到相关的output
			if !tx.IsCoinbase() {
				for _, in := range tx.Tins {
					// 验证公钥hash
					publicKeyHash := Base58Decode([]byte(addr))
					ripemd_sha := publicKeyHash[1 : len(publicKeyHash)-Utils.CHECKSUMLEN]
					if in.UnLockWithRipemd_SHA(ripemd_sha) {
						// 添加到已花费输出map中
						key := hex.EncodeToString(in.Tx_hash)
						spentTXOutput[key] = append(spentTXOutput[key], in.Index_out)
					}
				}
			}
		}
		// 退出条件
		var hashInt big.Int
		hashInt.SetBytes(block.Pre_Hash)
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}
	// 1. 遍历区块链得到所有的交易
	// 得到迭代器
	blockIterator = bc.Iterator()

	for {
		block := blockIterator.Next()
		for _, tx := range block.Txs { // 遍历区块中的交易
			//// 先查找输入
			//// 再通过索引找到相关的output
			//if !tx.IsCoinbase(){
			//	for _, in := range tx.Tins{
			//		// 验证地址
			//		if in.UnLockWithAddr(addr){
			//			// 添加到已花费输出map中
			//			key := hex.EncodeToString(in.Tx_hash)
			//			spentTXOutput[key] = append(spentTXOutput[key], in.Index_out)
			//		}
			//	}
			//}
			// 查找数据库输出

			// 在output链中，index也是从0开始
		work:
			for index, tout := range tx.Touts {
				var utxo *UTXO
				// 地址验证(钱是不是输出给指定传入的地址)
				if tout.UnLockPubKeyWithAddr(addr) {
					// 判断是否是未花费的输出
					// 如果已花费输出不为空
					if len(spentTXOutput) != 0 {
						var isSpentTXOutput bool
						for txHash, indexArray := range spentTXOutput {
							for _, i := range indexArray {
								if txHash == hex.EncodeToString(tx.Tx_hash) && i == index {
									// 已被花费
									isSpentTXOutput = true
									// 只要满足即跳到下一个output而不是遍历spentTXoutput中是否有还有
									continue work
								}
							}
						}
						// 遍历完所有都没有被花费，才能加入到utxos中
						if isSpentTXOutput == false {
							utxo = &UTXO{tx.Tx_hash, index, tout}
							utxos = append(utxos, utxo)
						}
					} else {
						utxo = &UTXO{tx.Tx_hash, index, tout}
						utxos = append(utxos, utxo)
					}
				}
			}
		}

		// 退出条件
		var hashInt big.Int
		hashInt.SetBytes(block.Pre_Hash)
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}

	}

	return utxos
}

// 通过指定地址查找余额
func (bc *BlockChain) GetBalance(addr string) int64 {
	utxos := bc.GetUTXOs(addr, []*TX{})
	var amount int64
	for _, utxo := range utxos {
		amount += utxo.Output.Value
	}
	return amount
}

// 转账
// 查找可用的UTXO，超过需要的资金即可
// 目标地址不用传，在cmdline中指明
func (bc *BlockChain) FindSpendableUTXO(from string, amount int64, txs []*TX) (int64, map[string][]int) {
	//查找的值总和: value
	var value int64
	// 可用的UTXO: spendableUTXO
	spendableUTXO := make(map[string][]int)
	// 获取所有的UTXO
	utxos := bc.GetUTXOs(from, txs)
	// 遍历
	for _, utxo := range utxos {
		// value >= amount
		value += utxo.Output.Value
		hash := hex.EncodeToString(utxo.Tx_hash)
		spendableUTXO[hash] = append(spendableUTXO[hash], utxo.Out_index)
		if value >= amount {
			break
		}
	}
	if value < amount {
		fmt.Printf("%s 余额不足. \n", from)
		os.Exit(1)
	}
	return value, spendableUTXO
}

// 查找指定的交易
// hash_id : input所引用的交易hash
func (bc *BlockChain) FindTX(hash_id []byte, txs []*TX) TX {
	// 查找缓存中是否有符合条件的关联交易
	for _, tx := range txs {
		if bytes.Compare(tx.Tx_hash, hash_id) == 0 {
			return *tx
		}
	}

	bcit := bc.Iterator()
	for {
		block := bcit.Next()
		for _, tx := range block.Txs {
			// 判断交易hash是否相等
			if bytes.Compare(tx.Tx_hash, hash_id) == 0 {
				return *tx
			}
		}
		var hashInt big.Int
		hashInt.SetBytes(block.Pre_Hash)
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}
	// 没找到，返回空的交易
	return TX{}
}

// 交易签名
func (bc *BlockChain) SignTX(tx *TX, priKey ecdsa.PrivateKey, txs []*TX) {
	// coinbase不需要签名
	if tx.IsCoinbase() {
		return
	}

	// 处理input, 查找tx的input所引用的output所属的交易
	preTXs := make(map[string]TX)
	for _, tin := range tx.Tins {
		// 查找所引用的交易
		preTX := bc.FindTX(tin.Tx_hash, txs)
		preTXs[hex.EncodeToString(preTX.Tx_hash)] = preTX
	}

	// 实现签名函数

	tx.Sign(priKey, preTXs)

}

func (bc *BlockChain) VerifyTX(tx *TX, txs []*TX) bool {
	// coinbase不需要签名
	if tx.IsCoinbase() {
		return true
	}

	// 查找指定交易的关联交易
	preTXs := make(map[string]TX)

	for _, tin := range tx.Tins {
		preTX := bc.FindTX(tin.Tx_hash, txs)
		preTXs[hex.EncodeToString(preTX.Tx_hash)] = preTX
	}

	return tx.Verify(preTXs)
}

func (bc *BlockChain) FindUTXOMap() map[string]*TxOutputs {
	bcit := bc.Iterator()

	// 存储所有已花费的UTXO
	// key : 指定交易hash
	// value: 代表所有引用了该指定交易的output的输入

	spentUTXOMap := make(map[string][]*TxInput)

	// 与上面对应
	utxoMap := make(map[string]*TxOutputs)

	for {
		block := bcit.Next()
		// 遍历每个区块中的交易
		for i := len(block.Txs) - 1; i >= 0; i-- {
			// 保存输出
			txOutputs := &TxOutputs{[]*UTXO{}}
			// 获取每一个交易
			tx := block.Txs[i]

			if !tx.IsCoinbase() {
				// 遍历输入
				for _, tin := range tx.Tins {
					// 当前输入引用的交易hash
					txhash := hex.EncodeToString(tin.Tx_hash)
					spentUTXOMap[txhash] = append(spentUTXOMap[txhash], tin)
				}
			}

			// 遍历输出
			txhash := hex.EncodeToString(tx.Tx_hash)
		WorkOutloop:
			for i, out := range tx.Touts {
				// 查找指定hash的关联输入
				txInputs := spentUTXOMap[txhash] //指定hash交易的关联输入
				if len(txInputs) > 0 {
					// 说明有output被花费了，但不一定是所有的output，因此还得遍历一下
					isSpent := false

					for _, in := range txInputs {
						outPubKey := out.PubkeyHash
						inPubkey := in.PublicKey
						// 检查是被哪个引用了, 或者没被引用
						if bytes.Compare(outPubKey, Ripemd160_SHA256(inPubkey)) == 0 &&
							i == in.Index_out {
							// 索引也必须相等，因为如果是拆分的两个不同的utxo属于同一个人，那么还得继续细分到底引用的是哪一个
							isSpent = true
							continue WorkOutloop
						}
					}

					if !isSpent {
						utxo := &UTXO{Tx_hash: tx.Tx_hash, Out_index: i, Output: out}
						txOutputs.UTXOS = append(txOutputs.UTXOS, utxo)
					}
				} else {
					// 说明都是未花费的输出
					utxo := &UTXO{Tx_hash: tx.Tx_hash, Out_index: i, Output: out}
					txOutputs.UTXOS = append(txOutputs.UTXOS, utxo)
				}

			}

			utxoMap[txhash] = txOutputs // 该交易的所有UTXO
		}

		var hashInt big.Int
		hashInt.SetBytes(block.Pre_Hash)
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}

	return utxoMap
}

// 获取区块hash列表
func (bc *BlockChain) GetBlockHashes() [][]byte {
	var blockHashed [][]byte
	bict := bc.Iterator()

	for {
		b := bict.Next()
		blockHashed = append(blockHashed, b.Self_Hash)

		var hashInt big.Int
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}

	return blockHashed
}

// 获取最新区块的高度
func (bc *BlockChain) GetHeight() int64 {
	return bc.Iterator().Next().Height
}

// 获取指定id区块信息
func (bc *BlockChain) GetBlockInfo(hash []byte) []byte {
	var blockBytes []byte

	err := bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(Utils.BLOCKTABLENAME))
		if b != nil {
			blockBytes = b.Get(hash)
		}

		return nil
	})
	if err != nil {
		log.Panicf("view the specified hash block failed. %v\n")
		return nil
	}
	return blockBytes
}

func (bc *BlockChain) AddBlock(block *Block) {
	err := bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(Utils.BLOCKTABLENAME))

		if b != nil {
			blockBytes := b.Get(block.Self_Hash)

			if blockBytes != nil {
				// 已存在不需要同步
				return nil
			}

			err := b.Put(block.Self_Hash, block.Serialize())
			if err != nil {
				log.Panicf("sync block failed. %v\n")
			}
			blockHash := b.Get([]byte(Utils.LATEST_HASH))
			latestBlock := b.Get(blockHash)
			blockDb := DeserializeBlock(latestBlock)

			if blockDb.Height < block.Height {
				b.Put([]byte(Utils.LATEST_HASH), block.Self_Hash)
				bc.Top = block.Self_Hash
			}
		}

		return nil
	})

	if nil != err {
		log.Panicf("addblock failed. %v\n", err)
	}

	fmt.Println("new block is added successfully.")
}

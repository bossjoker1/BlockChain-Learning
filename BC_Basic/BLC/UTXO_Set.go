package BLC

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

// 持久化保存

// UTXOSet结构
// 保存指定区块链中所有的utxos
type UTXOSet struct {
	BlockChain *BlockChain
}

// 重置utxo
func (us *UTXOSet) ResetUTXOSet() {
	// 更新utxo table
	// 覆盖

	err := us.BlockChain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(UTXOTABLENAME))
		if b != nil {
			_ = tx.DeleteBucket([]byte(UTXOTABLENAME))
		}
		c, _ := tx.CreateBucketIfNotExists([]byte(UTXOTABLENAME))

		if c != nil {
			// 先查找未花费的输出
			txOutputs := us.BlockChain.FindUTXOMap()
			for hash, touts := range txOutputs {

				// 存入表
				txhash, _ := hex.DecodeString(hash)
				err := c.Put(txhash, touts.SerializeUTXOSet())

				if err != nil {
					log.Panicf("put the txoutputs failed. %v\n", err)
				}
			}
		}
		return nil
	})

	if err != nil {
		log.Panicf("update the utxoset table failed. %v\n", err)
	}
}

// 序列化utxoSet

func (txOutputs *TxOutputs) SerializeUTXOSet() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	if err := encoder.Encode(txOutputs); err != nil {
		log.Panicf("encode txoutputs failed. %v\n", err)
	}

	return res.Bytes()
}

// 反序列化
func DeserializeUTXOSet(utxoSetBytes []byte) *TxOutputs {
	var t TxOutputs
	decoder := gob.NewDecoder(bytes.NewReader(utxoSetBytes))
	// 从输入流中读取进b
	if err := decoder.Decode(&t); err != nil {
		log.Panicf("deserialize the []byte to block failed. %v\n", err)
	}
	return &t
}

// 获取余额
func (us *UTXOSet) GetBalance(addr string) int64 {
	// 获取指定地址得utxo
	utxos := us.FindUTXOWithAddr(addr)

	var amount int64

	for _, utxo := range utxos {
		fmt.Printf("\t utxo-hash: %x\n", utxo.Tx_hash)
		fmt.Printf("\t utxo-index: %x\n", utxo.Out_index)
		fmt.Printf("\t utxo-pubkeyhash: %x\n", utxo.Output.PubkeyHash)
		fmt.Printf("\t utxo-value: %d\n", utxo.Output.Value)
		amount += utxo.Output.Value
	}

	return amount
}

// 查找特定地址的utxo

func (us *UTXOSet) FindUTXOWithAddr(addr string) []*UTXO {
	var utxos []*UTXO

	// 查找数据库表

	err := us.BlockChain.DB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(UTXOTABLENAME))
		c := b.Cursor() // 获取游标

		// 遍历每个utxo
		for k, v := c.First(); k != nil; k, v = c.Next() {
			// k -> hash
			// v -> txOutputs []byte
			txOutputs := DeserializeUTXOSet(v)

			for _, utxo := range txOutputs.UTXOS {
				if utxo.Output.UnLockPubKeyWithAddr(addr) {
					utxos = append(utxos, utxo)
				}
			}

		}
		return nil
	})

	if err != nil {
		log.Panicf("get the utxo_set failed. %v\n", err)
	}

	return utxos
}

// 查找可花费的utxo

func (us *UTXOSet) FindSpendableUTXO(from string, amount int64, txs []*TX) (int64, map[string][]int) {
	// 查找整个utxo表

	// 先直接考虑未打包的utxo，不够再去表中查
	// k -> txhash
	// v -> index slice
	spentableUTXO := make(map[string][]int)

	// 查找未打包交易中的utxo
	unPkgUTXOs := us.FindUnPackageSpendableUTXOs(from, txs)

	var value int64 = 0

	for _, utxo := range unPkgUTXOs {
		value += utxo.Output.Value
		spentableUTXO[hex.EncodeToString(utxo.Tx_hash)] = append(spentableUTXO[hex.EncodeToString(utxo.Tx_hash)], utxo.Out_index)
		if value >= amount {
			// 钱够了
			return value, spentableUTXO
		}
	}

	// 在获取到未打包交易后，钱还是不够，再从utxo表中试图获取(也可能还是不够，则余额不足，转账失败)
	_ = us.BlockChain.DB.View(func(tx *bolt.Tx) error {
		// 先获取表
		b := tx.Bucket([]byte(UTXOTABLENAME))
		if b != nil {
			cursor := b.Cursor() // 游标有序遍历
			for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
				txOutputs := DeserializeUTXOSet(v)
				for _, utxo := range txOutputs.UTXOS {

					if utxo.Output.UnLockPubKeyWithAddr(from) {

						value += utxo.Output.Value
						spentableUTXO[hex.EncodeToString(utxo.Tx_hash)] = append(spentableUTXO[hex.EncodeToString(utxo.Tx_hash)], utxo.Out_index)
						if value >= amount {
							break
						}
					}
				}
			}
		}

		return nil
	})
	if value < amount {
		log.Panicf("[%v]余额不足\n", from)
	}
	return value, spentableUTXO
}

// 从未打包的交易中进行查找
func (us *UTXOSet) FindUnPackageSpendableUTXOs(from string, txs []*TX) []*UTXO {
	var unUtxos []*UTXO

	spendTXOutputs := make(map[string][]int) // 每个交易中已经被引用(花费)的output的索引

	for _, tx := range txs {
		// 排除coinbase

		if !tx.IsCoinbase() {
			for _, tin := range tx.Tins {
				publicKeyHash := Base58Decode([]byte(from)) // 获取公钥hash
				ripemd_sha := publicKeyHash[1 : len(publicKeyHash)-CHECKSUMLEN]
				if tin.UnLockWithRipemd_SHA(ripemd_sha) {
					// 添加到已花费输出map中
					key := hex.EncodeToString(tin.Tx_hash)
					spendTXOutputs[key] = append(spendTXOutputs[key], tin.Index_out)
				}
			}
		}
	}

	for _, tx := range txs {
	UnUTXOLoop:
		for i, tout := range tx.Touts {
			if tout.UnLockPubKeyWithAddr(from) {
				if len(spendTXOutputs) == 0 {
					// 没有包含已花费输出
					utxo := &UTXO{tx.Tx_hash, i, tout}
					unUtxos = append(unUtxos, utxo)
				} else {
					for hash, index := range spendTXOutputs {
						if hash == hex.EncodeToString(tx.Tx_hash) {
							// 当前交易是否包含已花费输出
							var isUnpkgUTXO bool
							for _, idx := range index {
								if idx == i {
									// 判断索引是否匹配
									// 如果相等则说明确实被花费
									isUnpkgUTXO = true
									continue UnUTXOLoop
								}
							}

							if !isUnpkgUTXO {
								utxo := &UTXO{tx.Tx_hash, i, tout}
								unUtxos = append(unUtxos, utxo)
							}
						} else {
							utxo := &UTXO{tx.Tx_hash, i, tout}
							unUtxos = append(unUtxos, utxo)
						}
					}
				}
			}
		}
	}

	return unUtxos
}

// 实现utxo table 实时更新
func (us *UTXOSet) Update() {
	//1. 找到需要删除的utxo
	// 获取最新区块
	latest_block := us.BlockChain.Iterator().Next()
	// 存放最新区块的所有输入
	inputs := []*TxInput{} // 即表中的哪些utxo被引用了

	// 需要新存入的outputs

	outsMap := make(map[string]*TxOutputs)

	// 2.查找需要删除的数据

	for _, tx := range latest_block.Txs {
		// 遍历输入

		for _, tin := range tx.Tins {
			inputs = append(inputs, tin)
		}
	}

	// 3.遍历当前区块
	// 找到并删掉
	for _, tx := range latest_block.Txs {
		var utxos []*UTXO
		for index, out := range tx.Touts {
			isSpent := false
			for _, in := range inputs {
				if in.Index_out == index && bytes.Compare(tx.Tx_hash, in.Tx_hash) == 0 {
					if bytes.Compare(out.PubkeyHash, Ripemd160_SHA256(in.PublicKey)) == 0 {
						isSpent = true
						continue
					}
				}
			}
			if !isSpent {
				utxo := &UTXO{tx.Tx_hash, index, out}
				utxos = append(utxos, utxo)
			}
		}
		if len(utxos) > 0 {
			outsMap[hex.EncodeToString(tx.Tx_hash)] = &TxOutputs{}
		}
	}

	// 4.更新
	err := us.BlockChain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(UTXOTABLENAME))
		if b != nil {
			// 删除已花费输出
			for _, in := range inputs {
				txOutBytes := b.Get(in.Tx_hash) // 查找需要改动的交易的hash

				if len(txOutBytes) == 0 {

					continue
				}

				UTXOS := []*UTXO{}
				// 反序列化
				txOutputs := DeserializeUTXOSet(txOutBytes)

				isDel := false

				for _, utxo := range txOutputs.UTXOS {
					// 判断具体是哪一个输出被引用了

					if in.Index_out == utxo.Out_index && bytes.Compare(utxo.Output.PubkeyHash, Ripemd160_SHA256(in.PublicKey)) == 0 {
						isDel = true
					} else {
						UTXOS = append(UTXOS, utxo)
					}
				}

				if isDel {
					// 先删除输出
					_ = b.Delete(in.Tx_hash)
					if len(UTXOS) > 0 {
						preTXOutputs := outsMap[hex.EncodeToString(in.Tx_hash)]
						preTXOutputs.UTXOS = append(preTXOutputs.UTXOS, UTXOS...)

						// 更新
						outsMap[hex.EncodeToString(in.Tx_hash)] = preTXOutputs
					}
				}
			}

			for hash, outputs := range outsMap {
				hashBytes, _ := hex.DecodeString(hash)
				_ = b.Put(hashBytes, outputs.SerializeUTXOSet())
			}
		}
		return nil
	})

	if err != nil {
		log.Panicf("update the utxodb failed. %v\n", err)
	}
}

// 挖出新的区块时，UTXO_Set被更新
// 去掉已花费的UTXO, 增加新挖出的未花费的UTXO
// 如果一笔交易中的不包含UTXO了，则该交易也从UTXOSet中删掉

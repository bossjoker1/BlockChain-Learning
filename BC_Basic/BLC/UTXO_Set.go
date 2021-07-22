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

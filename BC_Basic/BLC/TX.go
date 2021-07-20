package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

type TX struct {
	// 交易唯一标识
	Tx_hash []byte
	// 输入
	Tins []*TxInput
	// 输出
	Touts []*TxOutput
}

// 生成交易hash
func (tx *TX) HashTX() {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	// 序列化
	err := encoder.Encode(tx)
	if err != nil {
		log.Panicf("tx hash generate failed. %v\n ", err)
	}

	hash := sha256.Sum256(res.Bytes())
	tx.Tx_hash = hash[:]
}

// 生成 coinbase
func NewCoinBaseTX(addr string) *TX {
	// 输入
	// 系统给的奖励，没有上一块的hash，也没有上一次输出的索引，这里置为 -1
	txInput := &TxInput{[]byte{}, -1, "Genesis Data"}
	// 输出
	txOutput := &TxOutput{MINEAWARD, addr}
	txCoinbase := &TX{nil, []*TxInput{txInput}, []*TxOutput{txOutput}}
	// hash
	txCoinbase.HashTX()

	return txCoinbase
}

// 生成 转账交易
func NewSimpleTX(from string, to string, amount int64, bc *BlockChain, txs []*TX) *TX {
	var (
		txInputs  []*TxInput
		txOutputs []*TxOutput
	)

	// 查找指定地址可用的UTXO
	money, spendableUTXO := bc.FindSpendableUTXO(from, amount, txs)
	fmt.Printf("from %s , money: %d\n", from, money)
	for txHash, indeArray := range spendableUTXO {
		txHashBytes, _ := hex.DecodeString(txHash)
		for _, index := range indeArray {
			// 此处的输出需要被花费，即被其它的交易所引用
			txInput := &TxInput{Tx_hash: txHashBytes, Index_out: index, ScriptSig: from}
			txInputs = append(txInputs, txInput)
		}
	}

	// 输入结构
	// 消费
	//txInput := &TxInput{[]byte("1"), 0, from}
	//txInputs = append(txInputs, txInput)
	//utxo := GetUTXOs(from)

	// 输出结构
	// 转账
	txOutput := &TxOutput{amount, to}
	txOutputs = append(txOutputs, txOutput)

	// 找零
	txOutput = &TxOutput{money - amount, from}
	txOutputs = append(txOutputs, txOutput)

	tx := &TX{nil, txInputs, txOutputs}
	tx.HashTX()

	return tx
}

// 判断是否为coinbase交易
func (tx *TX) IsCoinbase() bool {
	return len(tx.Tins[0].Tx_hash) == 0 && tx.Tins[0].Index_out == -1
}

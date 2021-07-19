package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
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
func NewSimpleTX(from string, to string, amount int64) *TX {
	var (
		txInputs  []*TxInput
		txOutputs []*TxOutput
	)

	// 输入结构
	// 消费
	txInput := &TxInput{[]byte("1"), 0, from}
	txInputs = append(txInputs, txInput)

	// 输出结构
	// 转账
	txOutput := &TxOutput{amount, to}
	txOutputs = append(txOutputs, txOutput)

	// 找零
	txOutput = &TxOutput{MINEAWARD - amount, from}
	txOutputs = append(txOutputs, txOutput)

	tx := &TX{nil, txInputs, txOutputs}
	tx.HashTX()

	return tx
}

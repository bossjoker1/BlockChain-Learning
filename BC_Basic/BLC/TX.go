package BLC

import (
	"BlockChain-Learning/BC_Basic/Utils"
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"time"
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
	tm := time.Now().Unix()
	txhashBytes := bytes.Join([][]byte{res.Bytes(), Utils.IntToHex(tm)}, []byte{})
	hash := sha256.Sum256(txhashBytes)
	tx.Tx_hash = hash[:]
}

// 生成 coinbase
func NewCoinBaseTX(addr string) *TX {
	// 输入
	// 系统给的奖励，没有上一块的hash，也没有上一次输出的索引，这里置为 -1
	// 创世区块没有hash
	// 添加时间参数
	txInput := &TxInput{[]byte{}, -1, nil, nil}
	// 输出
	txOutput := NewTXOutput(Utils.MINEAWARD, addr)
	txCoinbase := &TX{nil, []*TxInput{txInput}, []*TxOutput{txOutput}}
	// hash
	txCoinbase.HashTX()

	return txCoinbase
}

// 生成 转账交易
func NewSimpleTX(from string, to string, amount int64, bc *BlockChain, txs []*TX, us *UTXOSet, node_id string) *TX {
	var (
		txInputs  []*TxInput
		txOutputs []*TxOutput
	)

	// 查找指定地址可用的UTXO
	money, spendableUTXO := us.FindSpendableUTXO(from, amount, txs)

	fmt.Printf("from %s , money: %d\n", from, money)

	wallets, _ := NewWallets(node_id)
	wallet := wallets.Wallets[from]
	pubKey := wallet.PublicKey

	for txHash, indeArray := range spendableUTXO {
		txHashBytes, _ := hex.DecodeString(txHash)
		for _, index := range indeArray {
			// 此处的输出需要被花费，即被其它的交易所引用
			txInput := &TxInput{Tx_hash: txHashBytes, Index_out: index, Signature: nil, PublicKey: pubKey}
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
	txOutput := NewTXOutput(amount, to)
	txOutputs = append(txOutputs, txOutput)

	// 找零
	txOutput = NewTXOutput(money-amount, from)
	txOutputs = append(txOutputs, txOutput)

	tx := &TX{nil, txInputs, txOutputs}
	tx.HashTX()

	// 用私钥对新生成的交易进行签名
	bc.SignTX(tx, wallet.PrivateKey, txs)

	return tx
}

// 判断是否为coinbase交易
func (tx *TX) IsCoinbase() bool {
	return tx.Tins[0].Index_out == -1
}

// 对指定交易签名

func (tx *TX) Sign(priKey ecdsa.PrivateKey, prevTXs map[string]TX) {
	// 判断是否是挖矿交易
	if tx.IsCoinbase() {
		return
	}

	for _, tin := range tx.Tins {
		if prevTXs[hex.EncodeToString(tin.Tx_hash)].Tx_hash == nil {
			log.Panicf("prev txs is not correct\n")
		}
	}

	// 提取需要签名的字段
	txCopy := tx.TXTrimmedCopy()

	for i, _ := range txCopy.Tins {
		//preTx := prevTXs[hex.EncodeToString(tin.Tx_hash)]
		txCopy.Tins[i].Signature = nil
		// 这里有问题
		//txCopy.Tins[i].PublicKey = preTx.Touts[tin.Index_out].PubkeyHash
		txCopy.Tx_hash = txCopy.Hash()
		txCopy.Tins[i].PublicKey = nil

		// 这里的hash表示对摘要后的hash值进行签名
		r, s, err := ecdsa.Sign(rand.Reader, &priKey, txCopy.Tx_hash)
		if err != nil {
			log.Panicf("sign to tx %x failed. %v\n ", tx.Tx_hash, err)
		}

		signature := append(r.Bytes(), s.Bytes()...)

		tx.Tins[i].Signature = signature
	}
}

// 设置用于签名的数据Hash
func (tx *TX) Hash() []byte {
	txCopy := tx
	txCopy.Tx_hash = []byte{}
	hash := sha256.Sum256(txCopy.Serialize())

	return hash[:]
}

func (tx *TX) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	// 将b编码后存进res输出流
	if err := encoder.Encode(tx); err != nil {
		log.Panicf("serialize the tx to []byte failed %v\n", err)
	}
	return res.Bytes()
}

// 获取交易的copy,返回需要签名的交易字段
// 把不需要签名的字段置为空
func (tx *TX) TXTrimmedCopy() TX {
	var (
		inputs  []*TxInput
		outPuts []*TxOutput
	)

	for _, tin := range tx.Tins {
		inputs = append(inputs, &TxInput{tin.Tx_hash, tin.Index_out, nil, tin.PublicKey})
	}

	for _, tout := range tx.Touts {
		outPuts = append(outPuts, &TxOutput{tout.Value, tout.PubkeyHash})
	}

	txCopy := TX{tx.Tx_hash, inputs, outPuts}

	return txCopy
}

func (tx *TX) Verify(prevTXs map[string]TX) bool {
	if tx.IsCoinbase() {
		return true
	}

	// 能否找到交易
	for _, tin := range tx.Tins {
		if prevTXs[hex.EncodeToString(tin.Tx_hash)].Tx_hash == nil {
			log.Panic("Error: TX is incorrect.\n")
			return false
		}
	}

	txCopy := tx.TXTrimmedCopy()

	// 获取密钥对
	// 使用相同的椭圆
	curve := elliptic.P256()

	for i, tin := range tx.Tins {
		// 获取签名时的数据，才能用于还原
		//preTx := prevTXs[hex.EncodeToString(tin.Tx_hash)]
		txCopy.Tins[i].Signature = nil
		//txCopy.Tins[i].PublicKey = preTx.Touts[tin.Index_out].PubkeyHash
		txCopy.Tx_hash = txCopy.Hash() // 要验证的数据
		txCopy.Tins[i].PublicKey = nil

		// 获取R, S
		r := big.Int{}
		s := big.Int{}

		sigLen := len(tin.Signature)

		r.SetBytes(tin.Signature[:(sigLen / 2)])
		s.SetBytes(tin.Signature[(sigLen / 2):])

		// 生成X, Y坐标数据
		x := big.Int{}
		y := big.Int{}
		pubKeyLen := len(tin.PublicKey)
		x.SetBytes(tin.PublicKey[:(pubKeyLen / 2)])
		y.SetBytes(tin.PublicKey[(pubKeyLen / 2):])

		// 生成签名公钥
		rawPubkey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}

		// 验证签名

		if !ecdsa.Verify(&rawPubkey, txCopy.Tx_hash, &r, &s) {
			return false
		}
	}

	return true

}

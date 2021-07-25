package BLC

import (
	"BlockChain-Learning/BC_Basic/TX"
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

// 基本区块结构
// 根据结构示意图定义

type Block struct {
	TimeStamp int64 // 时间戳：区块产生的时间
	Height    int64 // 区块高度(代表区块唯一号码)
	Pre_Hash  []byte
	Self_Hash []byte
	Txs       []*TX.TX
	//Data      []byte // 交易结构，之后定义
	//随机数r & merkle root_hash
	Nonce int64 // 记录碰撞成功后结束的随机数值
}

// 创建新区块
func NewBlock(height int64, pre_hash []byte, txs []*TX.TX) *Block {
	block := &Block{Height: height, Pre_Hash: pre_hash, Txs: txs, TimeStamp: time.Now().Unix()}
	//block.SetHash()
	pow := NewPoW(block)
	hash, nonce := pow.Run()
	block.Self_Hash = hash
	block.Nonce = nonce
	return block
}

// 计算区块hash值
// 定义成方法
//func (b *Block) SetHash() {
//	heightBytes := IntToHex(b.Height)
//	timeStampBytes := IntToHex(b.TimeStamp)
//	// 拼接 生成hash
//	blockBytes := bytes.Join([][]byte{heightBytes, timeStampBytes, b.Pre_Hash, b.Data}, []byte{})
//	hash := sha256.Sum256(blockBytes)
//	b.Self_Hash = hash[:]
//}

// 生成创世区块
func CreateGenesisBlock(txs []*TX.TX) *Block {
	return NewBlock(1, nil, txs)
}

// 序列化，将区块结构序列化为[]byte
func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	// 将b编码后存进res输出流
	if err := encoder.Encode(b); err != nil {
		log.Panicf("serialize the block to []byte failed %v\n", err)
	}
	return res.Bytes()
}

// 反序列化
func DeserializeBlock(blockBytes []byte) *Block {
	var b Block
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	// 从输入流中读取进b
	if err := decoder.Decode(&b); err != nil {
		log.Panicf("deserialize the []byte to block failed. %v\n", err)
	}
	return &b
}

// 将区块中的交易结构转化为[]byte
func (b *Block) HashTransactions() []byte {
	var txs [][]byte

	for _, tx := range b.Txs {
		// 交易数据
		txs = append(txs, tx.Serialize())
	}

	mTree := NewMerkleTree(txs)

	//txHash := sha256.Sum256(bytes.Join(txHashes, []byte{}))

	// 改成merkle的根hash
	// 区块的hash为merkle树的根hash
	return mTree.Root.Data
}

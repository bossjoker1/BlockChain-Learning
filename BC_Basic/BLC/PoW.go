package BLC

import (
	"BlockChain-Learning/BC_Basic/Utils"
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

// 目标难度值
// 生成hash值前缀0的位数
// 先写死
const targetBit = 16

type PoW struct {
	Block  *Block
	target *big.Int
}

// 创建新的PoW对象
func NewPoW(block *Block) *PoW {
	//左移决定比较难度
	//0000 0001
	//0100 0000
	//比它小的肯定前缀有两个0
	target := big.NewInt(1)
	target = target.Lsh(target, 256-targetBit)
	return &PoW{block, target}
}

// 工作量证明

func (poW *PoW) Run() ([]byte, int64) {
	var (
		nonce   = 0      // 从0开始碰撞
		hash    [32]byte // 生成的hash值
		hashInt big.Int  // 用于数据比较
	)
	for {
		dataBytes := poW.PrepareData(nonce)
		hash = sha256.Sum256(dataBytes)
		hashInt.SetBytes(hash[:])
		//fmt.Printf("hash : \r %x", hash)
		// 难度比较

		if poW.target.Cmp(&hashInt) == 1 {
			break
		}
		nonce++
	}
	fmt.Println("\n碰撞次数：", nonce)
	return hash[:], int64(nonce)
}

// 随机数碰撞尝试
func (poW *PoW) PrepareData(nonce int) []byte {
	data := bytes.Join([][]byte{
		poW.Block.Pre_Hash,
		poW.Block.HashTransactions(),
		Utils.IntToHex(poW.Block.TimeStamp),
		Utils.IntToHex(poW.Block.Height),
		Utils.IntToHex(int64(nonce)),
		Utils.IntToHex(targetBit),
	}, []byte{})

	return data
}

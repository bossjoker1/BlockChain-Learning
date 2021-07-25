package Server

import (
	"BlockChain-Learning/BC_Basic/BLC"
	"BlockChain-Learning/BC_Basic/UTXO"
	"BlockChain-Learning/BC_Basic/Utils"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

//1. `version` : 验证当前节点的末端区块是不是最新区块，可以通过 `Height`来判断。

func HandleVersion(req []byte, bc *BLC.BlockChain) {
	fmt.Println("handle version in func")
	var buff bytes.Buffer
	var data Version

	// 解析req中的数据
	dataBytes := req[Utils.CMDLENGTH:]
	buff.Write(dataBytes)
	decoder := gob.NewDecoder(&buff)
	err := decoder.Decode(&data)
	if err != nil {
		log.Panicf("decode the version data failed. %v\n", err)
	}

	// 获取区块高度来判断是否需要更新
	baseHeight := bc.GetHeight()
	// 把自身版本信息发送给对方，看对方是否要同步版本
	if data.Height < baseHeight {
		SendVersion(data.AddrFrom, bc)
	} else if data.Height > baseHeight {
		// 收到req的节点需要同步
		sendGetBlocks(data.AddrFrom)
	}
}

//2. `GetBlocks` : 从最长的链上面获取区块

func HandleGetBlocks(req []byte, bc *BLC.BlockChain) {
	var buff bytes.Buffer
	var data Get_Blocks

	// 解析req中的数据
	dataBytes := req[Utils.CMDLENGTH:]
	buff.Write(dataBytes)
	decoder := gob.NewDecoder(&buff)
	err := decoder.Decode(&data)
	if err != nil {
		log.Panicf("decode the Getblocks data failed. %v\n", err)
	}
	hashes := bc.GetBlockHashes()
	SendInv(data.AddrFrom, Utils.BLOCK_TYPE, hashes)
}

//3. `Inv` : 向其他节点展示当前节点有哪些区块

func HandleInv(req []byte, bc *BLC.BlockChain) {
	var buff bytes.Buffer
	var data Inv

	// 解析req中的数据
	dataBytes := req[Utils.CMDLENGTH:]
	buff.Write(dataBytes)
	decoder := gob.NewDecoder(&buff)
	err := decoder.Decode(&data)
	if err != nil {
		log.Panicf("decode the inv failed. %v\n", err)
	}
	// 把需要展示的区块发回
	// SendGetData
	blockHash := data.Hashes[0]
	SendGetData(data.AddrFrom, Utils.BLOCK_TYPE, blockHash)
}

//4. `GetData` : 请求一个指定的区块

func HandleGetData(req []byte, bc *BLC.BlockChain) {
	var buff bytes.Buffer
	var data GetData

	// 解析req中的数据
	dataBytes := req[Utils.CMDLENGTH:]
	buff.Write(dataBytes)
	decoder := gob.NewDecoder(&buff)
	err := decoder.Decode(&data)
	if err != nil {
		log.Panicf("decode the inv failed. %v\n", err)
	}

	// 获取指定id区块信息
	// GetBlockInfo()
	block := bc.GetBlockInfo(data.ID)

	SendBlock(data.AddrFrom, block)
}

//5. `block` : 接受到新区块的时候，进行处理

func HandleBlock(req []byte, bc *BLC.BlockChain) {
	var buff bytes.Buffer
	var data BlockData

	// 解析req中的数据
	dataBytes := req[Utils.CMDLENGTH:]
	buff.Write(dataBytes)
	decoder := gob.NewDecoder(&buff)
	err := decoder.Decode(&data)
	if err != nil {
		log.Panicf("decode the inv failed. %v\n", err)
	}
	blockBytes := data.Block
	b := BLC.DeserializeBlock(blockBytes)
	// 上传bc数据库
	bc.AddBlock(b)

	// 更新utxo表

	us := &UTXO.UTXOSet{bc}

	us.Update()
}

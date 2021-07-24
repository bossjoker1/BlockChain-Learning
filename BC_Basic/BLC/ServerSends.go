package BLC

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

// P2P 节点(客户端)向节点(服务端)发送请求
func SendRequest(to string, msg []byte) {
	fmt.Println("send request to server.")
	conn, err := net.Dial(PROTOCOL, to)
	if err != nil {
		log.Panicf("dail to  %s failed. %v\n", to, err)
	}
	defer conn.Close()

	// 要发送的指令发送到请求中
	_, _ = io.Copy(conn, bytes.NewReader(msg))

}

// 数据同步
func SendVersion(to string, bc *BlockChain) {
	// 在比特币中，消息是底层的bit序列，前12个字节指定命令名
	// 后面的字节是gob编码的消息结构
	height := bc.GetHeight()
	data := GobEncode(Version{NODE_VERSION, int64(height), nodeAddr})

	request := append(CommandToBytes(VERSION_NUM), data...)
	SendRequest(to, request)
}

// 向其他节点显示区块信息
func SendInv(to string, kind string, hashes [][]byte) {
	data := GobEncode(Inv{Hashes: hashes, AddrFrom: nodeAddr, Type: kind})
	req := append(CommandToBytes(CMD_INV), data...)

	SendRequest(to, req)
}

// 从指定节点同步数据

func sendGetBlocks(to string) {
	data := GobEncode(Get_Blocks{nodeAddr})
	req := append(CommandToBytes(CMD_GETBLOCKS), data...)

	SendRequest(to, req)
}

// 向其他人展示区块信息或交易
func SendGetData(to string, kind string, hash []byte) {
	data := GobEncode(GetData{AddrFrom: nodeAddr, ID: hash, Type: BLOCK_TYPE})
	req := append(CommandToBytes(CMD_GETDATA), data...)

	SendRequest(to, req)
}

// 传输Block数据

func SendBlock(to string, block []byte) {
	data := GobEncode(BlockData{AddrFrom: nodeAddr, Block: block})
	req := append(CommandToBytes(CMD_BLOCK), data...)

	SendRequest(to, req)
}

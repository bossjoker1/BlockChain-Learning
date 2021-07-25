package Server

import (
	"BlockChain-Learning/BC_Basic/BLC"
	"BlockChain-Learning/BC_Basic/Utils"
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

// 3000作为主节点
var knownNodes = []string{"localhost:3000"}

// 服务处理文件
var nodeAddr string // 节点地址

// 启动服务器
func StartServer(nodeId string) {
	nodeAddr = fmt.Sprintf("localhost:%s", nodeId)
	fmt.Printf("服务节点 [%s] 启动...\n", nodeAddr)
	// 监听节点
	listen, err := net.Listen(Utils.PROTOCOL, nodeAddr)
	if err != nil {
		log.Panicf("listen addr of %s failed. %v\n", nodeAddr, err)
	}
	defer listen.Close()

	bc := BLC.GetBCObject(nodeId)

	if nodeAddr != knownNodes[0] {
		// 不是主节点, 需要向主节点发送请求同步数据
		// SendVersion()
		// SendRequest(knownNodes[0], []byte(nodeAddr))
		SendVersion(knownNodes[0], bc)
	}

	// 主节点接收请求
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Panicf("connect to addr of %s failed. %v\n", nodeAddr, err)
		}
		req, _ := ioutil.ReadAll(conn)
		//fmt.Println("firt source : ", conn)
		//fmt.Println("first req : ", req)
		cmd := Utils.BytesToCommand(req[:Utils.CMDLENGTH])
		fmt.Printf("Receive a request on server : %s\n", cmd)
		// 分出单独的goroutine并发的对请求进行处理
		go HandleRequest(req, bc)

		//go PrintHandle(conn)

		defer conn.Close()
	}

}

//func PrintHandle(conn net.Conn)  {
//	fmt.Println("second source : ", conn)
//	//fmt.Println("second req : ", req)
//	req, err := ioutil.ReadAll(conn)
//	if err != nil{
//		log.Panicf("err: %v\n", err)
//	}
//	cmd := bytesToCommand(req[:CMDLENGTH])
//	fmt.Println("second req : ", req)
//	fmt.Printf("2 cmd : %s\n", cmd)
//
//}

// 处理请求
func HandleRequest(req []byte, bc *BLC.BlockChain) {
	//req, _ := ioutil.ReadAll(conn)
	cmd := Utils.BytesToCommand(req[:Utils.CMDLENGTH])
	//fmt.Printf("Receive a request: %s\n", cmd)

	// 对命令进行判断
	switch cmd {
	case Utils.CMD_VERSION:
		fmt.Println("handle version")
		HandleVersion(req, bc)
	case Utils.CMD_GETDATA:
		HandleGetData(req, bc)
	case Utils.CMD_GETBLOCKS:
		HandleGetBlocks(req, bc)
	case Utils.CMD_INV:
		HandleInv(req, bc)
	case Utils.CMD_BLOCK:
		HandleBlock(req, bc)
	default:
		fmt.Println("unknown cmd!")
	}
	//defer conn.Close()
}

package BLC

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct {
	BC *BlockChain
}

// show tips

func PrintUsage() {
	fmt.Println("Usage: ")
	fmt.Printf("\tcreateblockchain	-- 创建区块链\n")
	fmt.Printf("\taddblock -data DATA	-- 交易数据\n")
	fmt.Printf("\tprintchain		-- 输出区块链信息\n")
}

// 判断输入参数是否规范
func IsValidArgs() {
	if len(os.Args) < 2 {
		PrintUsage() // 打印用法
		os.Exit(1)   // 退出程序
	}
}

// 添加区块
func (cli *CLI) AddBlock(data string) {
	if dbExists() == false {
		fmt.Println("database not exist.")
		os.Exit(1)
	}
	blockChain := GetBCObject() // 获得区块链对象
	defer blockChain.DB.Close()
	blockChain.AddBlock([]byte(data))
}

// 输出区块链信息
func (cli *CLI) PrintChain() {
	if dbExists() == false {
		fmt.Println("database not exist.")
		os.Exit(1)
	}
	blockChain := GetBCObject() // 获得区块链对象
	defer blockChain.DB.Close()
	blockChain.PrintChainInfo()
}

// 创建区块链
func (cli *CLI) CreateBlockchainWithGenesis() {
	CreateBlockChainWithGenesisBlock()
}

// 运行函数
func (cli *CLI) Run() {
	// 检测参数数量
	IsValidArgs()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createblcCmd := flag.NewFlagSet("createbc", flag.ExitOnError)

	// 获取命令行参数

	flagAddBlockArg := addBlockCmd.String("data", "send 0 BTC to sb.", "交易数据")
	switch os.Args[1] {
	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panicf("parse addblock failed. %v\n", err)
		}
	case "printchain":
		err := printCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panicf("parse printchain failed. %v\n", err)
		}
	case "createblockchain":
		err := createblcCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panicf("parse createblockchain failed. %v\n", err)
		}
	default:
		PrintUsage()
		os.Exit(1)
	}

	// Parsed()判断是否解析成功

	// 添加区块
	if addBlockCmd.Parsed() {
		if *flagAddBlockArg == "" {
			PrintUsage()
			os.Exit(1)
		}
		cli.AddBlock(*flagAddBlockArg)
	}
	// 输出信息
	if printCmd.Parsed() {
		cli.PrintChain()
	}
	// 创建区块链
	if createblcCmd.Parsed() {
		cli.CreateBlockchainWithGenesis()
	}

}

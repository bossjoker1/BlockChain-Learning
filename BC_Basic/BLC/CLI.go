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
	fmt.Printf("\tcreateblockchain -addr [address]	-- 地址\n")
	fmt.Printf("\taddblock -data [DATA]		-- 交易数据\n")
	fmt.Printf("\tprintchain			-- 输出区块链信息\n")
	fmt.Printf("\tsend -from [addr1] -to [addr2] -amount [value] -- 转账\n")
}

// 判断输入参数是否规范
func IsValidArgs() {
	if len(os.Args) < 2 {
		PrintUsage() // 打印用法
		os.Exit(1)   // 退出程序
	}
}

// 发送交易
func (cli *CLI) Send(from, to, amount []string) {
	if dbExists() == false {
		fmt.Println("database not exist")
		os.Exit(1)
	}
	bc := GetBCObject() //获得区块链对象
	defer bc.DB.Close()
	bc.MineNewBlock(from, to, amount)

}

// 添加区块
func (cli *CLI) AddBlock(txs []*TX) {
	if dbExists() == false {
		fmt.Println("database not exist.")
		os.Exit(1)
	}
	blockChain := GetBCObject() // 获得区块链对象
	defer blockChain.DB.Close()
	blockChain.AddBlock(txs)
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
func (cli *CLI) CreateBlockchainWithGenesis(addr string) {
	CreateBlockChainWithGenesisBlock(addr)
}

// 运行函数
func (cli *CLI) Run() {
	// 检测参数数量
	IsValidArgs()

	// 新建命令
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createblcCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)

	// 获取命令行参数
	flagCreateBlockChainAddr := createblcCmd.String("addr", "aa", "地址")
	flagAddBlockArg := addBlockCmd.String("data", "send 0 BTC to sb.", "交易数据")

	// 转账参数
	flagFromAddr := sendCmd.String("from", "", "转账源地址")
	flagToAddr := sendCmd.String("to", "", "转账目的地址")
	flagAmount := sendCmd.String("amount", "", "转账金额")

	switch os.Args[1] {
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panicf("parse send failed. %v\n", err)
		}

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

	// 转账
	if sendCmd.Parsed() {
		if *flagFromAddr == "" || *flagToAddr == "" {
			fmt.Println("源地址或目的地址不能为空")
			PrintUsage()
			os.Exit(1)
		}
		if *flagAmount == "" {
			fmt.Println("金额不能为空")
			PrintUsage()
			os.Exit(1)
		}
		// 打印转账信息

		fmt.Printf("From [%s] to [%s]   value [%s]\n", JsonToArray(*flagFromAddr), JsonToArray(*flagToAddr), JsonToArray(*flagAmount))
		cli.Send(JsonToArray(*flagFromAddr), JsonToArray(*flagToAddr), JsonToArray(*flagAmount)) // 发送交易
	}

	// 添加区块
	if addBlockCmd.Parsed() {
		if *flagAddBlockArg == "" {
			PrintUsage()
			os.Exit(1)
		}

		// 暂时传的空值
		cli.AddBlock([]*TX{})
	}
	// 输出信息
	if printCmd.Parsed() {
		cli.PrintChain()
	}
	// 创建区块链
	if createblcCmd.Parsed() {
		if *flagCreateBlockChainAddr == "" {
			PrintUsage()
			os.Exit(1)
		}
		//fmt.Println(*flagCreateBlockChainAddr)
		cli.CreateBlockchainWithGenesis(*flagCreateBlockChainAddr)
	}

}

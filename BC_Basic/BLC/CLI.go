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
	fmt.Printf("\ttest	-- 测试程序代码\n ")
	fmt.Printf("\tcreatewallets 	-- 创建钱包\n")
	fmt.Printf("\taddrlists 	-- 获取钱包地址列表\n")
	fmt.Printf("\tcreateblockchain -addr [address]	-- 地址\n")
	// 通过交易生成区块
	// fmt.Printf("\taddblock -data [DATA]		-- 交易数据\n")
	fmt.Printf("\tprintchain			-- 输出区块链信息\n")
	fmt.Printf("\tsend -from [addr1] -to [addr2] -amount [value] -- 转账\n")
	fmt.Printf("\tgetbalance -from [addr] 	-- 查询余额")

}

// 判断输入参数是否规范
func IsValidArgs() {
	if len(os.Args) < 2 {
		PrintUsage() // 打印用法
		os.Exit(1)   // 退出程序
	}
}

func (cli *CLI) GetBalance(from string) {
	// 获取指定地址的余额
	//outPuts := GetUTXOs(from)
	//fmt.Printf("unUTXO: %v\n", outPuts)
	bc := GetBCObject()
	defer bc.DB.Close()
	//amount := bc.GetBalance(from)
	//fmt.Printf("\t addr : %s , balance : %d \n", from, amount)

	utxoSet := &UTXOSet{bc}
	amount := utxoSet.GetBalance(from)
	fmt.Printf("\t addr : %s , balance : %d \n", from, amount)
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

//
//// 添加区块
//func (cli *CLI) AddBlock(txs []*TX) {
//	if dbExists() == false {
//		fmt.Println("database not exist.")
//		os.Exit(1)
//	}
//	blockChain := GetBCObject() // 获得区块链对象
//	defer blockChain.DB.Close()
//	blockChain.AddBlock(txs)
//}

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
	blockChain := CreateBlockChainWithGenesisBlock(addr)
	defer blockChain.DB.Close()

	// 设置UTXOSet操作

	utxoSet := &UTXOSet{BlockChain: blockChain}
	// 更新
	utxoSet.ResetUTXOSet()
}

// 创建钱包集合
func (cli *CLI) CreateWallets() {
	wallets, _ := NewWallets()
	wallets.CreateWallet()
	fmt.Printf("wallet: %v\n", wallets)
}

func (cli *CLI) GetAddrLists() {
	fmt.Println("print all wallets' address.")
	wallets, _ := NewWallets()
	for addr, _ := range wallets.Wallets {
		fmt.Printf("address : [%s]\n", addr)
	}
}

func (cli *CLI) TestMethod() {
	blockchain := GetBCObject()
	defer blockchain.DB.Close()

	utxoMap := blockchain.FindUTXOMap()
	for key, value := range utxoMap {
		fmt.Printf("key : [%x]\n", key)
		for _, utxo := range value.UTXOS {
			fmt.Printf("\thash : [%x], out_index : [%v], txOutput: [%v]\n", utxo.Tx_hash, utxo.Out_index, utxo.Output)
		}
	}
	fmt.Println("-----------------")
}

func (cli *CLI) TestResetUTXO() {
	bc := GetBCObject()
	defer bc.DB.Close()

	utxoSet := &UTXOSet{BlockChain: bc}
	utxoSet.ResetUTXOSet() // 重置UTXO
}

// 运行函数
func (cli *CLI) Run() {
	// 检测参数数量
	IsValidArgs()

	// 新建命令
	//addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)

	createWalletsCmd := flag.NewFlagSet("createwallets", flag.ExitOnError)

	getAddrListsCmd := flag.NewFlagSet("addrlists", flag.ExitOnError)

	printCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createblcCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)

	// 发送交易
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)

	// 查询余额
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)

	// 获取命令行参数

	flagCreateBlockChainAddr := createblcCmd.String("addr", "aa", "地址")
	// flagAddBlockArg := addBlockCmd.String("data", "send 0 BTC to sb.", "交易数据")

	// 转账参数
	flagFromAddr := sendCmd.String("from", "", "转账源地址")
	flagToAddr := sendCmd.String("to", "", "转账目的地址")
	flagAmount := sendCmd.String("amount", "", "转账金额")

	// 查询余额参数
	flagBalanceArg := getBalanceCmd.String("from", "", "查询地址")

	// 测试cmd line
	testCmd := flag.NewFlagSet("test", flag.ExitOnError)

	switch os.Args[1] {
	case "test":
		err := testCmd.Parse(os.Args[2:])
		if nil != err {
			log.Panicf("parse cmd of test failed! %v\n", err)
		}
	case "createwallets":
		err := createWalletsCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panicf("parse createwallets failed. %v\n", err)
		}
	case "addrlists":
		err := getAddrListsCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panicf("parse addrlists failed. %v\n", err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panicf("parse send failed. %v\n", err)
		}

	//case "addblock":
	//	err := addBlockCmd.Parse(os.Args[2:])
	//	if err != nil {
	//		log.Panicf("parse addblock failed. %v\n", err)
	//	}
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
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panicf("get balance failed. %v\n", err)
		}
	default:
		PrintUsage()
		os.Exit(1)
	}

	// Parsed()判断是否解析成功

	// 测试
	if testCmd.Parsed() {
		cli.TestMethod()
		cli.TestResetUTXO()
	}

	// 创建钱包
	if createWalletsCmd.Parsed() {
		cli.CreateWallets()
	}

	// 获得钱包地址
	if getAddrListsCmd.Parsed() {
		cli.GetAddrLists()
	}

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

	//// 添加区块
	//if addBlockCmd.Parsed() {
	//	if *flagAddBlockArg == "" {
	//		PrintUsage()
	//		os.Exit(1)
	//	}
	//
	//	// 暂时传的空值
	//	cli.AddBlock([]*TX{})
	//}
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

	// 添加余额查询命令
	if getBalanceCmd.Parsed() {
		if *flagBalanceArg == "" {
			fmt.Println("未指定查询地址...")
			PrintUsage()
			os.Exit(1)
		}
		cli.GetBalance(*flagBalanceArg)
	}

}

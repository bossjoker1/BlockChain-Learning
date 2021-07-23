package main

import "BlockChain-Learning/BC_Basic/BLC"

func main() {
	//blockChain := BLC.CreateBlockChainWithGenesisBlock()
	//blockChain.AddBlock([]byte("Send 100 BTC to James"))
	//blockChain.AddBlock([]byte("Send 10 BTC to Amy"))
	//blockChain.AddBlock([]byte("Send 50 BTC to Bob"))
	//err := blockChain.DB.View(func(tx *bolt.Tx) error {
	//	b := tx.Bucket([]byte(BLC.BLOCKTABLENAME))
	//	if b != nil {
	//		hash := b.Get([]byte("latest"))
	//		latest_block := BLC.DeserializeBlock(b.Get(hash))
	//		fmt.Printf("latest block hash height: %v\n", latest_block.Height)
	//	}
	//	return nil
	//})
	//
	//if err != nil {
	//	log.Panicf("Get the bucket failed. %v\n", err)
	//}
	//blockChain.PrintChainInfo()
	//defer blockChain.DB.Close()
	//BLC.PrintUsage()

	// 命令行测试G
	cli := BLC.CLI{}
	cli.Run()

	//base58编码测试
	//test := BLC.Base58Encode([]byte("fuck bug"))
	//fmt.Printf("%s\n", BLC.Base58Decode(test))

	// 钱包获取地址测试
	//wallet := BLC.NewWallet()
	//addr := wallet.GetAddr()
	//fmt.Printf("addr: %s\n", addr)
	//fmt.Println("isValid? = ", BLC.IsValidforAddr(addr))
}

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
	cli := BLC.CLI{}
	cli.Run()
}

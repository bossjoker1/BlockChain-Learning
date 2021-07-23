package BLC

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// 钱包集合结构定义
type Wallets struct {
	// 地址是无序的
	// key : addr
	Wallets map[string]*Wallet
}

// 创建或者得到钱包集合
func NewWallets(node_id string) (*Wallets, error) {
	// 文件操作
	// 判断文件是否存在
	walletFile := fmt.Sprintf(WALLETFILEPATH, node_id)
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		wallets := &Wallets{}
		wallets.Wallets = make(map[string]*Wallet)
		return wallets, err
	}

	f, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panicf("read file failed. %v\n", err)
	}

	var wallets Wallets

	//  包含interface的解析
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(f))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panicf("decode wallet file failed. %v\n", err)
	}

	return &wallets, nil
}

// 在钱包集合中创建新的钱包
func (wallets *Wallets) CreateWallet(nodeid string) {
	wallet := NewWallet()
	wallets.Wallets[string(wallet.GetAddr())] = wallet

	// 把新创建的钱包存储在文件中
	wallets.SaveWallets(nodeid)
}

// 持久化钱包信息
func (wallets *Wallets) SaveWallets(nodeId string) {
	var content bytes.Buffer
	// 注册
	gob.Register(elliptic.P256())
	// 序列化钱包数据
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(&wallets)
	if err != nil {
		log.Panicf("encode the wallets failed. %v\n", err)
	}
	walletFile := fmt.Sprintf(WALLETFILEPATH, nodeId)
	// 写入文件
	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)

	if err != nil {
		log.Panicf("write file failed. %v\n", err)
	}
}

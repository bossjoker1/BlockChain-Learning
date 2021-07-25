package TX

import (
	"BlockChain-Learning/BC_Basic/BLC"
	"BlockChain-Learning/BC_Basic/Utils"
	"bytes"
)

// 交易输出
type TxOutput struct {
	// 金额总量
	Value int64
	// 钱是谁的
	PubkeyHash []byte
}

// output 身份验证
func (tout *TxOutput) UnLockPubKeyWithAddr(addr string) bool {
	hash160 := Lock(addr)
	return bytes.Compare(tout.PubkeyHash, hash160) == 0
}

// 锁定
func Lock(addr string) []byte {
	pubKey_hash := BLC.Base58Decode([]byte(addr))
	hash160 := pubKey_hash[1 : len(pubKey_hash)-Utils.CHECKSUMLEN]
	//fmt.Printf("hash160 : %x\n", hash160)
	return hash160
}

// 创建output对象
func NewTXOutput(value int64, addr string) *TxOutput {
	txOutput := &TxOutput{}
	hash160 := Lock(addr)
	txOutput.Value = value
	txOutput.PubkeyHash = hash160
	return txOutput
}

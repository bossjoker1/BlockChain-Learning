package TX

import (
	"BlockChain-Learning/BC_Basic/Wallet"
	"bytes"
)

// 交易输入
type TxInput struct {
	// 引用的上一笔交易的hash
	Tx_hash []byte
	// 引用的上一笔交易的output索引
	Index_out int

	// 数字签名
	Signature []byte

	// 公钥
	PublicKey []byte
}

func (in *TxInput) UnLockWithRipemd_SHA(ripemd_sha []byte) bool {
	// 获取双hash值
	pubKey_hash := Wallet.Ripemd160_SHA256(in.PublicKey)
	return bytes.Compare(pubKey_hash, ripemd_sha) == 0
}

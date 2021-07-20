package BLC

// 交易输入
type TxInput struct {
	// 引用的上一笔交易的hash
	Tx_hash []byte
	// 引用的上一笔交易的output索引
	Index_out int
	// 锁定脚本 用户名
	ScriptSig string
}

// 判断能否引用指定地址的OUTPUT
func (in *TxInput) UnLockWithAddr(addr string) bool {
	return in.ScriptSig == addr
}

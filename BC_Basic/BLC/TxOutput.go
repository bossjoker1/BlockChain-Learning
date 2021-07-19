package BLC

// 交易输出
type TxOutput struct {
	// 金额总量
	Value int64
	// 钱是谁的
	ScriptPubkey string
}

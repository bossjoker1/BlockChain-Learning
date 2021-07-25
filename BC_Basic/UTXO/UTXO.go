package UTXO

import (
	"BlockChain-Learning/BC_Basic/TX"
)

type UTXO struct {
	// 对应的交易hash
	Tx_hash []byte
	// 该交易中的index 即OUPUT对应的索引
	Out_index int
	// OUTPUT结构
	Output *TX.TxOutput
}

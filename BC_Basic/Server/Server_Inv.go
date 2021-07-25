package Server

// 展示当前节点有哪些信息

type Inv struct {
	Hashes   [][]byte
	AddrFrom string
	Type     string // 类型（交易或者区块），这里其实用不到
}

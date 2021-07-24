package BLC

type BlockData struct {
	AddrFrom string // 来自哪个节点
	Block    []byte // 序列化的区块数据
}

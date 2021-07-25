package Server

// 代表当前节点版本信息，决定是否同步
type Version struct {
	Version  int    // 版本号
	Height   int64  // 当前节点区块的高度
	AddrFrom string // 当前节点的地址
}

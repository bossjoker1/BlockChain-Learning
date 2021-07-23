package BLC

import "crypto/sha256"

// Merkle树定义

type MerkleTree struct {
	// 根节点
	Root *TreeNode
}

// Merkle节点
type TreeNode struct {
	// 左子节点
	Left *TreeNode
	// 右子节点
	Right *TreeNode
	// 交易数据 | hash指针
	Data []byte
}

// 创建节点
func NewMerkleNode(left, right *TreeNode, data []byte) *TreeNode {
	node := &TreeNode{}
	if left == nil && right == nil {
		// 叶子节点
		hash := sha256.Sum256(data)
		node.Data = hash[:]
	} else {
		// 非叶子节点, 将左右节点的hash合并
		preHash := append(left.Data, right.Data...)
		hash := sha256.Sum256(preHash)
		node.Data = hash[:]

	}
	node.Left = left
	node.Right = right
	return node

}

// 创建树即根节点
// [][]byte表示区块中的所有交易
func NewMerkleTree(datas [][]byte) *MerkleTree {
	var nodes []TreeNode
	// 奇数条需要我们复制最后一条
	if len(datas)%2 != 0 {
		datas = append(datas, datas[len(datas)-1])
	}

	// merkle_tree自底向上建立树
	// 创建叶子节点
	for _, data := range datas {
		node := NewMerkleNode(nil, nil, data)
		nodes = append(nodes, *node)
	}

	// 创建非叶子节点
	for i := 0; i < len(datas)/2; i++ {
		var newNodes []TreeNode
		for j := 0; j < len(nodes); j += 2 {
			node := NewMerkleNode(&nodes[j], &nodes[j+1], nil)
			newNodes = append(newNodes, *node)
		}
		if len(newNodes)%2 == 0 {
			newNodes = append(newNodes, newNodes[len(newNodes)-1])
		}
		nodes = newNodes
	}

	mtree := MerkleTree{&nodes[0]}

	return &mtree
}

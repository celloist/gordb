package main

import "bytes"

const HEADER = 4

const BTREE_PAGE_SIZE = 4096
const BTREE_MAX_KEY_SIZE = 1000
const BTREE_MAX_VAL_SIZE = 3000

type BTree struct {
	root uint64
	get  func(uint64) BNode
	new  func(BNode) uint64
	del  func(uint64)
}

func init() {
	node1max := HEADER + 8 + 2 + 4 + BTREE_MAX_KEY_SIZE + BTREE_MAX_VAL_SIZE
	assertComparison("check nodemax", (node1max <= BTREE_PAGE_SIZE))
}

func treeInsert(tree *BTree, node BNode, key,val []byte)BNode {
	new := BNode{data: make([]byte, 2*BTREE_PAGE_SIZE)}

	idx := nodeLookupLE(node,key)

	switch node.btype() {
	case BNODE_LEAF:
		if bytes.Equal(key, node.getKey(idx)) {
			leafUpdate(new,node,idx,key,val)
		} else {
			leafInsert(new,node,idx+1,key,val)
		}
	case BNODE_NODE:
		nodeInsert(tree,new,node,idx,key,val)

	default:
		panic("bad node!")
	}
	return new

}

func nodeInsert(tree *BTree,new BNode, node BNode,idx uint16,key []byte, val []byte) {
	kptr := node.getPtr(idx)
	knode := tree.get(kptr)
	tree.del(kptr)
	knode =treeInsert(tree,knode,key,val)
	nsplit , splited := nodeSplit3(knode)

	nodeReplaceKidN(tree,new,node,idx,splited[:nsplit]...)
}

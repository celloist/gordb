package main

import "encoding/binary"

const (
	BNODE_NODE = 1
	BNODE_LEAF = 1
)

type BNode struct {
	data []byte
}

func (node BNode) btype() uint16 {
	return binary.LittleEndian.Uint16(node.data)
}

func (node BNode) nkeys() uint16 {
	return binary.BigEndian.Uint16(node.data[2:4])
}

func (node BNode) setHeader(btype, nkeys uint16) {
	binary.LittleEndian.PutUint16(node.data[0:2], btype)
	binary.LittleEndian.PutUint16(node.data[2:4], nkeys)
}

func(node BNode) getPtr(idx uint16) uint64 {
	assertComparison("getPtr",(idx < node.nkeys()))
	pos := HEADER +8*idx
	return binary.LittleEndian.Uint64(node.data[pos:])
}
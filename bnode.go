package main

import (
	"bytes"
	"encoding/binary"
)

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

func (node BNode) getPtr(idx uint16) uint64 {
	assertComparison("getPtr", (idx < node.nkeys()))
	pos := HEADER + 8*idx
	return binary.LittleEndian.Uint64(node.data[pos:])
}

func (node BNode) setPtr(idx uint16, val uint64) {
	assertComparison("setPtr", (idx < node.nkeys()))
	pos := HEADER + 8*idx
	binary.LittleEndian.PutUint64(node.data[pos:], val)
}

func (node BNode) getOffset(idx uint16) uint16 {
	if idx == 0 {
		return 0
	}
	return binary.LittleEndian.Uint16(node.data[offsetPos(node, idx):])
}

func (node BNode) setOffset(idx uint16, offset uint16) {
	binary.LittleEndian.PutUint16(node.data[offsetPos(node, idx):], offset)
}

func (node BNode) kvPos(idx uint16) uint16 {
	assertComparison("kvPos", (idx <= node.nkeys()))
	return HEADER + 10*node.nkeys() + node.getOffset(idx) //this is modified originall 8*x + 2*x
}

func (node BNode) getKey(idx uint16) []byte {
	assertComparison("getKey", (idx < node.nkeys()))
	pos := node.kvPos(idx)
	klen := binary.LittleEndian.Uint16(node.data[pos:])
	return node.data[pos+4:][:klen]
}

func (node BNode) getVal(idx uint16) []byte {
	assertComparison("getVal", (idx < node.nkeys()))
	pos := node.kvPos(idx)
	klen := binary.LittleEndian.Uint16(node.data[pos+0:])
	vlen := binary.LittleEndian.Uint16(node.data[pos+2:])
	return node.data[pos+4+klen:][:vlen]
}

func (node BNode) nbytes() uint16 {
	return node.kvPos(node.nkeys())
}

func offsetPos(node BNode, idx uint16) uint16 {
	assertComparison("offsetPos", (1 <= idx && idx <= node.nkeys()))
	return HEADER + 8*node.nkeys() + 2*(idx-1)
}

func nodeLookupLE(node BNode, key []byte) uint16 {
	nkeys := node.nkeys()
	found := uint16(0)

	for i := uint16(1); i < nkeys; i++ {
		cmp := bytes.Compare(node.getKey(i), key)
		if cmp <= 0 {
			found = i
		}
		if cmp >= 0 {
			break
		}
	}
	return found
}

func leafInsert(new BNode, old BNode, idx uint16, key []byte, val []byte) {
	new.setHeader(BNODE_LEAF, old.nkeys()+1)
	nodeAppendRange(new, old, 0, 0, idx)
	nodeAppendKv(new, idx, 0, key, val)
	nodeAppendRange(new, old, idx+1, idx, old.nkeys()-idx)
}

// todo
func leafUpdate(new BNode, old BNode, idx uint16, key []byte, val []byte) {
	new.setHeader(BNODE_LEAF, old.nkeys()+1)
	nodeAppendRange(new, old, 0, 0, idx)
	nodeAppendKv(new, idx, 0, key, val)
	nodeAppendRange(new, old, idx+1, idx, old.nkeys()-idx)
}

// copy multiple kv into positions
func nodeAppendRange(new, old BNode, dstNew, srcOld, n uint16) {
	assertComparison("nodeAppendRange check old", (srcOld+n <= old.nkeys()))
	assertComparison("nodeAppendRange check new", (dstNew+n <= new.nkeys()))

	if n == 0 {
		return
	}

	for i := uint16(0); i < n; i++ {
		new.setPtr(dstNew+i, old.getPtr(srcOld+i))
	}

	dstBegin := new.getOffset(dstNew)
	srcBegin := old.getOffset(srcOld)

	for i := uint16(1); i <= n; i++ {
		offset := dstBegin + old.getOffset(srcOld+1) - srcBegin
		new.setOffset(dstNew+i, offset)
	}

	begin := old.kvPos(srcOld)
	end := old.kvPos(srcOld + n)
	copy(new.data[new.kvPos(dstNew):], old.data[begin:end])
}

// copy kv into position
func nodeAppendKv(new BNode, idx uint16, ptr uint64, key, val []byte) {
	new.setPtr(idx, ptr)

	pos := new.kvPos(idx)
	binary.LittleEndian.PutUint16(new.data[pos+0:], uint16(len(key)))
	binary.LittleEndian.PutUint16(new.data[pos+2:], uint16(len(val)))

	copy(new.data[pos+4:], key)
	copy(new.data[pos+4+uint16(len(key)):], val)
	new.setOffset(idx+1, new.getOffset(idx)+4+uint16(len(key)+len(val)))
}

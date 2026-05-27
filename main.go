package main

import "encoding/binary"

// Fixed-size pages make space allocation and reuse easier because all deleted nodes are inter-changeable,
// which can be managed with a free list rather than reinventing malloc().
const BTREE_PAGE_SIZE = 4096

const BTREE_MAX_KEY_SIZE = 1000
const BTREE_MAX_VAL_SIZE = 3000

const PTR_SIZE = 8
const METADATA_FIELD_SIZE = 2
const HEADER_SIZE = 4

func main() {

}

// Simple format so no need to deserialize
// Also struct can't be variable-length, but byte array can
type BNode []byte

// Node getters (type and # of keys)
func (node BNode) btype() uint16 {
	return binary.LittleEndian.Uint16(node[0:2])
}

func (node BNode) nkeys() uint16 {
	return binary.LittleEndian.Uint16(node[2:4])
}

func (node BNode) getKey(idx uint16) []byte {
	assert(idx < node.nkeys())
	pos := node.kvPos(idx)
	klen := node.getKeyLen(pos)
	// starting from the position of the kv entry, skip the header and then take klen bytes.
	return node[pos + HEADER_SIZE:][:klen]
}

func (node BNode) getVal(idx uint16) []byte {
	assert (idx < node.nkeys())
	pos := node.kvPos(idx)
	klen := node.getKeyLen(pos)
	vlen := node.getValLen(pos)
	// skip header and key to access value
	return node[pos + HEADER_SIZE + klen:][:vlen]
}

// Node Setters
func (node BNode) SetHeader(btype uint16, nkeys uint16) {
	binary.LittleEndian.PutUint16(node[0:2], btype)
	binary.LittleEndian.PutUint16(node[2:4], nkeys)
}

// Node Get / Set pointers
func (node BNode) getPtr(idx uint16) uint64 {
	assert(idx < node.nkeys())
	pos := ptrPos(idx)
	// Uint64 always returns exactly 8 bytes, starting from the position we found.
	return binary.LittleEndian.Uint64(node[pos:])
}

func (node BNode) setPtr(idx uint16, val uint64) {
	assert(idx < node.nkeys())
	pos := ptrPos(idx)
	binary.LittleEndian.PutUint64(node[pos:], val)
}

// Read the 'offsets' array
func (node BNode) getOffset(idx uint16) uint16 {
	if idx == 0 {
		return 0
	}
	pos := node.offsetPos(idx)
	// Uint16 always returns exactly 2 bytes, starting from the position we found.
	return binary.LittleEndian.Uint16(node[pos:])
}
func (node BNode) getKeyLen(pos uint16) uint16 {
	// the length of the key is stored in the first two bytes of the kv entry (first field of kv entry header)
	return binary.LittleEndian.Uint16(node[pos:])
}
func (node BNode) getValLen(pos uint16) uint16 {
	// second field in kv entry header
	return binary.LittleEndian.Uint16(node[pos + METADATA_FIELD_SIZE:])
}





// POSITION FINDERS

// Since this node is simply an array of bytes, the offsets are stored as 2-byte integers (uint16)
// These come sequentially after the section which comprises pointers to values (keys)
// The pointers are all 8 bits, so we multiply the number of keys by 8 to reach the end of that section of the array
// Then we find the correct position within the virtual subarray of 2-bit offsets

func ptrPos(idx uint16) uint16 {
	return HEADER_SIZE + (PTR_SIZE * idx)
}

func (node BNode) offsetPos(idx uint16) uint16 {
	return HEADER_SIZE + (PTR_SIZE * node.nkeys()) + (METADATA_FIELD_SIZE * (idx - 1))
}

func (node BNode) kvPos(idx uint16) uint16 {
	assert(idx <= node.nkeys())
	return HEADER_SIZE + (PTR_SIZE * node.nkeys()) + node.getOffset(idx)
}


// Utilities

func assert(condition bool) {
	if !condition {
		panic("assertion failed")
	}
}
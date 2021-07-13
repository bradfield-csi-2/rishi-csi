package main

import (
	"math/rand"
)

const (
	MAX_LEVEL  = 16
	MAX_STRING = "~"
)

var NIL = &skipListNode{item: Item{Key: MAX_STRING}}

type skipListNode struct {
	item    Item
	forward [MAX_LEVEL]*skipListNode
}

type skipListOC struct {
	header *skipListNode
	level  int
	p      float32
}

func newSkipListOC() *skipListOC {
	return &skipListOC{
		p:      0.25,
		level:  1,
		header: newSkipListNode("", ""),
	}
}

func newSkipListNode(key, value string) *skipListNode {
	node := &skipListNode{item: Item{key, value}}
	for i := 0; i < MAX_LEVEL; i++ {
		node.forward[i] = NIL
	}
	return node
}

func (o *skipListOC) search(key string) (*skipListNode, [MAX_LEVEL]*skipListNode) {
	update := [MAX_LEVEL]*skipListNode{}
	x := o.header
	for i := o.level - 1; i >= 0; i-- {
		for x.forward[i].item.Key < key {
			x = x.forward[i]
		}
		update[i] = x
	}
	x = x.forward[0]
	return x, update
}

func (o *skipListOC) Get(key string) (string, bool) {
	x, _ := o.search(key)
	if x.item.Key == key {
		return x.item.Value, true
	}
	return "", false
}

func (o *skipListOC) Put(key, value string) bool {
	x, update := o.search(key)

	// Key found, update
	if x.item.Key == key {
		x.item.Value = value
		return false
	}

	// Key not found, insert
	level := o.randomLevel()
	if level > o.level {
		for i := o.level; i < level; i++ {
			update[i] = o.header
		}
		o.level = level
	}
	x = newSkipListNode(key, value)
	for i := 0; i < level; i++ {
		x.forward[i] = update[i].forward[i]
		update[i].forward[i] = x
	}
	return true
}

func (o *skipListOC) Delete(key string) bool {
	x, update := o.search(key)

	// Key not found, bail
	if x.item.Key != key {
		return false
	}

	// Key found, delete
	for i := 0; i < o.level; i++ {
		if update[i].forward[i] != x {
			break
		}
		update[i].forward[i] = x.forward[i]
	}
	for o.level > 1 && o.header.forward[o.level-1] == NIL {
		o.level--
	}
	return true
}

func (o *skipListOC) randomLevel() int {
	level := 1
	for rand.Float32() < o.p && level < MAX_LEVEL {
		level++
	}
	return level
}

func (o *skipListOC) RangeScan(startKey, endKey string) Iterator {
	return &skipListOCIterator{}
}

type skipListOCIterator struct {
}

func (iter *skipListOCIterator) Next() {
}

func (iter *skipListOCIterator) Valid() bool {
	return false
}

func (iter *skipListOCIterator) Key() string {
	return ""
}

func (iter *skipListOCIterator) Value() string {
	return ""
}

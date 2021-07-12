package main

import (
	"math/rand"
)

const MAX_LEVEL = 16

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
	header := &skipListNode{forward: [MAX_LEVEL]*skipListNode{newSkipListNode(1, "~", "~")}}
	return &skipListOC{
		p:      0.25,
		level:  1,
		header: header,
	}
}

func newSkipListNode(level int, key, value string) *skipListNode {
	return &skipListNode{
		item:    Item{Key: key, Value: value},
		forward: [MAX_LEVEL]*skipListNode{},
	}
}

func (o *skipListOC) search(key string) (*skipListNode, bool) {
	x := o.header
	for i := o.level - 1; i >= 0; i-- {
		for x.forward[i].item.Key < key {
			x = x.forward[i]
		}
	}
	x = x.forward[0]
	if x.item.Key == key {
		return x, true
	}
	return nil, false
}

func (o *skipListOC) Get(key string) (string, bool) {
	x, ok := o.search(key)
	if ok {
		return x.item.Value, true
	}
	return "", false
}

func (o *skipListOC) Put(key, value string) bool {
	update := [16]*skipListNode{}

	// Search for key
	x := o.header
	for i := o.level - 1; i >= 0; i-- {
		for x.forward[i].item.Key < key {
			x = x.forward[i]
		}
		update[i] = x
	}
	x = x.forward[0]

	// Key found, overwrite
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
	x = newSkipListNode(level, key, value)
	for i := 0; i < level; i++ {
		x.forward[i] = update[i].forward[i]
		update[i].forward[i] = x
	}
	return true
}

func (o *skipListOC) Delete(key string) bool {
	update := make([]*skipListNode, 10)

	// Search for key
	x := o.header
	for i := o.level - 1; i >= 0; i-- {
		for x.forward[i].item.Key < key {
			x = x.forward[i]
		}
		update[i] = x
	}
	x = x.forward[0]

	// Key nout found, bail
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
	for o.level > 1 && o.header.forward[o.level] == nil {
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

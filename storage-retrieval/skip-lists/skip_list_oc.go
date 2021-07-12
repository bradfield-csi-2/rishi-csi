package main

type skipListOC struct {
}

func newSkipListOC() *skipListOC {
	return &skipListOC{}
}

func (o *skipListOC) Get(key string) (string, bool) {
	return "", false
}

func (o *skipListOC) Put(key, value string) bool {
	return false
}

func (o *skipListOC) Delete(key string) bool {
	return false
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

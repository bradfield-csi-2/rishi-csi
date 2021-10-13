package main

type HashJoinOperator struct {
	outerChild Operator
	innerChild Operator
	key        string
	hashTable  map[string][]Tuple
	results    []Tuple
	idx        int
}

func NewHashJoinOperator(outer, inner Operator, key string) Operator {
	return &HashJoinOperator{
		outerChild: outer,
		innerChild: inner,
		key:        key,
		hashTable:  make(map[string][]Tuple),
		results:    []Tuple{},
		idx:        -1,
	}
}

func (n *HashJoinOperator) Init() {
	// Build hash table of outer
	for n.outerChild.Next() {
		o := n.outerChild.Execute()
		val := getVal(o, n.key)
		if _, ok := n.hashTable[val]; !ok {
			n.hashTable[val] = []Tuple{}
		}
		n.hashTable[val] = append(n.hashTable[val], o)
	}

	// Probe hash table with inner
	for n.innerChild.Next() {
		i := n.innerChild.Execute()
		val := getVal(i, n.key)
		tuples, ok := n.hashTable[val]
		if !ok {
			continue
		}
		for _, tuple := range tuples {
			n.results = append(n.results, combineTuples(tuple, i))
		}
	}
}

func (n *HashJoinOperator) Next() bool {
	n.idx++
	return n.idx < len(n.results)
}

func (n *HashJoinOperator) Execute() Tuple {
	return n.results[n.idx]
}

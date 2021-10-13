package main

type NestedLoopsJoinOperator struct {
	outerChild Operator
	innerChild Operator
	key        string
	results    []Tuple
	idx        int
}

func NewNestedLoopsJoinOperator(outer, inner Operator, key string) Operator {
	return &NestedLoopsJoinOperator{
		outerChild: outer,
		innerChild: inner,
		key:        key,
		results:    []Tuple{},
		idx:        -1,
	}
}

func (n *NestedLoopsJoinOperator) Init() {
	for n.outerChild.Next() {
		for n.innerChild.Next() {
			o := n.outerChild.Execute()
			i := n.innerChild.Execute()

			outerVal := getVal(o, n.key)
			innerVal := getVal(i, n.key)

			if outerVal == innerVal {
				n.results = append(n.results, combineTuples(o, i))
			}
		}
		n.innerChild.Init()
	}
}

func (n *NestedLoopsJoinOperator) Next() bool {
	n.idx++
	return n.idx < len(n.results)
}

func (n *NestedLoopsJoinOperator) Execute() Tuple {
	return n.results[n.idx]
}

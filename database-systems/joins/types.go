package main

// Operator is the interface implemented by all operators.
type Operator interface {
	// Init intializes the operator  in preparation for calling Next/Execute
	Init()
	// Next returns a boolean indicating whether the operator has more work to do.
	Next() bool
	// Execute executes the operation and returns the resulting tuple. Should only
	// be called if Next() returns true.
	Execute() Tuple
}

// Tuple represents a tuple (row) of values.
type Tuple struct {
	Values []Value
}

// Value represents a value and its associated name.
type Value struct {
	Name        string
	StringValue string
}

func combineTuples(a, b Tuple) Tuple {
	vals := []Value{}
	for _, v := range append(a.Values, b.Values...) {
		vals = append(vals, v)
	}
	return Tuple{Values: vals}
}

func getVal(tuple Tuple, col string) string {
	for _, v := range tuple.Values {
		if v.Name == col {
			return v.StringValue
		}
	}
	return ""
}

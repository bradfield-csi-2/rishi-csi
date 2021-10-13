package main

import (
	"testing"
)

func TestNestedLoopsJoinNoMatches(t *testing.T) {
	tuples1 := []Tuple{
		newTuple(
			"id", "bradfieldstudent1",
			"gender", "male"),
		newTuple(
			"id", "bradfieldstudent2",
			"gender", "female"),
		newTuple(
			"id", "bradfieldstudent3",
			"gender", "female"),
	}

	tuples2 := []Tuple{
		newTuple(
			"id", "bradfieldstudent4",
			"name", "Herman"),
		newTuple(
			"id", "bradfieldstudent5",
			"gender", "Maya"),
	}
	r1 := NewScanOperator(tuples1)
	r2 := NewScanOperator(tuples2)

	nlj := NewNestedLoopsJoinOperator(r1, r2, "id")
	nlj.Init()
	assertEq(t, false, nlj.Next())
}

func TestNestedLoopsJoinWithTuples(t *testing.T) {
	tuples1 := []Tuple{
		newTuple(
			"id", "bradfieldstudent1",
			"gender", "male"),
		newTuple(
			"id", "bradfieldstudent2",
			"gender", "female"),
		newTuple(
			"id", "bradfieldstudent3",
			"gender", "female"),
	}

	tuples2 := []Tuple{
		newTuple(
			"id", "bradfieldstudent1",
			"name", "Herman"),
		newTuple(
			"id", "bradfieldstudent2",
			"gender", "Maya"),
		newTuple(
			"id", "bradfieldstudent3",
			"gender", "Emily"),
	}
	r1 := NewScanOperator(tuples1)
	r2 := NewScanOperator(tuples2)

	nlj := NewNestedLoopsJoinOperator(r1, r2, "id")
	nlj.Init()

	for i := 0; i < len(tuples1); i++ {
		tuple := combineTuples(tuples1[i], tuples2[i])
		assertEq(t, true, nlj.Next())
		assertEq(t, tuple, nlj.Execute())
	}
	assertEq(t, false, nlj.Next())
}

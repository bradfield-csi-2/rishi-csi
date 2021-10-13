package main

import (
	"testing"
)

func TestHashJoinNoMatches(t *testing.T) {
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

	hj := NewHashJoinOperator(r1, r2, "id")
	hj.Init()
	assertEq(t, false, hj.Next())
}

func TestHashJoinWithTuples(t *testing.T) {
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

	hj := NewHashJoinOperator(r1, r2, "id")
	hj.Init()

	for i := 0; i < len(tuples1); i++ {
		tuple := combineTuples(tuples1[i], tuples2[i])
		assertEq(t, true, hj.Next())
		assertEq(t, tuple, hj.Execute())
	}
	assertEq(t, false, hj.Next())
}

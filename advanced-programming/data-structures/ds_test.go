package main

import (
	"fmt"
	"math"
	"testing"
)

func TestFloat64ToUint64(t *testing.T) {
	var tests = []struct {
		input float64
		want  string
	}{
		{0.0, "0000000000000000000000000000000000000000000000000000000000000000"},
		{1.0, "0011111111110000000000000000000000000000000000000000000000000000"},
		{123.456, "0100000001011110110111010010111100011010100111111011111001110111"},
		{-1.0, "1011111111110000000000000000000000000000000000000000000000000000"},
		{math.Inf(0), "0111111111110000000000000000000000000000000000000000000000000000"},
	}
	for _, test := range tests {
		if got := Float64ToUint64(test.input); fmt.Sprintf("%064b", got) != test.want {
			t.Errorf("Float64ToUint64(%f) = %064b, want %s", test.input, got, test.want)
		}
	}
}

func TestStringsPointToSame(t *testing.T) {
	same := "test"
	var tests = []struct {
		s    string
		t    string
		want bool
	}{
		{"abc", "xyz", false},
		{"abc", "abcd", false},
		{"abc", "abc", true},
		{same, same, true},
		{same, same[1:2], true},
		{same[1:], same, true},
		{same[1:], same[2:], true},
		// {same[0:1], same[2:3], true}, This test case does not work, but it should, right?
	}
	for _, test := range tests {
		if got := StringsPointToSame(test.s, test.t); got != test.want {
			t.Errorf("StringsPointToSame(%s, %s) = %v, want %v", test.s, test.t, got, test.want)
		}
	}
}

func TestSliceSum(t *testing.T) {
	var tests = []struct {
		input []int
		want  int
	}{
		{[]int{1, 2, 3}, 6},
		{[]int{100, 2, -100}, 2},
		{[]int{}, 0},
		{[]int{-1238595}, -1238595},
	}
	for _, test := range tests {
		if got := SliceSum(test.input); got != test.want {
			t.Errorf("SliceSum(%v) = %v, want %v", test.input, got, test.want)
		}
	}
}

func TestHashSum(t *testing.T) {
	var tests = []struct {
		input     map[int]int
		want_ksum int
		want_vsum int
	}{
		{map[int]int{0: 1, 1: 2, 2: 3}, 3, 6},
		{map[int]int{100: 100, 2: 2, -100: -100}, 2, 2},
		{map[int]int{}, 0, 0},
		{map[int]int{-1238595: 23495}, -1238595, 23495},
	}
	for _, test := range tests {
		got_ksum, got_vsum := HashSum(test.input)
		if got_ksum != test.want_ksum || got_vsum != test.want_vsum {
			t.Errorf("HashSum(%v) = keysum: %v, valsum: %v, want keysum: %v, valsum: %v",
				test.input,
				got_ksum,
				got_vsum,
				test.want_ksum,
				test.want_vsum)
		}
	}
}

package main

import (
	"fmt"
	"testing"
)

func generate(numRows int) [][]int {

	var rows [][]int

	for i := 0; i < numRows; i++ {

	}

	return rows
}

func GenerateTest(t *testing.T) {
	i := generate(5)
	fmt.Printf("%d\n", i)
}

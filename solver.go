package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func Solver(graph [][]int, squares int, colors int, constraints map[int]int) ([]int, error) {
	nodes := make([]int, 0)
	if squares <= 0 || colors < 0 {
		return nodes, nil
	}
	nodes = append(nodes, -1)
	for len(nodes) != 0 {
		nodeid := len(nodes) - 1
		if constraint, ok := constraints[nodeid]; ok {
			nodes[nodeid] = constraint
			if len(nodes) == squares {
				return nodes, nil
			}
			nodes = append(nodes, -1)
			continue
		}
		nodes[nodeid]++
		if nodes[nodeid] < colors {
			unique := true
			for _, neighbor := range graph[nodeid] {
				if neighbor < len(nodes) && nodes[neighbor] == nodes[nodeid] {
					unique = false
					break
				} else if constraint, ok := constraints[neighbor]; ok && constraint == nodes[nodeid] {
					unique = false
					break
				}
			}
			if unique && len(nodes) == squares {
				return nodes, nil
			}
			if unique {
				nodes = append(nodes, -1)
			}
		} else {
			// strip this node off
			nodes = nodes[:nodeid]
			// continue to strip off nodes while there are still nodes and
			// the back most node exists in the constraints map.
			for len(nodes) > 0 {
				if _, ok := constraints[nodes[len(nodes)-1]]; ok {
					nodes = nodes[:len(nodes)-1]
				} else {
					break
				}
			}
		}
	}
	return nil, errors.New("No solution possible")
}

func byteConvert(val byte) int {
	if val >= 49 && val <= 57 {
		return int(val - 49)
	} else if val >= 65 && val <= 90 {
		return int(val-65) + 9
	} else if val >= 97 && val <= 122 {
		return int(val-97) + 9 + 26
	} else {
		return -1
	}
}

func main() {
	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Println("usage: sudoko-solver <input_file string> <order int>\norder is bounded by 49 as an upperlimit")
		os.Exit(1)
	}

	order, ordererr := strconv.Atoi(args[1])
	if ordererr != nil || order > 49 {
		fmt.Println("order is bounded by 49 as an upper limit and must be a valid integer")
		os.Exit(1)
	}
	sqrtOrder := 1
	for sqrtOrder < 8 && sqrtOrder*sqrtOrder != order {
		sqrtOrder++
	}
	if sqrtOrder == 8 {
		fmt.Println("order must be a perfect square.")
		os.Exit(1)
	}
	inputFile, openerr := os.Open(args[0])
	if openerr != nil {
		fmt.Println(openerr)
		os.Exit(1)
	}
	fileBytes, readerr := ioutil.ReadAll(inputFile)
	if readerr != nil {
		fmt.Println(readerr)
		os.Exit(1)
	}

	puzzle := make([][]int, 1)
	puzzle[0] = make([]int, 0)
	for _, b := range fileBytes {
		if b == 10 || b == 0 { // new line
			puzzle = append(puzzle, make([]int, 0))
		} else {
			puzzle[len(puzzle)-1] = append(puzzle[len(puzzle)-1], byteConvert(b))
		}
	}
	if len(puzzle) != order {
		fmt.Printf("input file formatted incorrectly, expected %d nodes, but found %d\n", order, len(puzzle))
		os.Exit(1)
	}

	constraints := make(map[int]int)
	for i, row := range puzzle {
		if len(row) != order {
			fmt.Printf("input file formatted incorrectly at row %d\n", i)
			os.Exit(1)
		}
		for j, sq := range row {
			if sq >= order {
				fmt.Printf("input file formatted incorrectly at row %d, column %d where %d is greater than or equal to order %d in value\n", i, j, sq, order)
			}
			if sq != -1 {
				constraints[i*order+j] = sq
			}
		}
	}

	graph := make([][]int, order*order)
	for i := range graph {
		row := i / order
		col := i % order
		section_row := (row / sqrtOrder) * sqrtOrder
		section_col := (col / sqrtOrder) * sqrtOrder
		graph[i] = make([]int, 3*(order-1))
		edges := 0
		for x := 0; x < order; x++ {
			if x != row {
				graph[i][edges] = x*order + col
				edges++
			}
			if x != col {
				graph[i][edges] = row*order + x
				edges++
			}
			if (x/sqrtOrder)+section_row != row || (x%sqrtOrder)+section_col != col {
				graph[i][edges] = ((x/sqrtOrder)+section_row)*order + ((x % sqrtOrder) + section_col)
				edges++
			}
		}
	}
	// for k, v := range constraints {
	// 	fmt.Printf("[%d] -> %d\n", k, v)
	// }
	solution, solvererr := Solver(graph, order*order, order, constraints)
	if solvererr != nil {
		fmt.Println(solvererr)
		return
	}

	fmt.Println("Solution:")
	for i := 0; i < order; i++ {
		if i%sqrtOrder == 0 {
			fmt.Print("+")
			fmt.Printf(strings.Repeat(strings.Repeat("-", (sqrtOrder*3)+sqrtOrder-1)+"+", sqrtOrder))
			fmt.Print("\n|")
		} else {
			fmt.Print("|")
			fmt.Printf(strings.Repeat("-", (order*3)+order-1))
			fmt.Print("|\n|")
		}
		for j := 0; j < order; j++ {
			fmt.Printf(" %d |", solution[i*order+j]+1) // + 1 to re-adjust for indexes
		}
		fmt.Print("\n")
	}
	fmt.Print("+")
	fmt.Printf(strings.Repeat(strings.Repeat("-", (sqrtOrder*3)+sqrtOrder-1)+"+", sqrtOrder))
	fmt.Print("\n")
}

/*

+------------
| 2 | 3 | 1 |
|------------
| 2 | 3 | 1 |
|------------
| 2 | 3 | 1 |
+------------

*/

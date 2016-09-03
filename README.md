# Sudoku Solver

Solves a sudoku puzzle that is encoded in a file and passed to the program. It outputs the completed puzzle solution to `stdout`. The solution treats the Sudoko puzzle as a graph coloring problem, where each unit is considered a node in an undirected graph. Each node is given an edge between itself and all other nodes in it's row and column, as well as every node in it's local square. The starting values of nodes are considered additional constraints that takes the possiblility of `(9^2)^9` combinations to try in a `9x9` square down to a manageable number of states that are practically possible to go through.

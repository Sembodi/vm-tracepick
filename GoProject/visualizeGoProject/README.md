# Once run using `make visualize` and `bin/visualizeGoProject`:

In the go file (visualizeGoProject.go), one can specify in the variable `rootDirs` which project components should be included in the graph. 

## Run dot on it:
`cd visualizeGoProject/output`
`dot -T png output.dot -o graph.png`

or:

`dot -T png visualizeGoProject/output/output.dot -o visualizeGoProject/output/graph.png`

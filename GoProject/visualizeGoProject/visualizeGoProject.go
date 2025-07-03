package main

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/token"
	// "go/ast"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"project1/tracepick/lib/own/owncheckers"
)

const (
	projectPrefix = "project1/tracepick/"
	outputFile    = "visualizeGoProject/output/output.dot"
)
var rootDirs []string  = []string{"lib/helpers/","lib/collectors/","lib/fileparsers/"} //,"cmd/"} // Root of the project

func main() {
	var err error
	deps := make(map[string]map[string]struct{}) // map[importingPackage]set[importedPackage]
	for _, rootDir := range rootDirs {
		err = filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") {
				return err
			}
			// if !strings.HasPrefix(path, "lib/") && !strings.HasPrefix(path, "cmd/") {
				// 	return nil
				// }

			dir := filepath.Dir(path)
			relDir, err := filepath.Rel(".", dir)
			if err != nil {
				return err
			}
			pkgPath := strings.TrimPrefix(relDir, "./")
			importingPkg := pkgPath

			fset := token.NewFileSet()
			node, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed parsing %s: %v\n", path, err)
				return nil
			}

			for _, imp := range node.Imports {
				impPath := strings.Trim(imp.Path.Value, `"`)
				if strings.HasPrefix(impPath, projectPrefix) {
					importedPkg := strings.TrimPrefix(impPath, projectPrefix)
					if importingPkg != importedPkg {
						if deps[importingPkg] == nil {
							deps[importingPkg] = make(map[string]struct{})
						}
						if owncheckers.StringHasOneOfPrefixes(importedPkg, rootDirs) {
							deps[importingPkg][importedPkg] = struct{}{}
						}
					}
				}
			}
			return nil
		})
	}
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	buf.WriteString("digraph G {\n")
	buf.WriteString("  rankdir=LR;\n")
	buf.WriteString("  node [shape=box, fontname=Helvetica];\n")

	// Nodes
	allPkgs := make(map[string]struct{})
	for from, tos := range deps {
		allPkgs[from] = struct{}{}
		for to := range tos {
			allPkgs[to] = struct{}{}
		}
	}
	for pkg := range allPkgs {
		fmt.Fprintf(&buf, "  \"%s\";\n", pkg)
	}

  clusters := make(map[string][]string)
  for pkg := range allPkgs {
  	parts := strings.Split(pkg, "/")
  	group := parts[0]
  	if len(parts) > 1 {
  		group = parts[0] + "/" + parts[1] // e.g., lib/helpers
  	}
  	clusters[group] = append(clusters[group], pkg)
  }

  // Write subgraphs for each cluster
  clusterCount := 0
  for group, pkgs := range clusters {
  	fmt.Fprintf(&buf, "  subgraph cluster_%d {\n", clusterCount)
  	fmt.Fprintf(&buf, "    label=\"%s\";\n", group)
  	for _, pkg := range pkgs {
  		fmt.Fprintf(&buf, "    \"%s\";\n", pkg)
  	}
  	buf.WriteString("  }\n")
  	clusterCount++
  }

	// Edges
	for from, tos := range deps {
		for to := range tos {
			fmt.Fprintf(&buf, "  \"%s\" -> \"%s\";\n", from, to)
		}
	}

	buf.WriteString("}\n")
	err = os.WriteFile(outputFile, buf.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Package import graph written to %s\n", outputFile)
}

package packageHelper

import (
  "fmt"
  "strings"
  "regexp"
  "project1/tracepick/lib/own/owntypes"
  "project1/tracepick/lib/own/ownformatters"
  "project1/tracepick/lib/fileparsers/lineparse"
  "project1/tracepick/lib/helpers/fsHelper"
  "project1/tracepick/lib/helpers/optionalHelpers/bindepHelper"
)

func TrimPackageName(pkg string) string {
  return strings.Split(pkg, "-")[0]
}

func TrimPackageSet(dir string, pkgSlice []string) []string {
  programName := "TrimPackageSet"
  result := pkgSlice

  filesMap := lineparse.LineSlicesFromFiles(programName, dir)

	// Escape each package name for regex
	escaped := make([]string, len(pkgSlice))
	for i, pkg := range pkgSlice {
		escaped[i] = regexp.QuoteMeta(pkg)
	}

	// Build pattern
  // pattern := `\b(` + strings.Join(escaped, "|") + `)\b\s*->\s*\b(` + strings.Join(escaped, "|") + `)\b`
  pattern := `"` + `(` + strings.Join(escaped, "|") + `)` + `"` + `\s*->\s*"` + `(` + strings.Join(escaped, "|") + `)` + `"`

	re := regexp.MustCompile(pattern)

  for _, lines := range filesMap {
    for _, line := range lines {
      matches := re.FindStringSubmatch(line)
      // fmt.Println("Matches: ", matches)
      if len(matches) > 2 {
        rmTarget := matches[2] // The captured word after " -> "
        // fmt.Println("package removed:", rmTarget)
        result = ownformatters.RemoveSliceElemByValue(result, rmTarget) //remove package from list
      }
		}
  }
  return result
}

func SelectPackages(pkgList owntypes.StringSet, item string) owntypes.StringSet {
  pkgSlice := ownformatters.FromStringSetToSlice(pkgList)
  rawPath := "outputs/rdepends/graph/raw/"
  graphPath := "outputs/rdepends/graph/png/"

  fsHelper.CleanPath(graphPath)
  bindepHelper.MakeRecursiveDepGraphs(pkgList)  // saves package dependencies to outputs/rdepends/graph/raw
  // bindepHelper.PrintDependencyGraph(combinedGraph)

  result := TrimPackageSet(rawPath, pkgSlice)
  fmt.Println("Selected packages:", result)
  return ownformatters.FromSliceToStringSet(result)
}

// combinedGraph, err := bindepHelper.BuildDependencyGraph(rawPath)
// if err != nil {
  // 	fmt.Println(programName, ": error parsing dotfiles:", err)
  //   return nil
  // }
  //
  // result := bindepHelper.FindMinimalInstallSet(pkgSlice, combinedGraph)
  //
  // fmt.Println("Minimal install set:")
  // for _, pkg := range result {
    // 	fmt.Println(pkg)
    // }

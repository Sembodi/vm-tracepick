package main

import (
	"fmt"
	"regexp"
	"strings"
)

func main() {
	packages := []string{
		"libmariadb3:amd64",
		"mariadb-client-10.6",
		"mariadb-client-core-10.6",
		"mariadb-common",
		"mariadb-server",
		"mariadb-server-10.6",
		"mariadb-server-core-10.6",
	}

	// Escape each package name for regex
	escaped := make([]string, len(packages))
	for i, pkg := range packages {
		escaped[i] = regexp.QuoteMeta(pkg)
	}

	// Regex: capture the word before '->' if followed by a known package name
	pattern := `\b(\w+)\b\s*->\s*(?:` + strings.Join(escaped, "|") + `)\b`
	re := regexp.MustCompile(pattern)

	// Example test lines
	lines := []string{
		"node1 -> mariadb-client-10.6",
		"node2 -> not-this-one",
		"db1 -> mariadb-server-core-10.6",
		"irrelevant line",
		"thing -> mariadb-server",
	}

	// Match and print only the captured "source" words
	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if len(matches) > 1 {
			fmt.Println(matches[1]) // The captured word before ->
		}
	}
}

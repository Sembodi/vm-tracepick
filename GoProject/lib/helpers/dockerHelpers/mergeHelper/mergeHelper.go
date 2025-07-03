package mergeHelper

import (
  "fmt"
  "os"
  "strings"
  "project1/tracepick/lib/own/owntypes"
  "project1/tracepick/lib/own/owncheckers"
  "project1/tracepick/lib/fileparsers/lineparse"
  "project1/tracepick/lib/helpers/fsHelper"
  "project1/tracepick/lib/editors/fileedit"
)

func FuseDockerfiles(inPath string, d1name string, d2name string) {
  var (
    outPath string = fmt.Sprintf("artifacts/output/containers/%[1]s_%[2]s/", d1name, d2name)
    d1path string = inPath + d1name
    d2path string = inPath + d2name
    fileName string = fmt.Sprintf("%sDockerfile", outPath)
  )

  result := mergeDockerfiles(inPath, d1path, d2path)

  fsHelper.CleanPath(outPath)
  file, err := os.Create(fileName)
  if err != nil {
    fmt.Println("Dockerfile creation not successful:", err.Error())
    os.Exit(1)
  }

  defer file.Close() // Make sure to close it when done

  fmt.Println("Unfixed Dockerfile:")
  for _, line := range result {
    fmt.Println(line)
    file.WriteString(line + "\n")
  }

  FixMerged(fileName, d1path, d2path, d1name, d2name, outPath)
}

// mergeDockerfiles merges two Dockerfile line slices d1 and d2
// It keeps unique lines only once, but preserves order and blocks between common "anchor" lines.
func mergeDockerfiles(inPath string, d1path string, d2path string) []string {
  var (
    programName string = "mergeDockerfiles"
    d1lines []string
    d2lines []string

    er error
  )

  if d1lines, er = lineparse.ReadFileLinesSlice(programName, d1path + "/Dockerfile"); er != nil {
    fmt.Println("Error reading", d1path + "/Dockerfile, aborting program...")
    os.Exit(1)
    return nil
  }

  if d2lines, er = lineparse.ReadFileLinesSlice(programName, d2path + "/Dockerfile"); er != nil {
    fmt.Println("Error reading", d2path + "/Dockerfile, aborting program...")
    os.Exit(1)
    return nil
  }

	anchors := findAnchors(d1lines, d2lines)
	d1Blocks := splitByAnchors(d1lines, anchors)
	d2Blocks := splitByAnchors(d2lines, anchors)

	var merged []string
	seen := owntypes.StringSet{}

	addLines := func(lines []string) {
		for _, line := range lines {
			if !seen[line] {
				merged = append(merged, line)
				seen[line] = true
			}
		}
	}

	for i := 0; i < len(anchors)+1; i++ {
		addLines(d1Blocks[i])
		addLines(d2Blocks[i])
		if i < len(anchors) {
			if !seen[anchors[i]] {
				merged = append(merged, anchors[i])
				seen[anchors[i]] = true
			}
		}
	}
  return merged
}

func FixMerged(fileName, d1path, d2path, d1name, d2name, outPath string) {
  var (
    programName string = "FixMerged"
  )


  // Copy root filesystems accordingly
  outFSname := fmt.Sprintf("myrootfs-%[1]s_%[2]s", d1name, d2name)
  outFSpath := outPath + outFSname
  fsHelper.CreatePath(outFSpath)
  fsHelper.CopyRootFS(d1path + "/myrootfs-" + d1name + "/.", outFSpath)
  fsHelper.CopyRootFS(d2path + "/myrootfs-" + d2name + "/.", outFSpath)
  fsHelper.CopyRootFS(d1path + "/helperFiles", outPath)


  copyFSPattern := "(?m)^COPY myrootfs.*$"
  replacement := strings.Join([]string{fmt.Sprintf("COPY %s/ /bbackup/", outFSname),
                                        "COPY helperFiles/ "}, "\n")
  fileedit.EditFile(fileName, copyFSPattern, replacement, "ReplaceAll")

  cmdPattern := "(?m)^CMD .*$"
  cmdMatches := fileedit.GetAllMatchingString(fileName, cmdPattern)

  var sb strings.Builder
  for _, line := range cmdMatches {
    sb.WriteString(dockerCMDtoStartCmd(line) + " &\n")
  }
  cmdScript := strings.TrimSuffix(sb.String(), "&\n")

  outStr := mergeCMDEntrypoints(programName, d1path, d2path, cmdScript)
  outFileName := outPath + "entrypoint.sh"
  fsHelper.CreateExecFile(outFileName, outStr)


  cmdReplacement := strings.Join([]string{`COPY entrypoint.sh /entrypoint.sh`,
                                          `ENTRYPOINT ["/entrypoint.sh"]`},
                                  "\n")
  fileedit.EditFile(fileName, cmdPattern, cmdReplacement, "ReplaceAll")
}

//finds the common lines between the two dockerfiles, makes sure that all blocks between the identical lines are still between the lines
func findAnchors(d1, d2 []string) []string {
	d2Set := owntypes.StringSet{}
	for _, line := range d2 {
		d2Set[line] = true
	}

	var anchors []string
	seen := owntypes.StringSet{}
	for _, line := range d1 {
		if d2Set[line] && !seen[line] {
			anchors = append(anchors, line)
			seen[line] = true
		}
	}
	return anchors
}

func splitByAnchors(lines []string, anchors []string) [][]string {
	var blocks [][]string
	block := []string{}
	anchorSet := owntypes.StringSet{}
	for _, a := range anchors {
		anchorSet[a] = true
	}

	for _, line := range lines {
		if anchorSet[line] {
			blocks = append(blocks, block)
			block = []string{}
		} else {
			block = append(block, line)
		}
	}
	blocks = append(blocks, block)
	return blocks
}

// groupByInstruction puts lines into a map by their Dockerfile instruction keyword
func groupByInstruction(lines []string) map[string][]string {
	grouped := make(map[string][]string)
	for _, line := range lines {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "#") {
			// Group comments under empty string key to keep them at end
			grouped[""] = append(grouped[""], line)
			continue
		}
		parts := strings.Fields(trim)
		if len(parts) == 0 {
			// grouped[""] = append(grouped[""], line)
			continue
		}
		instr := strings.ToUpper(parts[0])
		grouped[instr] = append(grouped[instr], line)
	}
	return grouped
}

func dockerCMDtoStartCmd(dockerCMD string) string {
  dockerCMD = strings.TrimSpace(dockerCMD)
  dockerCMD = strings.TrimPrefix(dockerCMD, "CMD [")
  dockerCMD = strings.TrimSuffix(dockerCMD, "]")
  dockerCMD = strings.ReplaceAll(dockerCMD, `"`, ``)
  parts := strings.Split(dockerCMD, ",")
  dockerCMD = strings.Join(parts, " ")
  return dockerCMD
}

func mergeCMDEntrypoints(programName string, d1path string, d2path string, cmdScript string) string {

  entrypointStr := "entrypoint.sh"
  script1path := d1path + "/" + entrypointStr
  script2path := d2path + "/" + entrypointStr

  has1, err1 := owncheckers.FileExists(script1path)
  if err1 != nil {
    return "err"
  }
  has2, err2 := owncheckers.FileExists(script2path)
  if err2 != nil {
    return "err"
  }

  var (
    lines1 []string
    lines2 []string
    result string = "#!/bin/bash \n"
  )
  if has1 {
    lines1, _ = lineparse.ReadFileLinesSlice(programName, script1path)
    lines1[len(lines1)-1] = lines1[len(lines1)-1] + " &"
  }
  if has2 {
    lines2, _ = lineparse.ReadFileLinesSlice(programName, script2path)
  }

  switch {
  case has1 && has2:
    allLines := append(lines1, lines2...)
    result = result + strings.Join(allLines, "\n")
  case has1 && !has2:
    result = result + strings.Join(append(lines1, cmdScript), "\n")
  case !has1 && has2:
    cmdScript = cmdScript + " &"
    result = result + strings.Join(append([]string{cmdScript}, lines2...), "\n")
  default:
    result = result + cmdScript
  }
  return result
}

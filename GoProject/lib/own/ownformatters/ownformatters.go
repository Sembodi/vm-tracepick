package ownformatters

import (
  "fmt"
  "strings"
  "strconv"
  "unicode"
  "project1/tracepick/lib/own/owntypes"
  "project1/tracepick/lib/helpers/fsHelper"
)

func FromStringToStringSet(item string) owntypes.StringSet {
  var out owntypes.StringSet = make(owntypes.StringSet)
  out[item] = true
  return out
}

func FromStringLinesToStringSet(linesStr string) owntypes.StringSet {
  var out owntypes.StringSet = make(owntypes.StringSet)
  lines := strings.Split(linesStr, "\n")
  for _, line := range lines {
    if line != "" { out[line] = true }
  }
  return out
}

func FromStringSetToString(items owntypes.StringSet) string {
  var out strings.Builder
  for item, _ := range items {
    out.WriteString(item)
    out.WriteString("\n")
  }
  return strings.TrimSpace(out.String())
}

func FromStringSetToSlice(items owntypes.StringSet) []string {
  var out []string
  for item, _ := range items {
    out = append(out, item)
  }
  return out
}

func FromSliceToStringSet(items []string) owntypes.StringSet {
  var out owntypes.StringSet = make(owntypes.StringSet)
  for _, item := range items {
    out[item] = true
  }
  return out
}

func FromStringSetToSpacedStrings(items owntypes.StringSet) string {
  var out strings.Builder
  for item, _ := range items {
    out.WriteString(item)
    out.WriteString(" ")
  }
  return strings.TrimSpace(out.String())
}

func FromCommandMapToString(items owntypes.CommandMap) string {
  var out strings.Builder
  for _, cmd := range items {
    out.WriteString(cmd)
    out.WriteString("\n")
  }
  return strings.TrimSpace(out.String())
}

func FromStringToCommandMap(cmd string, fileName string) owntypes.CommandMap {
  var out owntypes.CommandMap = make(owntypes.CommandMap)
  out[fileName] = cmd
  return out
}

// func concatCommands(startCmd string, comms owntypes.StringSet, endCmd string) string {
func ConcatCommands(comms owntypes.CommandMap, fileName string) owntypes.CommandMap {
  var (
    allCmds strings.Builder
    result owntypes.CommandMap = make(owntypes.CommandMap)
  )

  // allCmds.WriteString(startCmd + "; ")
  for _, comm := range comms {
    allCmds.WriteString(comm + "; ")
  }
  // allCmds.WriteString(endCmd + ";")
  result[fileName] = allCmds.String()

  return result
}

func FormatCommands(fcomm string, items owntypes.StringSet) owntypes.CommandMap {
  var (
    cmd string
    allCmds owntypes.CommandMap = make(owntypes.CommandMap)
  )

  for item, _ := range items {
    cmd = fmt.Sprintf(fcomm, item)
    allCmds[item] = cmd
  }
  return allCmds
}

func MergeStringSets(ss1 owntypes.StringSet, ss2 owntypes.StringSet) owntypes.StringSet {
  resultSS := ss1
  for item, _ := range ss2 {
    resultSS[item] = true
  }
  return resultSS
}

func ReverseStringSetMap(mp map[string]owntypes.StringSet) map[string]owntypes.StringSet {
  visitedStrings := make(owntypes.StringSet)
  result := make(map[string]owntypes.StringSet)
  for key, ss := range mp {
    for item, _ := range ss {
      //guarantee that stringset value exists
      if !visitedStrings[item] {
        result[item] = make(owntypes.StringSet)
        visitedStrings[item] = true
      }
      result[item][key] = true
    }
  }
  return result
}

// ---- filesystem operations

// Formats filenames into CommandMap
func MapExecFromFiles(programName string, path string, formatCmd string) owntypes.CommandMap {
  files := fsHelper.GetFiles(programName, path)

  items := fsHelper.GetFileNames(files)

  fcomms := FormatCommands(formatCmd, items)
  return fcomms
}

// ----

//if below does not work, just return cmd to revert changes
func MakeSudo(cmd string) string {
  if owntypes.UseSudo { return fmt.Sprintf("sudo %s", cmd) }
  return cmd
}

// func Pop[T any](items *[]T) (T, bool) {
//   var zero T
//   if len(*items) == 0 {
//       return zero, false
//   }
//
//   result := (*items)[0]
//   *items = (*items)[1:]
//   return result, true
// }


// Strips version numbers from a set of service strings:
func StripServiceSet(items owntypes.StringSet) owntypes.StringSet {
  result := make(owntypes.StringSet)
  for item, _ := range items {
    result[StripServiceString(item)] = true
  }
  return result
}

// Strips version number from a service string
func StripServiceString(svc string) string {
  for i, c := range svc {
    if unicode.IsDigit(c) { return svc[:i] }
  }
  return svc
}

func FromMBtoBytes(mbStr string) int {
  numStr := strings.TrimSuffix(strings.TrimSpace(mbStr),"MB")

  num, err := strconv.Atoi(numStr)

  if err != nil {
    fmt.Println("MBtoByte conversion from string to int unsuccessful: ", err)
    return 0
  }

  return num * 1048576
}

func RemoveSliceElemByValue(slice []string, value string) []string {
	for i, v := range slice {
		if v == value {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice // value not found; return original
}

package owncheckers

import (
    "os"
    "strings"
    "project1/tracepick/lib/own/owntypes"
)

func StringHasOneOfPrefixes(target string, prefixes []string) bool {
  for _, prefix := range prefixes {
    if strings.HasPrefix(target, prefix) { return true }
  }
  return false
}

//target: bpfCMD, slice: services, result: service StringSet
func StringPartOfElementOrViceVersa(target string, slice owntypes.StringSet) (owntypes.StringSet, bool) {
  var (
    resultSet owntypes.StringSet = make(owntypes.StringSet)
    resultBool bool = false
  )

  targetLower := strings.ToLower(target)
  for item, _ := range slice {
      itemLower := strings.ToLower(item)
      if strings.Contains(itemLower, targetLower) || strings.Contains(targetLower, itemLower) {
        resultSet[item] = true
        resultBool = true
      }
  }
  return resultSet, resultBool
}

func StringSetElemPartOfString(items owntypes.StringSet, target string) (owntypes.StringSet, bool) {
  var (
    resultSet owntypes.StringSet = make(owntypes.StringSet)
    resultBool bool = false
  )

  targetLower := strings.ToLower(target)
  for item, _ := range items {
      if strings.Contains(targetLower, strings.ToLower(item)) {
        resultSet[item] = true
        resultBool = true
      }
  }
  return resultSet, resultBool
}

func StringPartOfSliceElem(target string, slice []string) bool {
  targetLower := strings.ToLower(target)
  for _, item := range slice {
      if strings.Contains(strings.ToLower(item), targetLower) {
          return true
      }
  }
  return false
}

// Helper to check if a file exists
func FileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	return !info.IsDir(), err
}

package lineparse

import (
  "fmt"
  "io/ioutil"
  "strings"
  "project1/tracepick/lib/own/owntypes"
  "project1/tracepick/lib/helpers/fsHelper"
)

func splitContent(programName string, fileName string) ([]string, error) {
  var (
    content string
    err error
  )
  // Read the file:
  if content, err = ReadContent(programName, fileName); err != nil { return nil, err }

  // Split the content by lines and add each line to the slice
	lines := strings.Split(string(content), "\n")
  return lines, nil
}


func ReadFileLines(programName string, fileName string) (owntypes.StringSet, error) {
  // Create a new StringSet for this file
  var (
    lines []string
    err error

    result owntypes.StringSet = make(owntypes.StringSet)
  )

  // Split the content by lines and add each line to the slice
	if lines, err = splitContent(programName, fileName); err != nil { return nil, err }

	for _, line := range lines {
		if line != "" { result[line] = true }
	}
	return result, nil
}

func ReadFileLinesSlice(programName string, fileName string) ([]string, error) {
  var (
    lines []string
    err error
    outSlice []string
  )

  if lines, err = splitContent(programName, fileName); err != nil { return nil, err }

  for _, line := range lines {
    if line != "" { outSlice = append(outSlice, line) }
  }

  return outSlice, nil
}

func ReadContent(programName string, fileName string) (string, error) {
  content, err := ioutil.ReadFile(fileName)
	if err != nil {
		// fmt.Printf("%s: Error reading file: %s", programName, err.Error())
		return "", err
	}
  return strings.TrimSpace(string(content)), err
}

// Puts file lines into StringSets (and maps them to the key filename)
func LineSetsFromFiles(programName string, path string) map[string]owntypes.StringSet {
  var (
    result map[string]owntypes.StringSet = make(map[string]owntypes.StringSet)
  )

  // Read the contents of the directory
	files := fsHelper.GetFiles(programName, path)

	// Iterate over the files in the directory
	for _, file := range files {
		// Skip directories (only process files)
		if file.IsDir() {
			continue
		}

		// Open the file and create the stringSet containing the lines
		filePath := path + file.Name()
    stringSet, parseErr := ReadFileLines(programName, filePath)

    if parseErr != nil {
      fmt.Println("Skipping: ", filePath)
      continue
    }

		// Add the StringSet to the result map with the filename as the key
		result[file.Name()] = stringSet
	}
  return result
}

// Puts file lines into slices (and maps them to the key filename)
func LineSlicesFromFiles(programName string, path string) map[string][]string {
  var (
    result map[string][]string = make(map[string][]string)
  )

  // Read the contents of the directory
	files := fsHelper.GetFiles(programName, path)

	// Iterate over the files in the directory
	for _, file := range files {
		// Skip directories (only process files)
		if file.IsDir() {
			continue
		}

		// Open the file and create the stringSet containing the lines
		filePath := path + file.Name()
    lines, parseErr := ReadFileLinesSlice(programName, filePath)

    if parseErr != nil {
      fmt.Println("Skipping: ", filePath)
      continue
    }

		// Add the []string to the result map with the filename as the key
		result[file.Name()] = lines
	}
  return result
}

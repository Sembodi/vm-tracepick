package fileedit

import(
  "fmt"
  "regexp"
  "os"
  "project1/tracepick/lib/fileparsers/lineparse"
  "project1/tracepick/lib/helpers/fsHelper"

)

//Insert before first appearance of pattern
func InsertBeforeString(content string, re *regexp.Regexp, added string) string {
  found := false
  newContent := re.ReplaceAllStringFunc(content, func(match string) string {
                                                    if found {
                                                      return match
                                                    }
                                                    found = true
                                                    return added + match
                                                  })
  return newContent
}

//Insert after first appearance of pattern
func InsertAfterString(content string, re *regexp.Regexp, added string) string {
  found := false
  newContent := re.ReplaceAllStringFunc(content, func(match string) string {
                                                    if found {
                                                      return match
                                                    }
                                                    found = true
                                                    return match + added
                                                  })
  return newContent
}

//Replace first appearance of pattern and remove all others
func ReplaceAllString(content string, re *regexp.Regexp, added string) string {
  found := false
  newContent := re.ReplaceAllStringFunc(content, func(match string) string {
                                                    if found {
                                                      return ""
                                                    }
                                                    found = true
                                                    return added
                                                  })
  return newContent
}


func GetAllMatchingString(fileName string, pattern string) []string {
  var matches []string
  content, re := ContentAndRegExp(fileName, pattern)
  re.ReplaceAllStringFunc(content, func(match string) string {
                                                    matches = append(matches, match)
                                                    return match //change nothing
                                                  })
  return matches
}

func ContentAndRegExp(fileName string, pattern string) (string, *regexp.Regexp) {
  var (
    programName string = "ContentAndRegExp"
  )

  content, err := lineparse.ReadContent(programName, fileName)
  if err != nil {
    fmt.Println(programName, ": Reading content failed:", err)
  }

  re := regexp.MustCompile(pattern)
  return content, re
}

func EditFile(fileName string, pattern string, added string, action string) error {
  var newContent string
  content, re := ContentAndRegExp(fileName, pattern)
  switch action {
  case "InsertBefore":
    newContent = InsertBeforeString(content, re, added)
  case "InsertAfter":
    newContent = InsertAfterString(content, re, added)
  case "ReplaceAll":
    newContent = ReplaceAllString(content, re, added)
  // case "PopAllMatching":
  //   newContent = PopAllMatchingString(content, re)
  default:
    fmt.Println("EditFile: specify one of the following actions: InsertBefore, InsertAfter, ReplaceAll")
  }
  fsHelper.RemovePath(fileName)
  file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
  file.WriteString(newContent)
  file.Close()

  return err
  //ioutil.WriteFile(filename, []byte(newContent), 0644)
}

package owndialogues

import (
  "fmt"
  "strings"
  "bufio"
  "regexp"
  "project1/tracepick/lib/own/owntypes"
  "project1/tracepick/lib/metrictrackers/timetracker"
)

func EditString(scanner *bufio.Scanner, input string, force bool) string {
  if owntypes.SkipDialogue && !force { return input }
  var (
    currentStr string = input

    pattern string
    replacement string
  )

  if currentStr == "" {
    currentStr = "nil"
  }

  timetracker.LogEvent("BeginDialogue")

  for true {
    fmt.Println("String: '", currentStr, "'")
    // Ask if the user wants to modify it
    fmt.Print("Do you want to change this string? (y/n): ")

    scanner.Scan()

    if strings.ToLower(scanner.Text()) == "y" {
      // Prompt the user for the new command
      fmt.Print("Enter the pattern you want to replace/remove: ")
      scanner.Scan()

      pattern = scanner.Text()

      fmt.Print("Enter the replacement (leave empty for removal): ")
      scanner.Scan()

      replacement = scanner.Text()

      // Replace the option with the new value
      re := regexp.MustCompile(`\b` + pattern + `\b`) // \b is for the boundaries (match whole word)
      currentStr = re.ReplaceAllString(currentStr, replacement)
      // startCmd = strings.Replace(startCmd, pattern, replacement, -1)
      } else {
        // Print the final command
        fmt.Println("Final string:")
        fmt.Println(currentStr)
        break
      }
  }

  timetracker.LogEvent("EndDialogue")
  return currentStr
}

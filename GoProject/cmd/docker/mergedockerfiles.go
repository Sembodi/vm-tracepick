package main

import (
  "fmt"
  "os"
  "bufio"
  "project1/tracepick/lib/helpers/dockerHelpers/mergeHelper"
  "project1/tracepick/lib/own/owntypes"
  "project1/tracepick/lib/own/owndialogues"
)

// MERGE TWO DOCKER FILES and save as "name1_name2"
func main() {
  fmt.Println("MAKE SURE to RUN this with ROOT PRIVILEGES . ")
  fmt.Println("MAKE SURE that the LATTER SERVICE is the one that should run in the FOREGROUND.\n",
              "And MAKE SURE that BOTH FOLDERS are located at the SAME PATH.")

  if len(os.Args) < 3 {
		fmt.Println("usage: merge [service 1] [service 2]")
		os.Exit(1)
	}

  var (
    inPath string = "artifacts/output/containers/"
    d1name string = os.Args[1]
    d2name string = os.Args[2]
  )

  if !owntypes.SkipDialogue { fmt.Println("Specify repository names:") }
  scanner := bufio.NewScanner(os.Stdin) //Going to need user input
  inPath = owndialogues.EditString(scanner, inPath, false)
  d1name = owndialogues.EditString(scanner, d1name, false)
  d2name = owndialogues.EditString(scanner, d2name, false)
  mergeHelper.FuseDockerfiles(inPath, d1name, d2name)
}

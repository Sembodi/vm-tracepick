package main

import (
  "fmt"
  "os"
  "os/exec"
  "bufio"
  "project1/tracepick/lib/own/owndialogues"
  "project1/tracepick/lib/own/owntypes"
)

func main() {
  if len(os.Args) < 2 {
		fmt.Println("usage: removedocker [service name]")
		os.Exit(1)
	}

  var (
    programName string = "removedocker.go"
    d1name string = os.Args[1]
    formatRemoveCmd string = "docker stop c-%[1]s && docker rm c-%[1]s" // && docker rmi i-%[1]s"
    cmd *exec.Cmd
  )

  if !owntypes.SkipDialogue { fmt.Println("Specify service name:") }
  scanner := bufio.NewScanner(os.Stdin) //Going to need user input
  d1name = owndialogues.EditString(scanner, d1name, false)

  removeCmdStr := fmt.Sprintf(formatRemoveCmd, d1name)
  cmd = exec.Command("bash", "-c", removeCmdStr)

  cmd.Stderr = os.Stderr
  cmd.Stdout = os.Stdout

  if err := cmd.Run(); err != nil {
    fmt.Println(programName, ": Docker remove failed:", err)
  }
}

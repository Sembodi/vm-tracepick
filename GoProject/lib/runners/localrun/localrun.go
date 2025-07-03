package localrun

import (
  "fmt"
  "strings"
  "os"
	"os/exec"

  "github.com/google/shlex" //external


  "project1/tracepick/lib/own/owntypes"
  "project1/tracepick/lib/helpers/fsHelper"
)

func RunCmd(cmd *exec.Cmd, programName string, fileName string) error {
  var (
    file   *os.File
    err     error
    )
  // fmt.Println("Executing command...")
  cmd.Stderr = os.Stderr
  if fileName != "" {
    cmd.Stdout, err = os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
      return returnLocalExecError(programName, cmd.String(), err)
    }
    defer file.Close()

    if err = cmd.Run(); err != nil {
      return returnLocalExecError(programName, cmd.String(), err)
    }

    // fmt.Println("Output redirected to", fileName)

  } else {
    //cmd.Stdout = os.Stdout
    if err = cmd.Run(); err != nil {
      return returnLocalExecError(programName, cmd.String(), err)
    }
    // fmt.Println("No output redirected")
  }
  return nil
}

func SingleLocalRun(programName string, cmdStr string, path string, fileName string) {
  var cmdMap owntypes.CommandMap = make(owntypes.CommandMap)

  fsHelper.CleanPath(path)

  cmdMap[fileName] = cmdStr

  LocalRun(programName, cmdMap, path, false)
}

func LocalRun(programName string, fcomms owntypes.CommandMap, outDir string, forceNoOutput bool) {
  var (
    mainprogram string = "LocalRun in "
    fileName string
    parts []string
    cmd *exec.Cmd
    err error
  )

  for item, fcomm := range fcomms {
    if forceNoOutput {item = ""}
    fileName = outDir + item

    parts, err = shlex.Split(fcomm) // strings.Fields(fcomm)

    cmd = exec.Command(parts[0], parts[1:]...)

    if err = RunCmd(cmd, mainprogram + programName, fileName); err != nil {
      fmt.Printf("Error: %s\nContinue run for other items...\n", err)
    }
  }
}


func returnLocalExecError(programName string, cmdStr string, err error) error {
  parts := strings.Split(cmdStr, " -")
  cmdName := parts[0]
  return fmt.Errorf("returnLocalExecError in program %s:\n%s: %s", programName, cmdName, err.Error())
}

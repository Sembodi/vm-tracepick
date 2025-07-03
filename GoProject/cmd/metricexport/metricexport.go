package main

import (
  "fmt"
  "os"
  "os/exec"
  "bufio"
  "project1/tracepick/lib/own/owntypes"
  "project1/tracepick/lib/own/owndialogues"
  "project1/tracepick/lib/helpers/fsHelper"
)

func main() {
  if len(os.Args) < 2 {
		fmt.Println("usage: export [VM name]")
		os.Exit(1)
	}
  var (
    programName string = "metricexport.go"
    metricFilePath string = owntypes.MetricPath + owntypes.MetricFileName
    outPath string = "keepfiles/metrics_dynamic/"
    vmName string = os.Args[1]
    cmd *exec.Cmd
  )

  if !owntypes.SkipDialogue { fmt.Println("Specify VM name:") }
  scanner := bufio.NewScanner(os.Stdin) //Going to need user input
  vmName = owndialogues.EditString(scanner, vmName, false)

  fsHelper.CreatePath(outPath)

  cmd = exec.Command("bash", "-c", fmt.Sprintf("cp %[1]s %[2]s/%[3]s.yaml", metricFilePath, outPath, vmName))

  if err := cmd.Run(); err != nil {
    fmt.Println(programName, ": export failed:", err)
  }

  ExportArtifacts(vmName)
}

func ExportArtifacts(vmName string) {
  sourcePath := "artifacts/output/containers/"
  targetPath := "keepfiles/artifacts_dynamic/"
  fsHelper.CreatePath(targetPath)
  fsHelper.CopyRootFS(sourcePath, targetPath + vmName)
}

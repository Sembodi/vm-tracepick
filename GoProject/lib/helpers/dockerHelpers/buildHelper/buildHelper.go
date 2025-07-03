package buildHelper

import (
  "fmt"
  "os"
  "os/exec"
)

func CachedBuild(inPath, d1name string) {
  var (
    programName string = "CachedBuild"
    repoPath string = inPath + d1name
    cmdStr string = fmt.Sprintf("docker build -t i-%[1]s %[2]s", d1name, repoPath)
    cmd *exec.Cmd = exec.Command("bash", "-c", cmdStr)
  )

  cmd.Stderr = os.Stderr
  cmd.Stdout = os.Stdout

  if err := cmd.Run(); err != nil {
    fmt.Println(programName, ": Docker build failed:", err)
  }
}

func NonCachedBuild(inPath, d1name string) {
  var (
    programName string = "NonCachedBuild"
    cmdStr string = "docker builder prune --all -f"
    cmd *exec.Cmd = exec.Command("bash", "-c", cmdStr)
  )
  if err := cmd.Run(); err != nil {
    fmt.Println(programName, ": Docker UNCACHED build failed:", err)
  }
  CachedBuild(inPath, d1name)
}

package sizetracker

import (
  "fmt"
  "strings"
  "os"
  "os/exec"
  "strconv"
  "project1/tracepick/lib/own/owntypes"
  "project1/tracepick/lib/own/ownformatters"
  "project1/tracepick/lib/runners/remoterun"
  "project1/tracepick/lib/external/scmd"
  "project1/tracepick/lib/external/yamlio"
)

func runSizeCmd(formatCmd string, item string) int {
  //make sure that this computes runtime format size
  cmd := exec.Command("bash", "-c",
    fmt.Sprintf(formatCmd, item),
  )

  out, err := cmd.CombinedOutput()
  if err != nil {
    fmt.Println("Size calculation failed:", err)
    return 0
  }

  size, err := strconv.Atoi(strings.TrimSpace(string(out)))
  if err != nil {
    fmt.Println("Size calculation failed:", err)
    return 0
  }

  return size
}

func TotalRuntimeSize(strippedItems owntypes.StringSet) {
  formatCmd := "sudo /usr/bin/du -s artifacts/output/containers/%[1]s/myrootfs-%[1]s | awk '{print $1}'"

  for item, _ := range strippedItems {
    // fmt.Println("TotalRuntimeSize calculation for", item)
    data := owntypes.EvaluationData.ArtifactDataMap[item]
    data.TotalRuntimeSize = runSizeCmd(formatCmd, item) * 512  // from unit-512 to Bytes
    owntypes.EvaluationData.ArtifactDataMap[item] = data
  }
}

// Execute next function at docker build, so make sure command string is included in dockerfile
func AddedRuntimeSizeCmd() string {
  //outputs which files/folders from 'bbackup/' are NOT located in /
  return "rsync -rcn --out-format='%l %n' /bbackup/ / | awk '{sum += $1} END {print sum}' > /added-runtime.log"
}

// Make sure all containers are running AS "c-<service>"
func AddedRuntimeSize(strippedItem string) {
  formatCmd := "docker exec -t c-%s /bin/cat /added-runtime.log"

  data := owntypes.EvaluationData.ArtifactDataMap[strippedItem]
  data.AddedRuntimeSize = runSizeCmd(formatCmd, strippedItem)
  owntypes.EvaluationData.ArtifactDataMap[strippedItem] = data
}

func TotalFSSize(strippedItem string) {
  formatCmd := `docker exec -t c-%s '/bin/bash' '-c' 'du -s / 2> /dev/null' | awk '{print $1}'`

    data := owntypes.EvaluationData.ArtifactDataMap[strippedItem]
    data.TotalFSSize = runSizeCmd(formatCmd, strippedItem) * 1024
    owntypes.EvaluationData.ArtifactDataMap[strippedItem] = data
}

// Only run when images are built
func ImageSize(strippedItem string) {
  formatCmd := "docker image inspect i-%s --format='{{.Size}}'"

  data := owntypes.EvaluationData.ArtifactDataMap[strippedItem]
  data.ImageSize = runSizeCmd(formatCmd, strippedItem)
  owntypes.EvaluationData.ArtifactDataMap[strippedItem] = data
}

func VMSize() {
  var (
    pmap remoterun.ConnectMap
    err  error
  )

  if pmap, err = yamlio.ReadYaml[remoterun.ConnectMap](os.Args[1]); err != nil {
    fmt.Println(err)
    return
  }

  cmdStr := ownformatters.MakeSudo("du -s / 2> /dev/null | awk '{print $1}'")

  var remoteCmd *exec.Cmd

  for _, form := range pmap {
    remoteCmd = scmd.RemoteCommand(form.Connection, form.Identity, cmdStr)
  }

  out, err := remoteCmd.CombinedOutput()

  size, err := strconv.Atoi(strings.TrimSpace(string(out)))
  if err != nil {
    fmt.Println("VMSize calculation failed:", err)
    return
  }

  owntypes.EvaluationData.GeneralData.VMSize = size * 1024
}

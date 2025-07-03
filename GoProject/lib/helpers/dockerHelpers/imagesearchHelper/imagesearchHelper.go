package imagesearchHelper

// TODO: MAKE SURE THAT docker search with no results does not crash the program (i.e. don't index the results or give it an if-statement before indexing)
import (
  "fmt"
  "strings"
  "os/exec"
  "project1/tracepick/lib/own/owntypes"
  "project1/tracepick/lib/own/ownformatters"
  "project1/tracepick/lib/helpers/fsHelper"
  "project1/tracepick/lib/fileparsers/lineparse"
  // "project1/tracepick/lib/runners/remoterun"
  "project1/tracepick/lib/runners/localrun"
)

func ListDockerImages(items owntypes.StringSet) {
  var (
    programName string = "ListDockerImages"
    formatSearchCmd string = "docker search '%s' --limit 3 --format '{{.Name}} {{.StarCount}}'"
    //      | sort -k 2 -nr | head -1
    outDir string = "outputs/dockersearch/"
    fcomms owntypes.CommandMap
  )
  fcomms = ownformatters.FormatCommands(formatSearchCmd, items)

  fsHelper.CleanPath(outDir)
  // remoterun.RemoteRun(programName, fcomms, outDir)

  localrun.LocalRun(programName, fcomms, outDir, false)
}

func ChooseDockerImages() owntypes.DockerMap {
  var (
    programName string = "ChooseDockerImages"
    path string = "outputs/dockersearch/"
    outfromfiles map[string][]string
    result owntypes.DockerMap = make(owntypes.DockerMap)
  )

  outfromfiles = lineparse.LineSlicesFromFiles(programName, path)

  for item, lines := range outfromfiles {
    fields := strings.Fields(lines[0])
    if len(fields) > 0 { result[item] = fields[0] }
  }
  return result
}

func PrintDockerMap(images owntypes.DockerMap) {
  fmt.Println("Found base images from all items:")
  for item, image := range images {
    fmt.Println(item, ": ", image)
  }
}

func FindClosestVersion(imageName string, version string) (string, error) {
  fmt.Println("Running shell script to check Docker tag...")
  out, err := exec.Command("bash", "bashscripts/checkversions.sh", imageName, version).Output()

	if err != nil { return "", err}

  return string(out), nil
}

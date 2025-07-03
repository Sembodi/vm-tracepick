package etccollector

import (
  // "fmt"
  "project1/tracepick/lib/fileparsers/pathparse"
  // "project1/tracepick/lib/helpers/fsHelper"
  // "project1/tracepick/lib/own/owntypes"
  // "project1/tracepick/lib/own/ownformatters"
  // "project1/tracepick/lib/runners/localrun"
  // "project1/tracepick/lib/runners/remoterun"
)

func EtcvarMinimalPaths() map[string][]string {
  var (
    path string = "outputs/etcvarusr/"
  )

  folderProfiles := pathparse.GetFolderProfiles(path)
  minPaths := pathparse.GetMinimalPathsMap(folderProfiles)

  return minPaths
}

package etcvarHelper

import (
  "project1/tracepick/lib/own/owntypes"
  "project1/tracepick/lib/own/ownformatters"
  "project1/tracepick/lib/fileparsers/lineparse"
  "project1/tracepick/lib/helpers/fsHelper"
  "project1/tracepick/lib/runners/remoterun"
)

//TODO: ADD /VAR to find command, but this should not take only the second folder for filetree analysis (pathparser)
func RetrieveConfigurationFiles(items owntypes.StringSet) {
  var (
    programName string    = "RetrieveConfigurationFiles"
    formatCmd string      = ownformatters.MakeSudo("find {/etc,/var,/opt,/srv} -type f -name '*%[1]s*'") //{/usr,/bin}
    outDir string         = "outputs/etcvarusr/"
    fcomms owntypes.CommandMap = make(owntypes.CommandMap)
  )

  fsHelper.CleanPath(outDir)

  fcomms = ownformatters.FormatCommands(formatCmd, ownformatters.StripServiceSet(items))

  remoterun.RemoteRun(programName, fcomms, outDir)
}


func ParseConfigurationFiles() map[string]owntypes.StringSet {
  var (
    programName string = "ParseConfigurationFiles"
    path string = "outputs/etcvarusr/"
  )
  return lineparse.LineSetsFromFiles(programName, path)
}

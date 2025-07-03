package bindepHelper

import (
  // "fmt"
  "project1/tracepick/lib/own/owntypes"
  "project1/tracepick/lib/own/ownformatters"
  "project1/tracepick/lib/fileparsers/lineparse"
  "project1/tracepick/lib/runners/localrun"
  "project1/tracepick/lib/runners/remoterun"
  "project1/tracepick/lib/helpers/fsHelper"
)


func RetrieveEssentialBinaries(items owntypes.StringSet) { //returns map of e.g. {nginx: /usr/bin/nginx}
  var (
    programName string    = "RetrieveEssentialBinaries"
    formatCmd string      = ownformatters.MakeSudo("which %[1]s")
    outDir string         = "outputs/which/"
    // fileName string       = "which"

    allComms owntypes.CommandMap = make(owntypes.CommandMap)
  )

  fsHelper.CleanPath(outDir)

  allComms = ownformatters.FormatCommands(formatCmd, items) //format first into  every separate command
  // allComms = ownformatters.ConcatCommands(allComms, fileName) //concatenate all commands so all output will be sent to the file in outputs/which

  remoterun.RemoteRun(programName, allComms, outDir)
}

func ParseEssentialBinaries() (owntypes.StringSet, error) {
  var (
    programName string = "ParseEssentialBinaries"
    fileName string = "outputs/which/which"
  )
  return lineparse.ReadFileLines(programName, fileName)
}

func RetrieveDynamicDependencies(items owntypes.StringSet) {
  var (
    programName string = "RetrieveDynamicDependencies"
    formatCmd string = ownformatters.MakeSudo("ldd $(which %[1]s) | awk ' { print $3 }'")
    outDir string = "outputs/ldd/"
    fcomms owntypes.CommandMap
  )

  fsHelper.CleanPath(outDir)

  fcomms = ownformatters.FormatCommands(formatCmd, items)

  remoterun.RemoteRun(programName, fcomms, outDir)
}

func ParseDynamicDependencies() map[string]owntypes.StringSet {
  var (
    programName string = "ParseDynamicDependencies"
    path string = "outputs/ldd/"
  )
  return lineparse.LineSetsFromFiles(programName, path)
}


//OS SPECIFIC
func RetrieveRecursiveDependencies(items owntypes.StringSet) {
  var (
    programName string = "RetrieveRecursiveDependencies"
    formatCmd string = "apt-rdepends %s | awk '{print $2 }'"
    outDir string = "outputs/rdepends/list/"
    fcomms owntypes.CommandMap
  )

  fsHelper.CleanPath(outDir)

  fcomms = ownformatters.FormatCommands(formatCmd, items)

  remoterun.RemoteRun(programName, fcomms, outDir)
}

//OS SPECIFIC
// Retrieves raw output for parsing a nice dependency graph:
func RetrieveRecursiveDepGraphs(items owntypes.StringSet) {
  var (
    programName string = "RetrieveRecursiveDepGraph"
    formatCmd string = "apt-rdepends -d %s"
    outDir string = "outputs/rdepends/graph/raw/"
    fcomms owntypes.CommandMap
  )

  fsHelper.CleanPath(outDir)

  fcomms = ownformatters.FormatCommands(formatCmd, items)

  remoterun.RemoteRun(programName, fcomms, outDir)
}

func DrawRecursiveDepGraphs() {
  var (
    programName string = "DrawRecursiveDepGraphs"
    basepath string = "outputs/rdepends/graph/"
    rawpath string = basepath + "raw/"
    pngpath string = basepath + "png/"
    formatCmd string = "dot -Tpng " + rawpath + "%[1]s -o " + pngpath + "%[1]s.png"

    fcomms owntypes.CommandMap
  )

  fcomms = ownformatters.MapExecFromFiles(programName, rawpath, formatCmd)

  localrun.LocalRun(programName, fcomms, owntypes.NoOutputFile, true)
}

func MakeRecursiveDepGraphs(input owntypes.StringSet) {
  RetrieveRecursiveDepGraphs(input)
  DrawRecursiveDepGraphs()
}




































//end

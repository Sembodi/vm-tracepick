package serviceHelper

import (
  "fmt"
  "project1/tracepick/lib/own/owntypes"
  "project1/tracepick/lib/own/ownformatters"
  "project1/tracepick/lib/fileparsers/lineparse"
  "project1/tracepick/lib/helpers/fsHelper"
  "project1/tracepick/lib/helpers/vmHelpers/osHelper"
  "project1/tracepick/lib/runners/remoterun"

)

//TODO make sure that mongod is found by service helper
var MainServiceStr string = `systemctl list-units --type=service | grep -v systemd | grep running | awk '{print $1}' | grep -v '@' | sed 's/\.service$//'`
var AltServiceStr string = `service --status-all | egrep '\[ \+ \]' | awk '{print $4}'`

func RetrieveUpServices(osName string) {
  var (
    programName string = "RetrieveUpServices"
    cmdStr string      = ownformatters.MakeSudo(`systemctl list-units --type=service | grep -v systemd | grep running | awk '{print $1}' | grep -v '@' | sed 's/\.service$//'`)
    outDir string      = "outputs/services/"
  )

  fsHelper.CleanPath(outDir)
  remoterun.SingleRemoteRun(programName, cmdStr, outDir, osName)
}

func GetAddedServices(osName string) owntypes.StringSet {
  var (
    programName string              = "GetAddedServices"
    defaultFilename string          = "defaultdata/defaultservices/" + osHelper.SimplifyOSName(osName)
    excludeFilename string          = "config/excludeservices"
    servicesFilename string         = "outputs/services/" + osName

    defaultServices owntypes.StringSet
    excludeServices owntypes.StringSet
    currentServices owntypes.StringSet

    result owntypes.StringSet = make(owntypes.StringSet)
    err error
  )


  if defaultServices, err = lineparse.ReadFileLines(programName, defaultFilename); err != nil {
    fmt.Println("Error reading ", defaultFilename, ": \n\n", err, "\n\n")
  }

  if excludeServices, err = lineparse.ReadFileLines(programName, excludeFilename); err != nil {
    fmt.Println("Error reading ", excludeFilename, ": \n\n", err, "\n\n")
  }

  if currentServices, err = lineparse.ReadFileLines(programName, servicesFilename); err != nil {
    fmt.Println("Error reading ", servicesFilename, ": \n\n", err, "\n\n")
  }



  for service, _ := range currentServices {
    _, existsDefault := defaultServices[service]
    _, existsExclude := excludeServices[service]
    if !existsDefault && !existsExclude {
      result[service] = true
    }
  }

  return result
}

func RestartService(item string) {
  var (
    programName string              = "RestartService"
    restartStr string               = ownformatters.MakeSudo("systemctl restart " + item)
  )
  fmt.Println(programName, item, "started...")

  remoterun.SingleRemoteRun(programName, restartStr, owntypes.NoOutputFile, owntypes.NoOutputFile)

  fmt.Println(programName, item, "done.")
}

func ShowUsersFromServices(items owntypes.StringSet) {
  var (
    programName string              = "ShowUsersFromServices"
    outPath string                  = "outputs/users/"
    formatCmd string                = ownformatters.MakeSudo("systemctl show -p User --value %s")
  )

  showUserCmdMap := ownformatters.FormatCommands(formatCmd, items)

  fsHelper.CleanPath(outPath)
  remoterun.RemoteRun(programName, showUserCmdMap, outPath)
}

func ParseUsersFromServices() map[string][]string {
  var (
    programName string              = "ParseUsersFromServices"
    path string                     = "outputs/users/"
  )

  return lineparse.LineSlicesFromFiles(programName, path)
}

func ShowDaemonsFromServices(items owntypes.StringSet) {
  var (
    programName string              = "ShowDaemonsFromServices"
    outPath string                  = "outputs/daemons/"
    formatCmd string                = ownformatters.MakeSudo("systemctl show -p ExecStart --value %s | awk '{print $2}' | xargs basename")
    // if empty output, keep old daemon name. If new output, take that

    // or use systemctl show -p ExecStart --value mariadb | awk '{print $2}' | grep -o '/.*' to get the path (bad idea)
  )



  showDaemonCmdMap := ownformatters.FormatCommands(formatCmd, items)


  fsHelper.CleanPath(outPath)
  remoterun.RemoteRun(programName, showDaemonCmdMap, outPath)
}

func ParseDaemonsFromServices() map[string][]string {
  var (
    programName string              = "ParseDaemonsFromServices"
    path string                     = "outputs/daemons/"
  )

  return lineparse.LineSlicesFromFiles(programName, path)
}

func ShowPreCommandsFromServices(items owntypes.StringSet) {
  var (
    programName string              = "ShowPreCommandsFromServices"
    outPath string                  = "outputs/precommands/"
    formatCmd string                = ownformatters.MakeSudo("systemctl show -p ExecStartPre --value %s")
  )

  showDaemonCmdMap := ownformatters.FormatCommands(formatCmd, items)

  fsHelper.CleanPath(outPath)
  remoterun.RemoteRun(programName, showDaemonCmdMap, outPath)
}

func ParsePreCommandsFromServices() map[string][]string {
  var (
    programName string              = "ParsePreCommandsFromServices"
    path string                     = "outputs/precommands/"
  )

  return lineparse.LineSlicesFromFiles(programName, path)
}


// UNUSED: ----------------------------------------------------------------------

func RestartServices(items owntypes.StringSet) {
  var (
    programName string              = "RestartServices"
    formatCmd string                = ownformatters.MakeSudo("systemctl restart %s")
  )

  restartCmdMap := ownformatters.FormatCommands(formatCmd, items)
  restartCmdMap = ownformatters.ConcatCommands(restartCmdMap, owntypes.NoOutputFile)  // run all commands, and create no output file

  fmt.Println("THE SERVICE RESTART COMMAND:")
  fmt.Println(ownformatters.FromCommandMapToString(restartCmdMap))

  remoterun.RemoteRun(programName, restartCmdMap, owntypes.NoOutputFile)
}

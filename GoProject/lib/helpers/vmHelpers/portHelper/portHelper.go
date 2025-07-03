package portHelper

import (
  // "fmt"
  "project1/tracepick/lib/own/owntypes"
  "project1/tracepick/lib/own/ownformatters"
  "project1/tracepick/lib/fileparsers/lineparse"
  "project1/tracepick/lib/helpers/fsHelper"
  "project1/tracepick/lib/runners/remoterun"
)


func RetrievePortNumbers(items owntypes.StringSet) {
  var (
    programName string = "GetPortNumbers"
    formatCmd string = ownformatters.MakeSudo(`ss -tulnp | grep -E "$(pids=$(pidof %[1]s); [[ -n "$pids" ]] && echo "$pids" | sed 's/ /|/g' || echo $(systemctl status %[1]s | grep 'Main PID' | awk '{print $3}') || echo %[1]s)" | awk '{print $5}' | sed -E 's/.*:([0-9]+)$/\1/' | sort -u`) //ownformatters.MakeSudo(`ss -tulnp | grep -E "$(pidof %[1]s | sed 's/ /|/g' || echo %[1]s)" | awk '{print $5}' | sed -E 's/.*:([0-9]+)$/\1/' | sort -u`)
    outDir = "outputs/ports/"
    fcomms owntypes.CommandMap
  )

    fsHelper.CleanPath(outDir)

    fcomms = ownformatters.FormatCommands(formatCmd, items)

    remoterun.RemoteRun(programName, fcomms, outDir)
}

func ParsePortMap(items owntypes.StringSet) owntypes.PortMap {
  var (
    programName string = "ParsePortMap"
    path = "outputs/ports/"

    result = make(owntypes.PortMap)
    portSet owntypes.StringSet
    err error
  )

  for item, _ := range items {
    if portSet, err = lineparse.ReadFileLines(programName, path + item); err != nil { continue }
    result[item] = portSet
  }

  return result
}

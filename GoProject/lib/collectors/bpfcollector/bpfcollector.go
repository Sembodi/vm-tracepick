package bpfcollector

import (
  "fmt"
  "project1/tracepick/lib/own/owntypes"
  "project1/tracepick/lib/fileparsers/bpfparse"
  "project1/tracepick/lib/fileparsers/pathparse"
)

func BpfMinimalPathsServiceMap(items owntypes.StringSet) (owntypes.StringSet, map[string][]string, map[string]owntypes.StringSet) {
	var (
    fileName string = "outputs/bpfout/bpfout_wrkCmd"

    events map[string]owntypes.OpenSnoopEvent

    err error
  )

  if events, err = bpfparse.GenerateEventsFromLog(fileName); err != nil {
    fmt.Println("Error (main): ", err)
    return nil, nil, nil
  }

  allComms := bpfparse.ListCommands(events)

  serviceMap, eventMap := bpfparse.MakeServicePathMap(events, items)
  eventFolderMap := bpfparse.MakeFolderMap(eventMap)
  eventMinPaths := pathparse.GetMinimalPathsMap(eventFolderMap)

  // fmt.Println(eventFolderMap)
  // owntypes.PrintStringSet(allComms)

  return allComms, eventMinPaths, serviceMap
}

func BpfServiceStartMinPaths(items owntypes.StringSet) map[string][]string {
  var (
    formatFileName string = "outputs/bpfout/bpfout_%s"

    events map[string]owntypes.OpenSnoopEvent
    err error

    folderProfiles map[string]*owntypes.Folder = make(map[string]*owntypes.Folder)
  )

  for item, _ := range items {
    fileName := fmt.Sprintf(formatFileName, item)
    if events, err = bpfparse.GenerateEventsFromLog(fileName); err != nil {
      fmt.Println("Error (main): ", err)
      return nil
    }


    paths := bpfparse.ListPaths(events)
    folderStructure := pathparse.BuildFolderStructure(paths)
    folderProfiles[item] = folderStructure
  }
  return pathparse.GetMinimalPathsMap(folderProfiles)
}

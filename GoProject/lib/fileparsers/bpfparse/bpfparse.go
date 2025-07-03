package bpfparse

import (
	"bufio"
	"fmt"
  "os"
  "errors"
	"strings"
	"strconv"
	"project1/tracepick/lib/own/owntypes"
	"project1/tracepick/lib/own/owncheckers"
	"project1/tracepick/lib/fileparsers/pathparse"
)

// Function to generate a unique key for deduplication
func generateKey(event owntypes.OpenSnoopEvent) string {
	return event.PID + "|" + event.Command + "|" + event.Path
}


func GenerateEventsFromLog(filePath string) (map[string]owntypes.OpenSnoopEvent, error) {
  file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error opening file: ", err))
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	events := make(map[string]owntypes.OpenSnoopEvent)


	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)

		// Ignore lines that don't have the expected number of columns
		if len(parts) < 5 {
			continue
		}

		// Extract relevant data
		pid := parts[0]
		comm := parts[1]
		bpfErr := parts[3]
		path := parts[len(parts)-1]

		// Ignore erroneous processes // comm == "irqbalance" ||
		// wrk command comes from running the html thread runner

		if bpfErr != "0" || comm == "wrk" {
			continue
		}

		if _, err = strconv.Atoi(pid); err != nil {
			continue  // skip line
		}

		blockList := []string{"/etc/apt", "/var/log", "/var/tmp"} //block /etc/apt folder to prevent repo security conflicts
		allowList := []string{"/etc", "/var", "/root", "/home", "/srv"}

		blocked := owncheckers.StringHasOneOfPrefixes(path, blockList)
		allowed := owncheckers.StringHasOneOfPrefixes(path, allowList)

		if blocked ||
			!allowed ||
			strings.Contains(path, ".git") { continue }

		event := owntypes.OpenSnoopEvent{PID: pid, Command: comm, Path: path}
		key := generateKey(event)

		// Store only unique entries
		events[key] = event
	}

  if err := scanner.Err(); err != nil {
    return nil, errors.New(fmt.Sprintf("Error reading file: ", err))
	}

  return events, nil
}


func makeList(events map[string]owntypes.OpenSnoopEvent, getValue func(owntypes.OpenSnoopEvent) string) owntypes.StringSet {
	visitedElems := make(owntypes.StringSet) //visited pids/commands/paths
  for _, event := range events {
		str := getValue(event)
    if _, exists := visitedElems[str]; !exists {
      visitedElems[str] = true
    }
  }
  return visitedElems
}

func EventPid(event owntypes.OpenSnoopEvent) string {
	return event.PID
}

func EventCommand(event owntypes.OpenSnoopEvent) string {
	return event.Command
}

func EventPath(event owntypes.OpenSnoopEvent) string {
	return event.Path
}

func ListPids(events map[string]owntypes.OpenSnoopEvent) owntypes.StringSet {
	return makeList(events, EventPid)
}

func ListCommands(events map[string]owntypes.OpenSnoopEvent) owntypes.StringSet {
	return makeList(events, EventCommand)
}

func ListPaths(events map[string]owntypes.OpenSnoopEvent) owntypes.StringSet {
	return makeList(events, EventPath)
}

//way to optimize the function below: make filterpathmap straight from file
func MakeServicePathMap(events map[string]owntypes.OpenSnoopEvent, items owntypes.StringSet) (map[string]owntypes.StringSet, map[string]owntypes.StringSet) {
	var (
		cmdStr string
		pathStr string
	)

	allVisitedComms := make(owntypes.StringSet) //visited commands
	serviceVisitedComms := make(owntypes.StringSet)
	resultServiceMap := make(map[string]owntypes.StringSet) //map from bpf command to StringSet of services
	resultPathMap := make(map[string]owntypes.StringSet)

	for _, event := range events {
		cmdStr = event.Command
		pathStr = event.Path

		if !allVisitedComms[cmdStr] {
			resultPathMap[cmdStr] = make(owntypes.StringSet)
			allVisitedComms[cmdStr] = true
		}

		resultPathMap[cmdStr][pathStr] = true

		if foundServices, containsservice := owncheckers.StringSetElemPartOfString(items, pathStr); containsservice {
			//guarantee the stringset exists if not visited already
			if !serviceVisitedComms[cmdStr] {
				resultServiceMap[cmdStr] = make(owntypes.StringSet)
				serviceVisitedComms[cmdStr] = true
			}
			for service, _ := range foundServices {
				resultServiceMap[cmdStr][service] = true
			}
		}
	}
	return resultServiceMap, resultPathMap
}

func MakeFolderMap(filteredEvents map[string]owntypes.StringSet) map[string]*owntypes.Folder {
	var (
		result map[string]*owntypes.Folder = make(map[string]*owntypes.Folder)
	)

	for cmd, paths := range filteredEvents {
    result[cmd] = pathparse.BuildFolderStructure(paths)

		// fmt.Println(fmt.Sprintf("BPF (%s):", cmd))
		// folder := pathparse.BuildFolderStructure(paths)
    // pathparse.PrintFolderStructure(folder, "  ")
  }

	return result
}

//In case we want to filter out some events beforehand (so we don't need to save them as well):
func FilterEvents(events map[string]owntypes.OpenSnoopEvent, include owntypes.StringSet, getValue func(owntypes.OpenSnoopEvent) string) map[string]owntypes.OpenSnoopEvent {
	result := make(map[string]owntypes.OpenSnoopEvent)
	for key, event := range events {
		if include[getValue(event)] {
			result[key] = event
		}
	}
	return result
}

func MakeIncludeFromExclude(allItems owntypes.StringSet, exclude owntypes.StringSet) owntypes.StringSet {
	result := make(owntypes.StringSet)
	for item, _ := range allItems {
		if exclude[item] { continue }
		result[item] = true
	}
	return result
}

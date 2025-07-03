package timetracker

import (
  "fmt"
  "strings"
  "strconv"
	"os"
	"time"
  "project1/tracepick/lib/own/owntypes"
  "project1/tracepick/lib/fileparsers/lineparse"
  "project1/tracepick/lib/helpers/fsHelper"
)

var startTime time.Time
var timetrackFile string


func Initialize(fileName string) {
  fsHelper.RemovePath(fileName) //remove file
  startTime = time.Now() // Start time reference
  timetrackFile = fileName
}

func LogEvent(label string) {
  LogElapsed(startTime, label, timetrackFile)
}

func LogElapsed(start time.Time, label string, filename string) error {
	// Calculate elapsed time in milliseconds
	elapsed := time.Since(start).Milliseconds()

	// Open the file in append mode
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write in the desired format: "label: milliseconds"
	logLine := fmt.Sprintf("%s: %d\n", label, elapsed)
	if _, err := f.WriteString(logLine); err != nil {
		return err
	}
	return nil
}

// ------------------------ PARSE TIME LOG FILE --------------------------------
func parseTimestamps(logFile string) []owntypes.TimingEvent {
  var (
    programName string = "SplitTimeLines"
    events []owntypes.TimingEvent
  )
  lines, err := lineparse.ReadFileLinesSlice(programName, logFile)
  if err != nil {
    fmt.Println(programName, ": no lines read from", logFile)
    return nil
  }

  for _, line := range lines {
    parts := strings.Split(line, ": ")

    if len(parts) > 1 {
      num, err := strconv.Atoi(parts[1])
      if err != nil {
        fmt.Println(programName, ": Conversion error:", err)
        return nil
      }
      event := owntypes.TimingEvent{parts[0], num} //label and number
      events = append(events, event)
    }
  }
  return events
}

// compute total time between START string and STOP string, minus the dialogue times
func computeTotalRunningTime(events []owntypes.TimingEvent, startStr string, stopStr string) int {
  var runtime int = 0

  for _, event := range events {
    switch event.Label {
    case startStr:
      runtime = -event.Time
    case stopStr:
      runtime = runtime + event.Time
    case "BeginDialogue":
      runtime = runtime + event.Time
    case "EndDialogue":
      runtime = runtime - event.Time //net difference between BeginDialogue and EndDialogue is subtracted from runtime
    // default:
    //    runtime = runtime // do nothing
    }
  }
  return runtime
}

func SaveMigrationTimeFromLog(logFile string, startLogStr string, stopLogStr string) {
  owntypes.EvaluationData.GeneralData.MigrationTime =
    computeTotalRunningTime(parseTimestamps(logFile), startLogStr, stopLogStr)
}

func SaveDeployTimeFromLog(strippedItem, logFile, startLogStr, stopLogStr string) {
  data := owntypes.EvaluationData.ArtifactDataMap[strippedItem]
  data.DeployTime = computeTotalRunningTime(parseTimestamps(logFile), startLogStr, stopLogStr)
  owntypes.EvaluationData.ArtifactDataMap[strippedItem] = data
}

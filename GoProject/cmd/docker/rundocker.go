package main

import (
  "fmt"
  "os"
  "bufio"
  "project1/tracepick/lib/helpers/fsHelper"
  "project1/tracepick/lib/helpers/dockerHelpers/runHelper"
  "project1/tracepick/lib/own/owntypes"
  "project1/tracepick/lib/own/owndialogues"
  "project1/tracepick/lib/metrictrackers/timetracker"
  "project1/tracepick/lib/metrictrackers/sizetracker"
  "project1/tracepick/lib/external/yamlio"

)

func main() {
  if len(os.Args) < 2 {
		fmt.Println("usage: rundocker [service name]")
		os.Exit(1)
	}

  var (
    inPath string = "artifacts/output/containers/"
    d1name string = os.Args[1]
    logEventFile string = "metricdata/logfiles/dockerrun.log"
    startLogStr string = "StartDockerRun"
    stopLogStr string = "StopDockerRun"
    err error
  )

  // Retrieve metric state from yaml file
  if owntypes.EvaluationData, err = yamlio.ReadYaml[owntypes.MetricMap](owntypes.MetricPath + owntypes.MetricFileName); err != nil {
    fmt.Println("Error retrieving metrics from artifact build at", owntypes.MetricPath)
  }


  timetracker.Initialize(logEventFile)
  timetracker.LogEvent(startLogStr) //StartDockerRun

  if !owntypes.SkipDialogue { fmt.Println("Specify repository names:") }
  scanner := bufio.NewScanner(os.Stdin) //Going to need user input
  inPath = owndialogues.EditString(scanner, inPath, false)
  d1name = owndialogues.EditString(scanner, d1name, false)


  runHelper.RunImage(inPath, d1name)

  timetracker.LogEvent(stopLogStr) //StopDockerRun
  timetracker.SaveDeployTimeFromLog(d1name, logEventFile, startLogStr, stopLogStr)


  sizetracker.AddedRuntimeSize(d1name)
  sizetracker.TotalFSSize(d1name)

  // Save metric state in yaml file
  fsHelper.CleanPath(owntypes.MetricPath)
  yamlio.WriteYaml[owntypes.MetricMap](owntypes.MetricPath + owntypes.MetricFileName, owntypes.EvaluationData)

  owntypes.PrintEvaluationData(owntypes.EvaluationData)
}

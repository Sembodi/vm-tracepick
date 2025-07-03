package main

import (
  "fmt"
  "os"
  "bufio"
  "project1/tracepick/lib/helpers/fsHelper"
  "project1/tracepick/lib/helpers/dockerHelpers/buildHelper"
  "project1/tracepick/lib/own/owntypes"
  "project1/tracepick/lib/own/owndialogues"
  "project1/tracepick/lib/metrictrackers/timetracker"
  "project1/tracepick/lib/metrictrackers/sizetracker"
  "project1/tracepick/lib/external/yamlio"

)

func main() {
  if len(os.Args) < 2 {
		fmt.Println("usage: builddocker [service name]")
		os.Exit(1)
	}
  var (
    inPath string = "artifacts/output/containers/"
    d1name string = os.Args[1]
    logEventFile string = "metricdata/logfiles/dockerbuild.log"
    startLogStr string = "StartDockerBuild"
    stopLogStr string = "StopDockerBuild"
    err error
  )

  // Retrieve metric state from yaml file
  if owntypes.EvaluationData, err = yamlio.ReadYaml[owntypes.MetricMap](owntypes.MetricPath + owntypes.MetricFileName); err != nil {
    fmt.Println("Error retrieving metrics from artifact build at", owntypes.MetricPath)
  }


  timetracker.Initialize(logEventFile)
  timetracker.LogEvent(startLogStr) //StartDockerBuild

  if !owntypes.SkipDialogue { fmt.Println("Specify repository names:") }
  scanner := bufio.NewScanner(os.Stdin) //Going to need user input
  inPath = owndialogues.EditString(scanner, inPath, false)
  d1name = owndialogues.EditString(scanner, d1name, false)

  if owntypes.CachedBuild {
    buildHelper.CachedBuild(inPath, d1name)
  } else {
    buildHelper.NonCachedBuild(inPath, d1name)
  }
  timetracker.LogEvent(stopLogStr) //StopDockerBuild
  timetracker.SaveDeployTimeFromLog(d1name, logEventFile, startLogStr, stopLogStr)


  sizetracker.ImageSize(d1name)



  // Save metric state in yaml file
  fsHelper.CleanPath(owntypes.MetricPath)
  yamlio.WriteYaml[owntypes.MetricMap](owntypes.MetricPath + owntypes.MetricFileName, owntypes.EvaluationData)

  owntypes.PrintEvaluationData(owntypes.EvaluationData)

}

package main

import (
  "fmt"
  "project1/tracepick/lib/own/owntypes"
  "project1/tracepick/lib/own/ownformatters"
  "project1/tracepick/lib/helpers/fsHelper"
  "project1/tracepick/lib/helpers/artifactHelper"
  "project1/tracepick/lib/metrictrackers/timetracker"
  "project1/tracepick/lib/metrictrackers/sizetracker"
  "project1/tracepick/lib/external/yamlio"
)

//TODO: make sure that if an error regarding /usr/share is given, make sure another package is installed on the container.
// For example, "adminer" is missing in this example, so it needs to be installed when this error is given.
// Command to find the proper package: find /usr/share/adminer -type f | head -n 1 | xargs dpkg -S
// Then, make sure install using RUN apt install adminer is added in the dockerfile in keepfiles/artifacts/final/LAMP_WP
// CHECK IF sources.list is included in copied files

func printTracePick() {
  fmt.Println(`
    _______                   _____
   |__   __|                 |  __ \*      _
      | |_ __ ____  ___ ___  | |_/ /_  ___| | __
      | | '__/ _  |/ __/ _ \ |  __/| |/ __| |/ /
      | | | | (_| | (__  __/ | |   | | (__|   <
      |_|_|  \__,_|\___\___| |_|   |_|\___|_|\_\
    `)
}


func generateArtifacts() {
    startLogStr := "StartArtifactGen"
    stopLogStr := "StopArtifactGen"
    logDir := "metricdata/logfiles/"
    logEventFile := logDir + "artifactgen.log"
    fsHelper.CleanPath(logDir)

    timetracker.Initialize(logEventFile)
    timetracker.LogEvent(startLogStr) //StartArtifactGen

    printTracePick()
    owntypes.UseSudo = artifactHelper.PrepareSudo()
    fsHelper.RemovePath(owntypes.MetricPath + owntypes.MetricFileName)

    var outputPath string = "artifacts/output/"

    // only clean containers path if no extra tracing is done
    if !owntypes.DoTracing {
      outputPath = outputPath + "containers/"
    }

    // specify VM name: var vmName string = os.Args[2]

    // CleanPath swipes the path's folder clean
    // by removing all of its components and recreating the folder specified
    fsHelper.CleanPath(outputPath)

    // Create the right dockerfile strings and write them to their dockerfile
    services, fromStr, runMap, workdirStr, copyMap, exposeMap, dockerCmdMap := artifactHelper.GetDockerStrings()
    // Write the dockerfile strings to one file located at containerPath

    if owntypes.DoTracing {
      outputPath = outputPath + "containers/"
    }

    artifactHelper.WriteDockerFile(services, outputPath, fromStr, runMap, workdirStr, copyMap, exposeMap, dockerCmdMap)
    timetracker.LogEvent(stopLogStr) //StopArtifactGen



    sizetracker.VMSize()
    timetracker.SaveMigrationTimeFromLog(logEventFile, startLogStr, stopLogStr)

    shortServices := ownformatters.StripServiceSet(services)
    sizetracker.TotalRuntimeSize(shortServices)

    owntypes.PrintEvaluationData(owntypes.EvaluationData)

    // Save metric state in yaml file
    fsHelper.CleanPath(owntypes.MetricPath)
    yamlio.WriteYaml[owntypes.MetricMap](owntypes.MetricPath + owntypes.MetricFileName, owntypes.EvaluationData)

}



//TODO:

// - make sure it can run in batches (based on VM names on different ports, save artifacts in same folder, but move them to keepfiles under VM name)
// - make sure that copying/replacing happens as wished for (all folders replaced, unless it's a protected folder (hence partial profile))

// - calculate migrationTime & save in GLOBAL datastructure
// - clean cache & save deployment time (with/without caching) at docker build
//   calculate deployTime by adding buildTime and runTime
// - compute vm size by running remote command 'du -hs /'
// - compute exported runtime size using 'du' for rootfs
// - incorporate runtime size calculation in docker file (can be done manually and exported to file at root directory)
// - incorporate container filesystem size calculation in dockerfile (and save to root dir file)
// - write down image size (using docker images | head -n1 | awk etc etc...)

// - add the logic that only if two files in a folder are used, the folder is taken (instead of just the largest non-protected folder)
// - make sure if file not found upon run, update image build

func main() {
  generateArtifacts()
}

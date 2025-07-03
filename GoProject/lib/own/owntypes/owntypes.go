package owntypes

import (
  "fmt"
  // "project1/tracepick/lib/external/yamlio"
)

//set of strings so no duplicate elements
type StringSet = map[string]bool

type ImageList = map[string]bool

//alias no output is empty string
const NoOutputFile string = ""

//for each outputfile string, show command string
type CommandMap map[string]string

type PortMap map[string]StringSet
type DockerMap map[string]string
type FormatMap map[string]string

//pathparse:
type Folder struct {
    Name string
    Folders map[string]*Folder
}

//bpfparse:
type OpenSnoopEvent struct {
	PID  string
	Command string
	Path string
}

type TimingEvent struct {
  Label string
  Time int
}

type GeneralData struct {
  MigrationTime     int `yaml:"MigrationTime"`
  VMSize            int `yaml:"VMSize"`
}

type ArtifactData struct {
    DeployTime        int `yaml:"deployTime"`
    TotalRuntimeSize  int `yaml:"totalRuntimeSize"`
    AddedRuntimeSize  int `yaml:"addedRuntimeSize"`
    TotalFSSize       int `yaml:"totalFSSize"`
    ImageSize         int `yaml:"imageSize"`
}

// Map service names to their metricdata:
type MetricMap struct {
  GeneralData     GeneralData             `yaml:"GeneralData"`
  ArtifactDataMap map[string]ArtifactData `yaml:"ArtifactDataMap"`
}

func PrintStringSet(ss StringSet) {
  fmt.Print("StringSet:[")
  for item, _ := range ss {
    fmt.Print(" ", item, " ")
  }
  fmt.Print("]\n")
}

func PrintStringStringSetMap(itemMap map[string]StringSet) {
  fmt.Println("String-StringSet Map: ")
  for s, ss := range itemMap {
    fmt.Println(s, ":")
    PrintStringSet(ss)
  }
}

// global variables

var UseSudo bool = true

var HasNoPackage StringSet = make(StringSet) //list of services with no matching packages

var DoTracing bool = true

var CachedBuild bool = false

var SkipDialogue bool = true //for batch runs

var EvaluationData MetricMap = MetricMap{
  ArtifactDataMap: make(map[string]ArtifactData),
}

var MetricPath string = "metricdata/metricstate/"
var MetricFileName string = "results.yaml"

func PrintEvaluationData(results MetricMap) {
  fmt.Println("GENERAL DATA:\n")
  fmt.Println("MigrationTime:", results.GeneralData.MigrationTime)
  fmt.Println("VMSize:", results.GeneralData.VMSize, "\n\n")

  fmt.Println("ARTIFACT DATA:\n")

  for service, data := range results.ArtifactDataMap {
    fmt.Println("Service:", service)
    fmt.Println("DeployTime:", data.DeployTime)
    fmt.Println("TotalRuntimeSize:", data.TotalRuntimeSize)
    fmt.Println("AddedRuntimeSize:", data.AddedRuntimeSize)
    fmt.Println("TotalFSSize:", data.TotalFSSize)
    fmt.Println("ImageSize:", data.ImageSize, "\n")
  }
}

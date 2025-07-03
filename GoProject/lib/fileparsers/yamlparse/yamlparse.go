package yamlparse

import (
  "fmt"
  "os"
  "io/ioutil"
  "log"
	"encoding/csv"
	"path/filepath"
  "project1/tracepick/lib/own/owntypes"
  "project1/tracepick/lib/helpers/fsHelper"
  "project1/tracepick/lib/external/yamlio"
)

func MakeCSVFromDir(dir string) {
  outDir := "keepfiles/csv/"
  outFilePath := outDir + "output.csv"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalf("Failed to read directory: %v", err)
	}

  fsHelper.CleanPath(outDir)
	outFile, err := os.Create(outFilePath)
	if err != nil {
		log.Fatalf("Failed to create CSV file: %v", err)
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	// Write CSV header
	writer.Write([]string{
		"NAME", "MIGRATIONTIME", "VMSIZE", "DEPLOYTIME", "RUNTIMESIZE", "ADDRUNTIMESIZE", "FSSIZE", "IMAGESIZE",
	})

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".yaml" || filepath.Ext(file.Name()) == ".yml" {
			fullPath := filepath.Join(dir, file.Name())

      parsed, err := yamlio.ReadYaml[owntypes.MetricMap](fullPath)

			if err != nil {
				log.Printf("Failed to parse YAML in file %s: %v", file.Name(), err)
				continue
			}

			for name, artifact := range parsed.ArtifactDataMap {
				record := []string{
					name,
					fmt.Sprint(parsed.GeneralData.MigrationTime),
					fmt.Sprint(parsed.GeneralData.VMSize),
					fmt.Sprint(artifact.DeployTime),
					fmt.Sprint(artifact.TotalRuntimeSize),
					fmt.Sprint(artifact.AddedRuntimeSize),
					fmt.Sprint(artifact.TotalFSSize),
					fmt.Sprint(artifact.ImageSize),
				}
				writer.Write(record)
			}
		}
	}

	fmt.Println("CSV generation complete.")
}

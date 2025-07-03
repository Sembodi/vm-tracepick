package fsHelper

import (
  "os"
  "os/exec"
  "io/ioutil"
  "fmt"
  "project1/tracepick/lib/own/owntypes"
)

func RemovePath(path string) {
  err := os.RemoveAll(path)
  if err != nil {
    fmt.Printf("Error removing directory %s: \n%s", path, err)
  }
  // else {
  //   fmt.Printf("Directory %s and all contents removed successfully.", path)
  // }
}

func CreatePath(path string) {
  err := os.MkdirAll(path, 0755)
  if err != nil {
    fmt.Printf("Error creating directory %s: \n%s", path, err)
  }
  // else {
  //   fmt.Printf("Directory %s and all contents created successfully.", path)
  // }
}

func CleanPath(path string) {
  RemovePath(path)
  CreatePath(path)
}


func GetFiles(programName string, path string) []os.FileInfo {
  var (
    files []os.FileInfo
    err error
  )

  if files, err = ioutil.ReadDir(path); err != nil {
    fmt.Printf("%s: Error reading directory: ", programName, err)
    os.Exit(2)
  }

  return files
}

func GetFileNames(files []os.FileInfo) owntypes.StringSet {
  var items owntypes.StringSet = make(owntypes.StringSet)

  for _, file := range files {
    if file.IsDir() {
      continue
    }
    items[file.Name()] = true
  }
  return items
}

func CopyRootFS(inDir string, outDir string) {
  fmt.Println("Copying folder", inDir, "to", outDir, "...")
  cmd := exec.Command("sudo", "/bin/cp", "-a", inDir, outDir)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error copying folder:", err)
		return
	} else {
		fmt.Println("Folder copied successfully.")
	}
}

func CreateExecFile(outFileName string, outStr string) {
  // fmt.Println("Creating file at:", outFileName, ":")

  file, err := os.OpenFile(outFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
  if err != nil {
    fmt.Println("entrypoint.sh creation not successful")
    os.Exit(1)
  }
  defer file.Close() // Make sure to close it when done
  file.WriteString(outStr)
}

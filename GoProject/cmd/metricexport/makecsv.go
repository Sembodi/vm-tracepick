package main

import (
  "project1/tracepick/lib/fileparsers/yamlparse"
)

func main() {
  var (
    dir string = "keepfiles/final_static/metrics/TracePick/"
  )

  yamlparse.MakeCSVFromDir(dir)
}

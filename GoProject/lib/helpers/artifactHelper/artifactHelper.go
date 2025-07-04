package artifactHelper

import (
  "fmt"
  "regexp"
  "strings"
  "bufio"
  "os"
  "os/exec"
  "project1/tracepick/lib/own/owntypes"
  "project1/tracepick/lib/own/ownformatters"
  "project1/tracepick/lib/own/owncheckers"
  "project1/tracepick/lib/own/owndialogues"
  "project1/tracepick/lib/fileparsers/lineparse"
  "project1/tracepick/lib/fileparsers/pathparse"
  "project1/tracepick/lib/helpers/fsHelper"
  "project1/tracepick/lib/helpers/vmHelpers/etcvarHelper"
  "project1/tracepick/lib/helpers/vmHelpers/osHelper"
  "project1/tracepick/lib/helpers/vmHelpers/portHelper"
  "project1/tracepick/lib/helpers/vmHelpers/packageHelper"
  "project1/tracepick/lib/helpers/vmHelpers/serviceHelper"
  "project1/tracepick/lib/helpers/vmHelpers/traceHelper"
  "project1/tracepick/lib/helpers/dockerHelpers/imagesearchHelper"
  "project1/tracepick/lib/collectors/bpfcollector"
  "project1/tracepick/lib/collectors/etccollector"
  "project1/tracepick/lib/runners/localrun"
  "project1/tracepick/lib/runners/remoterun"
  "project1/tracepick/lib/external/yamlio"
  "project1/tracepick/lib/metrictrackers/timetracker"

  // "github.com/google/shlex" //external

)


// Set whether prefix 'sudo' is needed and optionally remove connect address from known hosts to prevent conflicts (not recommended for prod)
func PrepareSudo() bool {
  var (
    pmap remoterun.ConnectMap
    err error
  )
  if pmap, err = yamlio.ReadYaml[remoterun.ConnectMap](os.Args[1]); err != nil {
		fmt.Println(err)
		return false
	}

  for _, form := range pmap { //only goes through first element of pmap (connect map)
    connectParts := strings.Split(form.Connection, "@")
    user := connectParts[0]
    // nodeParts := strings.Split(connectParts[1], ":")
    // node := nodeParts[len(nodeParts) - 1]
    // RemoveFromKnownHosts(node)

    if user == "root" { return false }
    return true
  }
  return true
}

func RemoveFromKnownHosts(pattern string) {
  programName := "RemoveFromKnownHosts"
  cmd := exec.Command("bash", "bashscripts/rmfromknownhosts.sh", pattern)
  if err := cmd.Run(); err != nil {
    fmt.Println(programName, fmt.Sprintf(": Error removing lines with '%s' from known hosts:", pattern), err)
  }
}


// Returns baseimageName and Dockerfile line,
// same OS for each service, so no map for this one:
func DockerFROM() (string, string) {
  var (
    programName string = "DockerFROM"

    baseimageName string
    tagStr string

    err error
  )

  osHelper.GetVMOS()
  inputOSNameStr := osHelper.ParseVMOS("os") //set of 1 element
  inputVersionStr := osHelper.ParseVMOS("version") //set of 1 element

  //find and choose base image NAME
  imagesearchHelper.ListDockerImages(ownformatters.FromStringToStringSet(inputOSNameStr))
  images := imagesearchHelper.ChooseDockerImages()


  baseimageName = images[inputOSNameStr]

  //find appropriate base image TAG
  if tagStr, err = imagesearchHelper.FindClosestVersion(baseimageName, inputVersionStr); err != nil {
    fmt.Println("ERROR ", programName, ": ", err)
    return "err", "err"
  }

  baseimage := baseimageName + ":" + tagStr

  return inputOSNameStr, "FROM " + baseimage
}

// In here, run app install items (nginx, python3)
// make sure that items are full (no php instead of php8.1-fpm)
func dockerRUN(items owntypes.StringSet, osName string) owntypes.CommandMap {
  var (
    programName string = "dockerRUN"

    searchCmds owntypes.CommandMap = make(owntypes.CommandMap) //for remoterun make command map that has key:php8.1-fpm value: <command with php8.1>
    searchOutDir string = "outputs/installs/"

    sb strings.Builder
    resultMap owntypes.CommandMap = make(owntypes.CommandMap)
  )

  runStr := "RUN "
  updateCmd, _ := osHelper.GetOSCommand(osName, osHelper.ReturnUpdateString)
  installCmd, _ := osHelper.GetOSCommand(osName, osHelper.ReturnInstallString)
  searchCmd, _ := osHelper.GetOSCommand(osName, osHelper.ReturnSearchPackageString)
  repoCmd, _ := osHelper.GetOSCommand(osName, osHelper.ReturnPackageRepoString)
  cleanCmd, _ := osHelper.GetOSCommand(osName, osHelper.ReturnCleanString)
  commentStr := "# "

  smallInstallCmd, _ := osHelper.GetOSCommand(osName, osHelper.ReturnSmallInstallString)

  // installAddApt, _ := osHelper.GetOSCommand(osName, osHelper.ReturnInstallAddAptString)
  // purgeAddApt, _ := osHelper.GetOSCommand(osName, osHelper.ReturnPurgeAddAptString)

  installRsync := fmt.Sprintf(smallInstallCmd, "rsync")

  //todo: in forloop below, also build the package: repo map
  for item, _ := range items {
     //todo error handling

    trimmedItem := packageHelper.TrimPackageName(item)
    fullSearchCmd := fmt.Sprintf(searchCmd, trimmedItem)

    searchCmds[item] = fullSearchCmd
  }

  fsHelper.CleanPath(searchOutDir)
  remoterun.RemoteRun(programName, searchCmds, searchOutDir)

  for item, _ := range items {
    sb.WriteString(runStr + updateCmd + "\n")
    // sb.WriteString(runStr + installAddApt + "\n")

    allPackages, err := lineparse.ReadFileLines(programName, searchOutDir + item)
    if err != nil {
      fmt.Println(programName, ": Reading file lines for item", item, "unsuccessful. Continuing for other items...")
      continue
    }
    if len(allPackages) == 0 {
      owntypes.HasNoPackage[item] = true // if this is true, then include all files related to the service name/user
    }

    selectedPackages := packageHelper.SelectPackages(allPackages, item)  // selects only those to cover all dependencies

    repoCmds := ownformatters.FormatCommands(repoCmd, allPackages)

    cmdMap := remoterun.GetRemoteRunMap(programName, repoCmds) //map of key: pkg and value:StringSet of repos

    for pkg, repoSet := range cmdMap {

      //if pkg is in selectedPackages, preStr="", otherwise preStr="# "
      preStr := commentStr
      if selectedPackages[pkg] { preStr = "" }

      repoSum := ownformatters.FromStringSetToSpacedStrings(repoSet)
      fullInstallCmd := fmt.Sprintf(installCmd, pkg, repoSum)

      sb.WriteString(preStr + runStr + fullInstallCmd + "\n")
    }

    //sb.WriteString(runStr + purgeAddApt + "\n")

    sb.WriteString(runStr + installRsync + "\n")
    sb.WriteString(runStr + cleanCmd)

    resultMap[item] = sb.String()
    sb.Reset()
  }

  //idea: RUN apt update once, make many lines for every repo hit

  //save in outputs: map of repo => repo link so it can be referred to if docker build fails
  //idea: RUN apt install [alle gevonden package names] || add-apt-repository repoNamesMap[iets] && install -y [alle gevonden package names]
  return resultMap
}

func dockerEXPOSE(items owntypes.StringSet) owntypes.CommandMap {
  var (
    formatCmd string = "EXPOSE %s"
    serviceSet owntypes.CommandMap = make(owntypes.CommandMap) //set of EXPOSE commands for ONE service
    resultSet owntypes.CommandMap = make(owntypes.CommandMap) //map from service key to string value of 1 or more EXPOSE cmds
  )

  portHelper.RetrievePortNumbers(items)
  portNumbers := portHelper.ParsePortMap(items)

  for name, portSet := range portNumbers {
    serviceSet = make(owntypes.CommandMap)
    for port, _ := range portSet {
      serviceSet[port] = fmt.Sprintf(formatCmd, port) //+ " # " + name
    }
    resultStr := ownformatters.FromCommandMapToString(serviceSet)
    resultSet[name] = resultStr
  }
  return resultSet
}

func dockerCOPY(items owntypes.StringSet) owntypes.CommandMap {
  var (
    programName string = "dockerCOPY"
    pullPath string = "artifacts/output/filesystem/"

    outDir string = "artifacts/output/containers/"
    formatRuntimeFolder string = "myrootfs-%s/"
    helperSource string = "defaultdata/helperFiles"
    formatCmd string = "tar -xpzf " + pullPath + "%s -C " + outDir

    bpfstartFormatFileName string = "bpfstartfiles_%s.tar.gz"
    bpfFormatFileName string = "bpffiles_%s.tar.gz"
    etcFormatFileName string = "etcfiles_%s.tar.gz"

    usedBpfCmdsServices map[string]owntypes.StringSet = make(map[string]owntypes.StringSet)
    result owntypes.CommandMap = make(owntypes.CommandMap)
  )

  if owntypes.DoTracing {
    traceHelper.TraceServiceRestart(items) // once run, for next run not necessary
    // traceHelper.TraceWrk()
  }

  items = ownformatters.StripServiceSet(items)

  // // EXTRACT general BPF FILES from wrk
  // bpfCmds, bpfneededPaths, bpfServiceMap := bpfcollector.BpfMinimalPathsServiceMap(items) //TODO if item name is found in the path, also include the service in the command, but do this when making the map already (maybe add a field that can be either nil OR have a list/stringset of services inside)

  if owntypes.DoTracing {
    etcvarHelper.RetrieveConfigurationFiles(items)

    fsHelper.CleanPath(pullPath)
    fmt.Println("Need files:")

    //EXTRACT SERVICE START BPF FILES
    bpfstartneededPaths := bpfcollector.BpfServiceStartMinPaths(items)

    for item, _ := range items {
      pathparse.PrintMinimalFolderList(fmt.Sprintf("From start BPF (service %s)", item), bpfstartneededPaths[item])
    }


    fmt.Println("Pulling bpfstartfiles...")
    for item, _ := range items {
      fileName := fmt.Sprintf(bpfstartFormatFileName, item)
      fmt.Println(fileName,"...")
      remoterun.RemotePull(bpfstartneededPaths[item], fileName)
    }

    for bpfCmd, _ := range bpfCmds {
      pathparse.PrintMinimalFolderList(fmt.Sprintf("From BPF (command %s)", bpfCmd), bpfneededPaths[bpfCmd])
    }

    // fmt.Println("Pulling bpffiles...")
    // for cmd, _ := range bpfCmds {
    //   if services, containsservice := owncheckers.StringPartOfElementOrViceVersa(cmd, items); containsservice {
    //     usedBpfCmdsServices[cmd] = ownformatters.MergeStringSets(services, bpfServiceMap[cmd]) //TODO if any item name is found in the path, also include the service in the command
    //     fmt.Println("FROM BPF USED: ", cmd, "FOR SERVICES: ")
    //     owntypes.PrintStringSet(usedBpfCmdsServices[cmd])
    //     fileName := fmt.Sprintf(bpfFormatFileName, cmd)
    //     remoterun.RemotePull(bpfneededPaths[cmd], fileName)
    //   }
    // }

    // EXTRACT usr/etc/var files from find-command
    etcneededPaths := etccollector.EtcvarMinimalPaths()
    for file, _ := range items {
      pathparse.PrintMinimalFolderList(fmt.Sprintf("From etc (command %s)", file), etcneededPaths[file])
    }

    fmt.Println("Pulling etcvarfiles...")
    for item, _ := range items {
      fileName := fmt.Sprintf(etcFormatFileName, item)
      remoterun.RemotePull(etcneededPaths[item], fileName)
    }
  }

  cmdMap := ownformatters.MapExecFromFiles(programName, pullPath, formatCmd)
  serviceBpfCmds := ownformatters.ReverseStringSetMap(usedBpfCmdsServices)

  for item, _ := range items {
    runtimeFolder := fmt.Sprintf(formatRuntimeFolder, item)
    subDir := item + "/" + runtimeFolder

    fsHelper.CreatePath(outDir + subDir)

    for _, formatFileName := range []string{bpfstartFormatFileName, etcFormatFileName} {
      fileName := fmt.Sprintf(formatFileName, item)
      newCmd := cmdMap[fileName] + subDir
      newCmdMap := ownformatters.FromStringToCommandMap(newCmd, owntypes.NoOutputFile)
      localrun.LocalRun(programName, newCmdMap, owntypes.NoOutputFile, true)
    }

    //traverse reversed bpfcmd-servicemap:
    for bpfCmd, _ := range serviceBpfCmds[item] {
      bpfFileName := fmt.Sprintf(bpfFormatFileName, bpfCmd)

      newCmd := cmdMap[bpfFileName] + subDir
      newCmdMap := ownformatters.FromStringToCommandMap(newCmd, owntypes.NoOutputFile)
      localrun.LocalRun(programName, newCmdMap, owntypes.NoOutputFile, true)
    }

    fsHelper.CopyRootFS(helperSource, outDir + item)

    result[item] = strings.Join([]string{"COPY " + runtimeFolder + " /bbackup/",
                                         `RUN rsync -rcn --out-format="%l %n" /bbackup/ / | awk '{sum += $1} END {print sum}' > /added-runtime.log`,
                                         `COPY helperFiles/ /helperFiles/`,
                                         `RUN chmod +x /helperFiles/copyAlgorithm.sh`,
                                         `RUN /helperFiles/copyAlgorithm.sh '/bbackup' '' '/helperFiles/protected_paths'`},
                                         // `RUN rm -r /bbackup/`,
                                         // `RUN rm -r /helperFiles/`},
                                "\n")
  }

  // owntypes.PrintStringSet(bpfCmds)
  // owntypes.PrintStringStringSetMap(usedBpfCmdsServices)

  return result
}

func dockerWORKDIR(items owntypes.StringSet) owntypes.CommandMap {
  var (
    programName string = "dockerWORKDIR"
    formatDirCmd string = `systemctl show %s --property=WorkingDirectory | sed '1s/^WorkingDirectory=//'`
    outDir string = "outputs/workdir/"
    formatDockerCmd string = "WORKDIR %s"

    result owntypes.CommandMap = make(owntypes.CommandMap)
  )

  fsHelper.CleanPath(outDir)
  dirCmdMap := ownformatters.FormatCommands(formatDirCmd, items)

  remoterun.RemoteRun(programName, dirCmdMap, outDir)

  workdirs := lineparse.LineSlicesFromFiles(programName, outDir)

  for item, lines := range workdirs {
    var path string = "/"
    if len(lines) > 0 { path = lines[0] }
    result[item] = fmt.Sprintf(formatDockerCmd, path)
  }

  return result
}

func dockerCMD(items owntypes.StringSet) (owntypes.CommandMap) { // , owntypes.CommandMap
  var (
    programName string = "dockerCMD"
    outDir string = "outputs/execstart/"
    formatCmd string = `bash -c 'awk "/^ExecStart=/, !/\\\\$/" /*/systemd/system/%s.service' | sed "1s/^ExecStart=//"` // "grep -A '^ExecStart=' /lib/systemd/system/%s.service | cut -d= -f2- | head -n1" //`systemctl show %[1]s --property=ExecStart | grep -o 'argv\[\]=[^;]*' | sed 's~argv\[\]=~~g'`
    startCmd string


    err error

    resultdockerCMD owntypes.CommandMap = make(owntypes.CommandMap)
    // resultStartCmd owntypes.CommandMap = make(owntypes.CommandMap)
  )

  fsHelper.CleanPath(outDir)

  cmdMap := ownformatters.FormatCommands(formatCmd, items)

  remoterun.RemoteRun(programName, cmdMap, outDir)

  serviceHelper.ShowUsersFromServices(items)
  userMap := serviceHelper.ParseUsersFromServices()

  serviceHelper.ShowDaemonsFromServices(items)
  daemonMap := serviceHelper.ParseDaemonsFromServices()


  // Define the patterns to look for (either "daemon" or "background")
	daemonPattern := regexp.MustCompile(`(?i)(daemon|background)`)

  scanner := bufio.NewScanner(os.Stdin) //Going to need user input

  for item, _ := range items {
    if startCmd, err = lineparse.ReadContent(programName, outDir + item); err != nil {
      fmt.Println("ERROR: ", programName, "Cannot find startcmd for item ", item)
    }

    // Check if the startCmd contains "daemon" or "background"
  	if daemonPattern.MatchString(startCmd) {
  		// Extract the part of the command that contains "daemon" or "background"
  		matches := daemonPattern.FindAllString(startCmd, -1)
  		if len(matches) > 0 {
  			// Show the user the full option with "daemon" or "background"
  			fmt.Printf("Detected background-related option(s): %v\n", matches) //always able to edit, but extra notification in case it's found (INFO)
      }
      // } else {
        // 	fmt.Println("No 'daemon' or 'background' options found in startCmd for " + item + ".")
        // }
      fmt.Println("PLEASE CHECK if the startup command runs on the FOREGROUND")
      fmt.Println("--nodaemon (nginx), -DFOREGROUND (apache2), etc...")
      fmt.Printf("Full %s StartCommand: ", item)

			// Ask if the user wants to modify it
			startCmd = owndialogues.EditString(scanner, startCmd, false)
  	}

    userStr := "root" //use root user as backup
    daemonStr := "" //use service name as backup

    addUserRunVarLib := ""

    if len(daemonMap[item]) > 0 {
      daemonStr = daemonMap[item][0]

      renameDaemonCmd := fmt.Sprintf(ownformatters.MakeSudo(`find /usr/sbin /usr/bin -type l -exec ls -l {} + | grep -o '[a-z]* \-> %s$' | awk '{print $1}'`), daemonStr)
      daemonPath := "outputs/daemons/"
      remoterun.SingleRemoteRun(programName, renameDaemonCmd, daemonPath, item)

      daemonLines, _ := lineparse.ReadFileLinesSlice(programName, daemonPath + item)

      if daemonLen := len(daemonLines); daemonLen > 0 {
        daemonStr = daemonLines[daemonLen - 1]
      }
    }

    if len(userMap[item]) > 0 {
      userStr = userMap[item][0]

      if daemonStr == "" {
        if userStr == "root" {
          daemonStr = item
        } else {
          daemonStr = userStr
        }
      }

      addUserRunVarLib = strings.Join([]string{fmt.Sprintf("RUN groupadd -r %[1]s && useradd -r -g %[1]s %[1]s || echo 'group and user already exist'", userStr),
        fmt.Sprintf("RUN mkdir -p /var/lib/%[1]s && chown -R %[1]s:%[1]s /var/lib/%[1]s", userStr),
        fmt.Sprintf("RUN mkdir -p /var/log/%[1]s && chown -R %[1]s:%[1]s /var/log/%[1]s", userStr),
        fmt.Sprintf("RUN mkdir -p /run/%s", daemonStr),
        // fmt.Sprintf("RUN mkdir -p /etc/%s", userStr),
        fmt.Sprintf("RUN chown %[1]s:%[1]s /run/%[2]s", userStr, daemonStr),
        // fmt.Sprintf("RUN chown %[1]s:%[1]s /etc/%[1]s", userStr),
        fmt.Sprintf("USER %s", userStr)},
        "\n")
    }

    entrypointPath := fmt.Sprintf("artifacts/output/containers/%s/", ownformatters.StripServiceString(item))
    entrypointStr := strings.Join([]string{ "#!/bin/bash",
                                            "set -e",
                                            fmt.Sprintf("exec %s", startCmd) },
                                  "\n")

    fsHelper.CreateExecFile(entrypointPath + "entrypoint.sh", entrypointStr)

    // TODO: add ExecStartPre commands in the entrypoint
  	// Build the CMD line
  	resultdockerCMD[item] = strings.Join([]string{addUserRunVarLib,
                                                  `COPY entrypoint.sh /entrypoint.sh`,
                                                  `CMD ["/entrypoint.sh"]`},  // strings.Join(quotedArgs, ", "))},
                                         "\n")
    // resultStartCmd[item] = startCmd
  }

  //  TODO write CMD in entrypoint file, copy to container location and simply call COPY entrypoint.sh /entrypoint.sh and CMD entrypoint.sh

  return resultdockerCMD // , resultStartCmd
}
// ---------------------------- COMBINE FUNCTIONS ------------------------------
func GetDockerStrings() (owntypes.StringSet, string, owntypes.CommandMap, owntypes.CommandMap, owntypes.CommandMap, owntypes.CommandMap, owntypes.CommandMap) {
  timetracker.LogEvent("Start FROM")
  osName, fromStr := DockerFROM()


  timetracker.LogEvent("Start RUN")
  fmt.Println("Retrieving services...")
  serviceHelper.RetrieveUpServices(osName)

  // List of services running on VM
  services := serviceHelper.GetAddedServices(osName)
  // owntypes.PrintStringSet(addServices)
  runMap := dockerRUN(services, osName)
  timetracker.LogEvent("Start WORKDIR")
  workdirMap := dockerWORKDIR(services)
  timetracker.LogEvent("Start COPY")
  copyMap := dockerCOPY(services)
  timetracker.LogEvent("Start EXPOSE")
  exposeMap := dockerEXPOSE(services)
  timetracker.LogEvent("Start CMD")
  dockerCMDMap := dockerCMD(services)

  return services, fromStr, runMap, workdirMap, copyMap, exposeMap, dockerCMDMap
}

func WriteDockerFile(items owntypes.StringSet, containerPath string,
                    fromStr string, runMap owntypes.CommandMap,
                    workdirMap owntypes.CommandMap, copyMap owntypes.CommandMap,
                    exposeMap owntypes.CommandMap, dockerCMDMap owntypes.CommandMap) {
  for service, _ := range items {
    serviceFolder := containerPath + ownformatters.StripServiceString(service) + "/"
    // Open or create the Dockerfile
    fsHelper.CreatePath(serviceFolder)
    file, err := os.Create(serviceFolder + "Dockerfile")
    if err != nil {
    	fmt.Println("Dockerfile creation not successful")
      os.Exit(1)
    }

    defer file.Close() // Make sure to close it when done

    file.WriteString(fromStr + "\n")
    file.WriteString(runMap[service] + "\n")
    file.WriteString(copyMap[ownformatters.StripServiceString(service)] + "\n")
    file.WriteString(exposeMap[service] + "\n")
    file.WriteString(workdirMap[service] + "\n")
    file.WriteString(dockerCMDMap[service] + "\n")
    fmt.Println("Dockerfile created successfully.")
  }
  fmt.Println("Artifacts generated successfully at artifacts/output")
  fmt.Println("Next Steps (EXAMPLE):")
  fmt.Println("1. Build the Docker image with: docker build -t custom-nginx ./output")
  fmt.Println("2. Run with: docker run -d -p 8080:80 custom-nginx")
}

// ---------------------------- PRINT/SHOW FUNCTIONS ---------------------------
func PrintDockerMaps(fromStr string, runMap owntypes.CommandMap,
                    workdirStr string, copyMap owntypes.CommandMap,
                    exposeMap owntypes.CommandMap, cmdMap owntypes.CommandMap) {
  fmt.Printf("\ndockerFROM:\n%s\n", fromStr)
  fmt.Printf("\ndockerRUN:\n%s\n", runMap)
  fmt.Printf("\ndockerWORKDIR:\n%s\n", workdirStr)
  fmt.Printf("\ndockerCOPY:\n%s\n", copyMap)
  fmt.Printf("\ndockerEXPOSE:\n%s\n", exposeMap)
  fmt.Printf("\ndockerCMD:\n%s\n", cmdMap)
}

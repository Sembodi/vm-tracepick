package remoterun

import (
  "fmt"
  // "io" //enable if we want to use MultiWriter
  "os"
	"os/exec"
	"strings"
  "project1/tracepick/lib/own/owntypes"
  "project1/tracepick/lib/own/ownformatters"
  "project1/tracepick/lib/runners/localrun"
  "project1/tracepick/lib/external/scmd"
  "project1/tracepick/lib/external/yamlio"
)


type ConnectForm struct {
	Connection string `yaml:"ssh_connection"`
	Identity   string `yaml:"ssh_identity"`
}

type ConnectMap map[string]ConnectForm


func saveCmdOutput(node string, form ConnectForm, cmdStr string, fileName string) error {
  var (
    cmd    *exec.Cmd

    programName string = fmt.Sprintf("{saveCmdOutput to node %s}", node)
  )

  cmd = scmd.RemoteCommand(form.Connection, form.Identity, cmdStr)

  return localrun.RunCmd(cmd, programName, fileName)
}

func returnRemoteExecError(node string, cmdStr string, err error) error {
  parts := strings.Split(cmdStr, " -")
  cmdName := parts[0]
  return fmt.Errorf("Node %s:\n%s: %s", node, cmdName, err.Error())
}

func SingleRemoteRun(programName string, cmdStr string, path string, fileName string) {
  var cmdMap owntypes.CommandMap = make(owntypes.CommandMap)


  // fmt.Println("filename: ", fileName, "\ncommand: ", cmdStr)

  // fmt.Println("Executing single remote run...")
  cmdMap[fileName] = cmdStr

  RemoteRun(programName, cmdMap, path)
  // fmt.Println("Single remote run completed.")
}

func RemoteRun(programName string, fcomms owntypes.CommandMap, path string) { //TODO: comms may just be one element, which is fine ALSO: revert comm, _ into: item, comm
  if len(os.Args) < 2 {
		fmt.Printf("usage: %s [yaml file]", programName)
		os.Exit(1)
	}
	var (
		pmap ConnectMap
		err  error
	)

  if pmap, err = yamlio.ReadYaml[ConnectMap](os.Args[1]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for n, f := range pmap {
    for item, fcomm := range fcomms {
      if err = saveCmdOutput(n, f, fcomm, path + item); err != nil {
  			fmt.Println("Error executing command: '", fcomm, "\n\nContinuing program nonetheless...")
  		}
    }
	}
}

//maybe we want to specify the path to pull it to as well (the outDir)
func RemotePull(paths []string, fileName string) {
  var (
    programName string = "RemotePull"
    outDir string = "artifacts/output/filesystem/"
    cmdMap owntypes.CommandMap = make(owntypes.CommandMap)
  )

  cmdMap[fileName] = ownformatters.MakeSudo(fmt.Sprintf("tar chpzf - --exclude='.git' --exclude='/var/run/*' %[1]s", strings.Join(paths, " "))) //

  RemoteRun(programName, cmdMap, outDir)
}

func GetRemoteRunMap(programName string, fcomms owntypes.CommandMap) map[string]owntypes.StringSet { //TODO: comms may just be one element, which is fine ALSO: revert comm, _ into: item, comm
  if len(os.Args) < 2 {
		fmt.Printf("usage: %s [yaml file]", programName)
		os.Exit(1)
	}

	var (
		pmap ConnectMap
		err  error
    result map[string]owntypes.StringSet = make(map[string]owntypes.StringSet)
	)

  if pmap, err = yamlio.ReadYaml[ConnectMap](os.Args[1]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, f := range pmap {
    user := strings.Split(f.Connection, "@")[0]
    if user == "root" { owntypes.UseSudo = false } else { owntypes.UseSudo = true }
    for item, fcomm := range fcomms {
      if err = buildCmdOutputMap(f, fcomm, item, result); err != nil {
  			fmt.Println("Error executing command: '", fcomm, "\n\nContinuing program nonetheless...")
  		}
    }
	}
  return result
}

func buildCmdOutputMap(form ConnectForm, cmdStr string, item string, mapbuilder map[string]owntypes.StringSet) error {
  var (
    programName string = "buildCmdOutputMap"
    cmd    *exec.Cmd

  )

  cmd = scmd.RemoteCommand(form.Connection, form.Identity, cmdStr)


  result, err := cmd.CombinedOutput()

  if err != nil {
    fmt.Println(programName, ": command unsuccessful")
    return err
  }

  mapbuilder[item] = ownformatters.FromStringLinesToStringSet(string(result))
  return nil
}

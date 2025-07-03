package runHelper

import (
  "fmt"
  "os"
  "os/exec"
  "bufio"
  "strings"
  "project1/tracepick/lib/own/owndialogues"
  "project1/tracepick/lib/editors/fileedit"
)


// "docker run -d --name c-redis-server -p 6379:6379 i-redis-server || docker run -d --name c-redis-server -p 6379:6379 -u root i-redis-server"

func RunImage(inPath, d1name string) {
  //idea: check exposed ports, and ask to which port on the local machine they should be forwarded
  var (
    programName string = "RunImage"
    dockerfilePath string = inPath + d1name + "/Dockerfile"
    portPattern string = "(?m)^EXPOSE.*$"
    formatRunCmd string = "docker run -d --name c-%[1]s %[2]s %[3]s --memory='1024m' --cpus='1' i-%[1]s && sleep 10"  //[1] d1name, [2] port options, [3] root user -- same resource usage as Vagrantfile VMs
    successProbe string = "docker inspect -f '{{.State.ExitCode}}' c-%s" // d1name
    cmd *exec.Cmd
  )

  exposeLines := fileedit.GetAllMatchingString(dockerfilePath, portPattern)
  var sb strings.Builder

  fmt.Println("EXPOSE lines found:")
  fmt.Println(exposeLines)

  repeat := true
  for _, line := range exposeLines {
    repeat = true
    for repeat {
      parts := strings.Split(line, " ")
      containerport := strings.TrimSpace(parts[1])
      localport, err := giveOpenPort(programName)
      if err != nil {
        fmt.Println(fmt.Sprintf("%[1]s: Finding open host port for port %[2]s of container %[3]s failed: %[4]s. Skipping this port.", programName, containerport, d1name, err))
        continue
      }
      localport = strings.TrimSpace(localport)
      setting := fmt.Sprintf(" -p %[1]s:%[2]s ", localport, containerport)
      sb.WriteString(setting)
      repeat = false
    }
  }

  portSettings := sb.String()

  fmt.Println("Current port settings:", portSettings)
  scanner := bufio.NewScanner(os.Stdin) //Going to need user input
  portSettings = owndialogues.EditString(scanner, portSettings, false) //set to FALSE if port exposure does not matter

  rootUser := ""
  runCmdStr := fmt.Sprintf(formatRunCmd, d1name, portSettings, rootUser)

  fmt.Println("Starting docker container (wait 10s)...")
  runDockerContainer(programName, runCmdStr)

  probeStr := fmt.Sprintf(successProbe, d1name)
  cmd = exec.Command("bash", "-c", probeStr)

  out, err := cmd.CombinedOutput()
  if err != nil { fmt.Println("Probing container failed:", err) }

  outStr := strings.TrimSpace(string(out))
  if outStr != "0" {
    // fmt.Println("EXIT CODE:", outStr)
    fmt.Println("Docker run for defined user failed, trying with root user...")
    rootUser = "-u root"
    rmCmdStr := fmt.Sprintf("docker stop c-%[1]s && docker rm c-%[1]s", d1name)
    runCmdStr = fmt.Sprintf(formatRunCmd, d1name, portSettings, rootUser)
    runDockerContainer(programName, rmCmdStr)
    runDockerContainer(programName, runCmdStr)
  }

  fmt.Println("Please check if container is running properly.")
}

func runDockerContainer(programName string, runCmdStr string) {
  cmd := exec.Command("bash", "-c", runCmdStr)

  if err := cmd.Run(); err != nil {
    fmt.Println(programName, ": Docker run failed:", err)
  }
}

func giveOpenPort(programName string) (string, error) {
  getOpenPortCmd := `python3 -c 'import socket; s=socket.socket(); s.bind(("", 0)); print(s.getsockname()[1]); s.close()'`
  cmd := exec.Command("bash", "-c", getOpenPortCmd)
  bytes, err := cmd.CombinedOutput()
  return string(bytes), err
}

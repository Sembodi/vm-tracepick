package traceHelper

import (
  "fmt"
  "sync"
  "time"
  "project1/tracepick/lib/own/owntypes"
  "project1/tracepick/lib/own/ownformatters"
  "project1/tracepick/lib/helpers/fsHelper"
  "project1/tracepick/lib/helpers/vmHelpers/serviceHelper"
  "project1/tracepick/lib/runners/remoterun"
)
// BEGIN { printf("PID     CMD              FD   ERR  PATH\n"); }


var bpfScript string = `sudo timeout --preserve-status --signal=SIGINT 10s bpftrace -e '


                        tracepoint:syscalls:sys_enter_open,
                        tracepoint:syscalls:sys_enter_openat,
                        tracepoint:syscalls:sys_enter_openat2 {
                          @filename[tid] = str(args->filename);
                        }

                        tracepoint:syscalls:sys_exit_open,
                        tracepoint:syscalls:sys_exit_openat,
                        tracepoint:syscalls:sys_exit_openat2 {
                          $fd = args->ret;
                          $tid = tid;

                          $name = @filename[$tid];
                          if ($name != "") {
                              printf("%-6d %-16s %-4d %-4d %s\n", pid, comm, $fd >= 0 ? $fd : -1, $fd < 0 ? $fd : 0, $name);
                              delete(@filename[$tid]);
                          }
                        } ' | awk '!seen[$5]++' > bpfout.log `

func RunBPFtrace(item string) {
  var (
    programName string = "RunBPFtrace"
    bpfRunStr string = ownformatters.MakeSudo(bpfScript)

    bpfOutStr string = ownformatters.MakeSudo("cat bpfout.log && rm bpfout.log")

    path string = "outputs/bpfout/"
    fileName string = "bpfout_" + item
  )

  remoterun.SingleRemoteRun(programName, bpfRunStr, owntypes.NoOutputFile, owntypes.NoOutputFile)

  fmt.Println("Now retrieving all bpf events...")
  remoterun.SingleRemoteRun(programName, bpfOutStr, path, fileName)
}

func TraceCommand(item string, cmdFunc func(string)) {
  var (
    wg sync.WaitGroup
  )

  fmt.Printf("Executing concurrent commands for %s...\n", item)

  wg.Add(2)


  go func() {  // first func lasts 10 seconds
        defer wg.Done()
        fmt.Println("Starting BPFtrace...")
        RunBPFtrace(ownformatters.StripServiceString(item))
        fmt.Println("BPFtrace done")
  }()

  go func() {
        defer wg.Done()

        time.Sleep(3 * time.Second)

        // fmt.Println("Starting cmd2...")

        cmdFunc(item)

        // fmt.Println("Command 2 done")
  }()

  wg.Wait()
  fmt.Println("Both threads finished")
}

func TraceServiceRestart(items owntypes.StringSet) {
  var path string = "outputs/bpfout/"
  fsHelper.CleanPath(path)
  for item, _ := range items {
    TraceCommand(item, serviceHelper.RestartService)
  }
}

func TraceWrk() {
  var (
    itemStr string = "wrkCmd"
  )

  TraceCommand(itemStr, RunWrk)
}

func RunWrk(item string) {
  var (
    programName string              = "RunWrk"
    cmdStr string                   = ownformatters.MakeSudo("wrk -t1 -c1 -d3s http://localhost/")
  )

  fmt.Println(programName, "started...")

  remoterun.SingleRemoteRun(programName, cmdStr, owntypes.NoOutputFile, owntypes.NoOutputFile)

  fmt.Println(programName, "done.")
}



func TraceFFUF() {
  var (
    itemStr string = "ffuf"
  )

  TraceCommand(itemStr, RunFFUF)
}

func RunFFUF(item string) {
  var (
    programName string              = "RunFFUF"
    restartStr string               = ownformatters.MakeSudo("ffuf -w SecLists/Discovery/Web-Content/common.txt -u http://localhost/FUZZ -mc 200,301,302")
  )

  fmt.Println(programName, "started...")

  remoterun.SingleRemoteRun(programName, restartStr, owntypes.NoOutputFile, owntypes.NoOutputFile)

  fmt.Println(programName, "done.")
}

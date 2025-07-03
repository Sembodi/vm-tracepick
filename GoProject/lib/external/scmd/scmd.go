package scmd

import (
	"fmt"
	"strings"
	"os/exec"
)

func FormatCreds(user, host string) string {
	return fmt.Sprint(user, "@", host)
}

/*
Return command to execute over SSH.
Discard the identity option by keeping it ""
*/

//TODO ADD -p OPTION IF CREDS HAVE A :{} suffix
func RemoteCommand(creds, identity, cmdStr string) *exec.Cmd {
	cmd := exec.Command(
		"ssh",
		"-o ConnectTimeout=10",
		"-o StrictHostKeyChecking=no", // don't do this in prod (exposure to MITM-attacks)!
	)

	credParts := strings.Split(creds,":")

	if len(credParts) > 1 {
		cmd.Args = append(cmd.Args, fmt.Sprintf("-p %s", credParts[1])) //add -p port option if necessary
	}

	if identity != "" {
		cmd.Args = append(cmd.Args, "-i")
		cmd.Args = append(cmd.Args, identity)
	}

	cmd.Args = append(cmd.Args, credParts[0])
	cmd.Args = append(cmd.Args, cmdStr)


	return cmd
}

// old ssh command script:
// func RemoteCommand1(creds, identity, cmd string) *exec.Cmd {
// 	if identity == "" {
// 		return exec.Command(
// 			"ssh",
// 			"-o ConnectTimeout=10",
// 			creds,
// 			cmd,
// 		)
// 	}
// 	return exec.Command(
// 		"ssh",
// 		"-o ConnectTimeout=10",
// 		"-i",
// 		identity,
// 		creds,
// 		cmd,
// 	)
// }

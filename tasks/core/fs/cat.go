package fs

import (
	"github.com/pygrum/Empress/transport"
	"os"
)

func CmdCat(req *transport.Request, response *transport.Response) {
	path := string(req.Args[0])
	bytes, err := os.ReadFile(path)
	if err != nil {
		transport.ResponseWithError(response, err)
		return
	}
	transport.AddResponse(response, transport.ResponseDetail{
		Status: transport.StatusSuccess,
		Dest:   transport.DestStdout,
		Data:   bytes,
	})
}

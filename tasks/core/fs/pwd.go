package fs

import (
	"github.com/pygrum/Empress/transport"
	"os"
)

func CmdPWD(_ *transport.Request, response *transport.Response) {
	dir, err := os.Getwd()
	if err != nil {
		transport.ResponseWithError(response, err)
		return
	}
	transport.AddResponse(response, transport.ResponseDetail{
		Status: transport.StatusSuccess,
		Dest:   transport.DestStdout,
		Data:   []byte(dir),
	})
}

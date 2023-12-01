package fs

import (
	"github.com/pygrum/Empress/transport"
	"os"
)

func CmdMkdir(req *transport.Request, response *transport.Response) {
	for _, arg := range req.Args {
		dir := string(arg)
		if err := os.Mkdir(dir, 0777); err != nil {
			transport.ResponseWithError(response, err)
			return
		}
	}
	transport.AddResponse(response, transport.ResponseDetail{
		Status: transport.StatusSuccess,
		Dest:   transport.DestNone,
	})
}

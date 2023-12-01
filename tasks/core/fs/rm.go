package fs

import (
	"github.com/pygrum/Empress/transport"
	"os"
)

func CmdRm(req *transport.Request, response *transport.Response) {
	for _, i := range req.Args {
		item := string(i)
		if err := os.Remove(item); err != nil {
			transport.ResponseWithError(response, err)
			continue
		}
	}
	transport.AddResponse(response, transport.ResponseDetail{
		Status: transport.StatusSuccess,
		Dest:   transport.DestNone,
	})
}

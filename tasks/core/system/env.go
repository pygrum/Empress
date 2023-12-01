package system

import (
	"github.com/pygrum/Empress/transport"
	"os"
)

func CmdEnv(req *transport.Request, response *transport.Response) {
	if len(req.Args) == 0 {
		for _, e := range os.Environ() {
			transport.AddResponse(response, transport.ResponseDetail{
				Status: transport.StatusSuccess,
				Dest:   transport.DestStdout,
				Data:   []byte(e),
			})
		}
		return
	}
	for _, arg := range req.Args {
		key := string(arg)
		val := os.Getenv(key)
		transport.AddResponse(response, transport.ResponseDetail{
			Status: transport.StatusSuccess,
			Dest:   transport.DestStdout,
			Data:   []byte(val),
		})
	}
}

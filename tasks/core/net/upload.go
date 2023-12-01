package net

import (
	"github.com/pygrum/Empress/transport"
	"os"
)

func CmdUpload(req *transport.Request, response *transport.Response) {
	for _, arg := range req.Args {
		fileName := string(arg)
		b, err := os.ReadFile(fileName)
		if err != nil {
			transport.ResponseWithError(response, err)
			continue
		}
		transport.AddResponse(response, transport.ResponseDetail{
			Status: transport.StatusSuccess,
			Dest:   transport.DestFile,
			Name:   fileName,
			Data:   b,
		})
	}
}

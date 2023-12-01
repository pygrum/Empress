package net

import (
	"encoding/binary"
	"errors"
	"github.com/pygrum/Empress/transport"
	"os"
)

func CmdDownload(req *transport.Request, response *transport.Response) {
	defer func() {
		if r := recover(); r != nil {
			transport.ResponseWithError(response, errors.New("specified file size exceeds packet length"))
		}
	}()
	for _, arg := range req.Args {
		filesize := binary.BigEndian.Uint64(arg[:8])
		fileBytes := arg[8 : filesize+8]
		fileName := string(arg[filesize+8:])
		if err := os.WriteFile(fileName, fileBytes, 0600); err != nil {
			transport.ResponseWithError(response, err)
			continue
		}
	}
	transport.AddResponse(response, transport.ResponseDetail{
		Status: transport.StatusSuccess,
		Dest:   transport.DestNone,
	})
}

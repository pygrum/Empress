package fs

import (
	"github.com/pygrum/Empress/transport"
	"os"
	"path/filepath"
)

func CmdCP(req *transport.Request, response *transport.Response) {
	src := string(req.Args[0])
	dst := string(req.Args[1])
	info, err := os.Stat(src)
	if err != nil {
		transport.ResponseWithError(response, err)
		return
	}
	oldPerm := info.Mode()
	b, err := os.ReadFile(src)
	if err != nil {
		transport.ResponseWithError(response, err)
		return
	}

	info, err = os.Stat(dst)
	if err == nil {
		if info.IsDir() {
			if err = os.WriteFile(filepath.Join(dst, filepath.Base(src)), b, oldPerm); err != nil {
				transport.ResponseWithError(response, err)
				return
			}
			transport.AddResponse(response, transport.ResponseDetail{
				Status: transport.StatusSuccess,
				Dest:   transport.DestNone,
			})
			return
		}
	}
	if err = os.WriteFile(dst, b, oldPerm); err != nil {
		transport.ResponseWithError(response, err)
		return
	}
	transport.AddResponse(response, transport.ResponseDetail{
		Status: transport.StatusSuccess,
		Dest:   transport.DestNone,
	})
}

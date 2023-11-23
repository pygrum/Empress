package commands

import (
	"fmt"
	"os"
)

func CmdLS(req *Request, response *Response) {
	path := string(req.Args[0])
	entries, err := os.ReadDir(path)
	if err != nil {
		response.Status = 1
		response.Responses = []ResponseDetail{
			{Dest: 1, Data: []byte(fmt.Sprintf("could not retrieve entries: %v", err))},
		}
		return
	}
	response.Status = 0
	for _, d := range entries {
		response.Responses = append(response.Responses, ResponseDetail{
			Dest: 1,
			Data: []byte(d.Name()),
		})
	}
}

type Request struct {
	AgentID   string
	RequestID string
	Opcode    int32
	Args      [][]byte
}

type ResponseDetail struct {
	Dest int32  // Where to send response to (file, stdout)
	Name string // Name of file to save, if applicable
	Data []byte // file or output data
}

type Response struct {
	AgentID   string
	RequestID string
	Status    int32
	Responses []ResponseDetail
}

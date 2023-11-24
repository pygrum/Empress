package commands

import (
	"os"
)

func CmdLS(req *Request, response *Response) {
	if len(req.Args) == 0 {
		req.Args = append(req.Args, []byte("."))
	}
	for _, path := range req.Args {
		entries, err := os.ReadDir(string(path))
		if err != nil {
			response.Status = 1
			response.Responses = []ResponseDetail{
				{Dest: 1, Data: []byte(err.Error())},
			}
			return
		}
		for _, d := range entries {
			response.Responses = append(response.Responses, ResponseDetail{
				Dest: 1,
				Data: []byte(d.Name()),
			})
		}
	}
}

type Request struct {
	AgentID   string   `json:"agent_id"`
	RequestID string   `json:"request_id"`
	Opcode    int32    `json:"opcode"`
	Args      [][]byte `json:"args"`
}

type ResponseDetail struct {
	Dest int32  `json:"dest"` // Where to send response to (file, stdout)
	Name string `json:"name"` // Name of file to save, if applicable
	Data []byte `json:"data"` // file or output data
}

type Response struct {
	AgentID   string           `json:"agent_id"`
	RequestID string           `json:"request_id"`
	Status    int32            `json:"status"`
	Responses []ResponseDetail `json:"responses"`
}

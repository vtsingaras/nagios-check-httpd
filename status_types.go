package main


import (
	"time"
	"github.com/simplereach/timeutils"
)


type Connections struct {
	Active uint
	Idle uint
}

type Server struct {
	Connections uint
	Built timeutils.Time
	Uptime	time.Duration
	Localtime timeutils.Time
	Extended bool
	Bytes uint
	Host string
	Version string
}

type Mpm struct {
	Type string
	ActiveServers uint
	MaxServers uint
	ThreadsPerChild uint
	Threaded bool
}

type Thread struct {
	LastUsed timeutils.Time `json:"last_used"`
	Request string
	Client string
	Cost uint
	Bytes uint
	Count uint
	Thread string
	Vhost string
}

type Process struct {
	Connections uint
	WorkerStates struct {
		Idle uint
		Reading uint
		Keepalive uint
		Graceful uint
		Closing uint
		Writing uint
	}
	Bytes uint
	Utime time.Duration
	Active bool
	Stime time.Duration
	Pid uint
	Threads []Thread `json:",omitempty"`
}

type HttpdStatus struct {
	Connections Connections
	Server Server
	Processes []Process
	Mpm Mpm
}
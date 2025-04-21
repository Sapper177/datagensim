package sim

import (
	"fmt"
	"net"
	"time"
)

type Config struct {
	BusName		string
	Interface   net.Interface
	BusType		string
	SrcHost		net.IP
	DestHost	net.IP
	SrcPort		int
	DestPort	int

	DbHost		string
	DbPort		string
	DbPassword	string
	DbNum		int
	DbReadTimeout time.Duration
	DbWriteTimeout time.Duration
	
	LogFile		string
	LogLevel	string
}

type LogLevels int

const (
	Debug LogLevels = iota
	Info
	Warn
	Error
)

func (l LogLevels) String() string {
	switch l {
	case Debug:
		return "debug"
	case Info:
		return "info"
	case Warn:
		return "warn"
	case Error:
		return "error"
	default:
		return "unknown"
	}
}

func (l *LogLevels) Set(value string) error {
	switch value {
	case "debug":
		*l = Debug
	case "info":
		*l = Info
	case "warn":
		*l = Warn
	case "error":
		*l = Error
	default:
		return fmt.Errorf("invalid log level: %s", value)
	}
	return nil
}

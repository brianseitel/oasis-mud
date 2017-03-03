package helpers

import (
	"github.com/davecgh/go-spew/spew"
)

func Dump(i interface{}) {
	spew.Dump(i)
}

const (
	Reset       = "\x1B[0m"
	Bold        = "\x1B[1m"
	Dim         = "\x1B[2m"
	Under       = "\x1B[4m"
	Reverse     = "\x1B[7m"
	Hide        = "\x1B[8m"
	Clearscreen = "\x1B[2J"
	Clearline   = "\x1B[2K"
	Black       = "\x1B[30m"
	Red         = "\x1B[31m"
	Green       = "\x1B[32m"
	Yellow      = "\x1B[33m"
	Blue        = "\x1B[34m"
	Magenta     = "\x1B[35m"
	Cyan        = "\x1B[36m"
	White       = "\x1B[37m"
	Bblack      = "\x1B[40m"
	Bred        = "\x1B[41m"
	Bgreen      = "\x1B[42m"
	Byellow     = "\x1B[43m"
	Bblue       = "\x1B[44m"
	Bmagenta    = "\x1B[45m"
	Bcyan       = "\x1B[46m"
	Bwhite      = "\x1B[47m"
	Newline     = "\r\n\x1B[0m"
)

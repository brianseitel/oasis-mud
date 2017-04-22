package mud

import (
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

func dump(i interface{}) {
	spew.Dump(i)
}

func dd(i interface{}) {
	spew.Dump(i)
	os.Exit(1)
}

func isSameGroup(p1 *mob, p2 *mob) bool {
	if p1.leader != nil {
		p1 = p1.leader
	}

	if p2.leader != nil {
		p2 = p2.leader
	}

	return p1 == p2
}

func matchesSubject(list string, s string) bool {
	for _, v := range strings.Split(strings.ToLower(list), " ") {
		if strings.HasPrefix(v, s) {
			return true
		}
	}

	return false
}

func transferItem(i int, from []*item, to []*item) ([]*item, []*item) {
	item := from[i]
	from = append(from[0:i], from[i+1:]...)
	to = append(to, item)

	return from, to
}

const (
	reset       = "\x1B[0m"
	bold        = "\x1B[1m"
	dim         = "\x1B[2m"
	under       = "\x1B[4m"
	reverse     = "\x1B[7m"
	hide        = "\x1B[8m"
	clearscreen = "\x1B[2J"
	clearline   = "\x1B[2K"
	black       = "\x1B[30m"
	red         = "\x1B[31m"
	green       = "\x1B[32m"
	yellow      = "\x1B[33m"
	blue        = "\x1B[34m"
	magenta     = "\x1B[35m"
	cyan        = "\x1B[36m"
	white       = "\x1B[37m"
	bBlack      = "\x1B[40m"
	bRed        = "\x1B[41m"
	bGreen      = "\x1B[42m"
	bYellow     = "\x1B[43m"
	bBlue       = "\x1B[44m"
	bMagenta    = "\x1B[45m"
	bCyan       = "\x1B[46m"
	bWhite      = "\x1B[47m"
	newline     = "\r\n\x1B[0m"
)

package mud

import (
	"os"
	"strings"
	"unicode"

	"github.com/davecgh/go-spew/spew"
)

func dump(i interface{}) {
	spew.Dump(i)
}

func dd(i interface{}) {
	spew.Dump(i)
	os.Exit(1)
}

func sliceContainsUint(s []uint, e uint) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func wordWrap(text string, lineWidth int) string {
	words := strings.Fields(strings.TrimSpace(text))
	if len(words) == 0 {
		return text
	}
	wrapped := words[0]
	spaceLeft := lineWidth - len(wrapped)
	for _, word := range words[1:] {
		if len(word)+1 > spaceLeft {
			wrapped += "\n" + word
			spaceLeft = lineWidth - len(word)
		} else {
			wrapped += " " + word
			spaceLeft -= 1 + len(word)
		}
	}

	return wrapped
}

func matchesSubject(list string, s string) bool {
	for _, v := range strings.Split(strings.ToLower(list), " ") {
		if strings.HasPrefix(v, s) {
			return true
		}
	}

	return false
}

// ToSnake convert the given string to snake case following the Golang format:
// acronyms are converted to lower-case and preceded by an underscore.
func toSnake(in string) string {
	runes := []rune(in)
	length := len(runes)

	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}

	return string(out)
}

func transferItem(i int, from []*item, to []*item) ([]*item, []*item) {
	item := from[i]
	from = append(from[0:i], from[i+1:]...)
	to = append(to, item)

	return from, to
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

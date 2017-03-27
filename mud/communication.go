package mud

import (
	"strings"

	"github.com/brianseitel/oasis-mud/helpers"
)

const (
	actToRoom = iota
	actToNotVict
	actToVict
	actToChar
)

func act(format string, player *mob, arg1 interface{}, arg2 interface{}, target int) {

	var (
		heShe  = []string{"it", "he", "she"}
		himHer = []string{"it", "him", "her"}
		hisHer = []string{"its", "his", "her"}
	)

	if len(format) == 0 {
		return
	}

	to := player.Room.Mobs

	var (
		victim   *mob
		obj1     *item
		obj2     *item
		message1 string
		message2 string
	)

	switch arg1.(type) {

	case *item:
		obj1 = arg1.(*item)
		break
	case string:
		message1 = arg1.(string)
	default:
		message1 = ""
	}

	switch arg2.(type) {
	case *mob:
		victim = arg2.(*mob)
		break
	case *item:
		obj2 = arg2.(*item)
		break
	case string:
		message2 = arg2.(string)
	default:
		message2 = ""
	}

	if target == actToVict {
		if victim == nil {
			return
		}

		to = victim.Room.Mobs
	}

	for _, m := range to {
		if m.client == nil || !m.isAwake() {
			continue
		}

		if target == actToChar && m != player {
			continue
		}

		if target == actToVict && (m != victim || m == player) {
			continue
		}

		if target == actToRoom && m == player {
			continue
		}

		if target == actToNotVict && (m == player || m == victim) {
			continue
		}

		if strings.Contains(format, "$t") {
			format = strings.Replace(format, "$t", message1, -1)
		}

		if strings.Contains(format, "$T") {
			format = strings.Replace(format, "$T", message2, -1)
		}

		if strings.Contains(format, "$n") {
			format = strings.Replace(format, "$n", m.looksAt(player), -1)
		}

		if strings.Contains(format, "$N") {
			format = strings.Replace(format, "$N", player.looksAt(victim), -1)
		}

		if strings.Contains(format, "$e") {
			format = strings.Replace(format, "$e", heShe[helpers.Range(0, player.Gender, 2)], -1)
		}
		if strings.Contains(format, "$E") {
			format = strings.Replace(format, "$E", heShe[helpers.Range(0, victim.Gender, 2)], -1)
		}
		if strings.Contains(format, "$m") {
			format = strings.Replace(format, "$m", himHer[helpers.Range(0, player.Gender, 2)], -1)
		}
		if strings.Contains(format, "$M") {
			format = strings.Replace(format, "$M", himHer[helpers.Range(0, victim.Gender, 2)], -1)
		}

		if strings.Contains(format, "$s") {
			format = strings.Replace(format, "$s", hisHer[helpers.Range(0, player.Gender, 2)], -1)
		}

		if strings.Contains(format, "$S") {
			format = strings.Replace(format, "$S", hisHer[helpers.Range(0, victim.Gender, 2)], -1)
		}

		if strings.Contains(format, "$p") {
			format = strings.Replace(format, "$p", m.looksAt(obj1), -1)
		}

		if strings.Contains(format, "$P") {
			format = strings.Replace(format, "$P", m.looksAt(obj2), -1)
		}

		if strings.Contains(format, "$d") {
			format = strings.Replace(format, "$d", "door", -1)
		}

		m.notify(format)
	}
}

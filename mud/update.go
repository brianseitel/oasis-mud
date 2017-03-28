package mud

import "github.com/brianseitel/oasis-mud/helpers"

const (
	pulsePerSecond = 4
	pulseViolence  = 3 * pulsePerSecond
	pulseMobile    = 4 * pulsePerSecond
	pulseTick      = 30 * pulsePerSecond
	pulseArea      = 60 * pulsePerSecond
)

var (
	pulseTimerArea     = pulseArea
	pulseTimerMobs     = pulseMobile
	pulseTimerViolence = pulseViolence
	pulseTimerPoint    = pulseTick
)

func objUpdate() {
	for e := itemList.Front(); e != nil; e = e.Next() {
		item := e.Value.(*item)

		if item.Timer <= 0 {
			continue
		}

		item.Timer--

		if item.Timer > 0 {
			continue
		}

		var message string
		switch item.ItemType {
		default:
			message = "$p vanishes."
			break
		case itemFountain:
			message = "$p dries up."
			break
		case itemCorpseNPC:
			message = "$p crumbles into dust."
			break
		case itemCorpsePC:
			message = "$p decays into dust."
			break
		case itemFood:
			message = "$p decomposes."
			break
		}

		if item.carriedBy != nil {
			act(message, item.carriedBy, item, nil, actToChar)
		} else if item.Room != nil && len(item.Room.Mobs) > 0 {
			act(message, nil, item, nil, actToRoom)
			act(message, nil, item, nil, actToChar)
		}
	}
}

func updateHandler() {

	pulseTimerArea--
	pulseTimerMobs--
	pulseTimerViolence--
	pulseTimerPoint--

	if pulseTimerArea <= 0 {
		pulseTimerArea = pulseArea
		// area_update()
	}

	if pulseTimerMobs <= 0 {
		pulseTimerMobs = pulseMobile
		// mob_update()
	}

	if pulseTimerViolence <= 0 {
		pulseTimerViolence = pulseViolence
		violenceUpdate()
	}

	if pulseTimerPoint <= 0 {
		pulseTimerPoint = dice().Intn(3*pulseTick/2) + (pulseTick / 2)
		// weather_update()
		// char_update()
		// obj_update()
	}

	// aggr_update()
	return
}

func wait(player *mob, npulse int) {
	player.wait = uint(helpers.Max(int(player.wait), npulse))
}

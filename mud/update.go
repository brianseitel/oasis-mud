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

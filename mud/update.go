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

func aggroUpdate() {
	for e := mobList.Front(); e != nil; e = e.Next() {
		ch := e.Value.(*mob)

		if ch.isNPC() || ch.isImmortal() || ch.Room == nil {
			continue
		}

		for _, m := range ch.Room.Mobs {
			var count int

			if !m.isNPC() || !helpers.HasBit(m.Act, actAggressive) || m.Fight != nil || !m.isAwake() || !m.canSee(ch) {
				continue
			}

			count = 0
			var victim *mob
			for _, v := range m.Room.Mobs {
				if !v.isNPC() && v.isImmortal() && v.canSee(m) {
					if dice().Intn(count) == 0 {
						victim = v
					}
					count++
				}
			}

			if victim == nil {
				continue
			}

			multiHit(m, victim, typeUndefined)
		}
	}
}

func mobUpdate() {
	for e := mobList.Front(); e != nil; e = e.Next() {
		ch := e.Value.(*mob)

		if !ch.isNPC() || ch.Room == nil {
			continue
		}

		if ch.Status != standing {
			continue
		}

		/* Scavenge */
		if helpers.HasBit(ch.Act, actScavenger) && len(ch.Room.Items) > 0 && dBits(2) == 0 {
			max := 1
			var objectBest *item
			objKey := 0
			for j, item := range ch.Room.Items {
				if item.canWear(itemTake) && item.Cost > max {
					objectBest = item
					objKey = j
					max = item.Cost
				}
			}

			if objectBest != nil {
				ch.Inventory, ch.Room.Items = transferItem(objKey, ch.Inventory, ch.Room.Items)
				act("$n gets $p.", ch, objectBest, nil, actToRoom)
			}
		}

		/* wander */
		if !helpers.HasBit(ch.Act, actSentinel) {
			ch.wander()
		}

		/* flee */
		door := dBits(3)
		var exit *exit
		if door <= 5 {
			exit = ch.Room.Exits[door]
		}

		if ch.Hitpoints < ch.MaxHitpoints/2 && exit != nil && !helpers.HasBit(exit.Room.RoomFlags, roomNoMob) {
			found := false
			for _, rch := range exit.Room.Mobs {
				if !rch.isNPC() {
					found = true
					break
				}
			}

			if !found {
				ch.move(exit)
			}
		}
	}
}

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
		mobUpdate()
	}

	if pulseTimerViolence <= 0 {
		pulseTimerViolence = pulseViolence
		violenceUpdate()
	}

	if pulseTimerPoint <= 0 {
		pulseTimerPoint = dice().Intn(3*pulseTick/2) + (pulseTick / 2)
		// weather_update()
		// char_update()
		objUpdate()
	}

	aggroUpdate()
	return
}

func wait(player *mob, npulse int) {
	player.wait = uint(helpers.Max(int(player.wait), npulse))
}

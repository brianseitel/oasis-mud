package mud

import (
	"fmt"

	"github.com/brianseitel/oasis-mud/helpers"
)

type fight struct {
	Mob1 *mob
	Mob2 *mob
}

func newFight(m1 *mob, m2 *mob) *fight {
	if m2.isSafe() {
		m1.notify("A voice from the heavens booms, 'Thou shalt not kill!'\r\n")
		return nil
	}

	m1.notify(fmt.Sprintf("You scream and attack %s!%s", m2.Name, helpers.Newline))

	m1.Status = fighting
	m2.Status = fighting

	f := &fight{
		Mob1: m1,
		Mob2: m2,
	}

	m1.Fight = f
	m2.Fight = f
	f.turn(m1)

	return f
}

func (f *fight) turn(m *mob) {
	if m.ID == f.Mob1.ID {
		fmt.Println("Mob1 attacks")
		m.attack(f.Mob2, f)
	} else if m.ID == f.Mob2.ID {
		fmt.Println("Mob2 attacks")
		m.attack(f.Mob1, f)
	}
}

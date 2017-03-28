package mud

type fight struct {
	Mob1 *mob
	Mob2 *mob
}

const (
	typeHit = iota
)

func newFight(m1 *mob, m2 *mob) *fight {
	if m2.isSafe() {
		m1.notify("A voice from the heavens booms, 'Thou shalt not kill!'\r\n")
		return nil
	}

	m1.notify("You scream and attack %s!", m2.Name)

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
		m.attack(f.Mob2, f)
	} else if m.ID == f.Mob2.ID {
		m.attack(f.Mob1, f)
	}
}

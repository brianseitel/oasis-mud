package mud

import (
	"fmt"

	"github.com/brianseitel/oasis-mud/helpers"

	"github.com/jinzhu/gorm"
)

type fight struct {
	gorm.Model

	Mob1   *mob `gorm:"ForeignKey:Mob1ID"`
	Mob1ID uint
	Mob2   *mob `gorm:"ForeignKey:Mob2ID"`
	Mob2ID uint
}

func newFight(m1 *mob, m2 *mob) *fight {
	m1.notify(fmt.Sprintf("You scream and attack %s!%s", m2.Name, helpers.Newline))

	m1.Status = fighting
	db.Save(&m1)
	m2.Status = fighting
	db.Save(&m2)

	f := &fight{
		Mob1: m1,
		Mob2: m2,
	}

	db.Save(&f)
	m1.setFight(f)
	m2.setFight(f)
	f.turn(m1)

	return f
}

func (f *fight) turn(m *mob) {
	db.Find(&f)
	var (
		m1 mob
		m2 mob
	)
	db.Model(&f).Related(&m1, "Mob1").Related(&m2, "Mob2")
	f.Mob1 = &m1
	f.Mob2 = &m2

	if m1.RoomID != m2.RoomID {
		return
	}

	if m.ID == m1.ID {
		fmt.Println("Mob1 attacks")
		m.attack(&m2, f)
	} else if m.ID == m2.ID {
		fmt.Println("Mob2 attacks")
		m.attack(&m1, f)
	} else {
		fmt.Println("what")
	}
}

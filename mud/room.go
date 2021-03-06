package mud

import (
	"fmt"
	"strings"
)

const (
	roomDark     = 1
	roomNoMob    = 2
	roomIndoors  = 4
	roomPrivate  = 8
	roomSafe     = 16
	roomSolitary = 32
	roomPetShop  = 64
	roomNoRecall = 128
)

const (
	sectorInside = iota
	sectorCity
	sectorField
	sectorForest
	sectorHills
	sectorMountain
	sectorWaterSwim
	sectorWaterNoSwim
	sectorUnused
	sectorAir
	sectorDesert
	sectorMax = 99999
)

type area struct {
	ID         int
	Name       string  `json:"name"`
	Rooms      []*room `json:"rooms",gorm:"-"`
	age        int
	numPlayers int
}

type room struct {
	ID   int
	Name string

	Area        area
	AreaID      int
	Description string
	Exits       []*exit `gorm:"many2many:room_exits;"`
	ItemIds     []int   `json:"items" gorm:"-"`
	Items       []*item `gorm:"many2many:room_items;"`
	Mobs        []*mob  `gorm:"many2many:room_mobs;"`
	MobIds      []int   `gorm:"-" json:"mobs"`

	Light      int `json:"light"`
	RoomFlags  int `json:"room_flags"`
	SectorType int `json:"sector_type"`
}

func (r *room) decayItems() {
	for j, item := range r.Items {
		if item.Timer == -1 {
			continue
		}

		if item.Timer <= 0 {
			r.Items = append(r.Items[0:j], r.Items[j+1:]...)
			for _, m := range r.Mobs {
				if r.ID == m.Room.ID {
					m.notify("Rats scurry forth and drag away %s!", item.Name)
				}
			}
			break
		}
		item.Timer--
	}
}

func getRoom(id int) *room {
	for e := roomList.Front(); e != nil; e = e.Next() {
		r := e.Value.(*room)
		if r.ID == id {
			return r
		}
	}
	return nil
}

func getItem(id int) *itemIndex {
	for e := itemIndexList.Front(); e != nil; e = e.Next() {
		i := e.Value.(itemIndex)
		if i.ID == id {
			return &i
		}
	}
	return nil
}

func getMob(id int) *mobIndex {
	for e := mobIndexList.Front(); e != nil; e = e.Next() {
		m := e.Value.(*mobIndex)
		if m.ID == id {
			return m
		}
	}
	return nil
}

func (r *room) isDark() bool {
	if r.Light > 0 {
		return false
	}

	if hasBit(r.RoomFlags, roomDark) {
		return true
	}

	if r.SectorType == sectorInside || r.SectorType == sectorCity {
		return true
	}

	return false
}

func (r *room) isPrivate() bool {
	count := len(r.Mobs)

	if hasBit(r.RoomFlags, roomPrivate) && count >= 2 {
		return true
	}

	if hasBit(r.RoomFlags, roomSolitary) && count >= 1 {
		return true
	}

	return false
}

func (r *room) notify(message string, except *mob) {
	for _, mob := range r.Mobs {
		if mob != except {
			mob.notify(message)
		}
	}
}

func (r *room) removeMob(m *mob) {
	for j, mob := range r.Mobs {
		if mob == m {
			r.Mobs = append(r.Mobs[0:j], r.Mobs[j+1:]...)
			break
		}
	}
}

func (r *room) removeObject(i *item) {
	for j, item := range r.Items {
		if item == i {
			r.Items = append(r.Items[0:j], r.Items[j+1:]...)
			return
		}
	}
}

func (r *room) showExits(player *mob) {
	var output string
	for _, e := range r.Exits {
		output = fmt.Sprintf("%s%s ", output, string(e.Dir))
	}

	player.notify(fmt.Sprintf("%s[%s]%s", white, strings.Trim(output, " "), reset))
}

func (r *room) findExit(arg string) *exit {
	var door string
	if strings.HasPrefix(arg, "n") || strings.HasPrefix(arg, "north") {
		door = "north"
	} else if strings.HasPrefix(arg, "e") || strings.HasPrefix(arg, "east") {
		door = "east"
	} else if strings.HasPrefix(arg, "s") || strings.HasPrefix(arg, "south") {
		door = "south"
	} else if strings.HasPrefix(arg, "w") || strings.HasPrefix(arg, "west") {
		door = "west"
	} else if strings.HasPrefix(arg, "u") || strings.HasPrefix(arg, "up") {
		door = "up"
	} else if strings.HasPrefix(arg, "d") || strings.HasPrefix(arg, "down") {
		door = "down"
	} else {
		return nil
	}

	for _, e := range r.Exits {
		if e.Dir == door {
			return e
		}
	}

	return nil
}

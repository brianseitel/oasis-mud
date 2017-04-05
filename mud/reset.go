package mud

import "fmt"

type resetData struct {
	AreaID            int                   `json:"area_id"`
	Mobs              []mobReset            `json:"mobs"`
	Items             []itemReset           `json:"objects"`
	ContainerItems    []containerItemReset  `json:"container_objects"`
	MobObjects        []mobObjectReset      `json:"mob_objects"`
	MobObjectsToEquip []mobObjectEquipReset `json:"mob_objects_equip"`
	Doors             []doorReset           `json:"doors"`
}

type mobReset struct {
	MobID  int `json:"mob_id"`
	Limit  int `json:"limit"`
	RoomID int `json:"room_id"`
}

type itemReset struct {
	ItemID int `json:"item_id"`
	Limit  int `json:"limit"`
	RoomID int `json:"room_id"`
}

type containerItemReset struct {
	ItemID      int `json:"item_id"`
	Limit       int `json:"limit"`
	ContainerID int `json:"container_id"`
}

type mobObjectReset struct {
	ItemID int `json:"item_id"`
	MobID  int `json:"mob_id"`
	RoomID int `json:"room_id"`
	Limit  int `json:"limit"`
}

type mobObjectEquipReset struct {
	ItemID       int `json:"item_id"`
	MobID        int `json:"mob_id"`
	WearLocation int `json:"wear_location"`
	RoomID       int `json:"room_id"`
}

type doorReset struct {
	RoomID    int    `json:"room_id"`
	Direction string `json:"direction"`
	State     string `json:"state"`
}

func resetArea(ar *area) {

	var reset *resetData
	for e := resetList.Front(); e != nil; e = e.Next() {
		r := e.Value.(*resetData)
		if r.AreaID == ar.ID {
			reset = r
			break
		}
	}

	if reset == nil {
		return
	}

	// reset mobs
	for _, mr := range reset.Mobs {
		mobIndex := getMob(mr.MobID)
		if mobIndex == nil {
			continue
		}

		roomIndex := getRoom(mr.RoomID)

		if roomIndex == nil {
			continue
		}

		here := false
		for _, m := range roomIndex.Mobs {
			if m.index.ID == mobIndex.ID {
				here = true
			}
		}

		if here {
			// already have it
			continue
		}

		mob := createMob(mobIndex)

		if roomIndex.isDark() {
			mob.AffectedBy = setBit(mob.AffectedBy, affectInfrared)
		}

		mob.Room = roomIndex
		roomIndex.Mobs = append(roomIndex.Mobs, mob)
		fmt.Println("RESET MOB", mob.Name)
	}

	// reset objects
	for _, mo := range reset.Items {
		ir := getItem(mo.ItemID)
		room := getRoom(mo.RoomID)

		if ir == nil || room == nil {
			continue
		}

		exists := false
		for _, it := range room.Items {
			if it.index.ID == ir.ID {
				exists = true
				break
			}
		}

		if exists {
			continue
		}

		object := createItem(ir)
		object.Cost = 0
		object.Room = room
		room.Items = append(room.Items, object)
	}

	// reset objects inside containers

	// reset objects carried by mobs
	for _, mo := range reset.MobObjects {
		ir := getItem(mo.ItemID)
		m := getMob(mo.MobID)
		room := getRoom(mo.RoomID)

		if ir == nil || m == nil || room == nil {
			continue
		}

		var olevel int
		if m.Shop != nil {

			switch ir.ItemType {
			default:
				olevel = 0
				break
			case itemPill:
				olevel = dice().Intn(10)
				break
			case itemPotion:
				olevel = dice().Intn(10)
				break
			case itemScroll:
				olevel = dice().Intn(10) + 5
				break
			case itemWand:
				olevel = dice().Intn(10) + 10
				break
			case itemStaff:
				olevel = dice().Intn(10) + 15
				break
			case itemWeapon:
			case itemArmor:
				olevel = dice().Intn(10) + 5
				break
			}
		}

		for _, mm := range room.Mobs {
			if mm.index.ID == m.ID {

				count := 0
				for _, i := range mm.Inventory {
					if i.index.ID == ir.ID {
						count++
					}
				}

				count = mo.Limit - count

				if count <= 0 {
					continue
				}

				for i := 0; i < count; i++ {
					item := createItem(ir)
					item.Level = olevel
					mm.Inventory = append(mm.Inventory, item)
				}
			}
		}
	}

	// reset objects equipped by mobs
	for _, mo := range reset.MobObjectsToEquip {
		ir := getItem(mo.ItemID)
		m := getMob(mo.MobID)
		room := getRoom(mo.RoomID)

		if ir == nil || m == nil || room == nil {

			dump("equipped")
			dump(mo)
			continue
		}

		var olevel int
		if m.Shop != nil {

			switch ir.ItemType {
			default:
				olevel = 0
				break
			case itemPill:
				olevel = dice().Intn(10)
				break
			case itemPotion:
				olevel = dice().Intn(10)
				break
			case itemScroll:
				olevel = dice().Intn(10) + 5
				break
			case itemWand:
				olevel = dice().Intn(10) + 10
				break
			case itemStaff:
				olevel = dice().Intn(10) + 15
				break
			case itemWeapon:
			case itemArmor:
				olevel = dice().Intn(10) + 5
				break
			}
		}

		item := createItem(ir)
		item.Level = olevel

		for _, mm := range room.Mobs {
			if mm.index.ID == m.ID {
				mm.Inventory = append(mm.Inventory, item)
				item.WearLocation = mo.WearLocation
				mm.wear(item, true)
				break
			}
		}
	}

	for _, d := range reset.Doors {
		room := getRoom(d.RoomID)
		if room == nil {
			continue
		}

		for _, x := range room.Exits {
			if x.Dir == d.Direction {
				switch d.State {
				case "open":
					x.Flags = removeBit(x.Flags, exitClosed)
					x.Flags = removeBit(x.Flags, exitLocked)
					break
				case "closed":
					x.Flags = setBit(x.Flags, exitClosed)
					x.Flags = removeBit(x.Flags, exitLocked)
					break
				case "locked":
					x.Flags = setBit(x.Flags, exitClosed)
					x.Flags = setBit(x.Flags, exitLocked)
					break
				}
			}
		}
	}
}

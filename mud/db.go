package mud

import (
	"container/list"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

var (
	areaList      list.List
	commandList   list.List
	itemList      list.List
	itemIndexList list.List
	jobList       list.List
	mobList       *list.List
	mobIndexList  list.List
	raceList      list.List
	roomList      list.List
	shopList      list.List
	skillList     list.List
	socialList    list.List
)

func areaUpdate() {
	for e := areaList.Front(); e != nil; e = e.Next() {
		area := e.Value.(*area)
		area.age++
		if area.age < 3 {
			continue
		}

		playerCount := 0
		for _, room := range area.Rooms {
			for _, mob := range room.Mobs {
				if mob.Playable && mob.client != nil {
					playerCount++
				}
			}
		}

		if playerCount == 0 || area.age >= 15 {
			resetArea(area)
			area.age = dice().Intn(3)
		}
	}
}

func bootDB() {
	loadCommands()
	loadSkills()
	loadSocials()
	loadJobs()
	loadRaces()
	loadItems()
	loadMobs()
	loadRooms()
	loadShops()

	areaUpdate()
}

func createItem(index *itemIndex) *item {
	var contained []*item
	item := &item{Name: index.Name, Description: index.Description, ShortDescription: index.ShortDescription, ItemType: index.ItemType, ExtraFlags: index.ExtraFlags, WearFlags: index.WearFlags, Weight: index.Weight, Value: index.Value, Timer: -1}
	for _, id := range index.ContainedIDs {
		i := getItem(id)
		contained = append(contained, createItem(i))
	}

	item.index = *index
	item.container = contained
	itemList.PushBack(item)
	return item
}

func createMob(index *mobIndex) *mob {

	m := &mob{}

	m.index = index
	m.Name = index.Name
	m.Title = index.Title
	m.Description = index.Description
	m.Affects = index.Affects
	m.AffectedBy = index.AffectedBy

	m.Skills = index.Skills

	for _, i := range index.ItemIds {
		m.Inventory = append(m.Inventory, createItem(getItem(i)))
	}
	for _, i := range index.EquippedIds {
		m.Equipped = append(m.Equipped, createItem(getItem(i)))
	}

	m.ExitVerb = index.ExitVerb
	m.Room = getRoom(index.RoomID)

	m.Hitpoints = index.Hitpoints
	m.MaxHitpoints = index.MaxHitpoints
	m.Mana = index.Mana
	m.MaxMana = index.MaxMana
	m.Movement = index.Movement
	m.MaxMovement = index.MaxMovement

	m.Armor = index.Armor
	m.Hitroll = index.Hitroll
	m.Damroll = index.Damroll

	m.Exp = index.Exp
	m.Level = index.Level
	m.Alignment = index.Alignment
	m.Practices = index.Practices
	m.Gold = index.Gold

	m.Carrying = index.Carrying
	m.CarryMax = index.CarryMax
	m.CarryWeight = index.CarryWeight
	m.CarryWeightMax = index.CarryWeightMax

	m.Job = getJob(uint(index.JobID))
	m.Race = getRace(uint(index.RaceID))
	m.Gender = index.Gender

	m.Attributes = index.Attributes
	m.ModifiedAttributes = index.ModifiedAttributes

	m.Status = index.Status
	m.Playable = index.Playable

	var skills []*mobSkill
	for _, s := range m.Skills {
		skill := getSkill(s.SkillID)
		skills = append(skills, &mobSkill{Skill: skill, SkillID: s.SkillID, Level: s.Level})
	}
	m.Skills = skills

	for _, i := range m.Equipped {
		m.wear(i, false)
	}

	if m.isNPC() && (m.Room != nil && m.Room.isDark()) {
		setBit(m.AffectedBy, affectInfrared)
	}
	return m
}

func loadCommands() {
	// Directions
	commandList.PushBack(&cmd{Name: "north", Trust: 0, Position: standing, Callback: func(player *mob, argument string) { doMove(player, "north") }})
	commandList.PushBack(&cmd{Name: "south", Trust: 0, Position: standing, Callback: func(player *mob, argument string) { doMove(player, "south") }})
	commandList.PushBack(&cmd{Name: "east", Trust: 0, Position: standing, Callback: func(player *mob, argument string) { doMove(player, "east") }})
	commandList.PushBack(&cmd{Name: "west", Trust: 0, Position: standing, Callback: func(player *mob, argument string) { doMove(player, "west") }})
	commandList.PushBack(&cmd{Name: "up", Trust: 0, Position: standing, Callback: func(player *mob, argument string) { doMove(player, "up") }})
	commandList.PushBack(&cmd{Name: "down", Trust: 0, Position: standing, Callback: func(player *mob, argument string) { doMove(player, "down") }})

	// Most common commands
	commandList.PushBack(&cmd{Name: "buy", Trust: 0, Position: standing, Callback: doBuy})
	commandList.PushBack(&cmd{Name: "cast", Trust: 0, Position: standing, Callback: doCast})
	commandList.PushBack(&cmd{Name: "scan", Trust: 0, Position: standing, Callback: doScan})
	commandList.PushBack(&cmd{Name: "get", Trust: 0, Position: standing, Callback: doGet})
	commandList.PushBack(&cmd{Name: "inventory", Trust: 0, Position: standing, Callback: doInventory})
	commandList.PushBack(&cmd{Name: "kill", Trust: 0, Position: standing, Callback: doKill})
	commandList.PushBack(&cmd{Name: "look", Trust: 0, Position: standing, Callback: doLook})

	// Informational commands
	commandList.PushBack(&cmd{Name: "consider", Trust: 0, Position: standing, Callback: doConsider})
	commandList.PushBack(&cmd{Name: "equipment", Trust: 0, Position: standing, Callback: doEquipment})
	// commandList.PushBack(&cmd{Name: "examine", Trust: 0, Position: standing, Callback: doExamine})
	// commandList.PushBack(&cmd{Name: "help", Trust: 0, Position: standing, Callback: doHelp})
	commandList.PushBack(&cmd{Name: "score", Trust: 0, Position: standing, Callback: doScore})
	commandList.PushBack(&cmd{Name: "socials", Trust: 0, Position: standing, Callback: doSocials})
	commandList.PushBack(&cmd{Name: "who", Trust: 0, Position: standing, Callback: doWho})

	// Communication commands
	commandList.PushBack(&cmd{Name: "answer", Trust: 0, Position: standing, Callback: doAnswer})
	commandList.PushBack(&cmd{Name: "auction", Trust: 0, Position: standing, Callback: doAuction})
	commandList.PushBack(&cmd{Name: "chat", Trust: 0, Position: standing, Callback: doChat})
	commandList.PushBack(&cmd{Name: "emote", Trust: 0, Position: standing, Callback: doEmote})
	commandList.PushBack(&cmd{Name: "music", Trust: 0, Position: standing, Callback: doMusic})
	commandList.PushBack(&cmd{Name: "question", Trust: 0, Position: standing, Callback: doQuestion})
	commandList.PushBack(&cmd{Name: "reply", Trust: 0, Position: standing, Callback: doReply})
	commandList.PushBack(&cmd{Name: "say", Trust: 0, Position: standing, Callback: doSay})
	commandList.PushBack(&cmd{Name: "tell", Trust: 0, Position: standing, Callback: doTell})
	// commandList.PushBack(&cmd{Name: "yell", Trust: 0, Position: standing, Callback: doYell})
	// commandList.PushBack(&cmd{Name: "shout", Trust: 0, Position: standing, Callback: doShout})
	// commandList.PushBack(&cmd{Name: "gtell", Trust: 0, Position: standing, Callback: doGroupTell})

	// Object manip commands
	// commandList.PushBack(&cmd{Name: "brandish", Trust: 0, Position: standing, Callback: doBrandish})
	// commandList.PushBack(&cmd{Name: "close", Trust: 0, Position: standing, Callback: doClose})
	// commandList.PushBack(&cmd{Name: "drink", Trust: 0, Position: standing, Callback: doDrink})
	commandList.PushBack(&cmd{Name: "drop", Trust: 0, Position: standing, Callback: doDrop})
	// commandList.PushBack(&cmd{Name: "eat", Trust: 0, Position: standing, Callback: doEat})
	// commandList.PushBack(&cmd{Name: "fill", Trust: 0, Position: standing, Callback: doFill})
	commandList.PushBack(&cmd{Name: "give", Trust: 0, Position: standing, Callback: doGive})
	// commandList.PushBack(&cmd{Name: "hold", Trust: 0, Position: standing, Callback: doHold})
	commandList.PushBack(&cmd{Name: "list", Trust: 0, Position: standing, Callback: doList})
	commandList.PushBack(&cmd{Name: "lock", Trust: 0, Position: standing, Callback: doLock})
	// commandList.PushBack(&cmd{Name: "open", Trust: 0, Position: standing, Callback: doOpen})
	commandList.PushBack(&cmd{Name: "pick", Trust: 0, Position: standing, Callback: doPick})
	commandList.PushBack(&cmd{Name: "put", Trust: 0, Position: standing, Callback: doPut})
	// commandList.PushBack(&cmd{Name: "quaff", Trust: 0, Position: standing, Callback: doQuaff})
	// commandList.PushBack(&cmd{Name: "recite", Trust: 0, Position: standing, Callback: doRecite})
	commandList.PushBack(&cmd{Name: "remove", Trust: 0, Position: standing, Callback: doRemove})
	commandList.PushBack(&cmd{Name: "sell", Trust: 0, Position: standing, Callback: doSell})
	// commandList.PushBack(&cmd{Name: "take", Trust: 0, Position: standing, Callback: doTake})
	commandList.PushBack(&cmd{Name: "sacrifice", Trust: 0, Position: standing, Callback: doSacrifice})
	commandList.PushBack(&cmd{Name: "unlock", Trust: 0, Position: standing, Callback: doUnlock})
	commandList.PushBack(&cmd{Name: "value", Trust: 0, Position: standing, Callback: doValue})
	commandList.PushBack(&cmd{Name: "wear", Trust: 0, Position: standing, Callback: doWear})
	// commandList.PushBack(&cmd{Name: "zap", Trust: 0, Position: standing, Callback: doZap})

	/* Combat Commands */
	commandList.PushBack(&cmd{Name: "backstab", Trust: 0, Position: standing, Callback: doBackstab})
	commandList.PushBack(&cmd{Name: "bs", Trust: 0, Position: standing, Callback: doBackstab})
	commandList.PushBack(&cmd{Name: "disarm", Trust: 0, Position: standing, Callback: doDisarm})
	commandList.PushBack(&cmd{Name: "flee", Trust: 0, Position: standing, Callback: doFlee})
	commandList.PushBack(&cmd{Name: "kick", Trust: 0, Position: standing, Callback: doKick})
	// commandList.PushBack(&cmd{Name: "murde", Trust: 0, Position: standing, Callback: doMurde})
	// commandList.PushBack(&cmd{Name: "murder", Trust: 0, Position: standing, Callback: doMurder})
	// commandList.PushBack(&cmd{Name: "rescue", Trust: 0, Position: standing, Callback: doRescue})

	/* Misc Commands */
	commandList.PushBack(&cmd{Name: "follow", Trust: 0, Position: standing, Callback: doFollow})
	// commandList.PushBack(&cmd{Name: "group", Trust: 0, Position: standing, Callback: doGroup})
	commandList.PushBack(&cmd{Name: "hide", Trust: 0, Position: standing, Callback: doHide})
	commandList.PushBack(&cmd{Name: "practice", Trust: 0, Position: standing, Callback: doPractice})
	commandList.PushBack(&cmd{Name: "qui", Trust: 0, Position: standing, Callback: doQui})
	commandList.PushBack(&cmd{Name: "quit", Trust: 0, Position: standing, Callback: doQuit})
	commandList.PushBack(&cmd{Name: "recall", Trust: 0, Position: standing, Callback: doRecall})
	// commandList.PushBack(&cmd{Name: "rent", Trust: 0, Position: standing, Callback: doRent})
	// commandList.PushBack(&cmd{Name: "save", Trust: 0, Position: standing, Callback: doSave})
	// commandList.PushBack(&cmd{Name: "sleep", Trust: 0, Position: standing, Callback: doSleep})
	commandList.PushBack(&cmd{Name: "sneak", Trust: 0, Position: standing, Callback: doSneak})
	// commandList.PushBack(&cmd{Name: "split", Trust: 0, Position: standing, Callback: doSplit})
	// commandList.PushBack(&cmd{Name: "steal", Trust: 0, Position: standing, Callback: doSteal})
	commandList.PushBack(&cmd{Name: "train", Trust: 0, Position: standing, Callback: doTrain})
	// commandList.PushBack(&cmd{Name: "visible", Trust: 0, Position: standing, Callback: doVisible})
	// commandList.PushBack(&cmd{Name: "wake", Trust: 0, Position: standing, Callback: doWake})
	commandList.PushBack(&cmd{Name: "where", Trust: 0, Position: standing, Callback: doWhere})

	/* Immortal commands */
	commandList.PushBack(&cmd{Name: "advance", Trust: 0, Position: standing, Callback: doAdvance})
	commandList.PushBack(&cmd{Name: "trust", Trust: 0, Position: standing, Callback: doTrust})

	// commandList.PushBack(&cmd{Name: "allow", Trust: 0, Position: standing, Callback: doAllow})
	// commandList.PushBack(&cmd{Name: "ban", Trust: 0, Position: standing, Callback: doBan})
	commandList.PushBack(&cmd{Name: "deny", Trust: 0, Position: standing, Callback: doDeny})
	commandList.PushBack(&cmd{Name: "disconnect", Trust: 0, Position: standing, Callback: doDisconnect})
	commandList.PushBack(&cmd{Name: "freeze", Trust: 0, Position: standing, Callback: doFreeze})
	commandList.PushBack(&cmd{Name: "reboo", Trust: 0, Position: standing, Callback: doReboo})
	commandList.PushBack(&cmd{Name: "reboot", Trust: 0, Position: standing, Callback: doReboot})
	commandList.PushBack(&cmd{Name: "shutdow", Trust: 0, Position: standing, Callback: doShutdow})
	commandList.PushBack(&cmd{Name: "shutdown", Trust: 0, Position: standing, Callback: doShutdown})
	// commandList.PushBack(&cmd{Name: "users", Trust: 0, Position: standing, Callback: doUsers})
	// commandList.PushBack(&cmd{Name: "wizlock", Trust: 0, Position: standing, Callback: doWizlock})

	commandList.PushBack(&cmd{Name: "force", Trust: 0, Position: standing, Callback: doForce})
	commandList.PushBack(&cmd{Name: "mload", Trust: 0, Position: standing, Callback: doMload})
	commandList.PushBack(&cmd{Name: "mset", Trust: 0, Position: standing, Callback: doMwhere})
	// commandList.PushBack(&cmd{Name: "noemote", Trust: 0, Position: standing, Callback: doNoEmote})
	// commandList.PushBack(&cmd{Name: "notell", Trust: 0, Position: standing, Callback: doNoTell})
	commandList.PushBack(&cmd{Name: "oload", Trust: 0, Position: standing, Callback: doOload})
	// commandList.PushBack(&cmd{Name: "oset", Trust: 0, Position: standing, Callback: doOset})
	commandList.PushBack(&cmd{Name: "pardon", Trust: 0, Position: standing, Callback: doPardon})
	commandList.PushBack(&cmd{Name: "purge", Trust: 0, Position: standing, Callback: doPurge})
	commandList.PushBack(&cmd{Name: "restore", Trust: 0, Position: standing, Callback: doRestore})
	// commandList.PushBack(&cmd{Name: "rset", Trust: 0, Position: standing, Callback: doRset})
	// commandList.PushBack(&cmd{Name: "silence", Trust: 0, Position: standing, Callback: doSilence})
	// commandList.PushBack(&cmd{Name: "sla", Trust: 0, Position: standing, Callback: doSla})
	// commandList.PushBack(&cmd{Name: "slay", Trust: 0, Position: standing, Callback: doSlay})
	// commandList.PushBack(&cmd{Name: "sset", Trust: 0, Position: standing, Callback: doSset})
	commandList.PushBack(&cmd{Name: "transfer", Trust: 0, Position: standing, Callback: doTransfer})

	commandList.PushBack(&cmd{Name: "at", Trust: 0, Position: standing, Callback: doAt})
	commandList.PushBack(&cmd{Name: "bamfin", Trust: 0, Position: standing, Callback: doBamfin})
	commandList.PushBack(&cmd{Name: "bamfout", Trust: 0, Position: standing, Callback: doBamfout})
	commandList.PushBack(&cmd{Name: "echo", Trust: 0, Position: standing, Callback: doEcho})
	commandList.PushBack(&cmd{Name: "goto", Trust: 0, Position: standing, Callback: doGoto})
	commandList.PushBack(&cmd{Name: "holylight", Trust: 0, Position: standing, Callback: doHolylight})
	commandList.PushBack(&cmd{Name: "invis", Trust: 0, Position: standing, Callback: doInvis})
	// commandList.PushBack(&cmd{Name: "log", Trust: 0, Position: standing, Callback: doLog})
	// commandList.PushBack(&cmd{Name: "memory", Trust: 0, Position: standing, Callback: doMemory})
	commandList.PushBack(&cmd{Name: "mfind", Trust: 0, Position: standing, Callback: doMfind})
	// commandList.PushBack(&cmd{Name: "mstat", Trust: 0, Position: standing, Callback: doMStat})
	commandList.PushBack(&cmd{Name: "mwhere", Trust: 0, Position: standing, Callback: doMwhere})
	commandList.PushBack(&cmd{Name: "ofind", Trust: 0, Position: standing, Callback: doOfind})
	// commandList.PushBack(&cmd{Name: "ostat", Trust: 0, Position: standing, Callback: doOsat})
	// commandList.PushBack(&cmd{Name: "owhere", Trust: 0, Position: standing, Callback: doOwhere})
	// commandList.PushBack(&cmd{Name: "peace", Trust: 0, Position: standing, Callback: doPeace})
	commandList.PushBack(&cmd{Name: "recho", Trust: 0, Position: standing, Callback: doRecho})
	commandList.PushBack(&cmd{Name: "return", Trust: 0, Position: standing, Callback: doReturn})
	commandList.PushBack(&cmd{Name: "rstat", Trust: 0, Position: standing, Callback: doRstat})
	// commandList.PushBack(&cmd{Name: "slookup", Trust: 0, Position: standing, Callback: doSlookup})
	commandList.PushBack(&cmd{Name: "snoop", Trust: 0, Position: standing, Callback: doSnoop})
	commandList.PushBack(&cmd{Name: "switch", Trust: 0, Position: standing, Callback: doSwitch})

	commandList.PushBack(&cmd{Name: "immtalk", Trust: 0, Position: standing, Callback: doImmtalk})
}

func loadItems() {
	itemFiles, _ := filepath.Glob("./data/items/*.json")

	for _, itemFile := range itemFiles {
		file, err := ioutil.ReadFile(itemFile)
		if err != nil {
			panic(err)
		}

		var list []itemIndex
		json.Unmarshal(file, &list)

		for _, it := range list {
			itemIndexList.PushBack(it)
		}

	}
}

func loadJobs() {
	var jobs []*job

	jobs = append(jobs, &job{ID: 1, Name: "Warrior", Abbr: "war", StartingWeapon: startingSword, PrimeAttribute: applyStrength, SkillAdept: 85, Thac0_00: 18, Thac0_32: 6, MinHitpoints: 11, MaxHitpoints: 15, GainsMana: false})
	jobs = append(jobs, &job{ID: 2, Name: "Mage", Abbr: "mag", StartingWeapon: startingStaff, PrimeAttribute: applyIntelligence, SkillAdept: 95, Thac0_00: 18, Thac0_32: 10, MinHitpoints: 6, MaxHitpoints: 8, GainsMana: true})
	jobs = append(jobs, &job{ID: 3, Name: "Cleric", Abbr: "cle", StartingWeapon: startingMace, PrimeAttribute: applyWisdom, SkillAdept: 95, Thac0_00: 18, Thac0_32: 12, MinHitpoints: 7, MaxHitpoints: 10, GainsMana: true})
	jobs = append(jobs, &job{ID: 4, Name: "Thief", Abbr: "thi", StartingWeapon: startingDagger, PrimeAttribute: applyDexterity, SkillAdept: 85, Thac0_00: 18, Thac0_32: 8, MinHitpoints: 8, MaxHitpoints: 8, GainsMana: false})
	jobs = append(jobs, &job{ID: 5, Name: "Ranger", Abbr: "ran", StartingWeapon: startingWhip, PrimeAttribute: applyConstitution, SkillAdept: 85, Thac0_00: 18, Thac0_32: 8, MinHitpoints: 10, MaxHitpoints: 14, GainsMana: false})
	jobs = append(jobs, &job{ID: 6, Name: "Bard", Abbr: "bar", StartingWeapon: startingDagger, PrimeAttribute: applyCharisma, SkillAdept: 85, Thac0_00: 18, Thac0_32: 11, MinHitpoints: 7, MaxHitpoints: 9, GainsMana: true})

	for _, j := range jobs {
		jobList.PushBack(j)
	}
}

func loadMobs() {
	mobList = list.New()

	mobFiles, _ := filepath.Glob("./data/mobs/*.json")

	for _, mobFile := range mobFiles {
		file, err := ioutil.ReadFile(mobFile)
		if err != nil {
			panic(err)
		}

		var list []*mobIndex
		err = json.Unmarshal(file, &list)
		if err != nil {
			panic(err)
		}

		for _, m := range list {
			mobIndexList.PushBack(m)
		}
	}

}

func loadRaces() {
	raceList.PushBack(&race{ID: 1, Name: "Human", Abbr: "hum"})
	raceList.PushBack(&race{ID: 2, Name: "Elf", Abbr: "elf"})
	raceList.PushBack(&race{ID: 3, Name: "Dwarf", Abbr: "dwf"})
	raceList.PushBack(&race{ID: 4, Name: "Dark Elft", Abbr: "drw"})
	raceList.PushBack(&race{ID: 5, Name: "Goblin", Abbr: "gob"})
	raceList.PushBack(&race{ID: 6, Name: "Dragon", Abbr: "drg"})
}

func loadRooms() {
	areaFiles, _ := filepath.Glob("./data/area/*.json")

	voidArea := &area{ID: 0, Name: "Limbo", age: 0, numPlayers: 0}
	void := &room{ID: 0, Exits: nil, ItemIds: nil, MobIds: nil, Name: "The Void", Description: "A dark, gaping void lies here."}

	voidArea.Rooms = append(voidArea.Rooms, void)
	areaList.PushBack(voidArea)

	for _, areaFile := range areaFiles {
		file, err := ioutil.ReadFile(areaFile)
		if err != nil {
			panic(err)
		}

		var a area
		json.Unmarshal(file, &a)
		if err != nil {
			panic(err)
		}

		area := &area{ID: a.ID, Name: a.Name, age: 0}

		for _, ro := range a.Rooms {
			ro.AreaID = int(a.ID)
			for _, i := range ro.ItemIds {
				index := getItem(i)
				item := createItem(index)
				ro.Items = append(ro.Items, item)
			}

			for _, i := range ro.MobIds {
				mob := createMob(getMob(i))
				ro.Mobs = append(ro.Mobs, mob)
			}

			roomList.PushBack(ro)

			area.Rooms = append(area.Rooms, ro)
		}

		exitsList := list.New()
		for e := roomList.Front(); e != nil; e = e.Next() {
			room := e.Value.(*room)
			for j, x := range room.Exits {
				room.Exits[j] = &exit{Dir: x.Dir, Room: getRoom(x.RoomID), RoomID: x.RoomID}
			}

			exitsList.PushBack(room)
		}

		roomList = *exitsList

		areaList.PushBack(area)
	}

	for e := roomList.Front(); e != nil; e = e.Next() {
		room := e.Value.(*room)
		for _, mob := range room.Mobs {
			mob.Room = room
		}
	}
}

func loadShops() {
	shopFiles, _ := filepath.Glob("./data/shops/*.json")

	for _, shopFile := range shopFiles {
		file, err := ioutil.ReadFile(shopFile)
		if err != nil {
			panic(err)
		}

		var list []*shop
		json.Unmarshal(file, &list)

		for _, sh := range list {
			sh.keeper = getMob(sh.KeeperID)
			shopList.PushBack(sh)
		}
	}
}

func loadSkills() {
	skillFiles, _ := filepath.Glob("./data/skills/*.json")

	for _, skillFile := range skillFiles {
		file, err := ioutil.ReadFile(skillFile)
		if err != nil {
			panic(err)
		}

		var list []*skill
		json.Unmarshal(file, &list)

		for _, sk := range list {
			skillList.PushBack(sk)
		}
	}
}

func loadSocials() {
	file, err := ioutil.ReadFile("./data/socials/socials.json")
	if err != nil {
		panic(err)
	}

	var list []*social
	err = json.Unmarshal(file, &list)
	if err != nil {
		panic(err)
	}

	for _, sk := range list {
		socialList.PushBack(sk)
	}
}

func resetArea(ar *area) {
	filename := toSnake(ar.Name)

	areaFile := fmt.Sprintf("./data/areas/%s.json", filename)
	file, err := ioutil.ReadFile(areaFile)
	if err != nil {
		panic(err)
	}

	var a *area
	json.Unmarshal(file, &a)
	if err != nil {
		panic(err)
	}

	var masterArea *area
	for e := areaList.Front(); e != nil; e = e.Next() {
		area := e.Value.(*area)
		if area.ID == a.ID {
			masterArea = area
			break
		}
	}

	var tempRooms []*room
	for _, ro := range a.Rooms {
		ro.AreaID = int(a.ID)
		for _, i := range ro.ItemIds {
			index := getItem(i)
			item := createItem(index)
			ro.Items = append(ro.Items, item)
		}

		for _, i := range ro.MobIds {
			mob := createMob(getMob(i))
			ro.Mobs = append(ro.Mobs, mob)
		}

		for e := roomList.Front(); e != nil; e = e.Next() {
			room := e.Value.(*room)
			if room.ID == ro.ID {
				room = ro
			}
		}

		tempRooms = append(tempRooms, ro)
	}

	masterArea.Rooms = tempRooms
	return
}

func extractMob(m *mob, pull bool) {
	if m.Room == nil {
		return
	}

	if pull {
		// TODO: m.dieFollower()
	}

	m.stopFighting(true)

	for _, i := range m.Inventory {
		extractObj(i)
	}

	m.Room.removeMob(m)

	if !pull {
		m.Room = getRoom(1)
		return
	}

	if m.isNPC() {
		m.index.count--
	}

	if m.client != nil && m.client.original != nil {
		doReturn(m, "")
	}

	for e := mobList.Front(); e != nil; e = e.Next() {
		wm := e.Value.(*mob)
		if wm.replyTarget == m {
			wm.replyTarget = nil
		}
	}

	for e := mobList.Front(); e != nil; e = e.Next() {
		wm := e.Value.(*mob)
		if wm == m {
			mobList.Remove(e)
		}
	}

	if m.client != nil {
		m.client.mob = nil
	}

	return
}

func extractObj(obj *item) {
	if obj.Room != nil {
		obj.Room.removeObject(obj)
	}
	if obj.carriedBy != nil {
		obj.carriedBy.removeItem(obj)
	}
	if obj.inObject != nil {
		obj.inObject.removeObject(obj)
	}

	for _, i := range obj.container {
		extractObj(i)
	}

	for e := itemList.Front(); e != nil; e = e.Next() {
		i := e.Value.(*item)
		if i == obj {
			itemList.Remove(e)
			break
		}
	}

	obj.Affected = nil
	obj.index.count--

	return
}

package mud

import (
	"container/list"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"path/filepath"
	"strings"
)

var (
	bonusTableStrength     = make(map[int]bonusStrength)
	bonusTableIntelligence = make(map[int]bonusIntelligence)
	bonusTableWisdom       = make(map[int]bonusWisdom)
	bonusTableDexterity    = make(map[int]bonusDexterity)
	bonusTableConstitution = make(map[int]bonusConstitution)
)

var (
	areaList      list.List
	banList       list.List
	commandList   list.List
	helpList      list.List
	itemList      list.List
	itemIndexList list.List
	jobList       list.List
	mobList       *list.List
	mobIndexList  list.List
	raceList      list.List
	resetList     list.List
	roomList      list.List
	shopList      list.List
	skillList     list.List
	socialList    list.List
)

func areaUpdate(now bool) {
	for e := areaList.Front(); e != nil; e = e.Next() {
		area := e.Value.(*area)
		area.age++
		if area.age < 3 && !now {
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

		if now || area.age >= 15 {
			resetArea(area)
			area.age = dice().Intn(3)
		}
	}
}

func bootDB() {
	loadBonuses()
	loadHelps()
	loadCommands()
	loadSkills()
	loadSocials()
	loadJobs()
	loadRaces()
	loadItems()
	loadMobs()
	loadRooms()
	loadShops()
	loadResets()

	areaUpdate(true)
}

func createItem(index *itemIndex) *item {
	var contained []*item
	item := &item{Name: index.Name, Description: index.Description, ShortDescription: index.ShortDescription, ItemType: index.ItemType, ExtraFlags: index.ExtraFlags, WearFlags: index.WearFlags, Weight: index.Weight, Value: index.Value, Timer: -1, Cost: index.Cost}
	for _, id := range index.ContainedIDs {
		i := getItem(id)
		contained = append(contained, createItem(i))
	}

	item.index = *index
	item.index.count++
	item.container = contained
	itemList.PushBack(item)

	return item
}

func createMob(index *mobIndex) *mob {

	m := &mob{}

	m.ID = index.ID
	m.SavedAt = index.SavedAt
	m.CreatedAt = index.CreatedAt
	m.LastSeenAt = index.LastSeenAt

	m.index = index
	m.Name = index.Name
	m.Password = index.Password
	m.Title = index.Title
	m.Description = index.Description
	m.LongDescription = index.LongDescription
	m.Affects = index.Affects

	var act int
	for _, a := range index.Act {
		act |= a
	}
	m.Act = act

	var affects int
	for _, a := range index.AffectedBy {
		affects |= a
	}
	m.AffectedBy = affects

	m.Skills = index.Skills

	m.Room = getRoom(index.RoomID)

	if m.Room == nil && m.WasInRoom != nil {
		m.Room = m.WasInRoom
	}

	if m.Room == nil {
		m.Room = getRoom(index.RecallRoomID)
	}

	if m.Room == nil {
		m.Room = getRoom(1) // start back at mudschool
	}

	for _, i := range index.Inventory {
		index := getItem(i.ID)
		i.index = *index
		if i.WearFlags == itemWearLight {
			m.Room.Light++
		}
		m.Inventory = append(m.Inventory, i)
	}

	for _, i := range index.Equipped {
		index := getItem(i.ID)
		i.index = *index
		if i.WearFlags == itemWearLight {
			m.Room.Light++
		}
		m.Inventory = append(m.Inventory, i)
		m.wear(i, false)
		fmt.Println("Wearing...")
	}

	m.ExitVerb = index.ExitVerb

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

	m.Job = getJob(index.JobID)
	m.Race = getRace(index.RaceID)
	m.Gender = index.Gender

	m.Attributes = index.Attributes
	m.ModifiedAttributes = &attributeSet{}

	m.Status = index.Status
	m.Playable = index.Playable

	var skills []*mobSkill
	for _, s := range m.Skills {
		skill := getSkill(s.SkillID)
		skills = append(skills, &mobSkill{Skill: skill, SkillID: s.SkillID, Level: s.Level})
	}
	m.Skills = skills

	if m.isNPC() && (m.Room != nil && m.Room.isDark()) {
		m.AffectedBy = setBit(m.AffectedBy, affectInfrared)
	}

	mobList.PushBack(m)

	return m
}

func loadBonuses() {
	bonusTableStrength[0] = bonusStrength{toHit: -5, toDamage: -4, toCarry: 0, wield: 0}
	bonusTableStrength[1] = bonusStrength{toHit: -4, toDamage: -4, toCarry: 0, wield: 0}
	bonusTableStrength[2] = bonusStrength{toHit: -3, toDamage: -4, toCarry: 0, wield: 0}
	bonusTableStrength[3] = bonusStrength{toHit: -2, toDamage: -4, toCarry: 0, wield: 0}
	bonusTableStrength[4] = bonusStrength{toHit: -2, toDamage: -4, toCarry: 0, wield: 0}
	bonusTableStrength[5] = bonusStrength{toHit: -1, toDamage: -4, toCarry: 0, wield: 0}
	bonusTableStrength[6] = bonusStrength{toHit: -1, toDamage: -4, toCarry: 0, wield: 0}
	bonusTableStrength[7] = bonusStrength{toHit: 0, toDamage: -4, toCarry: 0, wield: 0}
	bonusTableStrength[8] = bonusStrength{toHit: 0, toDamage: -4, toCarry: 0, wield: 0}
	bonusTableStrength[9] = bonusStrength{toHit: 0, toDamage: -4, toCarry: 0, wield: 0}
	bonusTableStrength[10] = bonusStrength{toHit: 0, toDamage: -4, toCarry: 0, wield: 0}
	bonusTableStrength[11] = bonusStrength{toHit: 0, toDamage: -4, toCarry: 0, wield: 0}
	bonusTableStrength[12] = bonusStrength{toHit: 0, toDamage: -4, toCarry: 0, wield: 0}
	bonusTableStrength[13] = bonusStrength{toHit: 0, toDamage: -4, toCarry: 0, wield: 0}
	bonusTableStrength[14] = bonusStrength{toHit: 1, toDamage: -4, toCarry: 0, wield: 0}
	bonusTableStrength[15] = bonusStrength{toHit: 1, toDamage: -4, toCarry: 0, wield: 0}
	bonusTableStrength[16] = bonusStrength{toHit: 2, toDamage: -4, toCarry: 0, wield: 0}
	bonusTableStrength[17] = bonusStrength{toHit: 2, toDamage: -4, toCarry: 0, wield: 0}
	bonusTableStrength[18] = bonusStrength{toHit: 3, toDamage: -4, toCarry: 0, wield: 0}
	bonusTableStrength[19] = bonusStrength{toHit: 3, toDamage: -4, toCarry: 0, wield: 0}
	bonusTableStrength[20] = bonusStrength{toHit: 4, toDamage: -4, toCarry: 0, wield: 0}
	bonusTableStrength[21] = bonusStrength{toHit: 5, toDamage: -4, toCarry: 0, wield: 0}
	bonusTableStrength[22] = bonusStrength{toHit: 5, toDamage: -4, toCarry: 0, wield: 0}
	bonusTableStrength[23] = bonusStrength{toHit: 6, toDamage: -4, toCarry: 0, wield: 0}
	bonusTableStrength[24] = bonusStrength{toHit: 8, toDamage: -4, toCarry: 0, wield: 0}
	bonusTableStrength[25] = bonusStrength{toHit: 10, toDamage: -4, toCarry: 0, wield: 0}

	bonusTableIntelligence[0] = bonusIntelligence{learn: 3}
	bonusTableIntelligence[1] = bonusIntelligence{learn: 5}
	bonusTableIntelligence[2] = bonusIntelligence{learn: 7}
	bonusTableIntelligence[3] = bonusIntelligence{learn: 8}
	bonusTableIntelligence[4] = bonusIntelligence{learn: 9}
	bonusTableIntelligence[5] = bonusIntelligence{learn: 10}
	bonusTableIntelligence[6] = bonusIntelligence{learn: 11}
	bonusTableIntelligence[7] = bonusIntelligence{learn: 12}
	bonusTableIntelligence[8] = bonusIntelligence{learn: 13}
	bonusTableIntelligence[9] = bonusIntelligence{learn: 15}
	bonusTableIntelligence[10] = bonusIntelligence{learn: 17}
	bonusTableIntelligence[11] = bonusIntelligence{learn: 19}
	bonusTableIntelligence[12] = bonusIntelligence{learn: 22}
	bonusTableIntelligence[13] = bonusIntelligence{learn: 25}
	bonusTableIntelligence[14] = bonusIntelligence{learn: 28}
	bonusTableIntelligence[15] = bonusIntelligence{learn: 31}
	bonusTableIntelligence[16] = bonusIntelligence{learn: 34}
	bonusTableIntelligence[17] = bonusIntelligence{learn: 37}
	bonusTableIntelligence[18] = bonusIntelligence{learn: 40}
	bonusTableIntelligence[19] = bonusIntelligence{learn: 44}
	bonusTableIntelligence[20] = bonusIntelligence{learn: 49}
	bonusTableIntelligence[21] = bonusIntelligence{learn: 55}
	bonusTableIntelligence[22] = bonusIntelligence{learn: 60}
	bonusTableIntelligence[23] = bonusIntelligence{learn: 70}
	bonusTableIntelligence[24] = bonusIntelligence{learn: 85}
	bonusTableIntelligence[25] = bonusIntelligence{learn: 99}

	bonusTableWisdom[0] = bonusWisdom{practice: 0}
	bonusTableWisdom[1] = bonusWisdom{practice: 0}
	bonusTableWisdom[2] = bonusWisdom{practice: 0}
	bonusTableWisdom[3] = bonusWisdom{practice: 0}
	bonusTableWisdom[4] = bonusWisdom{practice: 0}
	bonusTableWisdom[5] = bonusWisdom{practice: 1}
	bonusTableWisdom[6] = bonusWisdom{practice: 1}
	bonusTableWisdom[7] = bonusWisdom{practice: 1}
	bonusTableWisdom[8] = bonusWisdom{practice: 1}
	bonusTableWisdom[9] = bonusWisdom{practice: 1}
	bonusTableWisdom[10] = bonusWisdom{practice: 2}
	bonusTableWisdom[11] = bonusWisdom{practice: 2}
	bonusTableWisdom[12] = bonusWisdom{practice: 2}
	bonusTableWisdom[13] = bonusWisdom{practice: 2}
	bonusTableWisdom[14] = bonusWisdom{practice: 2}
	bonusTableWisdom[15] = bonusWisdom{practice: 3}
	bonusTableWisdom[16] = bonusWisdom{practice: 3}
	bonusTableWisdom[17] = bonusWisdom{practice: 4}
	bonusTableWisdom[18] = bonusWisdom{practice: 5}
	bonusTableWisdom[19] = bonusWisdom{practice: 5}
	bonusTableWisdom[20] = bonusWisdom{practice: 5}
	bonusTableWisdom[21] = bonusWisdom{practice: 6}
	bonusTableWisdom[22] = bonusWisdom{practice: 6}
	bonusTableWisdom[23] = bonusWisdom{practice: 6}
	bonusTableWisdom[24] = bonusWisdom{practice: 6}
	bonusTableWisdom[25] = bonusWisdom{practice: 7}

	bonusTableDexterity[0] = bonusDexterity{defensive: 60}
	bonusTableDexterity[1] = bonusDexterity{defensive: 50}
	bonusTableDexterity[2] = bonusDexterity{defensive: 50}
	bonusTableDexterity[3] = bonusDexterity{defensive: 40}
	bonusTableDexterity[4] = bonusDexterity{defensive: 30}
	bonusTableDexterity[5] = bonusDexterity{defensive: 20}
	bonusTableDexterity[6] = bonusDexterity{defensive: 10}
	bonusTableDexterity[7] = bonusDexterity{defensive: 0}
	bonusTableDexterity[8] = bonusDexterity{defensive: 0}
	bonusTableDexterity[9] = bonusDexterity{defensive: 0}
	bonusTableDexterity[10] = bonusDexterity{defensive: 0}
	bonusTableDexterity[11] = bonusDexterity{defensive: 0}
	bonusTableDexterity[12] = bonusDexterity{defensive: 0}
	bonusTableDexterity[13] = bonusDexterity{defensive: 0}
	bonusTableDexterity[14] = bonusDexterity{defensive: 0}
	bonusTableDexterity[15] = bonusDexterity{defensive: -10}
	bonusTableDexterity[16] = bonusDexterity{defensive: -15}
	bonusTableDexterity[17] = bonusDexterity{defensive: -20}
	bonusTableDexterity[18] = bonusDexterity{defensive: -30}
	bonusTableDexterity[19] = bonusDexterity{defensive: -40}
	bonusTableDexterity[20] = bonusDexterity{defensive: -50}
	bonusTableDexterity[21] = bonusDexterity{defensive: -60}
	bonusTableDexterity[22] = bonusDexterity{defensive: -75}
	bonusTableDexterity[23] = bonusDexterity{defensive: -90}
	bonusTableDexterity[24] = bonusDexterity{defensive: -105}
	bonusTableDexterity[25] = bonusDexterity{defensive: -120}

	bonusTableConstitution[0] = bonusConstitution{hitpoints: -4, shock: 20}
	bonusTableConstitution[1] = bonusConstitution{hitpoints: -3, shock: 25}
	bonusTableConstitution[2] = bonusConstitution{hitpoints: -2, shock: 30}
	bonusTableConstitution[3] = bonusConstitution{hitpoints: -1, shock: 35}
	bonusTableConstitution[4] = bonusConstitution{hitpoints: -1, shock: 40}
	bonusTableConstitution[5] = bonusConstitution{hitpoints: -1, shock: 45}
	bonusTableConstitution[6] = bonusConstitution{hitpoints: 0, shock: 50}
	bonusTableConstitution[7] = bonusConstitution{hitpoints: 0, shock: 55}
	bonusTableConstitution[8] = bonusConstitution{hitpoints: 0, shock: 60}
	bonusTableConstitution[9] = bonusConstitution{hitpoints: 0, shock: 65}
	bonusTableConstitution[10] = bonusConstitution{hitpoints: 0, shock: 70}
	bonusTableConstitution[11] = bonusConstitution{hitpoints: 0, shock: 75}
	bonusTableConstitution[12] = bonusConstitution{hitpoints: 0, shock: 80}
	bonusTableConstitution[13] = bonusConstitution{hitpoints: 0, shock: 85}
	bonusTableConstitution[14] = bonusConstitution{hitpoints: 0, shock: 90}
	bonusTableConstitution[15] = bonusConstitution{hitpoints: 1, shock: 95}
	bonusTableConstitution[16] = bonusConstitution{hitpoints: 2, shock: 99}
	bonusTableConstitution[17] = bonusConstitution{hitpoints: 2, shock: 99}
	bonusTableConstitution[18] = bonusConstitution{hitpoints: 3, shock: 99}
	bonusTableConstitution[19] = bonusConstitution{hitpoints: 3, shock: 99}
	bonusTableConstitution[20] = bonusConstitution{hitpoints: 4, shock: 99}
	bonusTableConstitution[21] = bonusConstitution{hitpoints: 4, shock: 99}
	bonusTableConstitution[22] = bonusConstitution{hitpoints: 5, shock: 99}
	bonusTableConstitution[23] = bonusConstitution{hitpoints: 6, shock: 99}
	bonusTableConstitution[24] = bonusConstitution{hitpoints: 7, shock: 29}
	bonusTableConstitution[25] = bonusConstitution{hitpoints: 8, shock: 99}
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
	commandList.PushBack(&cmd{Name: "areas", Trust: 0, Position: dead, Callback: doAreas})
	commandList.PushBack(&cmd{Name: "commands", Trust: 0, Position: dead, Callback: doCommands})
	commandList.PushBack(&cmd{Name: "compare", Trust: 0, Position: resting, Callback: doCompare})
	commandList.PushBack(&cmd{Name: "consider", Trust: 0, Position: resting, Callback: doConsider})
	commandList.PushBack(&cmd{Name: "equipment", Trust: 0, Position: dead, Callback: doEquipment})
	commandList.PushBack(&cmd{Name: "examine", Trust: 0, Position: resting, Callback: doExamine})
	commandList.PushBack(&cmd{Name: "help", Trust: 0, Position: dead, Callback: doHelp})
	commandList.PushBack(&cmd{Name: "score", Trust: 0, Position: dead, Callback: doScore})
	commandList.PushBack(&cmd{Name: "socials", Trust: 0, Position: dead, Callback: doSocials})
	commandList.PushBack(&cmd{Name: "who", Trust: 0, Position: dead, Callback: doWho})

	// Communication commands
	commandList.PushBack(&cmd{Name: "answer", Trust: 0, Position: sleeping, Callback: doAnswer})
	commandList.PushBack(&cmd{Name: "auction", Trust: 0, Position: sleeping, Callback: doAuction})
	commandList.PushBack(&cmd{Name: "chat", Trust: 0, Position: sleeping, Callback: doChat})
	commandList.PushBack(&cmd{Name: "emote", Trust: 0, Position: resting, Callback: doEmote})
	commandList.PushBack(&cmd{Name: "music", Trust: 0, Position: sleeping, Callback: doMusic})
	commandList.PushBack(&cmd{Name: "question", Trust: 0, Position: sleeping, Callback: doQuestion})
	commandList.PushBack(&cmd{Name: "reply", Trust: 0, Position: resting, Callback: doReply})
	commandList.PushBack(&cmd{Name: "say", Trust: 0, Position: resting, Callback: doSay})
	commandList.PushBack(&cmd{Name: "tell", Trust: 0, Position: resting, Callback: doTell})
	commandList.PushBack(&cmd{Name: "yell", Trust: 0, Position: resting, Callback: doYell})
	commandList.PushBack(&cmd{Name: "shout", Trust: 0, Position: resting, Callback: doShout})
	commandList.PushBack(&cmd{Name: "gtell", Trust: 0, Position: standing, Callback: doGroupTell})

	// Object manip commands
	commandList.PushBack(&cmd{Name: "brandish", Trust: 0, Position: resting, Callback: doBrandish})
	commandList.PushBack(&cmd{Name: "close", Trust: 0, Position: resting, Callback: doClose})
	// commandList.PushBack(&cmd{Name: "drink", Trust: 0, Position: resting, Callback: doDrink})
	commandList.PushBack(&cmd{Name: "drop", Trust: 0, Position: resting, Callback: doDrop})
	commandList.PushBack(&cmd{Name: "eat", Trust: 0, Position: resting, Callback: doEat})
	// commandList.PushBack(&cmd{Name: "fill", Trust: 0, Position: resting, Callback: doFill})
	commandList.PushBack(&cmd{Name: "give", Trust: 0, Position: resting, Callback: doGive})
	commandList.PushBack(&cmd{Name: "hold", Trust: 0, Position: resting, Callback: doWear})
	commandList.PushBack(&cmd{Name: "list", Trust: 0, Position: resting, Callback: doList})
	commandList.PushBack(&cmd{Name: "lock", Trust: 0, Position: resting, Callback: doLock})
	commandList.PushBack(&cmd{Name: "open", Trust: 0, Position: resting, Callback: doOpen})
	commandList.PushBack(&cmd{Name: "pick", Trust: 0, Position: resting, Callback: doPick})
	commandList.PushBack(&cmd{Name: "put", Trust: 0, Position: resting, Callback: doPut})
	commandList.PushBack(&cmd{Name: "quaff", Trust: 0, Position: resting, Callback: doQuaff})
	commandList.PushBack(&cmd{Name: "recite", Trust: 0, Position: resting, Callback: doRecite})
	commandList.PushBack(&cmd{Name: "remove", Trust: 0, Position: resting, Callback: doRemove})
	commandList.PushBack(&cmd{Name: "sell", Trust: 0, Position: resting, Callback: doSell})
	commandList.PushBack(&cmd{Name: "take", Trust: 0, Position: resting, Callback: doGet})
	commandList.PushBack(&cmd{Name: "sacrifice", Trust: 0, Position: resting, Callback: doSacrifice})
	commandList.PushBack(&cmd{Name: "unlock", Trust: 0, Position: resting, Callback: doUnlock})
	commandList.PushBack(&cmd{Name: "value", Trust: 0, Position: resting, Callback: doValue})
	commandList.PushBack(&cmd{Name: "wear", Trust: 0, Position: resting, Callback: doWear})
	commandList.PushBack(&cmd{Name: "zap", Trust: 0, Position: resting, Callback: doZap})

	/* Combat Commands */
	commandList.PushBack(&cmd{Name: "backstab", Trust: 0, Position: standing, Callback: doBackstab})
	commandList.PushBack(&cmd{Name: "bs", Trust: 0, Position: standing, Callback: doBackstab})
	commandList.PushBack(&cmd{Name: "disarm", Trust: 0, Position: fighting, Callback: doDisarm})
	commandList.PushBack(&cmd{Name: "flee", Trust: 0, Position: fighting, Callback: doFlee})
	commandList.PushBack(&cmd{Name: "kick", Trust: 0, Position: fighting, Callback: doKick})
	commandList.PushBack(&cmd{Name: "rescue", Trust: 0, Position: fighting, Callback: doRescue})

	/* Misc Commands */
	commandList.PushBack(&cmd{Name: "follow", Trust: 0, Position: resting, Callback: doFollow})
	commandList.PushBack(&cmd{Name: "group", Trust: 0, Position: sleeping, Callback: doGroup})
	commandList.PushBack(&cmd{Name: "hide", Trust: 0, Position: resting, Callback: doHide})
	commandList.PushBack(&cmd{Name: "practice", Trust: 0, Position: sleeping, Callback: doPractice})
	commandList.PushBack(&cmd{Name: "qui", Trust: 0, Position: dead, Callback: doQui})
	commandList.PushBack(&cmd{Name: "quit", Trust: 0, Position: dead, Callback: doQuit})
	commandList.PushBack(&cmd{Name: "recall", Trust: 0, Position: fighting, Callback: doRecall})
	// commandList.PushBack(&cmd{Name: "rent", Trust: 0, Position: dead, Callback: doRent})
	commandList.PushBack(&cmd{Name: "rest", Trust: 0, Position: sleeping, Callback: doRest})
	commandList.PushBack(&cmd{Name: "save", Trust: 0, Position: dead, Callback: doSave})
	commandList.PushBack(&cmd{Name: "sleep", Trust: 0, Position: sleeping, Callback: doSleep})
	commandList.PushBack(&cmd{Name: "sneak", Trust: 0, Position: standing, Callback: doSneak})
	// commandList.PushBack(&cmd{Name: "split", Trust: 0, Position: resting, Callback: doSplit})
	commandList.PushBack(&cmd{Name: "steal", Trust: 0, Position: standing, Callback: doSteal})
	commandList.PushBack(&cmd{Name: "train", Trust: 0, Position: resting, Callback: doTrain})
	commandList.PushBack(&cmd{Name: "visible", Trust: 0, Position: sleeping, Callback: doVisible})
	commandList.PushBack(&cmd{Name: "wake", Trust: 0, Position: sleeping, Callback: doWake})
	commandList.PushBack(&cmd{Name: "where", Trust: 0, Position: resting, Callback: doWhere})

	/* Immortal commands */
	commandList.PushBack(&cmd{Name: "advance", Trust: 98, Position: dead, Callback: doAdvance})
	commandList.PushBack(&cmd{Name: "trust", Trust: 98, Position: dead, Callback: doTrust})

	commandList.PushBack(&cmd{Name: "allow", Trust: 97, Position: dead, Callback: doAllow})
	commandList.PushBack(&cmd{Name: "ban", Trust: 97, Position: dead, Callback: doBan})
	commandList.PushBack(&cmd{Name: "deny", Trust: 97, Position: dead, Callback: doDeny})
	commandList.PushBack(&cmd{Name: "disconnect", Trust: 97, Position: dead, Callback: doDisconnect})
	commandList.PushBack(&cmd{Name: "freeze", Trust: 97, Position: dead, Callback: doFreeze})
	commandList.PushBack(&cmd{Name: "reboo", Trust: 97, Position: dead, Callback: doReboo})
	commandList.PushBack(&cmd{Name: "reboot", Trust: 97, Position: dead, Callback: doReboot})
	commandList.PushBack(&cmd{Name: "shutdow", Trust: 97, Position: dead, Callback: doShutdow})
	commandList.PushBack(&cmd{Name: "shutdown", Trust: 97, Position: dead, Callback: doShutdown})
	commandList.PushBack(&cmd{Name: "users", Trust: 97, Position: dead, Callback: doUsers})
	commandList.PushBack(&cmd{Name: "wizlock", Trust: 97, Position: dead, Callback: doWizlock})

	commandList.PushBack(&cmd{Name: "force", Trust: 96, Position: dead, Callback: doForce})
	commandList.PushBack(&cmd{Name: "mload", Trust: 96, Position: dead, Callback: doMload})
	commandList.PushBack(&cmd{Name: "mset", Trust: 96, Position: dead, Callback: doMwhere})
	commandList.PushBack(&cmd{Name: "noemote", Trust: 96, Position: dead, Callback: doNoEmote})
	commandList.PushBack(&cmd{Name: "notell", Trust: 96, Position: dead, Callback: doNoTell})
	commandList.PushBack(&cmd{Name: "oload", Trust: 96, Position: dead, Callback: doOload})
	// commandList.PushBack(&cmd{Name: "oset", Trust: 96, Position: dead, Callback: doOset})
	commandList.PushBack(&cmd{Name: "pardon", Trust: 96, Position: dead, Callback: doPardon})
	commandList.PushBack(&cmd{Name: "purge", Trust: 96, Position: dead, Callback: doPurge})
	commandList.PushBack(&cmd{Name: "restore", Trust: 96, Position: dead, Callback: doRestore})
	// commandList.PushBack(&cmd{Name: "rset", Trust: 96, Position: dead, Callback: doRset})
	commandList.PushBack(&cmd{Name: "silence", Trust: 96, Position: dead, Callback: doSilence})
	commandList.PushBack(&cmd{Name: "sla", Trust: 96, Position: dead, Callback: doSlayIncomplete})
	commandList.PushBack(&cmd{Name: "slay", Trust: 96, Position: dead, Callback: doSlay})
	// commandList.PushBack(&cmd{Name: "sset", Trust: 96, Position: dead, Callback: doSset})
	commandList.PushBack(&cmd{Name: "transfer", Trust: 96, Position: dead, Callback: doTransfer})

	commandList.PushBack(&cmd{Name: "at", Trust: 95, Position: dead, Callback: doAt})
	commandList.PushBack(&cmd{Name: "bamfin", Trust: 95, Position: dead, Callback: doBamfin})
	commandList.PushBack(&cmd{Name: "bamfout", Trust: 95, Position: dead, Callback: doBamfout})
	commandList.PushBack(&cmd{Name: "echo", Trust: 95, Position: dead, Callback: doEcho})
	commandList.PushBack(&cmd{Name: "goto", Trust: 95, Position: dead, Callback: doGoto})
	commandList.PushBack(&cmd{Name: "holylight", Trust: 95, Position: dead, Callback: doHolylight})
	commandList.PushBack(&cmd{Name: "invis", Trust: 95, Position: dead, Callback: doInvis})
	commandList.PushBack(&cmd{Name: "log", Trust: 95, Position: dead, Callback: doLog})
	commandList.PushBack(&cmd{Name: "memory", Trust: 95, Position: dead, Callback: doMemory})
	commandList.PushBack(&cmd{Name: "mfind", Trust: 95, Position: dead, Callback: doMfind})
	commandList.PushBack(&cmd{Name: "mstat", Trust: 95, Position: dead, Callback: doMstat})
	commandList.PushBack(&cmd{Name: "mwhere", Trust: 95, Position: dead, Callback: doMwhere})
	commandList.PushBack(&cmd{Name: "ofind", Trust: 95, Position: dead, Callback: doOfind})
	commandList.PushBack(&cmd{Name: "ostat", Trust: 95, Position: dead, Callback: doOstat})
	commandList.PushBack(&cmd{Name: "peace", Trust: 95, Position: dead, Callback: doPeace})
	commandList.PushBack(&cmd{Name: "recho", Trust: 95, Position: dead, Callback: doRecho})
	commandList.PushBack(&cmd{Name: "return", Trust: 95, Position: dead, Callback: doReturn})
	commandList.PushBack(&cmd{Name: "rstat", Trust: 95, Position: dead, Callback: doRstat})
	commandList.PushBack(&cmd{Name: "slookup", Trust: 95, Position: dead, Callback: doSlookup})
	commandList.PushBack(&cmd{Name: "snoop", Trust: 95, Position: dead, Callback: doSnoop})
	commandList.PushBack(&cmd{Name: "switch", Trust: 95, Position: dead, Callback: doSwitch})
	commandList.PushBack(&cmd{Name: "wizhelp", Trust: 95, Position: dead, Callback: doWizhelp})

	commandList.PushBack(&cmd{Name: "immtalk", Trust: 95, Position: dead, Callback: doImmtalk})
}

func loadHelps() {
	// TODO
}

func loadItems() {
	itemFiles, _ := filepath.Glob(fmt.Sprintf("%s%s", gameServer.BasePath, "./data/items/*.json"))

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

	jobs = append(jobs, &job{ID: 1, Name: "Soldier", Abbr: "sol", StartingWeapon: startingGun, PrimeAttribute: applyStrength, SkillAdept: 85, Thac0_00: 18, Thac0_32: 6, MinHitpoints: 11, MaxHitpoints: 15, GainsMana: false})
	jobs = append(jobs, &job{ID: 2, Name: "Merchant", Abbr: "mer", StartingWeapon: startingStaff, PrimeAttribute: applyWisdom, SkillAdept: 95, Thac0_00: 18, Thac0_32: 10, MinHitpoints: 6, MaxHitpoints: 8, GainsMana: true})
	jobs = append(jobs, &job{ID: 3, Name: "Spy", Abbr: "spy", StartingWeapon: startingDagger, PrimeAttribute: applyIntelligence, SkillAdept: 95, Thac0_00: 18, Thac0_32: 12, MinHitpoints: 7, MaxHitpoints: 10, GainsMana: true})
	jobs = append(jobs, &job{ID: 4, Name: "Sailor", Abbr: "sai", StartingWeapon: startingStaff, PrimeAttribute: applyDexterity, SkillAdept: 85, Thac0_00: 18, Thac0_32: 8, MinHitpoints: 8, MaxHitpoints: 8, GainsMana: false})
	jobs = append(jobs, &job{ID: 5, Name: "Guard", Abbr: "law", StartingWeapon: startingHammer, PrimeAttribute: applyConstitution, SkillAdept: 85, Thac0_00: 18, Thac0_32: 8, MinHitpoints: 10, MaxHitpoints: 14, GainsMana: false})

	for _, j := range jobs {
		jobList.PushBack(j)
	}
}

func loadMobs() {
	mobList = list.New()

	mobFiles, _ := filepath.Glob(fmt.Sprintf("%s%s", gameServer.BasePath, "./data/mobs/*.json"))

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
	raceList.PushBack(&race{ID: 1, Name: "West Indies", Adjective: "West Indian", Abbr: "wes", Stats: raceStats{Hitpoints: 100, Mana: 100, Movement: 100, Strength: 12, Intelligence: 12, Wisdom: 12, Charisma: 12, Dexterity: 12, Constitution: 12}})
	raceList.PushBack(&race{ID: 2, Name: "England", Adjective: "English", Abbr: "eng", Stats: raceStats{Hitpoints: 100, Mana: 100, Movement: 100, Strength: 12, Intelligence: 12, Wisdom: 12, Charisma: 12, Dexterity: 12, Constitution: 12}})
	raceList.PushBack(&race{ID: 3, Name: "France", Adjective: "French", Abbr: "fra", Stats: raceStats{Hitpoints: 100, Mana: 100, Movement: 100, Strength: 12, Intelligence: 12, Wisdom: 12, Charisma: 12, Dexterity: 12, Constitution: 12}})
	raceList.PushBack(&race{ID: 4, Name: "Spain", Adjective: "Spanish", Abbr: "spa", Stats: raceStats{Hitpoints: 100, Mana: 100, Movement: 100, Strength: 12, Intelligence: 12, Wisdom: 12, Charisma: 12, Dexterity: 12, Constitution: 12}})
	raceList.PushBack(&race{ID: 5, Name: "Portugal", Adjective: "Portuguese", Abbr: "por", Stats: raceStats{Hitpoints: 100, Mana: 100, Movement: 100, Strength: 12, Intelligence: 12, Wisdom: 12, Charisma: 12, Dexterity: 12, Constitution: 12}})
}

func loadResets() {
	resetFiles, _ := filepath.Glob(fmt.Sprintf("%s%s", gameServer.BasePath, "./data/resets/*.json"))

	for _, resetFile := range resetFiles {
		if strings.HasSuffix(resetFile, "sample.json") {
			continue
		}
		file, err := ioutil.ReadFile(resetFile)
		if err != nil {
			panic(err)
		}

		var res *resetData
		err = json.Unmarshal(file, &res)

		if err != nil {
			panic(err)
		}
		resetList.PushBack(res)
	}
}

func loadRooms() {
	areaFiles, _ := filepath.Glob(fmt.Sprintf("%s%s", gameServer.BasePath, "./data/area/*.json"))

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

		areaList.PushBack(area)

	}

	for e := roomList.Front(); e != nil; e = e.Next() {
		room := e.Value.(*room)
		for _, mob := range room.Mobs {
			mob.Room = room
		}

		for j, x := range room.Exits {
			if x.RoomID > 0 {
				room.Exits[j] = &exit{Dir: x.Dir, Room: getRoom(x.RoomID), RoomID: x.RoomID}
			}
		}
	}
}

func loadShops() {
	shopFiles, _ := filepath.Glob(fmt.Sprintf("%s%s", gameServer.BasePath, "./data/shops/*.json"))

	for _, shopFile := range shopFiles {
		file, err := ioutil.ReadFile(shopFile)
		if err != nil {
			panic(err)
		}

		var list []*shop
		json.Unmarshal(file, &list)

		for _, sh := range list {
			sh.keeper = getMob(sh.KeeperID)
			sh.keeper.Shop = sh
			shopList.PushBack(sh)
		}
	}
}

func loadSkills() {
	skillFiles, _ := filepath.Glob(fmt.Sprintf("%s%s", gameServer.BasePath, "./data/skills/*.json"))

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
	file, err := ioutil.ReadFile(fmt.Sprintf("%s%s", gameServer.BasePath, "./data/socials/socials.json"))
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

func extractMob(m *mob, pull bool) {
	if m.Room == nil {
		return
	}

	if pull {
		m.dieFollower()
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

	if !m.isNPC() {
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
}

func mockObject(name string, id int) *item {
	index := &itemIndex{}

	index.ID = 1
	index.Name = name
	index.Description = ""
	index.ShortDescription = ""
	index.ItemType = 1
	index.ContainedIDs = []int{}
	index.Affected = []*affect{}
	index.ExtraFlags = 0
	index.WearFlags = 0
	index.Weight = 0
	index.Value = 0
	index.Min = 0
	index.Max = 0
	index.Level = 1
	index.Cost = 0
	index.SkillID = 0
	index.Charges = 0

	return createItem(index)
}

func mockPlayer(name string) *mob {
	player := mockMob(name)
	player.Playable = true

	s, c := net.Pipe()
	s.Close()
	player.client = &connection{conn: c}
	gameServer.connections = append(gameServer.connections, *player.client)
	player.Room = mockRoom()
	player.Room.Mobs = append(player.Room.Mobs, player)

	return player
}

func mockAffect(name string) *affect {
	return &affect{affectType: mockSkill(name), bitVector: affectBlind}
}

func mockSkill(name string) *mobSkill {
	return &mobSkill{Skill: &skill{Name: name}, Level: 1}
}

func mockMob(name string) *mob {
	index := &mobIndex{}
	index.ID = 1
	index.Name = name
	index.Password = ""
	index.Description = ""
	index.LongDescription = ""
	index.Title = ""
	index.Affects = nil
	index.AffectedBy = []int{0}
	index.Act = []int{0}

	index.Skills = nil
	index.Inventory = nil
	index.Equipped = nil
	index.RoomID = 0

	index.ExitVerb = ""
	index.Bamfin = ""
	index.Bamfout = ""

	index.Hitpoints = 100
	index.MaxHitpoints = 100
	index.Mana = 100
	index.MaxMana = 100
	index.Movement = 100
	index.MaxMovement = 100

	index.Armor = 1
	index.Hitroll = 1
	index.Damroll = 1

	index.Exp = 0
	index.Level = 1
	index.Alignment = 0
	index.Practices = 0
	index.Gold = 0
	index.Trust = 0

	index.Carrying = 0
	index.CarryMax = 0
	index.CarryWeight = 0
	index.CarryWeightMax = 0

	index.JobID = 1
	index.RaceID = 1
	index.Gender = 0

	index.Attributes = &attributeSet{
		Strength:     12,
		Intelligence: 12,
		Wisdom:       12,
		Dexterity:    12,
		Charisma:     12,
		Constitution: 12,
	}
	index.ModifiedAttributes = &attributeSet{
		Strength:     0,
		Intelligence: 0,
		Wisdom:       0,
		Dexterity:    0,
		Charisma:     0,
		Constitution: 0,
	}

	index.Status = standing
	index.Shop = nil

	index.Playable = false

	m := createMob(index)
	m.Room = mockRoom()
	return m
}

func mockRoom() *room {
	r := &room{Mobs: []*mob{}, Items: []*item{}, Exits: []*exit{}}

	r.Name = "Some Room"
	r.RoomFlags = 0
	r.Light = 0
	r.SectorType = sectorHills
	return r
}

func resetTest() *mob {
	gameServer.connections = []connection{}
	var l list.List
	mobList = &l
	// roomList.Init()

	loadJobs()
	loadRaces()

	player := mockPlayer("player")
	player.Room = mockRoom()
	player.Room.Mobs = append(player.Room.Mobs, player)
	return player
}

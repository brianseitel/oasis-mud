package mud

type command string

const (
	cLook         command = "look"
	cNorth        command = "north"
	cSouth        command = "south"
	cEast         command = "east"
	cWest         command = "west"
	cUp           command = "up"
	cDown         command = "down"
	cGet          command = "get"
	cDrop         command = "drop"
	cWear         command = "wear"
	cRemove       command = "remove"
	cKill         command = "kill"
	cFlee         command = "flee"
	cQuit         command = "quit"
	cInventory    command = "inventory"
	cEquipment    command = "equipment"
	cScore        command = "score"
	cScan         command = "scan"
	cRecall       command = "recall"
	cSkill        command = "skills"
	cTrip         command = "trip"
	cTrain        command = "train"
	cCast         command = "cast"
	cAffect       command = "affects"
	cChat         command = "chat"
	cChatAuction  command = "auction"
	cChatMusic    command = "music"
	cChatQuestion command = "question"
	cChatAnswer   command = "answer"
	cChatImmtalk  command = "immtalk"
	cSay          command = "say"
	cTell         command = "tell"
	cReply        command = "reply"
	cPut          command = "put"
	cGive         command = "give"
	cConsider     command = "consider"
	cSneak        command = "sneak"
	cKick         command = "kick"
	cDisarm       command = "disarm"
	cBackstab     command = "backstab"
	cSteal        command = "steal"
	cPractice     command = "practice"
	cWho          command = "who"
	cWhere        command = "where"
	cHide         command = "hide"
	cFollow       command = "follow"
	cLock         command = "lock"
	cUnlock       command = "unlock"
	cPick         command = "pick"
	cNoop         command = "noop"
)

var commands []command

func init() {
	commands = []command{cLook, cNorth, cSouth, cEast, cWest, cUp, cDown, cWho, cFollow, cWhere, cHide, cGive, cPractice, cSneak, cConsider, cPut, cSay, cTell, cReply, cSteal, cChat, cChatAuction, cChatMusic, cChatQuestion, cChatAnswer, cChatImmtalk, cAffect, cCast, cGet, cDrop, cWear, cRemove, cKill, cKick, cBackstab, cDisarm, cFlee, cInventory, cScore, cFlee, cWear, cEquipment, cScan, cRecall, cQuit, cSkill, cUnlock, cLock, cPick, cTrip, cTrain}
}

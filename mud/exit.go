package mud

const (
	exitDoor = 1 << iota
	exitClosed
	exitLocked
	exitPickProof
)

type exit struct {
	ID          int
	Keyword     string `json:"keyword"`
	Description string `json:"description"`
	Dir         string `json:"direction"`
	Room        *room
	RoomID      int `json:"room_id"`
	Key         int
	Flags       int
}

func (e *exit) hasDoor() bool {
	return hasBit(e.Flags, exitDoor)
}

func (e *exit) isClosed() bool {
	return hasBit(e.Flags, exitClosed)
}

func (e *exit) isOpen() bool {
	return !e.isClosed()
}

func (e *exit) isLocked() bool {
	return hasBit(e.Flags, exitLocked)
}

func (e *exit) isUnlocked() bool {
	return !e.isLocked()
}

func (e *exit) isPickProof() bool {
	return hasBit(e.Flags, exitPickProof)
}

func (e *exit) isPickable() bool {
	return !e.isPickProof()
}

func reverseDirection(dir string) string {
	switch dir {
	case "east":
		return "west"
	case "west":
		return "east"
	case "up":
		return "down"
	case "down":
		return "up"
	case "north":
		return "south"
	case "south":
		return "north"
	default:
		return "oops"
	}
}

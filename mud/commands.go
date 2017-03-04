package mud

type command string

const (
	cLook      command = "look"
	cNorth     command = "north"
	cSouth     command = "south"
	cEast      command = "east"
	cWest      command = "west"
	cUp        command = "up"
	cDown      command = "down"
	cGet       command = "get"
	cDrop      command = "drop"
	cWear      command = "wear"
	cRemove    command = "remove"
	cKill      command = "kill"
	cFlee      command = "flee"
	cQuit      command = "quit"
	cInventory command = "inventory"
	cNoop      command = "noop"
)

var commands []command

func init() {
	commands = []command{cLook, cNorth, cSouth, cEast, cWest, cUp, cDown, cGet, cDrop, cWear, cRemove, cKill, cFlee, cInventory, cQuit}
}

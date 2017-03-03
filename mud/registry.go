package mud

type registry struct {
	rooms   []Room
	players []Player
	items   []Item
}

var Registry registry

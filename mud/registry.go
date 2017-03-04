package mud

type registry struct {
	rooms   map[int]Room
	players []Player
	items   []Item
}

var Registry registry

package mud

type registry struct {
	rooms   map[int]Room
	players []Mob
	mobs    []Mob
	items   []Item
}

var Registry registry

package mud

type status int

const (
	dead status = iota
	incapacitated
	stunned
	sleeping
	sitting
	fighting
	standing
)

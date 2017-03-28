package mud

type status int

const (
	dead status = iota
	mortal
	incapacitated
	stunned
	sleeping
	sitting
	resting
	fighting
	standing
)

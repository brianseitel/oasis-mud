package mud

type bonusStrength struct {
	toHit    int
	toDamage int
	toCarry  int
	wield    int
}

type bonusIntelligence struct {
	learn int
}

type bonusWisdom struct {
	practice int
}

type bonusDexterity struct {
	defensive int
}

type bonusConstitution struct {
	hitpoints int
	shock     int
}

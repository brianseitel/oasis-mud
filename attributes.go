package main

const (
	STRENGTH = iota
	HEALTH
	AGILITY
	MAXHITPOINTS
	ACCURACY
	DODGING
	STRIKEDAMAGE
	DAMAGEABSORB
	HPREGEN
	NUMATTRIBUTES
)

type ItemAttributeSet struct {
	Strength     int
	Health       int
	Agility      int
	MaxHitPoints int
	Accuracy     int
	Dodging      int
	StrikeDamage int
	DamageAbsorb int
	HPRegen      int
}

package mud

import "fmt"

type affect struct {
	affectType *mobSkill
	duration   uint
	location   int
	modifier   int
}

const (
	applyAC = iota
	applyStrength
	applyDexterity
	applyIntelligence
	applyWisdom
	applyConstitution
	applyCharisma
	applySex
	applyClass
	applyLevel
	applyMana
	applyHitpoints
	applyMovement
	applyGold
	applyExp
	applyHitroll
	applyDamroll
	applySavingParalysis
	applySavingRod
	applySavingPetrify
	applySavingBreath
	applySavingSpell
	applyNone
)

func affectModify(player *mob, paf *affect, add bool) {
	var mod int

	mod = paf.modifier

	if add {
		player.addAffect(paf)
	} else {
		player.removeAffect(paf)
	}

	if player.isNPC() {
		return
	}

	switch paf.location {
	default:
		fmt.Printf("affectModify: unknown location %s\r\n", paf.location)
		return

	case applySex:
	case applyClass:
	case applyLevel:
	case applyGold:
	case applyExp:
	case applySavingParalysis:
	case applySavingRod:
	case applySavingPetrify:
	case applySavingBreath:
	case applySavingSpell:
	case applyNone:
		break
	case applyStrength:
		player.ModifiedAttributes.Strength += mod
		break
	case applyDexterity:
		player.ModifiedAttributes.Dexterity += mod
		break
	case applyIntelligence:
		player.ModifiedAttributes.Intelligence += mod
		break
	case applyWisdom:
		player.ModifiedAttributes.Wisdom += mod
		break
	case applyCharisma:
		player.ModifiedAttributes.Charisma += mod
		break
	case applyConstitution:
		player.ModifiedAttributes.Constitution += mod
		break
	case applyMana:
		player.MaxMana += mod
		break
	case applyHitpoints:
		player.Hitpoints += mod
		break
	case applyMovement:
		player.Movement += mod
		break
	case applyAC:
		player.Armor += mod
		break
	case applyHitroll:
		player.Hitroll += mod
		break
	case applyDamroll:
		player.Damroll += mod
		break
	}

	return
}

func isAffected(player *mob, aff *affect) bool {
	for _, affect := range player.Affects {
		if affect.affectType.Skill.Name == aff.affectType.Skill.Name {
			return true
		}
	}
	return false
}

package mud

import "fmt"

type affect struct {
	affectType *mobSkill
	duration   int
	location   int
	modifier   int
	bitVector  int
}

const (
	applyArmorClass = iota
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

const (
	affectBlind           = 1
	affectInvisible       = 2
	affectDetectEvil      = 4
	affectDetectInvisible = 8
	affectDetectMagic     = 16
	affectDetectHidden    = 32
	affectHold            = 64 /* Unused	*/
	affectSanctuary       = 128
	affectFaerieFire      = 256
	affectInfrared        = 512
	affectCurse           = 1024
	affectFlaming         = 2048 /* Unused	*/
	affectPoison          = 4096
	affectProtect         = 8192
	affectParalysis       = 16384 /* Unused	*/
	affectSneak           = 32768
	affectHide            = 65536
	affectSleep           = 131072
	affectCharm           = 262144
	affectFlying          = 524288
	affectPassDoor        = 1048576
)

func affectModify(player *mob, paf *affect, add bool) {
	var mod int

	mod = paf.modifier

	if add {
		player.Affects = append(player.Affects, paf)
		player.AffectedBy = setBit(player.AffectedBy, paf.bitVector)
	} else {
		for j, aff := range player.Affects {
			if paf == aff {
				player.Affects = append(player.Affects[0:j], player.Affects[j+1:]...)
				player.notify("%s", paf.affectType.Skill.MessageOff)
				player.AffectedBy = removeBit(player.AffectedBy, paf.bitVector)
				return
			}
		}
	}

	if player.isNPC() {
		return
	}

	switch paf.location {
	default:
		fmt.Printf("affectModify: unknown location %d\r\n", paf.location)
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
		player.Mana += mod
		break
	case applyHitpoints:
		player.Hitpoints += mod
		break
	case applyMovement:
		player.Movement += mod
		break
	case applyArmorClass:
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

func getAffectNames(bit int) []string {
	var results []string

	if bit == 0 {
		return results
	}

	if hasBit(affectBlind, bit) {
		results = append(results, "blind")
	}
	if hasBit(affectInvisible, bit) {
		results = append(results, "invis")
	}
	if hasBit(affectDetectEvil, bit) {
		results = append(results, "detect evil")
	}
	if hasBit(affectDetectInvisible, bit) {
		results = append(results, "detect invis")
	}
	if hasBit(affectDetectMagic, bit) {
		results = append(results, "detect magic")
	}
	if hasBit(affectDetectHidden, bit) {
		results = append(results, "hidden")
	}
	if hasBit(affectHold, bit) {
		results = append(results, "hold")
	}
	if hasBit(affectSanctuary, bit) {
		results = append(results, "sanctuary")
	}
	if hasBit(affectFaerieFire, bit) {
		results = append(results, "faerie fire")
	}
	if hasBit(affectInfrared, bit) {
		results = append(results, "infrared")
	}
	if hasBit(affectFlaming, bit) {
		results = append(results, "flaming")
	}
	if hasBit(affectPoison, bit) {
		results = append(results, "poison")
	}
	if hasBit(affectProtect, bit) {
		results = append(results, "protect")
	}
	if hasBit(affectParalysis, bit) {
		results = append(results, "paralysis")
	}
	if hasBit(affectSneak, bit) {
		results = append(results, "sneak")
	}
	if hasBit(affectHide, bit) {
		results = append(results, "hide")
	}
	if hasBit(affectSleep, bit) {
		results = append(results, "sleep")
	}
	if hasBit(affectCharm, bit) {
		results = append(results, "charm")
	}
	if hasBit(affectFlying, bit) {
		results = append(results, "flying")
	}
	if hasBit(affectPassDoor, bit) {
		results = append(results, "pass door")
	}

	return results
}

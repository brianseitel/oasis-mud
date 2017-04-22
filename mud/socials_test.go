package mud

import "testing"

func TestCheckSocial(t *testing.T) {
	player := resetTest()
	loadSocials()

	// fake
	if checkSocial(player, "not a real social", "") == true {
		t.Error("Found invalid social")
	}

	// can't do it if dead
	player.Status = dead
	checkSocial(player, "dance", "")

	// can't do it if stunned or incapacitated
	player.Status = incapacitated
	checkSocial(player, "dance", "")

	player.Status = mortal
	checkSocial(player, "dance", "")

	// can only snore while asleep
	player.Status = sleeping
	checkSocial(player, "dance", "")

	// snore
	checkSocial(player, "dance", "")

	player.Status = standing

	// dance with someone else not here
	checkSocial(player, "dance", "bob")

	vic := mockPlayer("vic")
	player.Room = mockRoom()
	player.Room.Mobs = append(player.Room.Mobs, vic)
	player.Room.Mobs = append(player.Room.Mobs, player)

	checkSocial(player, "dance", "vic")

	// dance with yourself
	checkSocial(player, "dance", player.Name)

	vic.Playable = false
	for i := 0; i < 1000; i++ {
		checkSocial(player, "dance", "vic")
	}
}

func TestDoSocials(t *testing.T) {
	player := resetTest()

	doSocials(player, "")
}

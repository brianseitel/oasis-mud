package mud

type job struct {
	ID             uint
	Name           string
	Abbr           string
	StartingWeapon int
	PrimeAttribute int
	SkillAdept     int
	Thac0_00       int
	Thac0_32       int
	MinHitpoints   int
	MaxHitpoints   int
	GainsMana      bool
}

func getJob(id uint) *job {
	for e := jobList.Front(); e != nil; e = e.Next() {
		j := e.Value.(*job)
		if j.ID == id {
			return j
		}
	}
	return nil
}

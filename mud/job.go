package mud

type job struct {
	ID   uint
	Name string
	Abbr string
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

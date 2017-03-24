package mud

import "container/list"

type job struct {
	ID   uint
	Name string
	Abbr string
}

var jobList list.List

func newJobDatabase() {
	jobs := make(map[string]string)
	jobs["war"] = "Warrior"
	jobs["mag"] = "Mage"
	jobs["cle"] = "Cleric"
	jobs["thi"] = "Thief"
	jobs["ran"] = "Ranger"
	jobs["bar"] = "Bard"

	for abbr, name := range jobs {
		j := &job{Name: name, Abbr: abbr}

		jobList.PushBack(j)
	}
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

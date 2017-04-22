package mud

import "testing"

func TestGetJob(t *testing.T) {
	loadJobs()

	job := getJob(1)

	if job.ID != 1 || job.Name != "Warrior" {
		t.Error("Did not find job.")
	}

	job = getJob(99)

	if job != nil {
		t.Error("Found fake job when it shouldn't")
	}
}

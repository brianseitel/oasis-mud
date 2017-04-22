package mud

import "testing"

func TestDump(t *testing.T) {
	t.Skip("Skipping dump() test")
}

func TestDD(t *testing.T) {
	t.Skip("Skipping dd() test")
}

func TestIsSameGroup(t *testing.T) {
	p1 := &mob{}
	p2 := &mob{}

	if isSameGroup(p1, p2) {
		t.Error("p1 and p2 are not in same group")
	}

	p1.leader = p2

	if isSameGroup(p1, p2) == false {
		t.Error("p1 and p2 are in same group")
	}

	p1.leader = nil
	p2.leader = p1

	if isSameGroup(p1, p2) == false {
		t.Error("p1 and p2 are in same group")
	}
}

func TestMatchesSubject(t *testing.T) {
	phrase := "There's a million things I haven't done."

	if matchesSubject(phrase, "million") == false {
		t.Error("Could not match subject")
	}

	if matchesSubject(phrase, "haven't") == false {
		t.Error("Could not match subject")
	}

	if matchesSubject(phrase, "Aaron Burr, Sir") == true {
		t.Error("Found subject it shouldn't have found")
	}
}

func TestTransferItem(t *testing.T) {
	var from []*item
	var to []*item

	from = append(from, &item{ID: 1})

	from, to = transferItem(0, from, to)

	if len(from) != 0 || len(to) != 1 {
		t.Error("Failed to transfer item")
	}
}

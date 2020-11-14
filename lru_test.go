package lru

import (
	"testing"
)

func TestMain(t *testing.T) {
	d := Lru(4)
	if d.Len() != 0 {
		t.FailNow()
	}
	d.Set(1, 1)
	d.Set(2, 2)
	d.Set(3, 3)
	d.Set(4, 4) // Snapshot: [4, 3, 2, 1]
	if d.Len() != 4 {
		t.FailNow()
	}
	d.Set(5, 5) // Snapshot: [5, 4, 3, 2]
	if d.Len() != 4 {
		t.FailNow()
	}
	if _, fit := d.Get(1); fit {
		t.FailNow()
	}
	d.Get(2)    // Snapshot: [2, 5, 4, 3]
	d.Set(6, 6) // Snapshot: [6, 2, 5, 4]
	if _, fit := d.Get(3); fit {
		t.FailNow()
	}
	if _, fit := d.Get(6); !fit {
		t.FailNow()
	}
	if _, fit := d.Get(2); !fit {
		t.FailNow()
	}
	if _, fit := d.Get(5); !fit {
		t.FailNow()
	}
	if _, fit := d.Get(4); !fit {
		t.FailNow()
	}
}

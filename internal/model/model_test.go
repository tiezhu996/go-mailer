package model

import "testing"

func TestHigherPriority(t *testing.T) {
	a := &Message{ID: "a", Priority: 5}
	b := &Message{ID: "b", Priority: 1}
	if !HigherPriority(b, a) {
		t.Fatal("b should be higher priority than a")
	}
	if HigherPriority(a, b) {
		t.Fatal("a should not be higher priority than b")
	}
	c := &Message{ID: "c", Priority: 1}
	if !HigherPriority(b, c) {
		t.Fatal("same priority should fall back to ID order")
	}
}

func TestSortByPriority(t *testing.T) {
	in := []*Message{{ID: "c", Priority: 2}, {ID: "a", Priority: 1}, {ID: "b", Priority: 1}}
	got := SortByPriority(in)
	want := []string{"a", "b", "c"}
	for i, id := range want {
		if got[i].ID != id {
			t.Fatalf("order=%v want %v", got, want)
		}
	}
}

func TestBuildBatchesFresh(t *testing.T) {
	in := []*Message{{ID: "1"}, {ID: "2"}, {ID: "3"}}
	b := BuildBatches(in, 2)
	if len(b) != 2 {
		t.Fatalf("len=%d want 2", len(b))
	}
	b[0][0] = &Message{ID: "x"}
	if in[0].ID != "1" {
		t.Fatal("mutating batch corrupted input")
	}
}

func TestMergeSummary(t *testing.T) {
	got := MergeSummary(Summary{Sent: 1}, Summary{Sent: 2, Failed: 3})
	if got.Sent != 3 || got.Failed != 3 {
		t.Fatalf("got %+v", got)
	}
}

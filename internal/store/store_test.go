package store

import (
	"testing"

	"mailer/internal/model"
)

func TestEnqueueGet(t *testing.T) {
	s := New()
	if err := s.Enqueue(&model.Message{ID: "1", Status: model.StatusQueued}); err != nil {
		t.Fatal(err)
	}
	if err := s.Enqueue(&model.Message{ID: "1"}); err != ErrAlreadyExists {
		t.Fatalf("dup err=%v", err)
	}
	if _, err := s.Get("1"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Get("nope"); err != ErrNotFound {
		t.Fatalf("err=%v", err)
	}
}

func TestPendingSorted(t *testing.T) {
	s := New()
	for _, m := range []*model.Message{
		{ID: "low", Priority: 9, Status: model.StatusQueued},
		{ID: "high", Priority: 0, Status: model.StatusQueued},
		{ID: "mid", Priority: 5, Status: model.StatusQueued},
		{ID: "done", Priority: 0, Status: model.StatusSent},
	} {
		if err := s.Enqueue(m); err != nil {
			t.Fatal(err)
		}
	}
	got := s.Pending()
	if len(got) != 3 {
		t.Fatalf("len=%d want 3", len(got))
	}
	want := []string{"high", "mid", "low"}
	for i, id := range want {
		if got[i].ID != id {
			t.Fatalf("pending=%v want %v", got, want)
		}
	}
}

func TestPendingIDsFresh(t *testing.T) {
	s := New()
	for _, id := range []string{"b", "a"} {
		if err := s.Enqueue(&model.Message{ID: id, Status: model.StatusQueued}); err != nil {
			t.Fatal(err)
		}
	}
	first := s.PendingIDs()
	first[0] = "x"
	after := s.PendingIDs()
	if after[0] != "b" {
		t.Fatalf("PendingIDs after=%v want [b a]", after)
	}
}

func TestMarkSentFailed(t *testing.T) {
	s := New()
	if err := s.Enqueue(&model.Message{ID: "1", Status: model.StatusQueued}); err != nil {
		t.Fatal(err)
	}
	if err := s.MarkSent("1"); err != nil {
		t.Fatal(err)
	}
	m, _ := s.Get("1")
	if m.Status != model.StatusSent {
		t.Fatalf("status=%q", m.Status)
	}
	if err := s.MarkFailed("nope"); err != ErrNotFound {
		t.Fatalf("err=%v", err)
	}
}

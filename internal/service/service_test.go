package service

import (
	"errors"
	"testing"

	"mailer/internal/model"
	"mailer/internal/store"
)

func TestSubmitGet(t *testing.T) {
	s := store.New()
	svc := New(s, 2)
	m, err := svc.Submit("1", "a@b", "hi", 1)
	if err != nil || m.Status != model.StatusQueued {
		t.Fatalf("m=%v err=%v", m, err)
	}
	if _, err := svc.Submit("2", "a@b", "", 1); !errors.Is(err, ErrEmptyBody) {
		t.Fatalf("empty body err=%v", err)
	}
	if _, err := svc.Get("1"); err != nil {
		t.Fatal(err)
	}
}

func TestPrepareBatchesOrder(t *testing.T) {
	s := store.New()
	svc := New(s, 2)
	for _, m := range []*model.Message{
		{ID: "low", Priority: 9},
		{ID: "high", Priority: 0},
		{ID: "mid", Priority: 5},
	} {
		if _, err := svc.Submit(m.ID, "a@b", "p", m.Priority); err != nil {
			t.Fatal(err)
		}
	}
	batches, err := svc.PrepareBatches()
	if err != nil {
		t.Fatal(err)
	}
	// high, mid, low 分成 [high,mid] 和 [low]
	if len(batches) != 2 {
		t.Fatalf("batches=%d want 2", len(batches))
	}
	if batches[0][0].ID != "high" || batches[0][1].ID != "mid" {
		t.Fatalf("batch0=%v,%v", batches[0][0].ID, batches[0][1].ID)
	}
	if batches[1][0].ID != "low" {
		t.Fatalf("batch1=%v", batches[1][0].ID)
	}
}

func TestMarkSentWraps(t *testing.T) {
	s := store.New()
	svc := New(s, 2)
	if err := svc.MarkSent("nope"); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("errors.Is=false err=%v", err)
	}
}

func TestMarkFailedWraps(t *testing.T) {
	s := store.New()
	svc := New(s, 2)
	if err := svc.MarkFailed("nope"); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("errors.Is=false err=%v", err)
	}
}

package worker

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"mailer/internal/model"
	"mailer/internal/service"
	"mailer/internal/store"
)

type flakyDispatcher struct {
	mu    sync.Mutex
	fail  map[string]int // id -> 前 N 次失败
	calls map[string]int
}

func newFlaky(fail map[string]int) *flakyDispatcher {
	return &flakyDispatcher{fail: fail, calls: map[string]int{}}
}

func (d *flakyDispatcher) Send(ctx context.Context, m *model.Message) error {
	d.mu.Lock()
	d.calls[m.ID]++
	n := d.calls[m.ID]
	f := d.fail[m.ID]
	d.mu.Unlock()
	if n <= f {
		return fmt.Errorf("flaky fail %s", m.ID)
	}
	return nil
}

func TestRunSummary(t *testing.T) {
	st := store.New()
	svc := service.New(st, 2)
	for i := 0; i < 10; i++ {
		if _, err := svc.Submit(fmt.Sprintf("m%d", i), "a@b", "p", i%3); err != nil {
			t.Fatal(err)
		}
	}
	// m0 永久失败（重试 1 次仍失败），其余成功
	fail := map[string]int{"m0": 999}
	pool := New(st, svc, newFlaky(fail), 4, 1)
	sum := pool.Run(context.Background())
	if sum.Sent != 9 {
		t.Fatalf("Sent=%d want 9", sum.Sent)
	}
	if sum.Failed != 1 {
		t.Fatalf("Failed=%d want 1", sum.Failed)
	}
	// 状态核对
	for i := 0; i < 10; i++ {
		m, _ := st.Get(fmt.Sprintf("m%d", i))
		if i == 0 && m.Status != model.StatusFailed {
			t.Fatalf("m0 status=%q want failed", m.Status)
		}
		if i != 0 && m.Status != model.StatusSent {
			t.Fatalf("m%d status=%q want sent", i, m.Status)
		}
	}
}

func TestRunPriorityOrder(t *testing.T) {
	st := store.New()
	svc := service.New(st, 1) // 每批 1 条，观察发送顺序
	var mu sync.Mutex
	var sentOrder []string
	d := &recordingDispatcher{onSend: func(id string) { mu.Lock(); sentOrder = append(sentOrder, id); mu.Unlock() }}
	// 优先级乱序提交
	priorities := map[string]int{"low": 9, "high": 0, "mid": 5, "top": -1}
	for id, p := range priorities {
		if _, err := svc.Submit(id, "a@b", "p", p); err != nil {
			t.Fatal(err)
		}
	}
	pool := New(st, svc, d, 1, 0)
	pool.Run(context.Background())
	want := []string{"top", "high", "mid", "low"}
	if len(sentOrder) != len(want) {
		t.Fatalf("order len=%d want %d: %v", len(sentOrder), len(want), sentOrder)
	}
	for i, id := range want {
		if sentOrder[i] != id {
			t.Fatalf("order=%v want %v", sentOrder, want)
		}
	}
}

type recordingDispatcher struct {
	onSend func(id string)
}

func (d *recordingDispatcher) Send(ctx context.Context, m *model.Message) error {
	if d.onSend != nil {
		d.onSend(m.ID)
	}
	return nil
}

package worker

import (
	"context"
	"sync"

	"mailer/internal/model"
	"mailer/internal/service"
	"mailer/internal/store"
)

type Dispatcher interface {
	Send(ctx context.Context, m *model.Message) error
}

type Pool struct {
	store   *store.Store
	svc     *service.Service
	disp    Dispatcher
	workers int
	retries int
}

func New(s *store.Store, svc *service.Service, d Dispatcher, workers, retries int) *Pool {
	if workers <= 0 {
		workers = 1
	}
	if retries < 0 {
		retries = 0
	}
	return &Pool{store: s, svc: svc, disp: d, workers: workers, retries: retries}
}

func (p *Pool) Run(ctx context.Context) model.Summary {
	batches, err := p.svc.PrepareBatches()
	if err != nil {
		return model.Summary{}
	}

	var wg sync.WaitGroup
	ch := make(chan []*model.Message, len(batches))
	results := make(chan model.Summary, len(batches))

	go func() {
		defer close(ch)
		for _, b := range batches {
			select {
			case <-ctx.Done():
				return
			case ch <- b:
			}
		}
	}()

	for i := 0; i < p.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for batch := range ch {
				var local model.Summary
				for _, m := range batch {
					if err := p.disp.Send(ctx, m); err != nil {
						p.retry(m, &local)
						continue
					}
					if err := p.store.MarkSent(m.ID); err != nil {
						local.Skipped++
						continue
					}
					local.Sent++
				}
				results <- local
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var sum model.Summary
	for s := range results {
		sum = model.MergeSummary(sum, s)
	}
	return sum
}

func (p *Pool) retry(m *model.Message, sum *model.Summary) {
	for i := 0; i < p.retries; i++ {
		if err := p.disp.Send(context.Background(), m); err == nil {
			if err := p.store.MarkSent(m.ID); err != nil {
				sum.Skipped++
				return
			}
			sum.Sent++
			return
		}
	}
	if err := p.store.MarkFailed(m.ID); err != nil {
		sum.Skipped++
		return
	}
	sum.Failed++
}

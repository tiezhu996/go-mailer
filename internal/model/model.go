package model

import "sort"

type Status string

const (
	StatusQueued  Status = "queued"
	StatusSending Status = "sending"
	StatusSent    Status = "sent"
	StatusFailed  Status = "failed"
)

type Message struct {
	ID       string
	To       string
	Body     string
	Priority int // 数值越小优先级越高
	Attempts int
	Status   Status
}

type Summary struct {
	Sent    int
	Failed  int
	Skipped int
}

// HigherPriority 报告 a 是否应排在 b 之前（先按 Priority，再按 ID）。
func HigherPriority(a, b *Message) bool {
	if a.Priority != b.Priority {
		return a.Priority < b.Priority
	}
	return a.ID < b.ID
}

// SortByPriority 原地按优先级排序，返回传入的切片。
func SortByPriority(msgs []*Message) []*Message {
	sort.SliceStable(msgs, func(i, j int) bool { return HigherPriority(msgs[i], msgs[j]) })
	return msgs
}

// BuildBatches 把 msgs 切成每组最多 size 个的独立批次。
func BuildBatches(msgs []*Message, size int) [][]*Message {
	if size <= 0 {
		size = 1
	}
	out := make([][]*Message, 0, (len(msgs)+size-1)/size)
	for i := 0; i < len(msgs); i += size {
		end := i + size
		if end > len(msgs) {
			end = len(msgs)
		}
		b := make([]*Message, end-i)
		copy(b, msgs[i:end])
		out = append(out, b)
	}
	return out
}

// MergeSummary 把 src 累加到 dst。
func MergeSummary(dst Summary, src Summary) Summary {
	dst.Sent += src.Sent
	dst.Failed += src.Failed
	return dst
}

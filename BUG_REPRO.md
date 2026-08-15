# BUG_REPRO

## Bug 是什么
优先级比较方向写反、Pending 不再按优先级排序、并发汇总时去掉锁且把 WaitGroup.Add 放进 goroutine 内，导致优先级顺序错误、发送计数丢失并有数据竞争。

## 如何触发
`go test -race ./...`，或 `go test ./internal/model ./internal/service ./internal/worker`。

## 错误信息
- TestHigherPriority / TestSortByPriority / TestPrepareBatchesOrder 失败（顺序相反）。
- `go test -race` 报 data race；发送成功计数可能小于实际值。

# BUG_REPRO

## Bug 是什么
BuildBatches / PendingIDs / PrepareBatches 返回共享底层数组的子切片，worker 又对批次做 `batch[:len(batch)-1]` 切掉最后一条，导致消息漏发与列表互相串改。

## 如何触发
`go test ./...`，或 `go test ./internal/model ./internal/store ./internal/service ./internal/worker`。

## 错误信息
- TestBuildBatchesFresh / TestPendingIDsFresh 失败（切片共享底层数组）。
- TestRunSummary / TestRunPriorityOrder 失败（Sent 计数不对、优先级顺序为空）。

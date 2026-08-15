# BUG_REPRO

## Bug 是什么
service 包装错误用 `%v` 丢掉错误链，store 返回普通错误而非哨兵错误，MergeSummary 丢掉 Skipped，worker 重试成功后不落库，导致 errors.Is 失效且重试成功消息状态不对。

## 如何触发
`go test ./...`，或 `go test ./internal/service ./internal/worker ./internal/store`。

## 错误信息
- TestMarkSentWraps / TestMarkFailedWraps 失败（errors.Is 为 false）。
- TestRetrySuccessStatus 失败（f1 status="queued" want sent）。

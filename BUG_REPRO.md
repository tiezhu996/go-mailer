# BUG_REPRO

## Bug 是什么
worker.retry 在函数顶部 `defer p.store.MarkFailed(m.ID)`，即使重试成功、已经 MarkSent，返回时 defer 仍把消息覆盖为 failed。

## 如何触发
`go test ./internal/worker -run TestRetrySuccessStatus`。

## 错误信息
`f1 status="failed" want sent`

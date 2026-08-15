# BUG_REPRO

## Bug 是什么
Store 构造时未初始化 dedup map，service.Submit 内部走 store.Enqueue，首次向 nil map 写入触发 panic。

## 如何触发
`go test ./internal/store -run TestEnqueueGet`，或调用 service.Submit。

## 错误信息
`panic: assignment to entry in nil map`

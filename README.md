# mailer

一个用 Go 写的内存邮件外发队列示例服务，演示 `handler → service → repository → model` 分层、出队投递与并发 worker 池。

## 功能

- 提交 / 查询邮件，按提交顺序返回待发 ID 列表
- 出队投递，成功后标记为 sent 并累计计数
- 并发 worker 池批量发送，支持 context 取消

## 目录结构

```
cmd/mailer/          程序入口
internal/config/     环境变量配置
internal/model/      模型定义与纯工具函数
internal/store/      内存存储（map + 顺序 + sent 计数）
internal/service/    业务逻辑（提交 / 查询 / 出队 / 标记）
internal/worker/     并发发送 worker 池
```

## 运行与测试

```bash
go build ./...          # 编译
go test ./...           # 全量测试
go run ./cmd/mailer     # 启动
```

## 环境变量

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `MAILER_WORKERS` | worker 数量 | `2` |

## 技术栈

- Go 1.22

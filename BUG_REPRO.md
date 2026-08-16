# BUG_REPRO

## Bug 是什么
工单回复把关闭判断写反；回复后状态推进错误；按 ID+用户查询漏了用户过滤；工单状态默认值写错。

## 如何触发
`cd backend && go test ./...`

## 错误信息
- service.TestTicketReplyChain 失败：待处理工单无法回复。

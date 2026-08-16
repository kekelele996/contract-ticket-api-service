# BUG_REPRO

## Bug 是什么
合同按用户查询漏了用户过滤；合同/工单列表排序写反；合同状态默认值写错。

## 如何触发
`cd backend && go test ./...`

## 错误信息
- service.TestContractListByUser 失败：用他人 ID 仍能查到合同。

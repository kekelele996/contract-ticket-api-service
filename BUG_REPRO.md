# BUG_REPRO

## Bug 是什么
合同提交的状态流转判断写反；签署时目标状态写成草稿；按 ID+用户查询漏了用户过滤；合同状态默认值写错。

## 如何触发
`cd backend && go test ./...`

## 错误信息
- service.TestContractSubmitAndSign 失败：草稿无法提交、签署无法完成。

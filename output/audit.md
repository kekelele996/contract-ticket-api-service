# 审计报告

## 项目信息

- 项目目录：`/Users/gaobo/repositories/gitlab/评审项目/0-1代码生成提示词/golang-改编提示词/生活服务主题项目提示词/cy-397`
- 项目类型：纯后端服务「合同模板生成与法律工单 API」
- 技术栈：Go 1.22 + Gin + GORM + MySQL 8.0 + JWT + wkhtmltopdf
- 提示词文件：`/Users/gaobo/repositories/gitlab/评审项目/0-1代码生成提示词/golang-改编提示词/生活服务主题项目提示词/cy-397/cy-397.md`
- 提示词 SHA-256：`ca4a5efdf01a7003a301d55b8adf89edbe4065a80d82bc8d3a980add265f8480`

审计过程中未修改提示词 `.md` 文件。

## 结论

最终状态：**通过**。后端可正常构建、测试、部署并完成全部主要 API 流程；发现并修复 1 个真实问题。

修复代码提交：`65d8d7d4d167d1dfc9017421d6766ade1bf0a2c2`（`fix: 修复空模板变量持久化为 null 的问题`）。

## 验证过程与证据

### 1. 静态检查

在 `backend/` 下执行：

```bash
go vet ./...
go build ./...
go test ./...
```

结果：全部通过。

在项目根目录执行：

```bash
docker compose --env-file .env config --quiet
```

结果：退出码 0，Compose 配置有效。

### 2. 运行时验证

执行：

```bash
docker compose --env-file .env up -d --build --wait
```

结果：

- `contractapi-mysql` 进入 `healthy`。
- `contractapi-backend` 成功启动并监听 `0.0.0.0:19412 -> 8080`。
- `GET /healthz` 返回 `200 OK`，响应体 `{"status":"ok"}`。

已验证的主要 API 流程：

| 模块 | 验证项 | 结果 |
| --- | --- | --- |
| 认证 | 注册 / 重复注册 / 登录 / 无 Token / 非法 Token | 通过 |
| 模板 | 列表 / 详情 / 收藏 / 收藏列表 / 取消收藏 / 新建 / 更新 / 删除 / code 冲突 | 通过 |
| 合同 | 填充生成 / 列表 / 详情 / 提交待签署 / 签署 / 过期 / 签署方列表 / 状态筛选 | 通过 |
| PDF | 合同导出 `application/pdf`，文件头 `%PDF-1.4`，共 1 页 | 通过 |
| 工单 | 创建 / 列表 / 详情 / 用户回复 / 律师回复 / 回复列表 / 状态流转 / 非法类型 | 通过 |
| FAQ | 关键词搜索 / 分类搜索 / 详情 / 新建 / 更新 / 删除 | 通过 |

关键错误处理验证：

- 缺少必填变量生成合同：`400`，`{"code":40001,"message":"missing required variables: ..."}`。
- 不存在模板生成合同：`404`，`{"code":40400,"message":"template not found"}`。
- 草稿状态直接签署：`422`，`{"code":42200,"message":"cannot sign contract in status \"draft\""}`。
- 非法合同状态筛选：`400`，`{"code":40001,"message":"invalid contract status"}`。
- 非法工单类型：`400`，`{"code":40001,"message":"invalid ticket type"}`。
- 非法工单状态流转：`422`。
- 重复模板 code：`409`。
- 重复收藏：`409`。
- 不存在收藏：`404`。
- 非法 Token：`401`，`{"code":40101,"message":"invalid or expired token"}`。

### 3. 发现的问题

#### 问题 1：空模板变量序列化不一致

- 位置：`backend/internal/model/model_types.go` 的 `TemplateVariables.Value()`。
- 现象：创建模板时传入 `"variables": []`，接口创建响应显示 `[]`，但 MySQL 实际存储为 JSON `null`，再次读取列表/详情时返回 `null`。
- 根因：原实现使用 `append(TemplateVariables(nil), v...)` 生成排序副本；当 `v` 为非 nil 空切片时，`append` 返回 nil 切片，`json.Marshal(nil)` 得到 `null`。
- 影响：空变量模板的变量字段持久化值与 API 响应不一致，属于真实数据序列化缺陷。
- 修复：改为 `make(TemplateVariables, len(v))` + `copy`，保证非 nil 空切片序列化为 `[]`。
- 回归测试：新增 `backend/internal/model/model_types_test.go`，覆盖 nil、非 nil 空切片、非空排序三种场景。

修复后运行时验证：

- 创建 `variables: []` 的模板，接口返回 `variables: []`。
- 再次查询该模板，接口仍返回 `variables: []`。
- MySQL 中实际存储为 `[]`，`JSON_TYPE` 为 `ARRAY`。

### 4. 清理

执行：

```bash
docker compose --env-file .env down -v --remove-orphans
```

验证：

- `lsof -nP -iTCP:19412 -sTCP:LISTEN`：无监听。
- `docker compose ps`：无运行服务。
- `docker ps -a --filter name=contractapi`：无相关容器。
- `docker volume ls --filter name=contractapi`：无相关卷。
- `docker network ls --filter name=contractapi`：无相关网络。

## 修复内容

- 修改文件：`backend/internal/model/model_types.go`
- 新增文件：`backend/internal/model/model_types_test.go`

修复后已重新执行：

```bash
go vet ./...
go build ./...
go test ./...
docker compose --env-file .env config --quiet
```

均通过，并重新构建容器完成运行时冒烟验证。
